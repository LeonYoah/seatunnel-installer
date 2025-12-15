package controlplane

import (
	"fmt"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/server"
	"github.com/spf13/cobra"
)

// NewServerCmd creates the server command for control plane
func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the control plane server",
		Long:  `Start the control plane server to provide web UI and REST API.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")

			// 创建并启动服务器
			srv, err := server.NewServer(port)
			if err != nil {
				return err
			}

			return srv.Start()
		},
	}

	cmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	cmd.Flags().StringP("config", "c", "config.yaml", "Configuration file path")

	return cmd
}

// NewMigrateCmd creates the migrate command
func NewMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations and create default admin user",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running database migrations...")

			// 创建服务器实例以初始化数据库连接
			srv, err := server.NewServer(8080)
			if err != nil {
				return fmt.Errorf("初始化服务器失败: %w", err)
			}

			// 创建默认管理员用户
			if err := srv.CreateDefaultAdminUser(); err != nil {
				return fmt.Errorf("创建默认管理员用户失败: %w", err)
			}

			fmt.Println("✅ 数据库迁移完成")
			fmt.Println("✅ 默认管理员用户已创建")
			fmt.Println("   用户名: admin")
			fmt.Println("   密码: admin123")
			fmt.Println("   请登录后立即修改密码！")

			return nil
		},
	}
}

// NewVersionCmd creates the version command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("SeaTunnel Control Plane v0.1.0")
		},
	}
}
