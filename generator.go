package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateHeader processes a .c source file and writes the corresponding .h file to includeDir.
func GenerateHeader(cFilePath, includeDir string) (string, error) {
	if filepath.Ext(cFilePath) != ".c" {
		return "", fmt.Errorf("%s: not a .c file", cFilePath)
	}

	res, err := ParseCFile(cFilePath)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(includeDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", includeDir, err)
	}

	baseName := strings.TrimSuffix(filepath.Base(cFilePath), filepath.Ext(cFilePath))
	outPath := filepath.Join(includeDir, baseName+".h")

	if err := os.WriteFile(outPath, []byte(res.Content), 0644); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", outPath, err)
	}

	return outPath, nil
}
