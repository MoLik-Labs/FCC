use oh_my_molik::workspace::Workspace;
use std::{
    fs,
    path::PathBuf,
    sync::atomic::{AtomicU64, Ordering},
};
static COUNT: AtomicU64 = AtomicU64::new(0);
fn temp() -> PathBuf {
    let path = std::env::temp_dir().join(format!(
        "molik-workspace-{}-{}",
        std::process::id(),
        COUNT.fetch_add(1, Ordering::Relaxed)
    ));
    fs::create_dir(&path).unwrap();
    path
}

#[test]
fn reads_bounded_files_and_rejects_escape_and_git() {
    let root = temp();
    fs::write(root.join("hello.txt"), "safe text").unwrap();
    fs::write(root.join("large"), vec![b'x'; 65537]).unwrap();
    fs::write(root.join("binary"), b"\0").unwrap();
    fs::create_dir(root.join(".git")).unwrap();
    assert!(Workspace::new(&root.join(".git")).is_err());
    let w = Workspace::new(&root).unwrap();
    assert_eq!(w.call("read_file", r#"{"path":"hello.txt"}"#), "safe text");
    for path in [
        "../outside",
        "/etc/passwd",
        ".git/config",
        ".GIT/config",
        ".git./config",
        ".git::$DATA",
        "large",
        "binary",
    ] {
        let args = serde_json::json!({"path":path}).to_string();
        assert!(w.call("read_file", &args).starts_with("error:"), "{path}");
    }
    assert!(w
        .call("read_file", r#"{"path":"hello.txt","extra":true}"#)
        .starts_with("error:"));
    let listing = w.call("list_files", r#"{"path":"."}"#);
    assert!(listing.contains("hello.txt"));
    assert!(!listing.contains(".git"));
    fs::remove_dir_all(root).unwrap();
}

#[cfg(unix)]
#[test]
fn refuses_symlink_paths() {
    use std::os::unix::fs::symlink;
    let root = temp();
    symlink("/etc/passwd", root.join("link")).unwrap();
    let w = Workspace::new(&root).unwrap();
    assert!(w
        .call("read_file", r#"{"path":"link"}"#)
        .starts_with("error:"));
    fs::remove_dir_all(root).unwrap();
}
