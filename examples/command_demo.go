package main

import (
	"fmt"
	"runtime"

	"github.com/seatunnel/enterprise-platform/pkg/utils"
)

func main() {
	fmt.Println("=== 命令执行演示 ===")
	fmt.Printf("操作系统: %s\n", runtime.GOOS)
	fmt.Printf("架构: %s\n\n", runtime.GOARCH)

	// 1. 直接执行命令（跨平台）
	fmt.Println("1. 直接执行 Go 命令:")
	result := utils.ExecuteCommand("go", "version")
	if result.Error == nil {
		fmt.Printf("   输出: %s\n", result.Stdout)
	} else {
		fmt.Printf("   错误: %v\n", result.Error)
	}

	// 2. 检查命令是否存在
	fmt.Println("\n2. 检查命令是否存在:")
	commands := []string{"go", "git", "docker", "kubectl"}
	for _, cmd := range commands {
		exists := utils.CommandExists(cmd)
		status := "❌ 不存在"
		if exists {
			status = "✅ 存在"
			if path, err := utils.GetCommandPath(cmd); err == nil {
				status += fmt.Sprintf(" (%s)", path)
			}
		}
		fmt.Printf("   %s: %s\n", cmd, status)
	}

	// 3. 平台特定命令
	fmt.Println("\n3. 平台特定命令:")
	if runtime.GOOS == "windows" {
		// Windows 命令
		fmt.Println("   执行 Windows 命令:")

		// 列出当前目录
		result := utils.ExecuteCommand("cmd", "/c", "dir")
		if result.Error == nil {
			fmt.Printf("   dir 命令输出前 200 字符:\n   %s\n",
				truncate(result.Stdout, 200))
		}

		// 查看环境变量
		result = utils.ExecuteShellCommand("echo %USERPROFILE%")
		if result.Error == nil {
			fmt.Printf("   用户目录: %s\n", result.Stdout)
		}

		// PowerShell 命令
		result = utils.ExecuteCommand("powershell", "-Command", "Get-Date")
		if result.Error == nil {
			fmt.Printf("   当前时间: %s\n", result.Stdout)
		}
	} else {
		// Unix 命令
		fmt.Println("   执行 Unix 命令:")

		// 列出当前目录
		result := utils.ExecuteCommand("ls", "-la")
		if result.Error == nil {
			fmt.Printf("   ls 命令输出前 200 字符:\n   %s\n",
				truncate(result.Stdout, 200))
		}

		// 查看环境变量
		result = utils.ExecuteShellCommand("echo $HOME")
		if result.Error == nil {
			fmt.Printf("   用户目录: %s\n", result.Stdout)
		}
	}

	// 4. Shell 命令（自动适配平台）
	fmt.Println("\n4. Shell 命令（自动适配）:")
	result = utils.ExecuteShellCommand("echo Hello from shell")
	if result.Error == nil {
		fmt.Printf("   输出: %s\n", result.Stdout)
	} else {
		fmt.Printf("   错误: %v\n", result.Error)
	}

	// 5. 命令构建和拆分
	fmt.Println("\n5. 命令工具函数:")
	cmdStr := "git clone https://github.com/example/repo.git"
	cmd, args := utils.SplitCommand(cmdStr)
	fmt.Printf("   原始命令: %s\n", cmdStr)
	fmt.Printf("   拆分后 - 命令: %s, 参数: %v\n", cmd, args)

	rebuilt := utils.BuildCommand(cmd, args...)
	fmt.Printf("   重建命令: %s\n", rebuilt)

	// 6. 错误处理示例
	fmt.Println("\n6. 错误处理:")
	result = utils.ExecuteCommand("nonexistent_command_xyz")
	if result.Error != nil {
		fmt.Printf("   预期的错误: %v\n", result.Error)
		fmt.Printf("   退出码: %d\n", result.ExitCode)
	}

	fmt.Println("\n=== 演示完成 ===")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
