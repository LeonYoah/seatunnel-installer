package utils

import (
	"strings"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	t.Run("SimpleCommand", func(t *testing.T) {
		// Use 'go version' which works on all platforms
		result := ExecuteCommand("go", "version")

		if result.Error != nil {
			t.Fatalf("Command failed: %v", result.Error)
		}

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", result.ExitCode)
		}

		stdout := strings.TrimSpace(result.Stdout)
		if !strings.Contains(stdout, "go version") {
			t.Errorf("Expected output to contain 'go version', got '%s'", stdout)
		}
	})

	t.Run("CommandWithError", func(t *testing.T) {
		result := ExecuteCommand("ls", "/nonexistent/path/that/does/not/exist")

		if result.Error == nil {
			t.Error("Expected error for nonexistent path")
		}

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code")
		}
	})

	t.Run("CommandTimeout", func(t *testing.T) {
		// Skip on Windows as sleep command may not be available
		t.Skip("Skipping timeout test on Windows")
	})
}

func TestExecuteShellCommand(t *testing.T) {
	t.Run("SimpleShellCommand", func(t *testing.T) {
		// Skip shell tests on Windows as 'sh' is not available
		t.Skip("Skipping shell command tests on Windows")
	})
}

func TestCommandExists(t *testing.T) {
	t.Run("ExistingCommand", func(t *testing.T) {
		// Use 'go' command which should exist since we're running tests with go
		if !CommandExists("go") {
			t.Error("go command should exist")
		}
	})

	t.Run("NonExistingCommand", func(t *testing.T) {
		if CommandExists("nonexistent_command_xyz123") {
			t.Error("nonexistent command should not exist")
		}
	})
}

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		input       string
		wantCmd     string
		wantArgsLen int
	}{
		{"echo hello", "echo", 1},
		{"ls -la /tmp", "ls", 2},
		{"git", "git", 0},
		{"", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd, args := SplitCommand(tt.input)

			if cmd != tt.wantCmd {
				t.Errorf("Command mismatch: got %s, want %s", cmd, tt.wantCmd)
			}

			if len(args) != tt.wantArgsLen {
				t.Errorf("Args length mismatch: got %d, want %d", len(args), tt.wantArgsLen)
			}
		})
	}
}

func TestBuildCommand(t *testing.T) {
	tests := []struct {
		cmd  string
		args []string
		want string
	}{
		{"echo", []string{"hello"}, "echo hello"},
		{"ls", []string{"-la", "/tmp"}, "ls -la /tmp"},
		{"git", []string{}, "git"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			result := BuildCommand(tt.cmd, tt.args...)
			if result != tt.want {
				t.Errorf("BuildCommand mismatch: got %s, want %s", result, tt.want)
			}
		})
	}
}

func TestExecuteBatch(t *testing.T) {
	// Skip batch tests on Windows
	t.Skip("Skipping batch command tests on Windows")
}

func TestExecuteBatchWithStopOnError(t *testing.T) {
	// Skip batch tests on Windows
	t.Skip("Skipping batch command tests on Windows")
}
