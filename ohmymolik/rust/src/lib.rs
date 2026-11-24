//! A transport-independent, bounded agent state machine.
//!
//! This crate does not perform I/O, enforce a sandbox, or authenticate the user
//! responding to approval requests. Its caller is responsible for those tasks.
#![forbid(unsafe_code)]
pub mod chat;
pub mod sse;
pub mod stream;
pub mod theme;
pub mod workspace;

use std::collections::{HashSet, VecDeque};
use std::fmt;

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct ToolCall {
    pub id: String,
    pub name: String,
    /// Opaque input. The tool adapter must parse and validate its own schema.
    pub arguments: String,
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum ModelReply {
    Final(String),
    Tools(Vec<ToolCall>),
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum Message {
    User(String),
    AssistantFinal(String),
    AssistantTools(Vec<ToolCall>),
    ToolResult { id: String, content: String },
}

/// A read-only declaration is an assertion by the tool integrator, not a sandbox.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum Access {
    ReadOnly,
    Mutating,
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Tool {
    pub name: String,
    pub access: Access,
}

#[derive(Clone, Copy, Debug)]
pub struct Limits {
    pub max_turns: usize,
    pub max_calls_per_turn: usize,
    pub max_text_bytes: usize,
    pub max_transcript_bytes: usize,
}

impl Default for Limits {
    fn default() -> Self {
        Self {
            max_turns: 8,
            max_calls_per_turn: 16,
            max_text_bytes: 64 * 1024,
            max_transcript_bytes: 2 * 1024 * 1024,
        }
    }
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum Action {
    /// Send the transcript to the configured model. The caller handles transport.
    AskModel,
    /// Obtain an explicit decision for this exact call before executing it.
    RequestApproval(ToolCall),
    /// The caller must atomically claim this call with `begin_execution`.
    Execute(ToolCall),
    /// The call was claimed; wait for its result, do not execute it again.
    AwaitResult(String),
    Finished(String),
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Error(&'static str);
impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.0)
    }
}
impl std::error::Error for Error {}

#[derive(Clone, Debug, PartialEq, Eq)]
enum State {
    Model,
    Pending,
    Authorized,
    Executing,
    Finished(String),
}

/// A single, sequential session. Dropping it cancels further dispatch; it does
/// not stop a tool that the caller has already started.
pub struct Session {
    limits: Limits,
    tools: Vec<Tool>,
    messages: Vec<Message>,
    bytes: usize,
    turns: usize,
    seen_ids: HashSet<String>,
    pending: VecDeque<ToolCall>,
    state: State,
}

impl Session {
    pub fn new(prompt: String, tools: Vec<Tool>, limits: Limits) -> Result<Self, Error> {
        if limits.max_turns == 0
            || limits.max_calls_per_turn == 0
            || limits.max_text_bytes == 0
            || limits.max_transcript_bytes == 0
        {
            return Err(Error("limits must be positive"));
        }
        if prompt.trim().is_empty()
            || prompt.len() > limits.max_text_bytes
            || prompt.len() > limits.max_transcript_bytes
        {
            return Err(Error("invalid prompt size"));
        }
        let mut names = HashSet::new();
        for tool in &tools {
            if tool.name.trim().is_empty()
                || tool.name.len() > limits.max_text_bytes
                || !names.insert(&tool.name)
            {
                return Err(Error("invalid or duplicate tool name"));
            }
        }
        let bytes = prompt.len();
        Ok(Self {
            limits,
            tools,
            messages: vec![Message::User(prompt)],
            bytes,
            turns: 0,
            seen_ids: HashSet::new(),
            pending: VecDeque::new(),
            state: State::Model,
        })
    }

    pub fn messages(&self) -> &[Message] {
        &self.messages
    }

    /// Returns the currently authorized action. Calling repeatedly never grants
    /// approval and never advances or executes anything.
    pub fn action(&self) -> Action {
        match &self.state {
            State::Model => Action::AskModel,
            State::Pending | State::Authorized => {
                let call = self.pending.front().expect("pending state has a call");
                let read_only = self
                    .tools
                    .iter()
                    .any(|tool| tool.name == call.name && tool.access == Access::ReadOnly);
                if read_only || self.state == State::Authorized {
                    Action::Execute(call.clone())
                } else {
                    Action::RequestApproval(call.clone())
                }
            }
            State::Executing => Action::AwaitResult(
                self.pending
                    .front()
                    .expect("executing state has a call")
                    .id
                    .clone(),
            ),
            State::Finished(text) => Action::Finished(text.clone()),
        }
    }

    /// Accepts a complete model reply, atomically validating the entire batch
    /// before any calls become dispatchable. Unknown tools are rejected.
    pub fn accept(&mut self, reply: ModelReply) -> Result<(), Error> {
        if self.state != State::Model {
            return Err(Error("not waiting for a model reply"));
        }
        if self.turns >= self.limits.max_turns {
            return Err(Error("turn limit reached"));
        }
        match reply {
            ModelReply::Final(text) => {
                if text.trim().is_empty() || text.len() > self.limits.max_text_bytes {
                    return Err(Error("invalid final response size"));
                }
                self.reserve(text.len())?;
                self.turns += 1;
                self.messages.push(Message::AssistantFinal(text.clone()));
                self.state = State::Finished(text);
            }
            ModelReply::Tools(calls) => {
                if self.turns + 1 >= self.limits.max_turns {
                    return Err(Error("no turn remains for a final response"));
                }
                if calls.is_empty() || calls.len() > self.limits.max_calls_per_turn {
                    return Err(Error("invalid tool call count"));
                }
                let mut ids = HashSet::new();
                let mut bytes = 0usize;
                for call in &calls {
                    if call.id.trim().is_empty()
                        || call.name.trim().is_empty()
                        || call.id.len() > self.limits.max_text_bytes
                        || call.name.len() > self.limits.max_text_bytes
                        || call.arguments.len() > self.limits.max_text_bytes
                    {
                        return Err(Error("invalid tool call size"));
                    }
                    if self.seen_ids.contains(&call.id) || !ids.insert(&call.id) {
                        return Err(Error("duplicate tool call id"));
                    }
                    if !self.tools.iter().any(|tool| tool.name == call.name) {
                        return Err(Error("unknown tool"));
                    }
                    bytes = bytes
                        .checked_add(call.id.len())
                        .and_then(|n| n.checked_add(call.name.len()))
                        .and_then(|n| n.checked_add(call.arguments.len()))
                        .ok_or(Error("transcript size overflow"))?;
                }
                self.reserve(bytes)?;
                self.turns += 1;
                self.seen_ids
                    .extend(calls.iter().map(|call| call.id.clone()));
                self.messages.push(Message::AssistantTools(calls.clone()));
                self.pending = calls.into();
                self.state = State::Pending;
            }
        }
        Ok(())
    }

    /// Explicitly approve or deny the current mutating call. A decision applies
    /// only to this call, never to subsequent calls with the same name.
    pub fn decide(&mut self, id: &str, approved: bool) -> Result<(), Error> {
        let Action::RequestApproval(call) = self.action() else {
            return Err(Error("no approval pending"));
        };
        if call.id != id {
            return Err(Error("approval id mismatch"));
        }
        if approved {
            self.state = State::Authorized;
            Ok(())
        } else {
            self.finish_call(id, "Tool call denied by user.".to_string())
        }
    }

    /// Claim an authorized call exactly once before dispatch. The caller still
    /// owns actual execution, cancellation, and any at-most-once guarantees.
    pub fn begin_execution(&mut self, id: &str) -> Result<ToolCall, Error> {
        let Action::Execute(call) = self.action() else {
            return Err(Error("tool call is not authorized or already claimed"));
        };
        if call.id != id {
            return Err(Error("tool call id mismatch"));
        }
        self.state = State::Executing;
        Ok(call)
    }

    /// Record only the result of the claimed call. Error outputs count too.
    pub fn record_result(&mut self, id: &str, content: String) -> Result<(), Error> {
        if self.state != State::Executing {
            return Err(Error("no tool execution in progress"));
        }
        self.finish_call(id, content)
    }

    fn finish_call(&mut self, id: &str, content: String) -> Result<(), Error> {
        let call = self.pending.front().ok_or(Error("no tool call pending"))?;
        if call.id != id {
            return Err(Error("tool result id mismatch"));
        }
        if content.len() > self.limits.max_text_bytes {
            return Err(Error("tool result too large"));
        }
        let size = id
            .len()
            .checked_add(content.len())
            .ok_or(Error("transcript size overflow"))?;
        self.reserve(size)?;
        self.messages.push(Message::ToolResult {
            id: id.to_owned(),
            content,
        });
        self.pending.pop_front();
        self.state = if self.pending.is_empty() {
            State::Model
        } else {
            State::Pending
        };
        Ok(())
    }

    fn reserve(&mut self, extra: usize) -> Result<(), Error> {
        let total = self
            .bytes
            .checked_add(extra)
            .ok_or(Error("transcript size overflow"))?;
        if total > self.limits.max_transcript_bytes {
            return Err(Error("transcript limit reached"));
        }
        self.bytes = total;
        Ok(())
    }
}
