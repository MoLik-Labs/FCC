use oh_my_molik::sse::{Decoder, Error, Event};

#[test]
fn every_chunk_boundary_preserves_unicode_and_crlf() {
    let input = "\u{feff}: comment\r\nevent: update\r\ndata: cáo 🦊\r\ndata: second\r\n\r\ndata: [DONE]\n\n".as_bytes();
    let expected = vec![
        Event {
            event: "update".into(),
            data: "cáo 🦊\nsecond".into(),
        },
        Event {
            event: "message".into(),
            data: "[DONE]".into(),
        },
    ];
    for chunk in 1..=input.len() {
        let mut d = Decoder::new(128, 128).unwrap();
        let mut out = Vec::new();
        for part in input.chunks(chunk) {
            out.extend(d.push(part).unwrap());
        }
        d.finish().unwrap();
        assert_eq!(out, expected, "chunk size {chunk}");
        assert_eq!(d.push(b"x"), Err(Error::Closed));
    }
}

#[test]
fn event_framing_supports_cr_and_empty_data() {
    let mut d = Decoder::new(100, 100).unwrap();
    assert_eq!(
        d.push(b"event: hello\rdata\r\revent: ignored\r\rdata: x\r\r")
            .unwrap(),
        vec![
            Event {
                event: "hello".into(),
                data: "".into()
            },
            Event {
                event: "message".into(),
                data: "x".into()
            }
        ]
    );
}

#[test]
fn incomplete_event_is_not_dispatched() {
    let mut d = Decoder::new(30, 30).unwrap();
    assert!(d.push(b"data: partial\n").unwrap().is_empty());
    d.finish().unwrap();
}

#[test]
fn malformed_and_oversized_streams_are_poisoned() {
    assert!(matches!(Decoder::new(0, 1), Err(Error::InvalidLimits)));
    let mut d = Decoder::new(3, 20).unwrap();
    assert_eq!(d.push(b"1234"), Err(Error::LineTooLong));
    assert_eq!(d.push(b"\n"), Err(Error::Closed));
    let mut d = Decoder::new(100, 3).unwrap();
    assert_eq!(d.push(b"data: abc\n"), Err(Error::EventTooLarge));
    let mut d = Decoder::new(100, 100).unwrap();
    assert_eq!(d.push(b"data: \xff\n"), Err(Error::InvalidUtf8));
    assert_eq!(d.finish(), Err(Error::Closed));
    let mut d = Decoder::new(100, 100).unwrap();
    d.push(b"data: \xf0\x9f").unwrap();
    assert_eq!(d.finish(), Err(Error::InvalidUtf8));
}
