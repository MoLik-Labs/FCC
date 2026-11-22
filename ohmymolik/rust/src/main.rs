#![forbid(unsafe_code)]

use oh_my_molik::theme::{find, THEMES};

fn help() {
    println!("Oh My Molik — experimental Rust agent core\n\nUsage: oh-my-molik [--help | --version | theme list | theme show NAME]\n\nThe Rust model transport and tools are not implemented yet.\nThe theme commands display design tokens; there is no themed TUI yet.\nFor the read-only Go prototype, see ../../ohmyfox/README.md.");
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    match args
        .iter()
        .map(String::as_str)
        .collect::<Vec<_>>()
        .as_slice()
    {
        [] | ["--help" | "-h"] => help(),
        ["--version" | "-V"] => println!(
            "oh-my-molik {} (experimental core)",
            env!("CARGO_PKG_VERSION")
        ),
        ["theme", "list"] => {
            for theme in THEMES {
                println!("{}", theme.name);
            }
        }
        ["theme", "show", name] => {
            let Some(t) = find(name) else {
                eprintln!("Unknown theme: {name}. Run oh-my-molik theme list.");
                std::process::exit(2);
            };
            println!("{}\nbackground  {}\nsurface     {}\ntext        {}\nmuted       {}\naccent      {}\nsuccess     {}\nwarning     {}\nerror       {}", t.name, t.background, t.surface, t.text, t.muted, t.accent, t.success, t.warning, t.error);
        }
        _ => {
            eprintln!("Unknown arguments. Run oh-my-molik --help.");
            std::process::exit(2);
        }
    }
}
