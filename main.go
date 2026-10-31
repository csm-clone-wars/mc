package main

import (
	"fmt"
	"os"
)

const (
	Version = "0.1.0"
	Usage   = `mc - meta-c: A pre-preprocessor for C and C++

Usage:
  mc [options] [files...]

Commands:
  mc                 Print mc version and help
  mc -version,  -v   Print mc version
  mc -help,     -h   Print mc commands list
  mc -watch,    -w   Watch for any file changes and regenerate header files
  mc <files...>      Generate header files in include/ folder for given source files (.c, .cpp, .c++, .cxx)
`
)

// Options holds CLI configurations parsed from command line arguments.
type Options struct {
	ShowVersion bool
	ShowHelp    bool
	Watch       bool
	Files       []string
}

func parseCommandLine(args []string) (Options, error) {
	var opts Options

	if len(args) == 0 {
		opts.ShowVersion = true
		opts.ShowHelp = true
		return opts, nil
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-v", "-version", "--version":
			opts.ShowVersion = true
		case "-h", "-help", "--help":
			opts.ShowHelp = true
		case "-w", "-watch", "--watch":
			opts.Watch = true
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return opts, fmt.Errorf("unknown option: %s", arg)
			}
			opts.Files = append(opts.Files, arg)
		}
	}

	return opts, nil
}

func printVersion() {
	fmt.Printf("mc version %s\n", Version)
}

func printHelp() {
	fmt.Print(Usage)
}

func main() {
	opts, err := parseCommandLine(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "mc: %v\n", err)
		printHelp()
		os.Exit(1)
	}

	if opts.ShowVersion && opts.ShowHelp {
		printVersion()
		fmt.Println()
		printHelp()
		os.Exit(0)
	}

	if opts.ShowVersion {
		printVersion()
		os.Exit(0)
	}

	if opts.ShowHelp {
		printHelp()
		os.Exit(0)
	}

	if opts.Watch {
		watchOpts := WatchOptions{
			Files:      opts.Files,
			IncludeDir: "include",
			RootDir:    ".",
		}
		if err := Watch(watchOpts); err != nil {
			fmt.Fprintf(os.Stderr, "mc: watch error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(opts.Files) == 0 {
		printVersion()
		fmt.Println()
		printHelp()
		os.Exit(0)
	}

	hasError := false
	for _, file := range opts.Files {
		outPath, err := GenerateHeader(file, "include")
		if err != nil {
			fmt.Fprintf(os.Stderr, "mc: error: %v\n", err)
			hasError = true
		} else {
			fmt.Printf("Generated %s\n", outPath)
		}
	}

	if hasError {
		os.Exit(1)
	}
}
