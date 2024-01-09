package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files to Uploadthing",
	Long: `Usage Example:
uploadthing upload ./image.jpg ./video.mp4
	`,
	Run: func(cmd *cobra.Command, args []string) {
		paths := []string{}

		for _, arg := range args {
			paths = append(paths, filepath.Join(".", arg))
		}

		u, err := UTApi.UploadFiles(paths)

		fmt.Println(u, err)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
