package genmodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "genmodel",
	Short: "read / collect content from up to date migrations to generate models in go code.",
}

func init() {
	rootCmd.AddCommand(genModelCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
