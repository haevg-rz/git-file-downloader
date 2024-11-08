package cli

import (
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/github"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/gitlab"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/version"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"github.com/spf13/cobra"
)

const (
	// SOURCE, DEST

	FlagOutPath    = "out-path"
	FlagRemotePath = "remote-path"

	FlagBranch         = "branch"
	FlagIncludePattern = "include"
	FlagExcludePattern = "exclude"
	FlagLogLevel       = "log-level"

	// API

	FlagAuthToken = "token"
	FlagUrl       = "url"
	FlagUserAgent = "user-agent"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:       "gdown",
	Short:     "gdown",
	Long:      "git-file-downloader",
	Args:      cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	ValidArgs: []string{"file", "folder"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// set global log level
		log.Level = options.Current.LogLevel

		// no flags needed for version subcommand
		if cmd.CalledAs() == "version" || cmd.CalledAs() == "help" {
			return nil
		}

		return validate.Flags(map[string]interface{}{
			FlagAuthToken:  options.Current.Api.Auth,
			FlagOutPath:    options.Current.OutPath,
			FlagRemotePath: options.Current.RemotePath,
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand..")
		return nil
	},
}

func Command() *cobra.Command {
	rootCmd.AddCommand(version.Command())
	rootCmd.AddCommand(github.Command())
	rootCmd.AddCommand(gitlab.Command())
	return rootCmd
}

func init() {
	// source, dest
	rootCmd.PersistentFlags().StringVar(&options.Current.OutPath, FlagOutPath, options.Current.OutPath, "Path to write file to disk")
	rootCmd.PersistentFlags().StringVar(&options.Current.RemotePath, FlagRemotePath, options.Current.RemotePath, "Path to file/folder from remote source")

	rootCmd.PersistentFlags().StringVar(&options.Current.Branch, FlagBranch, options.Current.Branch, "Branch name")
	rootCmd.PersistentFlags().StringVar(&options.Current.IncludePattern, FlagIncludePattern, options.Current.IncludePattern, "Include this regex pattern")
	rootCmd.PersistentFlags().StringVar(&options.Current.ExcludePattern, FlagExcludePattern, options.Current.ExcludePattern, "Exclude this regex pattern")
	rootCmd.PersistentFlags().IntVar(&options.Current.LogLevel, FlagLogLevel, options.Current.LogLevel, "Higher loglevel leads to more verbose logging. Set log level to 0 if you dont want any logging.")

	// api flag def
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.Auth, FlagAuthToken, options.Current.Api.Auth, "Private-Token with access right for \"api\" and \"read_repository\", role must be minimum \"Reporter\"")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.UserAgent, FlagUserAgent, options.Current.Api.UserAgent, "User agent")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.BaseUrl, FlagUrl, options.Current.Api.BaseUrl, "url to Api v4, like https://my-git-lab-server.local/api/v4/")
}
