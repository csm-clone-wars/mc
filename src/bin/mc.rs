use std::process::exit;

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    let options = mc::parse(&args);

    if !options.errors.is_empty() {
        for error in options.errors {
            println!("error: {}", error);
        }
        exit(1);
    }

    if args.len() == 0 {
        println!("{}", MC_VERSION);
        println!("{}", MC_HELP);
        exit(0);
    }

    if args.len() == 1 {
        if options.print_help && options.print_version {
            println!("Malformed options!");
            exit(1);
        }

        if options.print_version {
            println!("{}", MC_VERSION);
            exit(0);
        } else if options.print_help {
            println!("{}", MC_VERSION);
            println!("{}", MC_HELP);
            exit(0);
        }
    }

    if options.out_dir_path.is_empty() {
    } else {
    }
}

const MC_VERSION: &str = "mc 1.0.0";
const MC_HELP: &str = r#"A pre-preprocessor for C.
- author: csm
- license: MIT
- overview:
	It helps in generating header files for C files using a custom macro syntax.
- commands:
	mc <no args>                  -> Print mc commands list
	mc -version, -v               -> Print mc version
	mc -help,    -h               -> Print mc commands list
	mc hello.c main.c             -> Generate the .h header files 
	mc hello.c -o include/        -> Generate the .h header files inside the specified folder
"#;
