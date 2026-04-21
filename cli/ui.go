package cli

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Open the portdor Web UI in your browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		url := "http://localhost:4242"
		fmt.Printf("Opening %s\n", url)
		var err error
		switch runtime.GOOS {
		case "darwin":
			err = exec.Command("open", url).Start()
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		default:
			fmt.Println("Open your browser and navigate to", url)
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
