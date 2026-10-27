package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println(VERSION)
		fmt.Println(HELP)
		os.Exit(0)
	} else if len(args) == 1 {
		arg := args[0]
		if arg == "-version" || arg == "-v" {
			fmt.Println(VERSION)
			os.Exit(0)
		} else if arg == "-help" || arg == "-h" {
			fmt.Println(VERSION)
			fmt.Println(HELP)
			os.Exit(0)
		} else if strings.HasSuffix(arg, ".c") {
			arg := args[0]
			filename := arg[:len(arg)-2]
			headerFileName := filename + ".h"
			fmt.Println(headerFileName)
		} else {
			fmt.Println("error: Unknown command!")
			os.Exit(1)
		}
	} else if len(args) > 1 {

	} else {
		fmt.Println("error: Unknown command!")
		os.Exit(1)
	}

}

const VERSION = "mc 1.0.0"
const HELP = `A pre-preprocessor for C.
- author: csm
- license: MIT
- overview:
	It helps in generating header files for C files using a custom macro syntax.
- commands:
	mc <no args>                  -> Print mc commands list
	mc -version, -v               -> Print mc version
	mc -help,    -h               -> Print mc commands list
	mc hello.c                    -> Generate .h header file inside 'include/' folder
	mc hello.c -o include/hello/  -> Generate the .h header file inside the specified folder`
