#![forbid(unsafe_code)]

use oh_my_molik::chat;
use oh_my_molik::theme::{find, THEMES};
use oh_my_molik::workspace::Workspace;
use oh_my_molik::{Access, Action, Limits, Session, Tool};
use std::io::{self, Read};
use std::path::Path;

fn help() {
    println!("Oh My Molik — experimental read-only agent\n\nUsage:\n  oh-my-molik run --endpoint URL --model MODEL --workspace DIR [--max-turns N] [--prompt-stdin | PROMPT]\n  oh-my-molik mcp serve --workspace DIR\n  oh-my-molik theme list\n  oh-my-molik theme show NAME\n  oh-my-molik [--help | --version]\n\nOptional bearer token: OMM_API_KEY environment variable.\nUses system curl for bounded, non-streaming Chat Completions.\nWorkspace checks are best-effort and are not a sandbox.\nNo file editing or command execution is implemented.");
}

fn clean_terminal(text: &str) -> String {
    text.chars()
        .map(|ch| {
            if ch == '\n' || ch == '\t' || !ch.is_control() {
                ch
            } else {
                '�'
            }
        })
        .collect()
}

fn run(args: &[String]) -> Result<(), String> {
    let (mut endpoint, mut model, mut workspace) = (None, None, None);
    let mut max_turns = 8usize;
    let (mut prompt, mut stdin) = (None, false);
    let mut i = 0;
    while i < args.len() {
        match args[i].as_str() {
            "--prompt-stdin" => {
                if stdin {
                    return Err("duplicate --prompt-stdin".into());
                }
                stdin = true;
            }
            "--endpoint" | "--model" | "--workspace" | "--max-turns" => {
                let value = args.get(i + 1).ok_or("missing option value")?.clone();
                match args[i].as_str() {
                    "--endpoint" => {
                        if endpoint.replace(value).is_some() {
                            return Err("duplicate endpoint".into());
                        }
                    }
                    "--model" => {
                        if model.replace(value).is_some() {
                            return Err("duplicate model".into());
                        }
                    }
                    "--workspace" => {
                        if workspace.replace(value).is_some() {
                            return Err("duplicate workspace".into());
                        }
                    }
                    _ => {
                        max_turns = value.parse().map_err(|_| "invalid max turns")?;
                    }
                }
                i += 1;
            }
            arg if arg.starts_with('-') => return Err("unknown option".into()),
            arg => {
                if prompt.replace(arg.to_owned()).is_some() {
                    return Err("provide exactly one prompt argument or use --prompt-stdin".into());
                }
            }
        }
        i += 1;
    }
    let endpoint = endpoint.ok_or("--endpoint is required")?;
    let model = model.ok_or("--model is required")?;
    let workspace = workspace.ok_or("--workspace is required")?;
    chat::validate_endpoint(&endpoint).map_err(str::to_owned)?;
    if !(1..=32).contains(&max_turns) {
        return Err("max turns must be between 1 and 32".into());
    }
    if stdin && prompt.is_some() {
        return Err("do not combine prompt argument and --prompt-stdin".into());
    }
    if stdin {
        let mut bytes = Vec::new();
        io::stdin()
            .take(65537)
            .read_to_end(&mut bytes)
            .map_err(|_| "cannot read prompt")?;
        if bytes.len() > 65536 {
            return Err("prompt exceeds size limit".into());
        }
        prompt = Some(String::from_utf8(bytes).map_err(|_| "prompt must be UTF-8")?);
    }
    let prompt = prompt.ok_or("prompt is required")?;
    let workspace = Workspace::new(Path::new(&workspace)).map_err(str::to_owned)?;
    let mut session = Session::new(
        prompt,
        vec![
            Tool {
                name: "read_file".into(),
                access: Access::ReadOnly,
            },
            Tool {
                name: "list_files".into(),
                access: Access::ReadOnly,
            },
        ],
        Limits {
            max_turns,
            ..Limits::default()
        },
    )
    .map_err(|e| e.to_string())?;
    let key = match std::env::var("OMM_API_KEY") {
        Ok(key) => key,
        Err(std::env::VarError::NotPresent) => String::new(),
        Err(std::env::VarError::NotUnicode(_)) => {
            return Err("OMM_API_KEY must be valid UTF-8".into())
        }
    };
    if endpoint.starts_with("http://") {
        eprintln!("Warning: HTTP does not encrypt prompts, files or the API key. Prefer HTTPS or a protected tunnel.");
    }
    loop {
        match session.action() {
            Action::AskModel => {
                let reply = chat::complete(&endpoint, &model, &key, session.messages())?;
                session.accept(reply).map_err(|e| e.to_string())?;
            }
            Action::Execute(call) => {
                let call = session
                    .begin_execution(&call.id)
                    .map_err(|e| e.to_string())?;
                let result = workspace.call(&call.name, &call.arguments);
                session
                    .record_result(&call.id, result)
                    .map_err(|e| e.to_string())?;
            }
            Action::Finished(text) => {
                println!("{}", clean_terminal(&text));
                return Ok(());
            }
            _ => return Err("unexpected session state".into()),
        }
    }
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    let words: Vec<&str> = args.iter().map(String::as_str).collect();
    match words.as_slice() {
        [] | ["--help" | "-h"] => help(),
        ["--version" | "-V"] => println!(
            "oh-my-molik {} (experimental read-only)",
            env!("CARGO_PKG_VERSION")
        ),
        ["theme", "list"] => {
            for theme in THEMES {
                println!("{}", theme.name);
            }
        }
        ["theme", "show", name] => {
            let Some(t) = find(name) else {
                eprintln!("Unknown theme. Run oh-my-molik theme list.");
                std::process::exit(2);
            };
            println!("{}\nbackground  {}\nsurface     {}\ntext        {}\nmuted       {}\naccent      {}\nsuccess     {}\nwarning     {}\nerror       {}", t.name, t.background, t.surface, t.text, t.muted, t.accent, t.success, t.warning, t.error);
        }
        ["mcp", "serve", "--workspace", root] => {
            let result = Workspace::new(Path::new(root))
                .map_err(str::to_owned)
                .and_then(|workspace| oh_my_molik::mcp::serve_stdio(&workspace));
            if let Err(error) = result {
                eprintln!("oh-my-molik: {}", clean_terminal(&error));
                std::process::exit(1);
            }
        }
        ["run", ..] => {
            if let Err(error) = run(&args[1..]) {
                eprintln!("oh-my-molik: {}", clean_terminal(&error));
                std::process::exit(1);
            }
        }
        _ => {
            eprintln!("Unknown arguments. Run oh-my-molik --help.");
            std::process::exit(2);
        }
    }
}
