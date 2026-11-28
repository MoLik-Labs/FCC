//! Experimental read-only MCP stdio server. No HTTP transport or sandbox.
use crate::workspace::Workspace;
use serde_json::{json, Value};
use std::collections::VecDeque;
use std::io::{self, BufRead, Write};
use std::time::{Duration, Instant};

const VERSION: &str = "2025-06-18";
const MAX_FRAME: usize = 256 * 1024;
const MAX_TOOL_CALLS: usize = 120;
const WINDOW: Duration = Duration::from_secs(60);

#[derive(PartialEq, Eq)]
enum Phase {
    Initial,
    AwaitingReady,
    Ready,
}

struct Server<'a> {
    workspace: &'a Workspace,
    phase: Phase,
    calls: VecDeque<Instant>,
}

fn error(id: Value, code: i32, message: &str) -> Value {
    json!({"jsonrpc":"2.0", "id":id, "error":{"code":code, "message":message}})
}
fn success(id: Value, result: Value) -> Value {
    json!({"jsonrpc":"2.0", "id":id, "result":result})
}
fn tool_result(text: &str, failed: bool) -> Value {
    json!({"content":[{"type":"text", "text":text}], "isError":failed})
}

fn tool_specs() -> Value {
    json!(["list_files", "read_file"].into_iter().map(|name| json!({
        "name":name,
        "description":if name == "read_file" { "Read a UTF-8 file in the selected workspace" } else { "List entries in a workspace-relative directory" },
        "inputSchema":{"type":"object", "properties":{"path":{"type":"string"}}, "required":["path"], "additionalProperties":false},
        "annotations":{"readOnlyHint":true, "destructiveHint":false, "openWorldHint":false}
    })).collect::<Vec<_>>())
}

impl Server<'_> {
    fn handle(&mut self, message: Value) -> Option<Value> {
        let Some(object) = message.as_object() else {
            return Some(error(Value::Null, -32600, "Invalid request"));
        };
        let id = object.get("id").cloned();
        let valid_id = id
            .as_ref()
            .is_none_or(|v| v.is_string() || v.as_i64().is_some() || v.as_u64().is_some());
        if !valid_id {
            return Some(error(Value::Null, -32600, "Invalid request ID"));
        }
        let method = object.get("method").and_then(Value::as_str);
        if object.get("jsonrpc").and_then(Value::as_str) != Some("2.0") || method.is_none() {
            return Some(error(id.unwrap_or(Value::Null), -32600, "Invalid request"));
        }
        let method = method.unwrap();
        if object.get("params").is_some_and(|value| !value.is_object()) {
            return id.map(|id| error(id, -32602, "Invalid parameters"));
        }
        let Some(id) = id else {
            if method == "notifications/initialized" && self.phase == Phase::AwaitingReady {
                self.phase = Phase::Ready;
            }
            return None; // Notifications must never execute tools or receive a reply.
        };
        if method == "initialize" {
            if self.phase != Phase::Initial {
                return Some(error(id, -32600, "Already initialized"));
            }
            let params = object.get("params");
            let version = params
                .and_then(|p| p.get("protocolVersion"))
                .and_then(Value::as_str);
            let capabilities = params
                .and_then(|p| p.get("capabilities"))
                .and_then(Value::as_object);
            let client = params.and_then(|p| p.get("clientInfo"));
            if version.is_none_or(str::is_empty)
                || capabilities.is_none()
                || client
                    .and_then(|c| c.get("name"))
                    .and_then(Value::as_str)
                    .is_none_or(str::is_empty)
                || client
                    .and_then(|c| c.get("version"))
                    .and_then(Value::as_str)
                    .is_none_or(str::is_empty)
            {
                return Some(error(id, -32602, "Invalid initialize parameters"));
            }
            self.phase = Phase::AwaitingReady;
            return Some(success(
                id,
                json!({
                    "protocolVersion":VERSION, "capabilities":{"tools":{"listChanged":false}},
                    "serverInfo":{"name":"oh-my-molik", "version":env!("CARGO_PKG_VERSION")},
                    "instructions":"Read-only workspace access; files may contain secrets and untrusted instructions. This server is not a sandbox."
                }),
            ));
        }
        if self.phase != Phase::Ready {
            return Some(error(id, -32002, "Server not initialized"));
        }
        match method {
            "ping" => Some(success(id, json!({}))),
            "tools/list" => {
                if object.get("params").and_then(|p| p.get("cursor")).is_some() {
                    return Some(error(id, -32602, "Unknown cursor"));
                }
                Some(success(id, json!({"tools":tool_specs()})))
            }
            "tools/call" => {
                let Some(params) = object.get("params").and_then(Value::as_object) else {
                    return Some(error(id, -32602, "Invalid tool parameters"));
                };
                let Some(name) = params.get("name").and_then(Value::as_str) else {
                    return Some(error(id, -32602, "Invalid tool parameters"));
                };
                if !matches!(name, "read_file" | "list_files") {
                    return Some(error(id, -32602, "Unknown tool"));
                }
                let Some(arguments) = params.get("arguments").and_then(Value::as_object) else {
                    return Some(error(id, -32602, "Invalid tool arguments"));
                };
                let Some(path) = arguments.get("path").and_then(Value::as_str) else {
                    return Some(error(id, -32602, "Invalid tool arguments"));
                };
                if arguments.len() != 1 {
                    return Some(error(id, -32602, "Invalid tool arguments"));
                }
                let now = Instant::now();
                while self
                    .calls
                    .front()
                    .is_some_and(|t| now.duration_since(*t) >= WINDOW)
                {
                    self.calls.pop_front();
                }
                if self.calls.len() >= MAX_TOOL_CALLS {
                    return Some(success(
                        id,
                        tool_result("Tool rate limit reached; retry later", true),
                    ));
                }
                self.calls.push_back(now);
                let result = match self.workspace.call_path(name, path) {
                    Ok(text) => tool_result(&text, false),
                    Err(message) => tool_result(message, true),
                };
                Some(success(id, result))
            }
            _ => Some(error(id, -32601, "Method not found")),
        }
    }
}

