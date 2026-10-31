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

// SupportedSourceExtensions maps supported source extensions to their corresponding header extension.
var SupportedSourceExtensions = map[string]string{
	".c":   ".h",
	".cpp": ".hpp",
	".c++": ".hpp",
	".cxx": ".hpp",
	".cc":  ".hpp",
}

// IsSupportedSource returns true if the file extension is supported.
func IsSupportedSource(filePath string) bool {
	ext := filepath.Ext(filePath)
	_, ok := SupportedSourceExtensions[ext]
	return ok
}

// HeaderExtForSource returns the header extension for a given source file, or empty string if unsupported.
func HeaderExtForSource(filePath string) string {
	ext := filepath.Ext(filePath)
	return SupportedSourceExtensions[ext]
}

// GenerateHeader processes a C or C++ source file and writes the corresponding header file (.h or .hpp) to includeDir.
func GenerateHeader(srcFilePath, includeDir string) (string, error) {
	headerExt := HeaderExtForSource(srcFilePath)
	if headerExt == "" {
		return "", fmt.Errorf("%s: not a supported source file (.c, .cpp, .c++, .cxx)", srcFilePath)
	}

	res, err := ParseFile(srcFilePath)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(includeDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", includeDir, err)
	}

	ext := filepath.Ext(srcFilePath)
	baseName := strings.TrimSuffix(filepath.Base(srcFilePath), ext)
	outPath := filepath.Join(includeDir, baseName+headerExt)

	headerContent := FormatHeader(res.Content)
	if err := os.WriteFile(outPath, []byte(headerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", outPath, err)
	}

	return outPath, nil
}
