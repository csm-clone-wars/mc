use std::{fs, path::PathBuf, process::exit};

use crate::macro_parser::parse_interface_macro;

pub fn generate_header_file(src_path: &PathBuf, _out_path: &PathBuf) {
    let src = read_file(src_path);
    let mut header_content = "".to_owned();

    let mut x:  = 0;
    while x < src.len() {
        if src[x] <= b' ' {
            // Skip ASCII whitespace and control characters
        } else if src[x] == b'/' {
            x += 1;
            if x + 11 < src.len() && &src[x..x + 11] == "/@interface" {
                parse_interface_macro(src_path, &src, x + 11, &mut header_content);
                break;
            } else if x < src.len() && src[x] == b'/' {
                // Skip single-line comment that appears before interface macro at the top of file
                let mut newline = false;

                let mut y = x + 1;
                while y < src.len() {
                    if src[y] == b'\n' {
                        newline = true;
                        break;
                    }
                    y += 1;
                }

                if !newline {
                    println!("error: Comment before interface macro needs to end with newline!");
                    println!("-> '{}'", src_path.display());
                    exit(1);
                }

                x = y + 1;
                continue;
            } else {
                interface_macro_error();
            }
        } else {
            interface_macro_error();
        }
        x += 1;
    }
}

pub fn read_file(path: &PathBuf) -> String {
    fs::read_to_string(path).unwrap_or_else(|err| {
        println!("error: Cannot read file '{}': {}", path.display(), err);
        exit(1);
    })
}

pub fn interface_macro_error() {
    println!(
        "error: Each .c file should have a valid interface macro, with valid opening and closing, at the top of file."
    );
    println!("Even if its empty, it should be present!");
    exit(1);
}
