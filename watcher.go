package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatchOptions configures the watch behavior.
type WatchOptions struct {
	Files      []string        // Specific files to watch, or empty for all .c files
	IncludeDir string          // Target directory for header files, default "include"
	RootDir    string          // Root directory to watch, default "."
	StopChan   <-chan struct{} // Optional channel to cleanly stop watcher (useful in tests)
}

// Watch watches for changes in .c files and regenerates .h files.
func Watch(opts WatchOptions) error {
	if opts.IncludeDir == "" {
		opts.IncludeDir = "include"
	}
	if opts.RootDir == "" {
		opts.RootDir = "."
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	watchSpecific := len(opts.Files) > 0
	targetFiles := make(map[string]bool)
	for _, f := range opts.Files {
		targetFiles[filepath.Clean(f)] = true
	}

	// Initial generation
	if watchSpecific {
		fmt.Println("[mc] generating headers for specified files...")
		for f := range targetFiles {
			outPath, err := GenerateHeader(f, opts.IncludeDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[mc] error in %s: %v\n", f, err)
			} else {
				fmt.Printf("[mc] generated %s\n", outPath)
			}
		}
	} else {
		fmt.Println("[mc] scanning directory and generating headers for all .c files...")
		cFiles, err := findCFiles(opts.RootDir)
		if err != nil {
			return fmt.Errorf("failed to scan for .c files: %w", err)
		}
		for _, f := range cFiles {
			outPath, err := GenerateHeader(f, opts.IncludeDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[mc] error in %s: %v\n", f, err)
			} else {
				fmt.Printf("[mc] generated %s\n", outPath)
			}
		}
	}

	// Add directories to watcher
	if watchSpecific {
		// Watch the directories containing the target files
		watchedDirs := make(map[string]bool)
		for f := range targetFiles {
			dir := filepath.Dir(f)
			if !watchedDirs[dir] {
				if err := watcher.Add(dir); err != nil {
					return fmt.Errorf("failed to watch directory %s: %w", dir, err)
				}
				watchedDirs[dir] = true
			}
		}
	} else {
		if err := addWatchDirectories(watcher, opts.RootDir); err != nil {
			return fmt.Errorf("failed to watch directories: %w", err)
		}
	}

	fmt.Println("[mc] watching for file changes... (Press Ctrl+C to stop)")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var (
		mu       sync.Mutex
		timers   = make(map[string]*time.Timer)
		debounce = 100 * time.Millisecond
	)

	for {
		select {
		case <-sigChan:
			fmt.Println("\n[mc] stopping watcher...")
			return nil

		case <-opts.StopChan:
			return nil

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "[mc] watcher error: %v\n", err)

		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// If a new directory is created, watch it
			if event.Op.Has(fsnotify.Create) {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					if !shouldSkipDir(event.Name) && !watchSpecific {
						_ = addWatchDirectories(watcher, event.Name)
					}
					continue
				}
			}

			// Only react to Write and Create events
			if !event.Op.Has(fsnotify.Write) && !event.Op.Has(fsnotify.Create) {
				continue
			}

			if filepath.Ext(event.Name) != ".c" {
				continue
			}

			cleanEventFile := filepath.Clean(event.Name)
			if watchSpecific && !targetFiles[cleanEventFile] {
				continue
			}

			// Debounce file changes
			mu.Lock()
			if t, exists := timers[cleanEventFile]; exists {
				t.Stop()
			}
			fileToRegen := cleanEventFile
			timers[cleanEventFile] = time.AfterFunc(debounce, func() {
				mu.Lock()
				delete(timers, fileToRegen)
				mu.Unlock()

				outPath, err := GenerateHeader(fileToRegen, opts.IncludeDir)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[mc] error in %s: %v\n", fileToRegen, err)
				} else {
					fmt.Printf("[mc] regenerated %s from %s\n", outPath, fileToRegen)
				}
			})
			mu.Unlock()
		}
	}
}

func findCFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if shouldSkipDir(path) && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".c" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func addWatchDirectories(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if shouldSkipDir(path) && path != root {
				return filepath.SkipDir
			}
			if err := watcher.Add(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func shouldSkipDir(path string) bool {
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return true
	}
	switch base {
	case "include", "build", "node_modules", "vendor":
		return true
	}
	return false
}
