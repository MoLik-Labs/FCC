#![forbid(unsafe_code)]

fn main() {
    let mut args = std::env::args().skip(1);
    match (args.next().as_deref(), args.next()) {
        (None | Some("--help") | Some("-h"), None) => {
            println!("Oh My Fox — experimental Rust agent core\n\nUsage: oh-my-fox [--help | --version]\n\nThe Rust model transport and tools are not implemented yet.\nFor the working read-only Go prototype, see ../README.md.");
        }
        (Some("--version" | "-V"), None) => println!(
            "oh-my-fox {} (experimental core)",
            env!("CARGO_PKG_VERSION")
        ),
        _ => {
            eprintln!("Unknown arguments. Run oh-my-fox --help.");
            std::process::exit(2);
        }
    }
}
