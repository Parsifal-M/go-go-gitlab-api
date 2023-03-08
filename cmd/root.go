package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-gitlab",
	Short: "Go Gitlab is a CLI for Gitlab written in Go",
	Long:  `Go Gitlab is a CLI for Gitlab written in Go to interact with Gitlab API`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
