/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addBlackCmd represents the addBlack command.
var addBlackCmd = &cobra.Command{
	Use:   "addBlack <ip>",
	Short: "addBlack ip",
	Long:  `Add to black list command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("usage: ab-client addBlack <ip>")
		} else {
			ret, err := AddToBlackList(ctx, client, args[0])
			CheckRetErr(ret, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addBlackCmd)
}
