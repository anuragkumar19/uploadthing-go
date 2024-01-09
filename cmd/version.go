package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uploadthing CLI v0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
