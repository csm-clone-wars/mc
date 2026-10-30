package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HeaderTemplate represents the header guard directive placed at the top of generated header files.
const HeaderTemplate = "#pragma once\n"

// FormatHeader prepares the content for a generated header file,
// ensuring the #pragma once guard is present even if the content is empty.
func FormatHeader(content string) string {
	if strings.TrimSpace(content) == "" {
		return HeaderTemplate
	}

	if strings.HasPrefix(strings.TrimSpace(content), "#pragma once") {
		return content
	}

	if strings.HasPrefix(content, "\n") {
		return HeaderTemplate + content
	}

	return HeaderTemplate + "\n" + content
}

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

	headerContent := FormatHeader(res.Content)
	if err := os.WriteFile(outPath, []byte(headerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", outPath, err)
	}

	return outPath, nil
}
