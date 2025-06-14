package cmd

import (
	"github.com/spf13/cobra"
	"ocserv-bakend/pkg/bootstrap"
)

var debug bool

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the api server",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.Serve(debug)
	},
}

func init() {
	serverCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	rootCmd.AddCommand(serverCmd)
}
