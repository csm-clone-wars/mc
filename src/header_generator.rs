use std::{fs, path::PathBuf, process::exit};

pub fn generate_header_file(src_path: &PathBuf, _out_path: &PathBuf) {
    let src = read_file(src_path);
    dbg!(&src);

    let mut x = 0;
    while x < src.len() {
        x += 1;
    }
}

pub fn read_file(path: &PathBuf) -> String {
    fs::read_to_string(path).unwrap_or_else(|err| {
        println!("error: Cannot read file '{}': {}", path.display(), err);
        exit(1);
    })
}
