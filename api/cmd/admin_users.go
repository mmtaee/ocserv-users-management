package cmd

import (
	"api/pkg/bootstrap"
	"github.com/spf13/cobra"
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
