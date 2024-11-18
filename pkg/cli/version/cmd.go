package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	AppName = "GitLab File Downloader"
	Version = "v1.0.0"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "version",
	Short: "version of git file downloader",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s %s\n", AppName, Version)
		return nil
	},
}

func Command() *cobra.Command {
	return rootCmd
}

func init() {
	// flag provider
}
