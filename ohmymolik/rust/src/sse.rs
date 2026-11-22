//! Bounded incremental Server-Sent Events framing, independent of HTTP and JSON.
//! This parses data and event fields only; id and retry are deliberately ignored.

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Event {
    pub event: String,
    pub data: String,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum Error {
    InvalidLimits,
    LineTooLong,
    EventTooLarge,
    InvalidUtf8,
    Closed,
}

/// One decoder per response. A framing error poisons the decoder: discard it.
/// EOF without a blank line does not dispatch an incomplete event.
pub struct Decoder {
    line: Vec<u8>,
    data: String,
    event: String,
    max_line: usize,
    max_event: usize,
    first_line: bool,
    skip_lf: bool,
    closed: bool,
}

impl Decoder {
    pub fn new(max_line: usize, max_event: usize) -> Result<Self, Error> {
        if max_line == 0 || max_event == 0 {
            return Err(Error::InvalidLimits);
        }
        Ok(Self {
            line: Vec::new(),
            data: String::new(),
            event: String::new(),
            max_line,
            max_event,
            first_line: true,
            skip_lf: false,
            closed: false,
        })
    }

    /// Feed raw HTTP body bytes. Delivers events in order after a blank line.
    /// Bytes must be bounded by the caller at the network layer as well.
    pub fn push(&mut self, bytes: &[u8]) -> Result<Vec<Event>, Error> {
        if self.closed {
            return Err(Error::Closed);
        }
        let mut out = Vec::new();
        for &byte in bytes {
            if self.skip_lf {
                self.skip_lf = false;
                if byte == b'\n' {
                    continue;
                }
            }
            if byte == b'\r' || byte == b'\n' {
                if let Err(error) = self.end_line(&mut out) {
                    self.closed = true;
                    return Err(error);
                }
                self.skip_lf = byte == b'\r';
            } else {
                if self.line.len() >= self.max_line {
                    self.closed = true;
                    return Err(Error::LineTooLong);
                }
                self.line.push(byte);
            }
        }
        Ok(out)
    }

    /// Validate a trailing partial line and discard an uncommitted event.
    pub fn finish(&mut self) -> Result<(), Error> {
        if self.closed {
            return Err(Error::Closed);
        }
        self.closed = true;
        if self.first_line && self.line.starts_with(&[0xef, 0xbb, 0xbf]) {
            self.line.drain(..3);
        }
        std::str::from_utf8(&self.line).map_err(|_| Error::InvalidUtf8)?;
        self.line.clear();
        self.data.clear();
        self.event.clear();
        Ok(())
    }

    fn end_line(&mut self, out: &mut Vec<Event>) -> Result<(), Error> {
        if self.first_line {
            self.first_line = false;
            if self.line.starts_with(&[0xef, 0xbb, 0xbf]) {
                self.line.drain(..3);
            }
        }
        let line = std::str::from_utf8(&self.line).map_err(|_| Error::InvalidUtf8)?;
        if line.is_empty() {
            if !self.data.is_empty() {
                self.data.pop(); // the last data field's newline
                out.push(Event {
                    event: if self.event.is_empty() {
                        "message".into()
                    } else {
                        std::mem::take(&mut self.event)
                    },
                    data: std::mem::take(&mut self.data),
                });
            }
            self.event.clear();
        } else if !line.starts_with(':') {
            let (field, value) = line.split_once(':').unwrap_or((line, ""));
            let value = value.strip_prefix(' ').unwrap_or(value);
            match field {
                "data" => {
                    let size = self
                        .data
                        .len()
                        .checked_add(value.len())
                        .and_then(|n| n.checked_add(1))
                        .ok_or(Error::EventTooLarge)?;
                    if size > self.max_event {
                        return Err(Error::EventTooLarge);
                    }
                    self.data.push_str(value);
                    self.data.push('\n');
                }
                "event" => {
                    if value.len() > self.max_event {
                        return Err(Error::EventTooLarge);
                    }
                    self.event.clear();
                    self.event.push_str(value);
                }
                _ => {}
            }
        }
        self.line.clear();
        Ok(())
    }
}
