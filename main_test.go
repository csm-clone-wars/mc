package main

import (
	"reflect"
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expected    Options
		expectError bool
	}{
		{
			name: "no arguments prints version and help",
			args: []string{},
			expected: Options{
				ShowVersion: true,
				ShowHelp:    true,
			},
		},
		{
			name: "flag -v",
			args: []string{"-v"},
			expected: Options{
				ShowVersion: true,
			},
		},
		{
			name: "flag -version",
			args: []string{"-version"},
			expected: Options{
				ShowVersion: true,
			},
		},
		{
			name: "flag -h",
			args: []string{"-h"},
			expected: Options{
				ShowHelp: true,
			},
		},
		{
			name: "flag -help",
			args: []string{"-help"},
			expected: Options{
				ShowHelp: true,
			},
		},
		{
			name: "flag -w",
			args: []string{"-w"},
			expected: Options{
				Watch: true,
			},
		},
		{
			name: "flag -watch",
			args: []string{"-watch"},
			expected: Options{
				Watch: true,
			},
		},
		{
			name: "files only",
			args: []string{"hi.c", "main.c"},
			expected: Options{
				Files: []string{"hi.c", "main.c"},
			},
		},
		{
			name: "watch with files",
			args: []string{"-w", "hi.c", "main.c"},
			expected: Options{
				Watch: true,
				Files: []string{"hi.c", "main.c"},
			},
		},
		{
			name:        "unknown flag",
			args:        []string{"-unknown"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := parseCommandLine(tt.args)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error for args %v, got nil", tt.args)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for args %v: %v", tt.args, err)
			}

			if opts.ShowVersion != tt.expected.ShowVersion {
				t.Errorf("ShowVersion: expected %v, got %v", tt.expected.ShowVersion, opts.ShowVersion)
			}
			if opts.ShowHelp != tt.expected.ShowHelp {
				t.Errorf("ShowHelp: expected %v, got %v", tt.expected.ShowHelp, opts.ShowHelp)
			}
			if opts.Watch != tt.expected.Watch {
				t.Errorf("Watch: expected %v, got %v", tt.expected.Watch, opts.Watch)
			}
			if !reflect.DeepEqual(opts.Files, tt.expected.Files) && (len(opts.Files) != 0 || len(tt.expected.Files) != 0) {
				t.Errorf("Files: expected %v, got %v", tt.expected.Files, opts.Files)
			}
		})
	}
}
