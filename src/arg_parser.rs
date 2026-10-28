use std::path::{Path, PathBuf};

pub struct Options {
    pub print_help: bool,
    pub print_version: bool,
    pub c_src_paths: Vec<PathBuf>,
    pub out_dir_path: PathBuf,
    pub out_header_paths: Vec<PathBuf>,
    pub errors: Vec<String>,
}

pub fn parse(args: &Vec<String>) -> Options {
    let mut options = Options {
        print_help: false,
        print_version: false,
        c_src_paths: vec![],
        out_dir_path: PathBuf::new(),
        out_header_paths: vec![],
        errors: vec![],
    };

    let mut x = 0;
    while x < args.len() {
        if args[x].ends_with(".c") {
            let path = PathBuf::from(&args[x]);
            if path.is_file() {
                options.c_src_paths.push(path);
            } else {
                let error = format!("Not a file '{}'", &args[x]);
                options.errors.push(error);
            }
        } else if args[x] == "-version" || args[x] == "-v" {
            options.print_version = true;
        } else if args[x] == "-help" || args[x] == "-h" {
            options.print_help = true;
        } else if args[x] == "-o" {
            x += 1;
            if x < args.len() && Path::new(&args[x]).is_dir() {
                options.out_dir_path = PathBuf::from(&args[x]);
                x += 1;
                continue;
            } else {
                let error = "-o flag needs a valid path to a directory!".to_owned();
                options.errors.push(error);
                break;
            }
        } else {
            let error = "Unknown command!".to_owned();
            options.errors.push(error);
            break;
        }
        x += 1;
    }

    if options.out_dir_path.is_empty() {
        for src_path in &options.c_src_paths {
            let Some(filename) = src_path.file_name() else {
                let error = format!("Missing filename for '{}'", src_path.display());
                options.errors.push(error);
                continue;
            };
            let out_header_path = options.out_dir_path.join(filename).with_extension("h");
            options.out_header_paths.push(out_header_path);
        }
    } else {
        for src_path in &options.c_src_paths {
            let Some(_filename) = src_path.file_name() else {
                let error = format!("Missing filename for '{}'", src_path.display());
                options.errors.push(error);
                continue;
            };
            let out_header_path = src_path.with_extension("h");
            options.out_header_paths.push(out_header_path);
        }
    }

    options
}
