use std::process::Command;

#[test]
fn help_version_and_invalid_arguments_are_truthful() {
    let bin = env!("CARGO_BIN_EXE_oh-my-fox");
    let help = Command::new(bin).arg("--help").output().unwrap();
    assert!(help.status.success());
    let text = String::from_utf8(help.stdout).unwrap();
    assert!(text.contains("not implemented yet"));
    let version = Command::new(bin).arg("--version").output().unwrap();
    assert!(version.status.success());
    assert!(String::from_utf8(version.stdout).unwrap().contains("0.1.0"));
    assert_eq!(
        Command::new(bin).arg("run").status().unwrap().code(),
        Some(2)
    );
}
