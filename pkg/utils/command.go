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

// CommandResult 命令执行结果
type CommandResult struct {
	Stdout   string // 标准输出
	Stderr   string // 标准错误输出
	ExitCode int    // 退出码
	Error    error  // 错误信息
}

// ExecuteCommand 执行shell命令并返回结果
func ExecuteCommand(command string, args ...string) *CommandResult {
	return ExecuteCommandWithTimeout(0, command, args...)
}

// ExecuteCommandWithTimeout 执行带超时的命令
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

// ExecuteCommandWithContext 使用context执行命令
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

// ExecuteShellCommand 执行shell命令字符串（例如："ls -la | grep test"）
// 在Windows上使用cmd.exe，在类Unix系统上使用sh
func ExecuteShellCommand(command string) *CommandResult {
	return ExecuteShellCommandWithTimeout(0, command)
}

// ExecuteShellCommandWithTimeout 执行带超时的shell命令
// 在Windows上使用cmd.exe /c，在类Unix系统上使用sh -c
func ExecuteShellCommandWithTimeout(timeout time.Duration, command string) *CommandResult {
	shell, shellArg := getShellCommand()
	return ExecuteCommandWithTimeout(timeout, shell, shellArg, command)
}

// CommandStreamConfig 流式命令执行配置
type CommandStreamConfig struct {
	Command       string        // 命令
	Args          []string      // 参数
	StdoutHandler func(string)  // 标准输出处理函数
	StderrHandler func(string)  // 标准错误输出处理函数
	Timeout       time.Duration // 超时时间
}

// ExecuteCommandWithStream 执行命令并将输出流式传输到处理函数
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

	// 设置标准输出管道
	if config.StdoutHandler != nil {
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("创建标准输出管道失败: %w", err)
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
						config.StdoutHandler(fmt.Sprintf("读取标准输出错误: %v\n", err))
					}
					break
				}
			}
		}()
	}

	// 设置标准错误输出管道
	if config.StderrHandler != nil {
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("创建标准错误输出管道失败: %w", err)
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
						config.StderrHandler(fmt.Sprintf("读取标准错误输出错误: %v\n", err))
					}
					break
				}
			}
		}()
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动命令失败: %w", err)
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("命令执行失败: %w", err)
	}

	return nil
}

// CommandExists 检查命令是否在PATH中可用
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// GetCommandPath 返回命令的完整路径
func GetCommandPath(command string) (string, error) {
	return exec.LookPath(command)
}

// ExecuteBatch 顺序执行多个命令
func ExecuteBatch(commands []string) []*CommandResult {
	results := make([]*CommandResult, len(commands))
	for i, cmd := range commands {
		results[i] = ExecuteShellCommand(cmd)
	}
	return results
}

// ExecuteBatchWithStopOnError 执行命令，遇到第一个错误时停止
func ExecuteBatchWithStopOnError(commands []string) ([]*CommandResult, error) {
	results := make([]*CommandResult, 0, len(commands))
	for _, cmd := range commands {
		result := ExecuteShellCommand(cmd)
		results = append(results, result)
		if result.Error != nil {
			return results, fmt.Errorf("命令失败: %s, 错误: %w", cmd, result.Error)
		}
	}
	return results, nil
}

// SplitCommand 将命令字符串分割为命令和参数
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

// BuildCommand 从命令和参数构建命令字符串
func BuildCommand(command string, args ...string) string {
	if len(args) == 0 {
		return command
	}
	return command + " " + strings.Join(args, " ")
}

// getShellCommand 返回当前平台适用的shell命令
// Windows: cmd.exe /c
// 类Unix系统: sh -c
func getShellCommand() (string, string) {
	// 通过查找常见的Windows命令来检查是否在Windows上
	if CommandExists("cmd.exe") || CommandExists("cmd") {
		return "cmd", "/c"
	}
	// 默认使用Unix shell
	return "sh", "-c"
}
