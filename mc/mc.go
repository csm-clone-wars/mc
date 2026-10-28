package mc

import (
	"fmt"
	"os"
	"strings"
)

func GenerateHeaderFile(src string, out_path string) {
	interfaceMacroFound := false
	output := "#pragma once\n"
	x := 0
	for x < len(src) {
		if src[x] <= ' ' {
			// Skip ASCII whitespace and control characters
		} else if x+1 < len(src) && src[x] == '/' && src[x+1] == '/' {
			foundNewline := false
			y := x + 2
			for y < len(src) {
				if src[y] == '\n' {
					foundNewline = true
					break
				}
				y += 1
			}

			if foundNewline {
				interfaceStart := 0
				interfaceEnd := 0
				if strings.TrimRight(src[x:y], " ") == "//@interface" {
					interfaceMacroFound = true
					interfaceStart = y + 1
					
					for y < len(src){
						if src[y] == '/' && y + {
							
						}
						y +=1
					}

				}
				x = y + 1
				continue
			} else {
				break
			}
		} else if src[x] == '/' {
			// Skip multi-line comments
			if x+1 < len(src) && src[x] == '*' {

			}

		} else {
			fmt.Println("error: Interface macro should be present at the top of each .c file!")
			os.Exit(1)
		}
		x += 1
	}

	if !interfaceMacroFound {
		fmt.Println("error: Interface macro should be present at the top of each .c file!")
		os.Exit(1)
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
