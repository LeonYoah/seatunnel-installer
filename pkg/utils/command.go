package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// CommandResult holds the result of a command execution
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

// ExecuteCommand executes a shell command and returns the result
func ExecuteCommand(command string, args ...string) *CommandResult {
	return ExecuteCommandWithTimeout(0, command, args...)
}

// ExecuteCommandWithTimeout executes a command with a timeout
func ExecuteCommandWithTimeout(timeout time.Duration, command string, args ...string) *CommandResult {
	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	return ExecuteCommandWithContext(ctx, command, args...)
}

// ExecuteCommandWithContext executes a command with a context
func ExecuteCommandWithContext(ctx context.Context, command string, args ...string) *CommandResult {
	cmd := exec.CommandContext(ctx, command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	result := &CommandResult{}

	err := cmd.Run()
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	if err != nil {
		result.Error = err
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	} else {
		result.ExitCode = 0
	}

	return result
}

// ExecuteShellCommand executes a shell command string (e.g., "ls -la | grep test")
// On Windows, it uses cmd.exe; on Unix-like systems, it uses sh
func ExecuteShellCommand(command string) *CommandResult {
	return ExecuteShellCommandWithTimeout(0, command)
}

// ExecuteShellCommandWithTimeout executes a shell command with timeout
// On Windows, it uses cmd.exe /c; on Unix-like systems, it uses sh -c
func ExecuteShellCommandWithTimeout(timeout time.Duration, command string) *CommandResult {
	shell, shellArg := getShellCommand()
	return ExecuteCommandWithTimeout(timeout, shell, shellArg, command)
}

// CommandStreamConfig holds configuration for streaming command execution
type CommandStreamConfig struct {
	Command       string
	Args          []string
	StdoutHandler func(string)
	StderrHandler func(string)
	Timeout       time.Duration
}

// ExecuteCommandWithStream executes a command and streams output to handlers
func ExecuteCommandWithStream(config *CommandStreamConfig) error {
	var ctx context.Context
	var cancel context.CancelFunc

	if config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), config.Timeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	cmd := exec.CommandContext(ctx, config.Command, config.Args...)

	// Set up stdout pipe
	if config.StdoutHandler != nil {
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to create stdout pipe: %w", err)
		}

		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stdoutPipe.Read(buf)
				if n > 0 {
					config.StdoutHandler(string(buf[:n]))
				}
				if err != nil {
					if err != io.EOF {
						config.StdoutHandler(fmt.Sprintf("Error reading stdout: %v\n", err))
					}
					break
				}
			}
		}()
	}

	// Set up stderr pipe
	if config.StderrHandler != nil {
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("failed to create stderr pipe: %w", err)
		}

		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stderrPipe.Read(buf)
				if n > 0 {
					config.StderrHandler(string(buf[:n]))
				}
				if err != nil {
					if err != io.EOF {
						config.StderrHandler(fmt.Sprintf("Error reading stderr: %v\n", err))
					}
					break
				}
			}
		}()
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

// CommandExists checks if a command is available in PATH
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// GetCommandPath returns the full path of a command
func GetCommandPath(command string) (string, error) {
	return exec.LookPath(command)
}

// ExecuteBatch executes multiple commands sequentially
func ExecuteBatch(commands []string) []*CommandResult {
	results := make([]*CommandResult, len(commands))
	for i, cmd := range commands {
		results[i] = ExecuteShellCommand(cmd)
	}
	return results
}

// ExecuteBatchWithStopOnError executes commands and stops on first error
func ExecuteBatchWithStopOnError(commands []string) ([]*CommandResult, error) {
	results := make([]*CommandResult, 0, len(commands))
	for _, cmd := range commands {
		result := ExecuteShellCommand(cmd)
		results = append(results, result)
		if result.Error != nil {
			return results, fmt.Errorf("command failed: %s, error: %w", cmd, result.Error)
		}
	}
	return results, nil
}

// SplitCommand splits a command string into command and arguments
func SplitCommand(command string) (string, []string) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", nil
	}
	if len(parts) == 1 {
		return parts[0], nil
	}
	return parts[0], parts[1:]
}

// BuildCommand builds a command string from command and arguments
func BuildCommand(command string, args ...string) string {
	if len(args) == 0 {
		return command
	}
	return command + " " + strings.Join(args, " ")
}

// getShellCommand returns the appropriate shell command for the current platform
// Windows: cmd.exe /c
// Unix-like: sh -c
func getShellCommand() (string, string) {
	// Check if we're on Windows by looking for common Windows commands
	if CommandExists("cmd.exe") || CommandExists("cmd") {
		return "cmd", "/c"
	}
	// Default to Unix shell
	return "sh", "-c"
}
