package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/anuragkumar19/uploadthing-go/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type CLIConfig struct {
	ApiKey  string
	TempDir string
}

var configFile string
var config *CLIConfig = nil
var UTApi *api.UploadthingApi

var rootCmd = &cobra.Command{
	Use:   "uploadthing",
	Short: "Uploadthing cli for managing files (not official)",
	Long: `File Uploads For Web Developers
Web dev is better than ever.
File uploads needed help catching up.

//TODO:add domain (For Go developers)
https://uploadthing.com/ (For Typescript developers)
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b, err := os.ReadFile(filepath.Join(home, ".uploadthing"))

	configFile = string(b)

	if err != nil || configFile == "" {
		apiKey := stringPrompt("Your Uploadthing api key/secret?")
		configFile = strings.Join([]string{apiKey, os.TempDir()}, ";")

		file, err := os.Create(filepath.Join(home, ".uploadthing"))
		if err != nil {
			log.Fatal("failed to create config file")
		}
		file.WriteString(configFile)
	}

	s := strings.Split(configFile, ";")

	if len(s) != 2 || s[0] == "" || s[1] == "" {
		log.Fatal("config file is not valid please reset config file: for help run `uploadthing help`")
	}

	config = &CLIConfig{
		ApiKey:  s[0],
		TempDir: s[1],
	}

	UTApi = api.NewWithConfig(&api.UploadthingApiConfig{
		ApiKey:  config.ApiKey,
		TempDir: config.TempDir,
	})
}

func stringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
