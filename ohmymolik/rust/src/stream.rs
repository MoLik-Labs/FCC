//! Assemble provider-decoded deltas into a complete, validated model reply.
//! JSON parsing, choice selection and HTTP framing belong to the adapter.
use crate::{ModelReply, ToolCall};
use std::collections::{BTreeMap, HashSet};

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum Error {
    InvalidLimits,
    Poisoned,
    TooManyTools,
    TooLarge,
    InvalidToolCalls,
    MixedContent,
    Truncated,
    Filtered,
    Empty,
    InvalidFinish,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Finish {
    Stop,
    ToolCalls,
    Length,
    ContentFilter,
}

/// An adapter must only feed deltas for one choice, in provider order. This
/// assembler cannot deduplicate a repeated network chunk or resume a response.
pub struct Assembler {
    text: String,
    calls: BTreeMap<usize, ToolCall>,
    total: usize,
    max_total: usize,
    max_field: usize,
    max_calls: usize,
    poisoned: bool,
}

impl Assembler {
    pub fn new(max_total: usize, max_field: usize, max_calls: usize) -> Result<Self, Error> {
        if max_total == 0 || max_field == 0 || max_calls == 0 {
            return Err(Error::InvalidLimits);
        }
        Ok(Self {
            text: String::new(),
            calls: BTreeMap::new(),
            total: 0,
            max_total,
            max_field,
            max_calls,
            poisoned: false,
        })
    }

    pub fn text(&mut self, delta: &str) -> Result<(), Error> {
        if self.poisoned {
            return Err(Error::Poisoned);
        }
        let result = self.append_text(delta);
        if result.is_err() {
            self.poisoned = true;
        }
        result
    }

    fn append_text(&mut self, delta: &str) -> Result<(), Error> {
        let next_field = self
            .text
            .len()
            .checked_add(delta.len())
            .ok_or(Error::TooLarge)?;
        let next_total = self.total.checked_add(delta.len()).ok_or(Error::TooLarge)?;
        if next_field > self.max_field || next_total > self.max_total {
            return Err(Error::TooLarge);
        }
        self.text.push_str(delta);
        self.total = next_total;
        Ok(())
    }

    /// Append fragments for a zero-based tool index. Only after `finish` can
    /// a caller submit the reply to Session; deltas never authorize execution.
    pub fn tool(
        &mut self,
        index: usize,
        id: &str,
        name: &str,
        arguments: &str,
    ) -> Result<(), Error> {
        if self.poisoned {
            return Err(Error::Poisoned);
        }
        let result = self.append_tool(index, id, name, arguments);
        if result.is_err() {
            self.poisoned = true;
        }
        result
    }

    fn append_tool(
        &mut self,
        index: usize,
        id: &str,
        name: &str,
        arguments: &str,
    ) -> Result<(), Error> {
        if index >= self.max_calls {
            return Err(Error::TooManyTools);
        }
        let current = self.calls.get(&index);
        for (existing, delta) in [
            (current.map_or(0, |c| c.id.len()), id),
            (current.map_or(0, |c| c.name.len()), name),
            (current.map_or(0, |c| c.arguments.len()), arguments),
        ] {
            if existing
                .checked_add(delta.len())
                .is_none_or(|n| n > self.max_field)
            {
                return Err(Error::TooLarge);
            }
        }
        let next_total = self
            .total
            .checked_add(id.len())
            .and_then(|n| n.checked_add(name.len()))
            .and_then(|n| n.checked_add(arguments.len()))
            .ok_or(Error::TooLarge)?;
        if next_total > self.max_total {
            return Err(Error::TooLarge);
        }
        let call = self.calls.entry(index).or_insert_with(|| ToolCall {
            id: String::new(),
            name: String::new(),
            arguments: String::new(),
        });
        call.id.push_str(id);
        call.name.push_str(name);
        call.arguments.push_str(arguments);
        self.total = next_total;
        Ok(())
    }

    pub fn finish(self, reason: Finish) -> Result<ModelReply, Error> {
        if self.poisoned {
            return Err(Error::Poisoned);
        }
        match reason {
            Finish::Length => return Err(Error::Truncated),
            Finish::ContentFilter => return Err(Error::Filtered),
            _ => {}
        }
        if !self.calls.is_empty() && !self.text.trim().is_empty() {
            return Err(Error::MixedContent);
        }
        match reason {
            Finish::Stop if self.calls.is_empty() => {
                if self.text.trim().is_empty() {
                    Err(Error::Empty)
                } else {
                    Ok(ModelReply::Final(self.text))
                }
            }
            Finish::ToolCalls if !self.calls.is_empty() => {
                let mut ids = HashSet::new();
                let mut calls = Vec::with_capacity(self.calls.len());
                for (expected, (index, call)) in self.calls.into_iter().enumerate() {
                    if index != expected
                        || call.id.trim().is_empty()
                        || call.name.trim().is_empty()
                        || !ids.insert(call.id.clone())
                    {
                        return Err(Error::InvalidToolCalls);
                    }
                    calls.push(call);
                }
                Ok(ModelReply::Tools(calls))
            }
            _ => Err(Error::InvalidFinish),
        }
    }
}
