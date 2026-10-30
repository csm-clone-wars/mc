package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateHeader_Success(t *testing.T) {
	tempDir := t.TempDir()
	cFile := filepath.Join(tempDir, "hello.c")
	includeDir := filepath.Join(tempDir, "include")

	cContent := `// Test hello
//@pub
void hello(void);
//@end

void hello(void) {}
`
	if err := os.WriteFile(cFile, []byte(cContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	outPath, err := GenerateHeader(cFile, includeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedOut := filepath.Join(includeDir, "hello.h")
	if outPath != expectedOut {
		t.Errorf("expected %q, got %q", expectedOut, outPath)
	}

	headerData, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read generated header: %v", err)
	}

	expectedHeader := "void hello(void);\n"
	if string(headerData) != expectedHeader {
		t.Errorf("expected %q, got %q", expectedHeader, string(headerData))
	}
}

func TestGenerateHeader_EmptyMacro(t *testing.T) {
	tempDir := t.TempDir()
	cFile := filepath.Join(tempDir, "main.c")
	includeDir := filepath.Join(tempDir, "include")

	cContent := `//@pub
//@end
int main() { return 0; }
`
	if err := os.WriteFile(cFile, []byte(cContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	outPath, err := GenerateHeader(cFile, includeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	headerData, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read generated header: %v", err)
	}

	if string(headerData) != "" {
		t.Errorf("expected empty header file, got %q", string(headerData))
	}
}

func TestGenerateHeader_NotCFile(t *testing.T) {
	tempDir := t.TempDir()
	txtFile := filepath.Join(tempDir, "test.txt")
	includeDir := filepath.Join(tempDir, "include")

	if err := os.WriteFile(txtFile, []byte("//@pub\n//@end"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := GenerateHeader(txtFile, includeDir)
	if err == nil {
		t.Fatal("expected error for non .c file, got nil")
	}
}
