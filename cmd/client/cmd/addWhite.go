/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addWhiteCmd represents the addWhite command
var addWhiteCmd = &cobra.Command{
	Use:   "addWhite <ip>",
	Short: "addWhite ip",
	Long:  `Add to white list command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("usage: ab-client addWhite <ip>")
		} else {
			ret, err := AddToWhiteList(ctx, client, args[0])
			CheckRetErr(ret, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addWhiteCmd)
}
