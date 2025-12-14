package main

import (
	"fmt"
	"os"

	"github.com/seatunnel/enterprise-platform/internal/controlplane"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seatunnel-control-plane",
	Short: "SeaTunnel Enterprise Platform Control Plane",
	Long:  `Control Plane provides web UI and REST API for managing SeaTunnel clusters.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(controlplane.NewServerCmd())
	rootCmd.AddCommand(controlplane.NewMigrateCmd())
	rootCmd.AddCommand(controlplane.NewVersionCmd())
}
