use serde_json::{json, Value};
use std::{
    fs,
    io::Write,
    process::{Command, Stdio},
    thread,
};

fn exchange(root: &std::path::Path, input: &str) -> (std::process::Output, Vec<Value>) {
    let mut child = Command::new(env!("CARGO_BIN_EXE_oh-my-molik"))
        .args(["mcp", "serve", "--workspace", root.to_str().unwrap()])
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .unwrap();
    let mut stdin = child.stdin.take().unwrap();
    let input = input.as_bytes().to_vec();
    let sender = thread::spawn(move || {
        let _ = stdin.write_all(&input);
    });
    let output = child.wait_with_output().unwrap();
    sender.join().unwrap();
    let responses = std::str::from_utf8(&output.stdout)
        .unwrap()
        .lines()
        .map(|line| serde_json::from_str(line).unwrap())
        .collect();
    (output, responses)
}

fn init() -> Value {
    json!({"jsonrpc":"2.0", "id":1, "method":"initialize", "params":{
        "protocolVersion":"2025-06-18", "capabilities":{}, "clientInfo":{"name":"test","version":"1"}
    }})
}
fn wire(messages: &[Value]) -> String {
    messages.iter().map(|x| format!("{x}\n")).collect()
}

#[test]
fn stdio_read_only_mcp_lifecycle_and_tool_errors() {
    let root = std::env::temp_dir().join(format!("molik-mcp-{}", std::process::id()));
    fs::create_dir(&root).unwrap();
    fs::write(
        root.join("example.txt"),
        "error: this is actual file content",
    )
    .unwrap();
    let input = wire(&[
        json!({"jsonrpc":"2.0","id":"early","method":"tools/list"}),
        init(),
        json!({"jsonrpc":"2.0","method":"notifications/initialized"}),
        json!({"jsonrpc":"2.0","id":2,"method":"tools/list"}),
        json!({"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"read_file","arguments":{"path":"example.txt"}}}),
        json!({"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"read_file","arguments":{"path":"../secret"}}}),
        json!({"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"shell","arguments":{"path":"."}}}),
        json!({"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"read_file","arguments":{"path":"example.txt","extra":true}}}),
        json!({"jsonrpc":"2.0","method":"tools/call","params":{"name":"read_file","arguments":{"path":"example.txt"}}}),
        json!({"jsonrpc":"2.0","id":7,"method":"unknown"}),
        json!({"jsonrpc":"2.0","id":8,"method":"ping"}),
    ]);
    let (out, reply) = exchange(&root, &input);
    fs::remove_dir_all(root).unwrap();
    assert!(
        out.status.success(),
        "{}",
        String::from_utf8_lossy(&out.stderr)
    );
    assert_eq!(reply.len(), 9);
    assert_eq!(reply[0]["error"]["code"], -32002);
    assert_eq!(reply[1]["result"]["protocolVersion"], "2025-06-18");
    assert!(reply[1]["result"]["capabilities"]["tools"].is_object());
    let tools = reply[2]["result"]["tools"].as_array().unwrap();
    assert_eq!(tools.len(), 2);
    assert_eq!(tools[0]["name"], "list_files");
    assert_eq!(tools[1]["name"], "read_file");
    assert_eq!(
        reply[3]["result"]["content"][0]["text"],
        "error: this is actual file content"
    );
    assert_ne!(reply[3]["result"]["isError"], true);
    assert_eq!(reply[4]["result"]["isError"], true);
    assert_eq!(reply[5]["error"]["code"], -32602);
    assert_eq!(reply[6]["error"]["code"], -32602);
    assert_eq!(reply[7]["error"]["code"], -32601);
    assert_eq!(reply[8]["result"], json!({}));
}

#[test]
fn refuses_oversized_frame_without_consuming_following_requests() {
    let root = std::env::temp_dir();
    let input = format!("{}\n{}\n", "x".repeat(256 * 1024 + 1), init());
    let (out, replies) = exchange(&root, &input);
    assert!(!out.status.success());
    assert!(replies.is_empty());
    assert!(String::from_utf8_lossy(&out.stderr).contains("size limit"));
}

#[test]
fn malformed_messages_and_negotiation_do_not_bypass_lifecycle() {
    let root = std::env::temp_dir();
    let input = format!(
        "{{not json}}\n{}",
        wire(&[
            json!({"jsonrpc":"2.0","id":null,"method":"ping"}),
            json!({"jsonrpc":"2.0","id":10,"method":"initialize","params":{}}),
            json!({"jsonrpc":"2.0","method":"notifications/initialized"}),
            json!({"jsonrpc":"2.0","id":11,"method":"tools/list"}),
            json!({"jsonrpc":"2.0","id":12,"method":"initialize","params":{"protocolVersion":"2099-01-01","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}),
            json!({"jsonrpc":"2.0","id":13,"method":"tools/list"}),
            json!({"jsonrpc":"2.0","method":"notifications/initialized","params":[]}),
            json!({"jsonrpc":"2.0","id":14,"method":"tools/list"}),
            json!({"jsonrpc":"2.0","method":"notifications/initialized"}),
            json!({"jsonrpc":"2.0","id":15,"method":"ping","params":[]}),
            json!({"jsonrpc":"2.0","id":16,"method":"initialize","params":{}}),
            json!({"jsonrpc":"2.0","id":17,"method":"ping"}),
            json!({"jsonrpc":"2.0","id":18,"method":"tools/list","params":{"cursor":"unknown"}}),
        ])
    );
    let (out, replies) = exchange(&root, &input);
    assert!(out.status.success());
    assert_eq!(replies.len(), 11);
    assert_eq!(replies[0]["error"]["code"], -32700);
    assert_eq!(replies[1]["error"]["code"], -32600);
    assert_eq!(replies[2]["error"]["code"], -32602);
    assert_eq!(replies[3]["error"]["code"], -32002);
    assert_eq!(replies[4]["result"]["protocolVersion"], "2025-06-18");
    assert_eq!(replies[5]["error"]["code"], -32002);
    assert_eq!(replies[6]["error"]["code"], -32002);
    assert_eq!(replies[7]["error"]["code"], -32602);
    assert_eq!(replies[8]["error"]["code"], -32600);
    assert_eq!(replies[9]["result"], json!({}));
    assert_eq!(replies[10]["error"]["code"], -32602);
}

#[test]
fn rate_limit_bounds_tool_calls_per_session() {
    let mut messages = vec![
        init(),
        json!({"jsonrpc":"2.0","method":"notifications/initialized"}),
    ];
    for id in 2..=122 {
        messages.push(json!({"jsonrpc":"2.0","id":id,"method":"tools/call","params":{"name":"list_files","arguments":{"path":"."}}}));
    }
    let (out, replies) = exchange(&std::env::temp_dir(), &wire(&messages));
    assert!(out.status.success());
    assert_eq!(replies.len(), 122);
    assert_eq!(replies[120]["result"]["isError"], false);
    assert_eq!(replies[121]["result"]["isError"], true);
    assert!(replies[121]["result"]["content"][0]["text"]
        .as_str()
        .unwrap()
        .contains("rate limit"));
}

#[test]
fn partial_frame_at_eof_is_not_executed() {
    let root = std::env::temp_dir();
    let (out, replies) = exchange(&root, &format!("{}", init()));
    assert!(!out.status.success());
    assert!(replies.is_empty());
    assert!(String::from_utf8_lossy(&out.stderr).contains("incomplete MCP frame"));
}

#[test]
fn supports_crlf_framing() {
    let root = std::env::temp_dir();
    let input = format!("{}\r\n", init());
    let (out, replies) = exchange(&root, &input);
    assert!(out.status.success());
    assert_eq!(replies[0]["result"]["protocolVersion"], "2025-06-18");
}

#[test]
fn git_metadata_is_not_exposed_through_mcp() {
    let root = std::env::temp_dir().join(format!("molik-mcp-git-{}", std::process::id()));
    fs::create_dir_all(root.join(".git")).unwrap();
    fs::write(root.join(".git/config"), "private").unwrap();
    let input = wire(&[
        init(),
        json!({"jsonrpc":"2.0","method":"notifications/initialized"}),
        json!({"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"read_file","arguments":{"path":".git/config"}}}),
        json!({"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_files","arguments":{"path":"."}}}),
    ]);
    let (out, replies) = exchange(&root, &input);
    fs::remove_dir_all(root).unwrap();
    assert!(out.status.success());
    assert_eq!(replies[1]["result"]["isError"], true);
    assert!(!replies[1].to_string().contains("private"));
    assert!(!replies[2]["result"]["content"][0]["text"]
        .as_str()
        .unwrap()
        .contains(".git"));
}
