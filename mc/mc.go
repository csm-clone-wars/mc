package mc

import (
	"fmt"
	"os"
)

func GenerateHeaderFile(src_content string, out_path string) {

}

func ReadFileToString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error: Cannot read file '%s'.\n", path)
		os.Exit(1)
	}
	return string(data)
}
