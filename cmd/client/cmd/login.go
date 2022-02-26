/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command.
var loginCmd = &cobra.Command{
	Use:   "login <login> <pass> <ip>",
	Short: "Check login, password and ip",
	Long:  `Check login, password and ip`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("usage: ab-client login <login> <pass> <ip>")
		} else {
			ret, err := Login(ctx, client, args[0], args[1], args[2])
			CheckRetErr(ret, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
