/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// delBlackCmd represents the delBlack command.
var delBlackCmd = &cobra.Command{
	Use:   "delBlack <ip>",
	Short: "delBlack ip",
	Long:  `Del from black list command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("usage: ab-client delBlack <ip>")
		} else {
			ret, err := DelFromBlackList(ctx, client, args[0])
			CheckRetErr(ret, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(delBlackCmd)
}
