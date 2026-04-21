package cli

import (
	"fmt"
	"os"

	"github.com/jmora/portdor/internal/registry"
	"github.com/jmora/portdor/internal/server"
	"github.com/jmora/portdor/internal/state"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the portdor server on :4242",
	RunE: func(cmd *cobra.Command, args []string) error {
		reg := registry.New()
		statePath := state.DefaultPath()

		if st, err := state.Load(statePath); err == nil {
			reg.Restore(st.Services)
			fmt.Fprintf(os.Stdout, "Restored %d services from state\n", len(st.Services))
		}

		st := &state.State{}
		srv := server.New(reg, st, statePath)

		fmt.Println("portdor server listening on :4242")
		return srv.ListenAndServe(":4242")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
