use std::{
    fs,
    io::{Read, Write},
    net::{TcpListener, TcpStream},
    process::{Command, Stdio},
    thread,
    time::{Duration, Instant},
};

fn accept_with_timeout(listener: &TcpListener) -> TcpStream {
    listener.set_nonblocking(true).unwrap();
    let deadline = Instant::now() + Duration::from_secs(15);
    loop {
        match listener.accept() {
            Ok((stream, _)) => return stream,
            Err(error)
                if error.kind() == std::io::ErrorKind::WouldBlock && Instant::now() < deadline =>
            {
                thread::sleep(Duration::from_millis(10))
            }
            Err(error) => panic!("failed to accept mock request: {error}"),
        }
    }
}

fn request(mut stream: &TcpStream) -> (String, Vec<u8>) {
    stream
        .set_read_timeout(Some(Duration::from_secs(20)))
        .unwrap();
    let mut raw = Vec::new();
    let mut buf = [0u8; 4096];
    let split = loop {
        let n = stream.read(&mut buf).unwrap();
        assert!(n > 0);
        raw.extend_from_slice(&buf[..n]);
        if let Some(i) = raw.windows(4).position(|b| b == b"\r\n\r\n") {
            break i + 4;
        }
    };
    let headers = String::from_utf8(raw[..split].to_vec()).unwrap();
    let length: usize = headers
        .lines()
        .find_map(|line| {
            line.to_ascii_lowercase()
                .strip_prefix("content-length:")
                .map(|s| s.trim().parse().unwrap())
        })
        .unwrap();
    while raw.len() - split < length {
        let n = stream.read(&mut buf).unwrap();
        assert!(n > 0);
        raw.extend_from_slice(&buf[..n]);
    }
    (headers, raw[split..split + length].to_vec())
}
fn respond(mut stream: &TcpStream, status: &str, body: &str, extra: &str) {
    write!(stream, "HTTP/1.1 {status}\r\nContent-Type: application/json\r\nContent-Length: {}\r\nConnection: close\r\n{extra}\r\n{body}", body.len()).unwrap();
}

#[test]
fn cli_executes_read_only_loop_against_mock_endpoint() {
    let root = std::env::temp_dir().join(format!("molik-run-{}", std::process::id()));
    fs::create_dir(&root).unwrap();
    fs::write(root.join("hello.txt"), "workspace secret sample").unwrap();
    let listener = TcpListener::bind("127.0.0.1:0").unwrap();
    let port = listener.local_addr().unwrap().port();
    let server = thread::spawn(move || {
        for turn in 0..2 {
            let stream = accept_with_timeout(&listener);
            let (headers, body) = request(&stream);
            assert!(headers
                .to_ascii_lowercase()
                .contains("authorization: bearer test-key"));
            let value: serde_json::Value = serde_json::from_slice(&body).unwrap();
            assert_eq!(value["model"], "example");
            assert_eq!(
                value["messages"][1]["content"],
                "inspect \"quoted\" \\ slash\nsecond line 🦊"
            );
            assert_eq!(value["stream"], false);
            let reply = if turn == 0 {
                serde_json::json!({"choices":[{"finish_reason":"tool_calls","message":{"role":"assistant","content":null,"tool_calls":[{"id":"call1","type":"function","function":{"name":"read_file","arguments":"{\"path\":\"hello.txt\"}"}}]}}]})
            } else {
                assert_eq!(value["messages"][3]["content"], "workspace secret sample");
                serde_json::json!({"choices":[{"finish_reason":"stop","message":{"role":"assistant","content":"Finished safely.\u{1b}[31m"}}]})
            };
            respond(&stream, "200 OK", &reply.to_string(), "");
        }
    });
    let mut child = Command::new(env!("CARGO_BIN_EXE_oh-my-molik"))
        .args([
            "run",
            "--endpoint",
            &format!("http://127.0.0.1:{port}/v1/chat/completions"),
            "--model",
            "example",
            "--workspace",
            root.to_str().unwrap(),
            "--prompt-stdin",
        ])
        .env("OMM_API_KEY", "test-key")
        .env("NO_PROXY", "127.0.0.1,localhost")
        .env("no_proxy", "127.0.0.1,localhost")
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .unwrap();
    child
        .stdin
        .take()
        .unwrap()
        .write_all("inspect \"quoted\" \\ slash\nsecond line 🦊".as_bytes())
        .unwrap();
    let output = child.wait_with_output().unwrap();
    server.join().unwrap();
    fs::remove_dir_all(root).unwrap();
    assert!(
        output.status.success(),
        "{}",
        String::from_utf8_lossy(&output.stderr)
    );
    assert_eq!(
        String::from_utf8(output.stdout).unwrap(),
        "Finished safely.�[31m\n"
    );
    assert!(!String::from_utf8_lossy(&output.stderr).contains("test-key"));
}

#[test]
fn redirect_does_not_leak_credentials() {
    let listener = TcpListener::bind("127.0.0.1:0").unwrap();
    let port = listener.local_addr().unwrap().port();
    let redirect = TcpListener::bind("127.0.0.1:0").unwrap();
    let redirect_port = redirect.local_addr().unwrap().port();
    redirect.set_nonblocking(true).unwrap();
    let server = thread::spawn(move || {
        let stream = accept_with_timeout(&listener);
        let _ = request(&stream);
        respond(
            &stream,
            "302 Found",
            "",
            &format!("Location: http://127.0.0.1:{redirect_port}/\r\n"),
        );
    });
    let out = Command::new(env!("CARGO_BIN_EXE_oh-my-molik"))
        .args([
            "run",
            "--endpoint",
            &format!("http://127.0.0.1:{port}/"),
            "--model",
            "m",
            "--workspace",
            ".",
            "test",
        ])
        .env("OMM_API_KEY", "test-key")
        .env("NO_PROXY", "127.0.0.1,localhost")
        .env("no_proxy", "127.0.0.1,localhost")
        .output()
        .unwrap();
    server.join().unwrap();
    assert!(!out.status.success());
    assert!(String::from_utf8_lossy(&out.stderr).contains("HTTP 302"));
    assert!(redirect.accept().is_err());
}
