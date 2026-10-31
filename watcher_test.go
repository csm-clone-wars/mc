package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatch_FileChanges(t *testing.T) {
	tempDir := t.TempDir()
	cFile := filepath.Join(tempDir, "sample.c")
	includeDir := filepath.Join(tempDir, "include")

	initialSrc := `//@pub
void sample_v1(void);
//@end
`
	if err := os.WriteFile(cFile, []byte(initialSrc), 0644); err != nil {
		t.Fatalf("failed to write initial c file: %v", err)
	}

	stopChan := make(chan struct{})
	watchDone := make(chan error, 1)

	go func() {
		watchDone <- Watch(WatchOptions{
			RootDir:    tempDir,
			IncludeDir: includeDir,
			StopChan:   stopChan,
		})
	}()

	// Wait for watcher to initialize and generate initial header
	time.Sleep(200 * time.Millisecond)

	headerPath := filepath.Join(includeDir, "sample.h")
	data, err := os.ReadFile(headerPath)
	if err != nil {
		t.Fatalf("expected header to exist initially: %v", err)
	}
	if string(data) != "#pragma once\n\nvoid sample_v1(void);\n" {
		t.Errorf("expected v1 header, got %q", string(data))
	}

	// Now modify the .c file
	updatedSrc := `//@pub
void sample_v2(int x);
//@end
`
	if err := os.WriteFile(cFile, []byte(updatedSrc), 0644); err != nil {
		t.Fatalf("failed to update c file: %v", err)
	}

	// Allow debounce and file event to process
	var updatedData []byte
	expectedV2 := "#pragma once\n\nvoid sample_v2(int x);\n"
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
		updatedData, _ = os.ReadFile(headerPath)
		if string(updatedData) == expectedV2 {
			break
		}
	}

	if string(updatedData) != expectedV2 {
		t.Fatalf("expected header to be updated to v2, got: %q", string(updatedData))
	}

	// Stop watcher
	close(stopChan)
	select {
	case err := <-watchDone:
		if err != nil {
			t.Errorf("watcher returned unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("timed out waiting for watcher to stop")
	}
}

func TestWatch_CppFileChanges(t *testing.T) {
	tempDir := t.TempDir()
	cppFile := filepath.Join(tempDir, "widget.cpp")
	includeDir := filepath.Join(tempDir, "include")

	initialSrc := `/* Widget initial */
//@pub
class WidgetV1 {};
//@end
`
	if err := os.WriteFile(cppFile, []byte(initialSrc), 0644); err != nil {
		t.Fatalf("failed to write initial cpp file: %v", err)
	}

	stopChan := make(chan struct{})
	watchDone := make(chan error, 1)

	go func() {
		watchDone <- Watch(WatchOptions{
			RootDir:    tempDir,
			IncludeDir: includeDir,
			StopChan:   stopChan,
		})
	}()

	// Wait for watcher to initialize and generate initial header
	time.Sleep(200 * time.Millisecond)

	headerPath := filepath.Join(includeDir, "widget.hpp")
	data, err := os.ReadFile(headerPath)
	if err != nil {
		t.Fatalf("expected header to exist initially: %v", err)
	}
	expectedV1 := "#pragma once\n\nclass WidgetV1 {};\n"
	if string(data) != expectedV1 {
		t.Errorf("expected v1 header, got %q", string(data))
	}

	// Now modify the .cpp file
	updatedSrc := `/* Widget updated */
//@pub
class WidgetV2 {};
//@end
`
	if err := os.WriteFile(cppFile, []byte(updatedSrc), 0644); err != nil {
		t.Fatalf("failed to update cpp file: %v", err)
	}

	// Allow debounce and file event to process
	var updatedData []byte
	expectedV2 := "#pragma once\n\nclass WidgetV2 {};\n"
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
		updatedData, _ = os.ReadFile(headerPath)
		if string(updatedData) == expectedV2 {
			break
		}
	}

	if string(updatedData) != expectedV2 {
		t.Fatalf("expected header to be updated to v2, got: %q", string(updatedData))
	}

	// Stop watcher
	close(stopChan)
	select {
	case err := <-watchDone:
		if err != nil {
			t.Errorf("watcher returned unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("timed out waiting for watcher to stop")
	}
}
