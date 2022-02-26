/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command.
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset login, password and ip",
	Long:  `Reset command - delete login, password and ip from all lists`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("usage: ab-client reset <login> <pass> <ip>")
		} else {
			ret, err := Reset(ctx, client, args[0], args[1], args[2])
			CheckRetErr(ret, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
