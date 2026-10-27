package mc

import (
	"fmt"
	"os"
)

func GenerateHeaderFile(src string, out_path string) {
	output := "#pragma once \n"
	x := 0
	for x < len(src) {
		if src[x] <= ' ' {
			//Skip ASCII whitespace and control characters
		} else if src[x] == '/' {
			// need to handle multi-line and single line commetns and interface macero

		} else {
			fmt.Println("error: Interface macro should be present at the top of each .c file!")
			os.Exit(1)
		}
		x += 1
	}

	CreateFile(out_path, output)
}

func ReadFileToString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error: Cannot read file '%s'.\n", path)
		os.Exit(1)
	}
	return string(data)
}

func CreateFile(path string, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Printf("error: Unable to write to file '%s'\n", path)
		os.Exit(1)
	}
}