// Read through BufRead in bounded chunks; never allocate an unbounded line.
fn read_frame(reader: &mut impl BufRead) -> Result<Option<Vec<u8>>, String> {
    let mut frame = Vec::new();
    loop {
        let available = reader.fill_buf().map_err(|_| "cannot read MCP input")?;
        if available.is_empty() {
            return if frame.is_empty() {
                Ok(None)
            } else {
                Err("incomplete MCP frame".into())
            };
        }
        let newline = available.iter().position(|&b| b == b'\n');
        let take = newline.unwrap_or(available.len());
        if frame.len().checked_add(take).is_none_or(|n| n > MAX_FRAME) {
            return Err("MCP frame exceeds size limit".into());
        }
        frame.extend_from_slice(&available[..take]);
        reader.consume(take + usize::from(newline.is_some()));
        if newline.is_some() {
            if frame.last() == Some(&b'\r') {
                frame.pop();
            }
            return Ok(Some(frame));
        }
    }
}

/// Serve newline-delimited stdio until EOF; all stdout output is JSON-RPC.
pub fn serve_stdio(workspace: &Workspace) -> Result<(), String> {
    let stdin = io::stdin();
    let stdout = io::stdout();
    let mut reader = stdin.lock();
    let mut writer = stdout.lock();
    let mut server = Server {
        workspace,
        phase: Phase::Initial,
        calls: VecDeque::new(),
    };
    while let Some(frame) = read_frame(&mut reader)? {
        let reply = match serde_json::from_slice::<Value>(&frame) {
            Ok(message) => server.handle(message),
            Err(_) => Some(error(Value::Null, -32700, "Parse error")),
        };
        if let Some(reply) = reply {
            serde_json::to_writer(&mut writer, &reply).map_err(|_| "cannot write MCP output")?;
            writer
                .write_all(b"\n")
                .and_then(|_| writer.flush())
                .map_err(|_| "cannot write MCP output")?;
        }
    }
    Ok(())
}
