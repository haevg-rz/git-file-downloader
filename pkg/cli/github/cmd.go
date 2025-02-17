package github

import (
	"strings"

	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/github/options"
	globalOptions "github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"github.com/haevg-rz/git-file-downloader/pkg/logic"

	"github.com/spf13/cobra"
)

const (
	FlagOwner      = "owner"
	FlagRepo       = "repo"
	FlagApiVersion = "api-version"

	Endpoint = "https://api.github.com"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "github <file|folder> <flags>",
	Short: "retrieves data from github",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validate.Flags(map[string]interface{}{
			FlagOwner: options.Current.Owner,
			FlagRepo:  options.Current.Repo,
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.V(1).Println("retrieving files from github")

		var gitApi api.GitApi = api.NewGitHubApi(
			&api.BaseConfig{
				Url:       Endpoint,
				Auth:      globalOptions.Current.Api.Auth,
				UserAgent: globalOptions.Current.Api.UserAgent,
			},
			&api.GitHubConfig{
				Owner:      options.Current.Owner,
				Repo:       options.Current.Repo,
				ApiVersion: options.Current.ApiVersion,
			})

		return logic.NewGitFileDownloader(gitApi).Handle(&logic.Context{
			OutPath:    globalOptions.Current.OutPath,
			RemotePath: globalOptions.Current.RemotePath,
			Branch:     globalOptions.Current.Branch,
			Patterns: &logic.RegexRules{
				Include: globalOptions.Current.IncludePattern,
				Exclude: globalOptions.Current.ExcludePattern,
			},
		}, strings.ToLower(args[0]))
	},
}

func Command() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.Flags().StringVar(&options.Current.Owner, FlagOwner, options.Current.Owner, "repo owner")
	rootCmd.Flags().StringVar(&options.Current.Repo, FlagRepo, options.Current.Repo, "repo name")
	rootCmd.Flags().StringVar(&options.Current.ApiVersion, FlagApiVersion, options.Current.ApiVersion, "api version")
}
