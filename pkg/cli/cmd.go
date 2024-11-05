package cli

import (
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/file"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/folder"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/version"
	"github.com/spf13/cobra"
)

const (
	FlagOutPath        = "out-path"
	FlagBranch         = "branch"
	FlagProjectNumber  = "project-number"
	FlagIncludePattern = "include-pattern"
	FlagExcludePattern = "exclude-pattern"

	FlagPrivateToken = "token"
	FlagUrl          = "url"
	FlagUserAgent    = "user-agent"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "gdown",
	Short: "git file downloader",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		return nil
	},
}

func Command() *cobra.Command {
	rootCmd.AddCommand(version.Command())
	rootCmd.AddCommand(file.Command())
	rootCmd.AddCommand(folder.Command())

	return rootCmd
}

func init() {
	// global flags
	rootCmd.PersistentFlags().StringVar(&options.Current.OutPath, FlagOutPath, options.Current.OutPath, "Path to write file to disk")
	rootCmd.PersistentFlags().StringVar(&options.Current.Branch, FlagBranch, options.Current.Branch, "branch name")
	rootCmd.PersistentFlags().StringVar(&options.Current.IncludePattern, FlagIncludePattern, options.Current.IncludePattern, "include this regex pattern")
	rootCmd.PersistentFlags().StringVar(&options.Current.ExcludePattern, FlagExcludePattern, options.Current.ExcludePattern, "exclude this regex pattern")

	// api flag def
	rootCmd.PersistentFlags().IntVar(&options.Current.Api.ProjectNumber, FlagProjectNumber, options.Current.Api.ProjectNumber, "The Project ID from your project")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.PrivateToken, FlagPrivateToken, options.Current.Api.PrivateToken, "Private-Token with access right for \"api\" and \"read_repository\", role must be minimum \"Reporter\"")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.UserAgent, FlagUserAgent, options.Current.Api.UserAgent, "user agent")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.ApiBaseUrl, FlagUrl, options.Current.Api.ApiBaseUrl, "Url to Api v4, like https://my-git-lab-server.local/api/v4/")
}
