//! Read-only best-effort workspace inspection. Not an OS sandbox.
use serde_json::{json, Value};
use std::fs::{self, File};
use std::io::Read;
use std::path::{Component, Path, PathBuf};

const MAX_FILE: u64 = 64 * 1024;
const MAX_ENTRIES: usize = 200;

fn is_git_metadata(name: &std::ffi::OsStr) -> bool {
    let name = name.to_string_lossy();
    let name = name
        .split(':')
        .next()
        .unwrap_or("")
        .trim_end_matches([' ', '.']);
    name.eq_ignore_ascii_case(".git")
}

pub struct Workspace {
    root: PathBuf,
}

impl Workspace {
    pub fn new(root: &Path) -> Result<Self, &'static str> {
        let root = fs::canonicalize(root).map_err(|_| "cannot resolve workspace")?;
        if root
            .components()
            .any(|part| matches!(part, Component::Normal(name) if is_git_metadata(name)))
        {
            return Err("workspace cannot be inside git metadata");
        }
        if !root.is_dir() {
            return Err("workspace must be an existing directory");
        }
        Ok(Self { root })
    }

    fn resolve(&self, path: &str) -> Result<PathBuf, &'static str> {
        if path.is_empty() || path.len() > 4096 || path.contains('\0') {
            return Err("invalid workspace path");
        }
        let mut result = self.root.clone();
        for part in Path::new(path).components() {
            match part {
                Component::CurDir => {}
                Component::Normal(name) => {
                    if is_git_metadata(name) {
                        return Err("git metadata is not available");
                    }
                    result.push(name);
                    let metadata =
                        fs::symlink_metadata(&result).map_err(|_| "path is unavailable")?;
                    if metadata.file_type().is_symlink() {
                        return Err("symlinks are not supported");
                    }
                }
                _ => return Err("path must be relative and cannot traverse parents"),
            }
        }
        Ok(result)
    }

    /// Model-visible failures are intentionally generic; file names and content
    /// returned by successful calls will be sent to the selected model endpoint.
    pub fn call(&self, name: &str, arguments: &str) -> String {
        self.try_call(name, arguments)
            .unwrap_or_else(|error| format!("error: {error}"))
    }

    /// Structured result for protocol adapters; unlike `call`, file content
    /// beginning with `error:` cannot be mistaken for an execution failure.
    pub fn try_call(&self, name: &str, arguments: &str) -> Result<String, &'static str> {
        let args: Value = serde_json::from_str(arguments).map_err(|_| "invalid tool arguments")?;
        let object = args.as_object().ok_or("invalid tool arguments")?;
        if object.len() != 1 {
            return Err("invalid tool arguments");
        }
        let path = object
            .get("path")
            .and_then(Value::as_str)
            .ok_or("invalid tool arguments")?;
        self.call_path(name, path)
    }

    pub(crate) fn call_path(&self, name: &str, path: &str) -> Result<String, &'static str> {
        let resolved = self.resolve(path)?;
        match name {
            "read_file" => {
                let info = fs::symlink_metadata(&resolved).map_err(|_| "path is unavailable")?;
                if !info.file_type().is_file() {
                    return Err("path is not a regular file");
                }
                let file = File::open(resolved).map_err(|_| "cannot open file")?;
                let mut bytes = Vec::new();
                file.take(MAX_FILE + 1)
                    .read_to_end(&mut bytes)
                    .map_err(|_| "cannot read file")?;
                if bytes.len() as u64 > MAX_FILE {
                    return Err("file exceeds size limit");
                }
                if bytes.contains(&0) {
                    return Err("file is not UTF-8 text");
                }
                String::from_utf8(bytes).map_err(|_| "file is not UTF-8 text")
            }
            "list_files" => {
                let dir = fs::read_dir(resolved).map_err(|_| "cannot read directory")?;
                let mut entries = Vec::new();
                let mut truncated = false;
                for entry in dir {
                    let entry = entry.map_err(|_| "cannot read directory")?;
                    if entries.len() >= MAX_ENTRIES {
                        truncated = true;
                        break;
                    }
                    let name = entry.file_name();
                    let name = name.to_string_lossy();
                    if is_git_metadata(entry.file_name().as_os_str()) {
                        continue;
                    }
                    let kind = entry.file_type().map_err(|_| "cannot inspect entry")?;
                    let kind = if kind.is_symlink() {
                        "symlink"
                    } else if kind.is_dir() {
                        "directory"
                    } else if kind.is_file() {
                        "file"
                    } else {
                        "other"
                    };
                    entries.push(json!({"name":name,"type":kind}));
                }
                entries.sort_by(|a, b| a["name"].as_str().cmp(&b["name"].as_str()));
                Ok(json!({"entries":entries,"truncated":truncated}).to_string())
            }
            _ => Err("unknown tool"),
        }
    }
}
