package cmd

import (
	"bytes"
	"encoding/json"
	"runtime"
	"testing"
)

func TestBuildInfoCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantJSON bool
	}{
		{"default output", []string{}, false},
		{"json output", []string{"--json"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetArgs(append([]string{"build-info"}, tt.args...))

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("build-info command failed: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Fatal("expected non-empty output")
			}

			if tt.wantJSON {
				var info BuildInfo
				if err := json.Unmarshal([]byte(output), &info); err != nil {
					t.Fatalf("invalid JSON output: %v\noutput: %s", err, output)
				}
				if info.Version != Version {
					t.Errorf("version = %q, want %q", info.Version, Version)
				}
				if info.GoVersion != runtime.Version() {
					t.Errorf("go_version = %q, want %q", info.GoVersion, runtime.Version())
				}
				if info.OS != runtime.GOOS {
					t.Errorf("os = %q, want %q", info.OS, runtime.GOOS)
				}
				if info.Arch != runtime.GOARCH {
					t.Errorf("arch = %q, want %q", info.Arch, runtime.GOARCH)
				}
			}
		})
	}
}
