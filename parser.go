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

// LexState represents the lexical state when scanning C/C++ source code.
type LexState int

const (
	lexNormal LexState = iota
	lexBlockComment
	lexLineComment
	lexRawString
	lexString
	lexChar
)

// CppLexer maintains the state of lexical analysis across lines.
type CppLexer struct {
	state      LexState
	rawDelim   string // e.g. ")delim\""
	lineContin bool   // true if previous line ended with backslash continuation
}

// NewCppLexer creates a new CppLexer initialized to normal state.
func NewCppLexer() *CppLexer {
	return &CppLexer{
		state: lexNormal,
	}
}

// isIdentChar checks if character is a valid identifier character.
func isIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// matchRawStringStart checks if line starting at index i begins a C++ raw string literal (e.g. R"(...)" or R"delim(...)delim").
func matchRawStringStart(line string, i int) (string, int, bool) {
	// Must not be immediately preceded by an identifier character (e.g. FOOR"...")
	if i > 0 && isIdentChar(line[i-1]) {
		return "", 0, false
	}

	pos := i
	// Check optional encoding prefix: u8, u, U, L
	if strings.HasPrefix(line[pos:], "u8") {
		pos += 2
	} else if strings.HasPrefix(line[pos:], "u") || strings.HasPrefix(line[pos:], "U") || strings.HasPrefix(line[pos:], "L") {
		pos += 1
	}

	// Must be followed by R"
	if !strings.HasPrefix(line[pos:], `R"`) {
		return "", 0, false
	}
	pos += 2 // skip R"

	// Find '(' within at most 16 characters (C++ standard limit for raw string delimiter)
	parenIdx := -1
	for j := pos; j < len(line) && j <= pos+16; j++ {
		c := line[j]
		if c == '(' {
			parenIdx = j
			break
		}
		// Delimiter cannot contain space, control chars, '(', ')', or '\'
		if c <= 32 || c == 127 || c == ')' || c == '\\' {
			return "", 0, false
		}
	}

	if parenIdx == -1 {
		return "", 0, false
	}

	delim := line[pos:parenIdx]
	return delim, parenIdx + 1, true
}

// CheckLineBeforePub checks a line when in stateBeforePub.
// It returns whether the line contains the start of @pub, unexpected @end, or unexpected non-comment code.
func (l *CppLexer) CheckLineBeforePub(line string) (isPub bool, isEnd bool, hasCode bool) {
	if l.lineContin {
		l.lineContin = false
		if l.state == lexLineComment {
			if strings.HasSuffix(strings.TrimRight(line, " \t\r"), "\\") {
				l.lineContin = true
			} else {
				l.state = lexNormal
			}
			return false, false, false
		}
	}

	i := 0
	n := len(line)

	for i < n {
		switch l.state {
		case lexBlockComment:
			if i+1 < n && line[i] == '*' && line[i+1] == '/' {
				l.state = lexNormal
				i += 2
				continue
			}
			i++

		case lexNormal:
			c := line[i]
			// Skip whitespace
			if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
				i++
				continue
			}

			// Single-line comment //
			if i+1 < n && line[i] == '/' && line[i+1] == '/' {
				remainder := strings.TrimSpace(line[i:])
				if isPubStart(remainder) {
					isPub = true
					return
				}
				if isPubEnd(remainder) {
					isEnd = true
					return
				}
				// Regular single-line comment
				if strings.HasSuffix(strings.TrimRight(line, " \t\r"), "\\") {
					l.lineContin = true
					l.state = lexLineComment
				}
				return // Rest of line is comment
			}

			// Block comment /*
			if i+1 < n && line[i] == '/' && line[i+1] == '*' {
				l.state = lexBlockComment
				i += 2
				continue
			}

			// Any other character in lexNormal before @pub is code
			hasCode = true
			return
		}
	}

	return
}

