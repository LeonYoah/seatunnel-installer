package main

import (
	"fmt"
	"os"

	"github.com/seatunnel/enterprise-platform/internal/agent"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seatunnel-agent",
	Short: "SeaTunnel Enterprise Platform Agent",
	Long:  `Agent runs on each node to manage SeaTunnel processes and report status.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(agent.NewStartCmd())
	rootCmd.AddCommand(agent.NewVersionCmd())
}
