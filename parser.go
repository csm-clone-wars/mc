package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ParseResult contains the extracted header content.
type ParseResult struct {
	Content string
}

// ParseCFile reads a .c file and extracts the content between //@pub and //@end.
func ParseCFile(path string) (*ParseResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", path, err)
	}
	defer file.Close()

	return ParseCReader(path, file)
}

// ParseCReader parses C source from an io.Reader.
func ParseCReader(filename string, r io.Reader) (*ParseResult, error) {
	reader := bufio.NewReader(r)

	const (
		stateBeforePub = iota
		stateInPub
		stateAfterPub
	)

	state := stateBeforePub
	var headerLines []string
	lineNum := 0
	pubStartLine := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("%s: read error at line %d: %w", filename, lineNum+1, err)
		}

		// Handle EOF with no content read
		if err == io.EOF && len(line) == 0 {
			break
		}

		lineNum++
		// Strip trailing CR and LF
		rawLine := strings.TrimRight(line, "\r\n")
		trimmed := strings.TrimSpace(rawLine)

		switch state {
		case stateBeforePub:
			if trimmed == "" {
				// Blank lines are allowed before //@pub
				continue
			}

			if isPubStart(trimmed) {
				state = stateInPub
				pubStartLine = lineNum
				continue
			}

			if isPubEnd(trimmed) {
				return nil, fmt.Errorf("%s:%d: unexpected //@end before //@pub macro", filename, lineNum)
			}

			if strings.HasPrefix(trimmed, "//") {
				// Single line comments are allowed before //@pub
				continue
			}

			// Any other content before //@pub is a violation
			return nil, fmt.Errorf("%s:%d: unexpected content before //@pub macro: each .c file must have @pub macro at the top (only whitespace and single-line comments allowed)", filename, lineNum)

		case stateInPub:
			if isPubEnd(trimmed) {
				state = stateAfterPub
				continue
			}

			if isPubStart(trimmed) {
				return nil, fmt.Errorf("%s:%d: nested //@pub macro is not allowed (already opened at line %d)", filename, lineNum, pubStartLine)
			}

			// Add line as-is into header content
			headerLines = append(headerLines, rawLine)

		case stateAfterPub:
			if isPubStart(trimmed) {
				return nil, fmt.Errorf("%s:%d: multiple //@pub macros are not allowed; only one @pub macro at the top of file is permitted", filename, lineNum)
			}
			if isPubEnd(trimmed) {
				return nil, fmt.Errorf("%s:%d: unexpected //@end macro after closed @pub block", filename, lineNum)
			}
			// Regular C code is ignored
		}

		if err == io.EOF {
			break
		}
	}

	if state == stateBeforePub {
		return nil, fmt.Errorf("%s: missing //@pub macro at the top of file", filename)
	}

	if state == stateInPub {
		return nil, fmt.Errorf("%s:%d: unclosed //@pub macro (missing //@end)", filename, pubStartLine)
	}

	var content string
	if len(headerLines) > 0 {
		content = strings.Join(headerLines, "\n") + "\n"
	}

	return &ParseResult{Content: content}, nil
}

func isPubStart(trimmed string) bool {
	return trimmed == "//@pub" || trimmed == "// @pub"
}

func isPubEnd(trimmed string) bool {
	return trimmed == "//@end" || trimmed == "// @end"
}
