use std::{path::PathBuf, process::exit};

pub fn parse_interface_macro(src_path: &PathBuf, src: &str, y: usize, header_content: &mut String) {
    let mut y = y;

    let content_start: usize;
    let content_end: usize;

    let mut interface_opening_newline = false;
    while y < src.len() {
        if src[y] == b' ' {
            // Skip ASCII whitespace character
        } else if src[y] == b'\n' {
            interface_opening_newline = true;
            break;
        } else {
            println!("error: Invalid characters found in interface macro opening!");
            println!("-> '{}'", src_path.display());
            exit(1);
        }
        y += 1;
    }
    content_start = y;
    y += 1;

    let mut interface_closing_newline = false;
    while y < src.len() {
        if src[y] == b'\'' {
            let mut end_single_quote = false;

            let mut z = y + 1;
            while z < src.len() {
                if src[z] == b'\'' && (z + 1 < src.len()) {
                    if z - 1 != 0 && src[z] == b' ' {}
                    end_single_quote = true;
                    break;
                }
                z += 1;
            }
            if !end_single_quote {
                println!("error: Unterminated C character literal inside interface macro!");
                println!("-> {}", src_path.display());
                exit(1);
            }

            y = z + 1;
            continue;
        }
        y += 1;
    }
}
