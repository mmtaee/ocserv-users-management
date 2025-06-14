package cmd

import (
	"github.com/spf13/cobra"
	"ocserv-bakend/pkg/bootstrap"
)

var (
	username string
	password string
)

var createAdminUserCmd = &cobra.Command{
	Use:   "create-admin",
	Short: "create a new admin user",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.CreateSuperAdmin(username, password)
	},
}

func init() {
	createAdminUserCmd.Flags().StringVarP(&username, "username", "u", "", "Username for admin")
	createAdminUserCmd.Flags().StringVarP(&password, "password", "p", "", "Password for admin")
	if err := createAdminUserCmd.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
	if err := createAdminUserCmd.MarkFlagRequired("password"); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(createAdminUserCmd)
}
