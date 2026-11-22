use oh_my_molik::{Access, Action, Limits, Message, ModelReply, Session, Tool, ToolCall};

fn tools() -> Vec<Tool> {
    vec![
        Tool {
            name: "read".into(),
            access: Access::ReadOnly,
        },
        Tool {
            name: "write".into(),
            access: Access::Mutating,
        },
    ]
}
fn call(id: &str, name: &str) -> ToolCall {
    ToolCall {
        id: id.into(),
        name: name.into(),
        arguments: "{}".into(),
    }
}
fn session() -> Session {
    Session::new("inspect".into(), tools(), Limits::default()).unwrap()
}

#[test]
fn ordered_loop_and_exact_approval() {
    let mut s = session();
    let calls = vec![call("1", "read"), call("2", "write"), call("3", "write")];
    s.accept(ModelReply::Tools(calls.clone())).unwrap();
    assert_eq!(s.action(), Action::Execute(calls[0].clone()));
    assert!(s.record_result("1", "unclaimed".into()).is_err());
    assert_eq!(s.begin_execution("1").unwrap(), calls[0]);
    assert_eq!(s.action(), Action::AwaitResult("1".into()));
    assert!(s.begin_execution("1").is_err());
    assert!(s.record_result("2", "wrong".into()).is_err());
    s.record_result("1", "data".into()).unwrap();
    assert_eq!(s.action(), Action::RequestApproval(calls[1].clone()));
    assert!(s.begin_execution("2").is_err());
    assert!(s.decide("3", true).is_err());
    s.decide("2", true).unwrap();
    assert_eq!(s.begin_execution("2").unwrap(), calls[1]);
    s.record_result("2", "saved".into()).unwrap();
    assert_eq!(s.action(), Action::RequestApproval(calls[2].clone()));
    s.decide("3", false).unwrap();
    assert_eq!(s.action(), Action::AskModel);
    assert_eq!(
        s.messages()[4],
        Message::ToolResult {
            id: "3".into(),
            content: "Tool call denied by user.".into()
        }
    );
    s.accept(ModelReply::Final("done".into())).unwrap();
    assert_eq!(s.action(), Action::Finished("done".into()));
    assert!(s.accept(ModelReply::Final("again".into())).is_err());
}

#[test]
fn invalid_batch_is_atomic_and_unknown_tools_are_denied() {
    let mut s = session();
    for calls in [
        vec![],
        vec![call("a", "read"), call("a", "write")],
        vec![call("a", "read"), call("b", "shell")],
        vec![call("", "read")],
    ] {
        assert!(s.accept(ModelReply::Tools(calls)).is_err());
        assert_eq!(s.action(), Action::AskModel);
        assert_eq!(s.messages().len(), 1);
    }
    s.accept(ModelReply::Tools(vec![call("a", "read")]))
        .unwrap();
    s.begin_execution("a").unwrap();
    s.record_result("a", "ok".into()).unwrap();
    assert!(s
        .accept(ModelReply::Tools(vec![call("a", "read")]))
        .is_err());
    assert_eq!(s.action(), Action::AskModel);
}

#[test]
fn invalid_results_do_not_advance_or_reauthorize() {
    let limits = Limits {
        max_text_bytes: 30,
        max_transcript_bytes: 34,
        ..Limits::default()
    };
    let mut s = Session::new("p".into(), tools(), limits).unwrap();
    s.accept(ModelReply::Tools(vec![call("a", "write")]))
        .unwrap();
    assert!(s.record_result("a", "fake".into()).is_err());
    // Denial cannot fit: failure must not open a path to record a forged result.
    assert!(s.decide("a", false).is_err());
    assert!(matches!(s.action(), Action::RequestApproval(_)));
    assert!(s.record_result("a", "fake".into()).is_err());
    s.decide("a", true).unwrap();
    s.begin_execution("a").unwrap();
    assert!(s.record_result("a", "x".repeat(31)).is_err());
    assert!(matches!(s.action(), Action::AwaitResult(_)));
    s.record_result("a", "ok".into()).unwrap();
}

#[test]
fn turn_and_transcript_limits_fail_before_dispatch() {
    let mut s = Session::new(
        "p".into(),
        tools(),
        Limits {
            max_turns: 1,
            ..Limits::default()
        },
    )
    .unwrap();
    assert!(s
        .accept(ModelReply::Tools(vec![call("a", "read")]))
        .is_err());
    s.accept(ModelReply::Final("answer".into())).unwrap();
    let mut s = Session::new(
        "p".into(),
        tools(),
        Limits {
            max_transcript_bytes: 5,
            ..Limits::default()
        },
    )
    .unwrap();
    assert!(s
        .accept(ModelReply::Tools(vec![call("a", "read")]))
        .is_err());
    assert_eq!(s.messages().len(), 1);
    s.accept(ModelReply::Final("done".into())).unwrap();
}

#[test]
fn validates_configuration_and_utf8_byte_budgets() {
    assert!(Session::new(" ".into(), tools(), Limits::default()).is_err());
    assert!(Session::new(
        "p".into(),
        vec![tools()[0].clone(), tools()[0].clone()],
        Limits::default()
    )
    .is_err());
    assert!(Session::new(
        "p".into(),
        tools(),
        Limits {
            max_turns: 0,
            ..Limits::default()
        }
    )
    .is_err());
    assert!(Session::new(
        "🦊🦊".into(),
        vec![],
        Limits {
            max_text_bytes: 7,
            ..Limits::default()
        }
    )
    .is_err());
    let mut s = session();
    assert!(s.accept(ModelReply::Final(" \n".into())).is_err());
    assert_eq!(s.action(), Action::AskModel);
}

#[test]
fn batch_limit_and_results_at_exact_byte_boundary() {
    let limits = Limits {
        max_calls_per_turn: 1,
        max_text_bytes: 8,
        max_transcript_bytes: 18,
        ..Limits::default()
    };
    let mut s = Session::new("p".into(), vec![tools()[0].clone()], limits).unwrap();
    assert!(s
        .accept(ModelReply::Tools(vec![
            call("a", "read"),
            call("b", "read")
        ]))
        .is_err());
    s.accept(ModelReply::Tools(vec![call("a", "read")]))
        .unwrap();
    assert_eq!(s.begin_execution("a").unwrap(), call("a", "read"));
    // 1 prompt + 7 call bytes + 1 result ID + 8 result bytes = 17.
    s.record_result("a", "🦊🦊".into()).unwrap();
    assert!(s.accept(ModelReply::Final("no".into())).is_err());
    s.accept(ModelReply::Final("!".into())).unwrap();
    assert_eq!(s.action(), Action::Finished("!".into()));
}
