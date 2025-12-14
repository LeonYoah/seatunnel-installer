package agent

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewStartCmd creates the start command for agent
func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the agent daemon",
		Long:  `Start the agent daemon to manage SeaTunnel processes on this node.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Starting SeaTunnel agent...")
			// TODO: Implement agent startup
			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "/etc/seatunnel/agent.yaml", "Configuration file path")
	cmd.Flags().String("control-plane", "http://localhost:8080", "Control plane URL")

	return cmd
}

// NewVersionCmd creates the version command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("SeaTunnel Agent v0.1.0")
		},
	}
}
