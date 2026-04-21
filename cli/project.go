package cli

import (
	"fmt"
	"os"

	"github.com/jmora/portdor/config"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Register all services from portdor.toml in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		configPath, err := config.FindProjectConfig(cwd)
		if err != nil {
			return fmt.Errorf("no portdor.toml found in current directory")
		}
		cfg, err := config.LoadProject(configPath)
		if err != nil {
			return err
		}

		for _, svc := range cfg.Services {
			err := apiClient.RegisterService(svc.Name, svc.Command, svc.Cwd, svc.Port, cfg.Project.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", svc.Name, err)
				continue
			}
			msg := fmt.Sprintf("  ✓ Registered: %s", svc.Name)
			if svc.Port > 0 {
				msg += fmt.Sprintf(" (port %d)", svc.Port)
			}
			fmt.Println(msg)
		}
		return nil
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and unregister all services from portdor.toml",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		configPath, err := config.FindProjectConfig(cwd)
		if err != nil {
			return fmt.Errorf("no portdor.toml found in current directory")
		}
		cfg, err := config.LoadProject(configPath)
		if err != nil {
			return err
		}

		for _, svc := range cfg.Services {
			apiClient.StopService(svc.Name)
			if err := apiClient.UnregisterService(svc.Name); err != nil {
				fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", svc.Name, err)
				continue
			}
			fmt.Printf("  ✓ Removed: %s\n", svc.Name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd, downCmd)
}
