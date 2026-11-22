use std::process::Command;

#[test]
fn help_version_and_invalid_arguments_are_truthful() {
    let bin = env!("CARGO_BIN_EXE_oh-my-molik");
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

#[test]
fn theme_commands_list_and_inspect_without_claiming_to_set_a_theme() {
    let bin = env!("CARGO_BIN_EXE_oh-my-molik");
    let output = Command::new(bin).args(["theme", "list"]).output().unwrap();
    assert!(output.status.success());
    let text = String::from_utf8(output.stdout).unwrap();
    assert!(text.lines().any(|line| line == "forest"));
    assert!(text.lines().any(|line| line == "catppuccin-mocha"));
    let output = Command::new(bin)
        .args(["theme", "show", "catppuccin-mocha"])
        .output()
        .unwrap();
    assert!(output.status.success());
    assert!(String::from_utf8(output.stdout)
        .unwrap()
        .contains("#1E1E2E"));
    assert_eq!(
        Command::new(bin)
            .args(["theme", "show", "missing"])
            .output()
            .unwrap()
            .status
            .code(),
        Some(2)
    );
    assert_eq!(
        Command::new(bin)
            .args(["theme", "show", "forest", "extra"])
            .output()
            .unwrap()
            .status
            .code(),
        Some(2)
    );
}
