package main

import (
	"fmt"
	"os"

	"github.com/seatunnel/enterprise-platform/internal/installer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seatunnel-installer",
	Short: "SeaTunnel Enterprise Platform Installer",
	Long:  `A comprehensive installer for SeaTunnel clusters with web-based management.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(installer.NewServerCmd())
	rootCmd.AddCommand(installer.NewVersionCmd())
}
