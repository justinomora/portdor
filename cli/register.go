package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register <name>",
	Short: "Register a service with portdor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		port, _ := cmd.Flags().GetInt("port")
		command, _ := cmd.Flags().GetString("cmd")
		cwd, _ := cmd.Flags().GetString("cwd")
		project, _ := cmd.Flags().GetString("project")

		if command == "" {
			return fmt.Errorf("--cmd is required")
		}

		if err := apiClient.RegisterService(name, command, cwd, port, project); err != nil {
			return err
		}
		fmt.Printf("✓ Registered '%s'", name)
		if port > 0 {
			fmt.Printf(" on port %d", port)
		}
		fmt.Println()
		return nil
	},
}

var unregisterCmd = &cobra.Command{
	Use:   "unregister <name>",
	Short: "Remove a service from the registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.UnregisterService(args[0]); err != nil {
			return err
		}
		fmt.Printf("✓ Unregistered '%s'\n", args[0])
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update a registered service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fields := map[string]any{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			fields["name"] = v
		}
		if cmd.Flags().Changed("project") {
			v, _ := cmd.Flags().GetString("project")
			fields["project"] = v
		}
		if cmd.Flags().Changed("port") {
			v, _ := cmd.Flags().GetInt("port")
			fields["port"] = v
		}
		if cmd.Flags().Changed("cmd") {
			v, _ := cmd.Flags().GetString("cmd")
			fields["command"] = v
		}
		if cmd.Flags().Changed("cwd") {
			v, _ := cmd.Flags().GetString("cwd")
			fields["cwd"] = v
		}
		if len(fields) == 0 {
			return fmt.Errorf("no fields specified to update")
		}
		if err := apiClient.UpdateService(args[0], fields); err != nil {
			return err
		}
		fmt.Printf("✓ Updated '%s'\n", args[0])
		return nil
	},
}

func init() {
	registerCmd.Flags().Int("port", 0, "Port the service listens on")
	registerCmd.Flags().String("cmd", "", "Command to run (required)")
	registerCmd.Flags().String("cwd", "", "Working directory (default: current dir)")
	registerCmd.Flags().String("project", "", "Project label for grouping")

	updateCmd.Flags().String("name", "", "New service name")
	updateCmd.Flags().String("project", "", "New project label")
	updateCmd.Flags().Int("port", 0, "New port")
	updateCmd.Flags().String("cmd", "", "New command")
	updateCmd.Flags().String("cwd", "", "New working directory")

	rootCmd.AddCommand(registerCmd, unregisterCmd, updateCmd)
}
