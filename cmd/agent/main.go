package main

import (
	"fmt"
	"os"

	"github.com/seatunnel/enterprise-platform/internal/agent"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seatunnel-agent",
	Short: "SeaTunnel 企业级平台统一代理",
	Long: `Agent是部署在每个节点上的统一守护进程，负责：
1. 安装管理：集群部署、卸载、升级、诊断
2. 进程管理：SeaTunnel进程生命周期管理、监控、日志收集
3. 运维管理：执行Control Plane下发的运维指令`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// 进程管理相关命令
	rootCmd.AddCommand(agent.NewStartCmd())  // 启动Agent守护进程
	rootCmd.AddCommand(agent.NewStopCmd())   // 停止Agent守护进程
	rootCmd.AddCommand(agent.NewStatusCmd()) // 查看Agent状态

	// 安装管理相关命令
	rootCmd.AddCommand(agent.NewInstallCmd())   // 安装SeaTunnel
	rootCmd.AddCommand(agent.NewUninstallCmd()) // 卸载SeaTunnel
	rootCmd.AddCommand(agent.NewUpgradeCmd())   // 升级SeaTunnel
	rootCmd.AddCommand(agent.NewPrecheckCmd())  // 环境预检查
	rootCmd.AddCommand(agent.NewDiagnoseCmd())  // 收集诊断信息

	// 通用命令
	rootCmd.AddCommand(agent.NewVersionCmd()) // 版本信息
}
