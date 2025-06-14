package cmd

import (
	"github.com/spf13/cobra"
	"ocserv-bakend/pkg/bootstrap"
)

var adminUsersCmd = &cobra.Command{
	Use:   "admins",
	Short: "List all admin users",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.AdminUsers()
	},
}

func init() {
	rootCmd.AddCommand(adminUsersCmd)
}