// ScanLine updates the lexer state by scanning characters in the line.
func (l *CppLexer) ScanLine(line string) {
	if l.lineContin {
		l.lineContin = false
		if l.state == lexLineComment {
			if strings.HasSuffix(strings.TrimRight(line, " \t\r"), "\\") {
				l.lineContin = true
			} else {
				l.state = lexNormal
			}
			return
		}
	}

	i := 0
	n := len(line)
	for i < n {
		switch l.state {
		case lexNormal:
			// Check for raw string literal
			if delim, contentStart, ok := matchRawStringStart(line, i); ok {
				l.state = lexRawString
				l.rawDelim = ")" + delim + "\""
				i = contentStart
				continue
			}

			// Check for line comment //
			if i+1 < n && line[i] == '/' && line[i+1] == '/' {
				if strings.HasSuffix(strings.TrimRight(line, " \t\r"), "\\") {
					l.lineContin = true
					l.state = lexLineComment
				}
				return
			}

			// Check for block comment /*
			if i+1 < n && line[i] == '/' && line[i+1] == '*' {
				l.state = lexBlockComment
				i += 2
				continue
			}

			// Check for string literal "
			if line[i] == '"' {
				l.state = lexString
				i++
				continue
			}

			// Check for character literal '
			if line[i] == '\'' {
				l.state = lexChar
				i++
				continue
			}

			i++

		case lexBlockComment:
			if i+1 < n && line[i] == '*' && line[i+1] == '/' {
				l.state = lexNormal
				i += 2
				continue
			}
			i++

		case lexRawString:
			if strings.HasPrefix(line[i:], l.rawDelim) {
				l.state = lexNormal
				i += len(l.rawDelim)
				l.rawDelim = ""
				continue
			}
			i++

		case lexString:
			if line[i] == '\\' {
				if i+1 < n {
					i += 2
				} else {
					i++
				}
				continue
			}
			if line[i] == '"' {
				l.state = lexNormal
				i++
				continue
			}
			i++

		case lexChar:
			if line[i] == '\\' {
				if i+1 < n {
					i += 2
				} else {
					i++
				}
				continue
			}
			if line[i] == '\'' {
				l.state = lexNormal
				i++
				continue
			}
			i++
		}
	}

	// At end of line: check if string or char literal had backslash line continuation
	if l.state == lexString || l.state == lexChar {
		if strings.HasSuffix(strings.TrimRight(line, " \t\r"), "\\") {
			l.lineContin = true
		} else {
			// Unescaped newline terminates the string/char in C/C++
			l.state = lexNormal
		}
	}
}

// ParseFile reads a C or C++ source file and extracts the content between //@pub and //@end.
func ParseFile(path string) (*ParseResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", path, err)
	}
	defer file.Close()

	return ParseReader(path, file)
}

// ParseCFile is an alias for ParseFile for backward compatibility.
func ParseCFile(path string) (*ParseResult, error) {
	return ParseFile(path)
}

// ParseReader parses C or C++ source from an io.Reader.
func ParseReader(filename string, r io.Reader) (*ParseResult, error) {
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
	lexer := NewCppLexer()

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
			isPub, isEnd, hasCode := lexer.CheckLineBeforePub(rawLine)
			if hasCode {
				return nil, fmt.Errorf("%s:%d: unexpected content before //@pub macro: each source file must have @pub macro at the top (only whitespace and comments allowed)", filename, lineNum)
			}
			if isEnd {
				return nil, fmt.Errorf("%s:%d: unexpected //@end before //@pub macro", filename, lineNum)
			}
			if isPub {
				state = stateInPub
				pubStartLine = lineNum
				lexer.state = lexNormal
			}

		case stateInPub:
			initialState := lexer.state
			if initialState == lexNormal {
				if isPubEnd(trimmed) {
					state = stateAfterPub
					continue
				}
				if isPubStart(trimmed) {
					return nil, fmt.Errorf("%s:%d: nested //@pub macro is not allowed (already opened at line %d)", filename, lineNum, pubStartLine)
				}
			}

			// Add line as-is into header content
			headerLines = append(headerLines, rawLine)
			lexer.ScanLine(rawLine)

		case stateAfterPub:
			initialState := lexer.state
			if initialState == lexNormal {
				if isPubStart(trimmed) {
					return nil, fmt.Errorf("%s:%d: multiple //@pub macros are not allowed; only one @pub macro at the top of file is permitted", filename, lineNum)
				}
				if isPubEnd(trimmed) {
					return nil, fmt.Errorf("%s:%d: unexpected //@end macro after closed @pub block", filename, lineNum)
				}
			}
			lexer.ScanLine(rawLine)
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

// ParseCReader is an alias for ParseReader for backward compatibility.
func ParseCReader(filename string, r io.Reader) (*ParseResult, error) {
	return ParseReader(filename, r)
}

func isPubStart(trimmed string) bool {
	if trimmed == "//@pub" {
		return true
	}
	if strings.HasPrefix(trimmed, "//") {
		return strings.TrimSpace(trimmed[2:]) == "@pub"
	}
	return false
}

func isPubEnd(trimmed string) bool {
	if trimmed == "//@end" {
		return true
	}
	if strings.HasPrefix(trimmed, "//") {
		return strings.TrimSpace(trimmed[2:]) == "@end"
	}
	return false
}
