use oh_my_molik::chat::{parse_response, request_body, validate_endpoint};
use oh_my_molik::ModelReply;
use serde_json::json;

#[test]
fn endpoint_policy_disallows_credentials_and_ambiguity_but_not_remote() {
    for url in [
        "http://localhost:8082/v1/chat/completions",
        "https://example.com/v1/chat/completions",
        "https://[::1]:8080/a",
    ] {
        assert!(validate_endpoint(url).is_ok());
    }
    for url in [
        "file:///etc/passwd",
        "https://user:pass@example.com/a",
        "https://example.com/?token=x",
        "https://example.com/#frag",
        "https://example.com\\@evil/a",
        "https://example.com/\nheader",
        "https:///missing",
    ] {
        assert!(validate_endpoint(url).is_err(), "{url}");
    }
}

#[test]
fn parse_checks_protocol_and_refuses_incomplete_or_mixed_results() {
    let parse = |value: serde_json::Value| parse_response(value.to_string().as_bytes());
    let choice =
        |reason: &str, message| json!({"choices":[{"finish_reason":reason,"message":message}]});
    assert_eq!(
        parse(choice("stop", json!({"role":"assistant","content":"ok"}))),
        Ok(ModelReply::Final("ok".into()))
    );
    for value in [
        json!({"choices":[]}),
        choice("length", json!({"role":"assistant","content":"partial"})),
        choice(
            "content_filter",
            json!({"role":"assistant","content":"filtered"}),
        ),
        choice("stop", json!({"role":"user","content":"bad"})),
        choice("stop", json!({"role":"assistant","content":""})),
        choice(
            "stop",
            json!({"role":"assistant","content":"x","tool_calls":[{}]}),
        ),
        choice(
            "tool_calls",
            json!({"role":"assistant","content":"text","tool_calls":[{}]}),
        ),
        choice(
            "tool_calls",
            json!({"role":"assistant","tool_calls":[{"id":"x","type":"bad"}]}),
        ),
    ] {
        assert!(parse(value).is_err());
    }
    assert!(request_body("", &[]).is_err());
}
