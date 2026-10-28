use std::path::Path;

pub struct Options {
    pub print_help: bool,
    pub print_version: bool,
    pub c_src_paths: Vec<String>,
    pub out_dir_path: String,
    pub errors: Vec<String>,
}

pub fn parse(args: &Vec<String>) -> Options {
    let mut options = Options {
        print_help: false,
        print_version: false,
        c_src_paths: vec![],
        out_dir_path: "".to_string(),
        errors: vec![],
    };

    let mut x = 0;
    while x < args.len() {
        if args[x].ends_with(".c") {
            options.c_src_paths.push(args[x].to_owned());
        } else if args[x] == "-version" || args[x] == "-v" {
            options.print_version = true;
        } else if args[x] == "-help" || args[x] == "-h" {
            options.print_help = true;
        } else if args[x] == "-o" {
            x += 1;
            if x < args.len() && Path::new(&args[x]).is_dir() {
                options.out_dir_path = args[x].to_owned();
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
    options
}
