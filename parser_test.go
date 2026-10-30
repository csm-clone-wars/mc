package main

import (
	"strings"
	"testing"
)

func TestParseCReader_HelloExample(t *testing.T) {
	src := `// Hello module example!

//@pub
void hello();
//@end

#include<stdio.h>

void hello(){
    puts("Hello world!");
}
`
	res, err := ParseCReader("hello.c", strings.NewReader(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "void hello();\n"
	if res.Content != expected {
		t.Errorf("expected %q, got %q", expected, res.Content)
	}
}

func TestParseCReader_EmptyPubMacro(t *testing.T) {
	src := `// Enven if empty each .c file should have the @pub macro at the top
//@pub
//@end

#include "hello.h"

int main(){
    hello();
    return 0;
}
`
	res, err := ParseCReader("main.c", strings.NewReader(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Content != "" {
		t.Errorf("expected empty content, got %q", res.Content)
	}
}

func TestParseCReader_SpacedSyntax(t *testing.T) {
	src := `   // single line comment with leading spaces
   // @pub
void foo(int a, int b);
   // @end
int x = 1;
`
	res, err := ParseCReader("spaced.c", strings.NewReader(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "void foo(int a, int b);\n"
	if res.Content != expected {
		t.Errorf("expected %q, got %q", expected, res.Content)
	}
}

func TestParseCReader_MultiLineHeaderContent(t *testing.T) {
	src := `//@pub
#ifndef MY_HEADER_H
#define MY_HEADER_H

typedef struct {
    int x;
    int y;
} Point;

Point create_point(int x, int y);

#endif
//@end

Point create_point(int x, int y) {
    Point p = {x, y};
    return p;
}
`
	res, err := ParseCReader("point.c", strings.NewReader(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `#ifndef MY_HEADER_H
#define MY_HEADER_H

typedef struct {
    int x;
    int y;
} Point;

Point create_point(int x, int y);

#endif
`
	if res.Content != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, res.Content)
	}
}

func TestParseCReader_Errors(t *testing.T) {
	tests := []struct {
		name        string
		src         string
		errContains string
	}{
		{
			name:        "empty file",
			src:         "",
			errContains: "missing //@pub macro",
		},
		{
			name:        "only comments and whitespace without pub",
			src:         "// comment\n\n   \n",
			errContains: "missing //@pub macro",
		},
		{
			name:        "code before pub",
			src:         "#include <stdio.h>\n//@pub\nvoid f();\n//@end\n",
			errContains: "unexpected content before //@pub macro",
		},
		{
			name:        "unclosed pub",
			src:         "//@pub\nvoid f();\n",
			errContains: "unclosed //@pub macro",
		},
		{
			name:        "unexpected end before pub",
			src:         "//@end\nvoid f();\n",
			errContains: "unexpected //@end before //@pub macro",
		},
		{
			name:        "nested pub",
			src:         "//@pub\n//@pub\nvoid f();\n//@end\n",
			errContains: "nested //@pub macro",
		},
		{
			name:        "multiple pub macros",
			src:         "//@pub\nvoid f1();\n//@end\nvoid f1(){}\n//@pub\nvoid f2();\n//@end\n",
			errContains: "multiple //@pub macros are not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCReader("test.c", strings.NewReader(tt.src))
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.errContains)
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("expected error to contain %q, but got: %v", tt.errContains, err)
			}
		})
	}
}
