/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// delWhiteCmd represents the delWhite command.
var delWhiteCmd = &cobra.Command{
	Use:   "delWhite <ip>",
	Short: "delWhite ip",
	Long:  `Del from white list command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("usage: ab-client delWhite <ip>")
		} else {
			ret, err := DelFromWhiteList(ctx, client, args[0])
			CheckRetErr(ret, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(delWhiteCmd)
}
