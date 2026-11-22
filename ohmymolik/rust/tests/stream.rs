use oh_my_molik::stream::{Assembler, Error, Finish};
use oh_my_molik::{Access, Action, Limits, ModelReply, Session, Tool, ToolCall};

#[test]
fn assembles_interleaved_calls_then_session_controls_dispatch() {
    let mut a = Assembler::new(256, 128, 2).unwrap();
    a.tool(1, "b", "wr", "{\"a\":").unwrap();
    a.tool(0, "a", "read", "{}").unwrap();
    a.tool(1, "", "ite", "1}").unwrap();
    let reply = a.finish(Finish::ToolCalls).unwrap();
    assert_eq!(
        reply,
        ModelReply::Tools(vec![
            ToolCall {
                id: "a".into(),
                name: "read".into(),
                arguments: "{}".into()
            },
            ToolCall {
                id: "b".into(),
                name: "write".into(),
                arguments: "{\"a\":1}".into()
            }
        ])
    );
    let mut s = Session::new(
        "hi".into(),
        vec![
            Tool {
                name: "read".into(),
                access: Access::ReadOnly,
            },
            Tool {
                name: "write".into(),
                access: Access::Mutating,
            },
        ],
        Limits::default(),
    )
    .unwrap();
    s.accept(reply).unwrap();
    s.begin_execution("a").unwrap();
    s.record_result("a", "ok".into()).unwrap();
    assert!(matches!(s.action(), Action::RequestApproval(_)));
}

#[test]
fn rejects_bad_finishes_gaps_duplicate_ids_and_mixed_content() {
    for reason in [Finish::Length, Finish::ContentFilter] {
        let mut a = Assembler::new(100, 100, 2).unwrap();
        a.tool(0, "a", "write", "{}").unwrap();
        assert!(a.finish(reason).is_err());
    }
    let mut a = Assembler::new(100, 100, 3).unwrap();
    a.tool(1, "a", "read", "{}").unwrap();
    assert_eq!(a.finish(Finish::ToolCalls), Err(Error::InvalidToolCalls));
    let mut a = Assembler::new(100, 100, 3).unwrap();
    a.tool(0, "a", "read", "{}").unwrap();
    a.tool(1, "a", "write", "{}").unwrap();
    assert_eq!(a.finish(Finish::ToolCalls), Err(Error::InvalidToolCalls));
    let mut a = Assembler::new(100, 100, 3).unwrap();
    a.tool(0, "a", "read", "{}").unwrap();
    a.text("hello").unwrap();
    assert_eq!(a.finish(Finish::ToolCalls), Err(Error::MixedContent));
    let mut a = Assembler::new(100, 100, 3).unwrap();
    a.tool(0, "a", "read", "{}").unwrap();
    assert_eq!(a.finish(Finish::Stop), Err(Error::InvalidFinish));
    let mut a = Assembler::new(100, 100, 3).unwrap();
    a.text("hello").unwrap();
    assert_eq!(a.finish(Finish::ToolCalls), Err(Error::InvalidFinish));
}

#[test]
fn overflow_errors_poison_assembly_so_truncated_calls_cannot_execute() {
    let mut a = Assembler::new(8, 5, 1).unwrap();
    a.tool(0, "a", "read", "").unwrap();
    assert_eq!(a.tool(0, "", "xx", ""), Err(Error::TooLarge));
    assert_eq!(a.tool(0, "", "", "{}"), Err(Error::Poisoned));
    assert_eq!(a.finish(Finish::ToolCalls), Err(Error::Poisoned));
    let mut a = Assembler::new(8, 5, 1).unwrap();
    assert_eq!(a.tool(1, "b", "x", ""), Err(Error::TooManyTools));
    assert_eq!(a.finish(Finish::ToolCalls), Err(Error::Poisoned));
    let mut a = Assembler::new(4, 4, 1).unwrap();
    a.text("🦊").unwrap();
    assert_eq!(a.text("x"), Err(Error::TooLarge));
    assert_eq!(a.finish(Finish::Stop), Err(Error::Poisoned));
}
