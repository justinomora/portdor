package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <name>",
	Short: "Stop a service gracefully (SIGTERM)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.StopService(args[0]); err != nil {
			return err
		}
		fmt.Printf("✓ Stopped '%s'\n", args[0])
		return nil
	},
}

var killCmd = &cobra.Command{
	Use:   "kill <name>",
	Short: "Force kill a service (SIGKILL)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.KillService(args[0]); err != nil {
			return err
		}
		fmt.Printf("✓ Force killed '%s'\n", args[0])
		return nil
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart <name>",
	Short: "Restart a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.RestartService(args[0]); err != nil {
			return err
		}
		fmt.Printf("✓ Restarted '%s'\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd, killCmd, restartCmd)
}
