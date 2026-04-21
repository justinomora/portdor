package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/jmora/portdor/internal/client"
	"github.com/spf13/cobra"
)

var apiClient *client.Client

var rootCmd = &cobra.Command{
	Use:   "portdor",
	Short: "Local dev service registry and manager",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "serve" {
			return nil
		}
		apiClient = client.New(client.DefaultAddr)
		if !apiClient.IsReachable() {
			if err := autoStartServer(); err != nil {
				return fmt.Errorf("could not start portdor server: %w", err)
			}
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func autoStartServer() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(self, "serve")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return err
	}
	apiClient = client.New(client.DefaultAddr)
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		if apiClient.IsReachable() {
			return nil
		}
	}
	return fmt.Errorf("server did not become ready in time")
}
