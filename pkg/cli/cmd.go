package cli

import (
	"github.com/haevg-rz/git-file-downloader/pkg/cli/file"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/folder"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/version"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"github.com/spf13/cobra"
)

const (
	FlagOutPath        = "out-path"
	FlagBranch         = "branch"
	FlagProjectNumber  = "project-number"
	FlagIncludePattern = "include-pattern"
	FlagExcludePattern = "exclude-pattern"
	FlagLogLevel       = "log-level"

	FlagPrivateToken = "token"
	FlagUrl          = "url"
	FlagUserAgent    = "user-agent"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "gdown",
	Short: "git file downloader",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// set global log level
		log.Level = options.Current.LogLevel

		return validate.Flags(map[string]interface{}{
			FlagPrivateToken:  options.Current.Api.PrivateToken,
			FlagProjectNumber: options.Current.Api.ProjectNumber,
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.V(1).Println("please use a subcommand...")
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
	// application relevant flags
	rootCmd.PersistentFlags().StringVar(&options.Current.OutPath, FlagOutPath, options.Current.OutPath, "Path to write file to disk")
	rootCmd.PersistentFlags().StringVar(&options.Current.Branch, FlagBranch, options.Current.Branch, "Branch name")
	rootCmd.PersistentFlags().StringVar(&options.Current.IncludePattern, FlagIncludePattern, options.Current.IncludePattern, "Include this regex pattern")
	rootCmd.PersistentFlags().StringVar(&options.Current.ExcludePattern, FlagExcludePattern, options.Current.ExcludePattern, "Exclude this regex pattern")

	// log level flag
	rootCmd.PersistentFlags().IntVar(&options.Current.LogLevel, FlagLogLevel, options.Current.LogLevel, "Higher loglevel leads to more verbose logging. Set log level to 0 if you dont want any logging.")

	// api flag def
	rootCmd.PersistentFlags().IntVar(&options.Current.Api.ProjectNumber, FlagProjectNumber, options.Current.Api.ProjectNumber, "The Project ID from your project")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.PrivateToken, FlagPrivateToken, options.Current.Api.PrivateToken, "Private-Token with access right for \"api\" and \"read_repository\", role must be minimum \"Reporter\"")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.UserAgent, FlagUserAgent, options.Current.Api.UserAgent, "User agent")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.ApiBaseUrl, FlagUrl, options.Current.Api.ApiBaseUrl, "Url to Api v4, like https://my-git-lab-server.local/api/v4/")
}
