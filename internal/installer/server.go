package installer

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewServerCmd creates the server command for installer
func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the installer API server",
		Long:  `Start the installer API server to provide web-based installation interface.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			fmt.Printf("Starting installer server on port %d...\n", port)
			// TODO: Implement server startup
			return nil
		},
	}

	cmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	cmd.Flags().StringP("config", "c", "config.yaml", "Configuration file path")

	return cmd
}

// NewVersionCmd creates the version command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("SeaTunnel Installer v0.1.0")
		},
	}
}
