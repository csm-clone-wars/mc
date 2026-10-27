package main

import (
	"fmt"
	"mc/mc"
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
			src_path := args[0]
			out_path := src_path[:len(src_path)-2] + ".h"
			mc.GenerateHeaderFile(src_path, out_path)
		} else {
			fmt.Println("error: Unknown command!")
			os.Exit(1)
		}
	} else if len(args) == 3 {
		fmt.Println("Under construction!")
		os.Exit(1)
		// src_path := ""
		// out_path := ""
		// y := 0
		// for y < len(args) {
		// 	if strings.HasSuffix(args[y], ".c") {
		// 		src_path = args[y]
		// 	} else if args[y] == "-o" {
		// 		y += 1
		// 		if y < len(args) {
		// 			out_path = args[y]
		// 			y += 1
		// 			continue
		// 		}
		// 	} else {
		// 		fmt.Println("error: Unknown command!")
		// 		os.Exit(1)
		// 	}
		// 	y += 1
		// }

		// if src_path == "" {
		// 	fmt.Println("error: No .c file provided!")
		// 	os.Exit(1)
		// } else if out_path == "" {
		// 	fmt.Println("error: No output path provided!")
		// 	os.Exit(1)
		// }

		// filename := src_path[:len(src_path)-2]
		// out_path += filename

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
	mc hello.c                    -> Generate the .h header file inside 'include/' folder
	mc hello.c -o include/hello/  -> Generate the .h header file inside the specified folder`
