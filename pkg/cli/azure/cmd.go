package azure

import (
	"strings"

	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/azure/options"
	globalOptions "github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"github.com/haevg-rz/git-file-downloader/pkg/logic"

	"github.com/spf13/cobra"
)

const (
	FlagOrganization = "organization"
	FlagProject      = "project"
	FlagRepo         = "repo"
	FlagApiVersion   = "api-version"
	Endpoint         = "https://dev.azure.com"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "azure <file|folder> <flags>",
	Short: "retrieves data from azure dev ops",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validate.Flags(map[string]interface{}{
			FlagOrganization: options.Current.Organization,
			FlagProject:      options.Current.Project,
			FlagRepo:         options.Current.Repo,
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.V(1).Println("retrieving files from azure dev ops")

		var gitApi api.GitApi = api.NewAzureGitApi(
			&api.BaseConfig{
				Url:       Endpoint,
				Auth:      globalOptions.Current.Api.Auth,
				UserAgent: globalOptions.Current.Api.UserAgent,
			},
			&api.AzureGitConfig{
				ApiVersion:   options.Current.ApiVersion,
				Organization: options.Current.Organization,
				Project:      options.Current.Project,
				Repo:         options.Current.Repo,
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
	rootCmd.Flags().StringVar(&options.Current.Organization, FlagOrganization, options.Current.Organization, "azure devops organization")
	rootCmd.Flags().StringVar(&options.Current.Project, FlagProject, options.Current.Project, "azure devops project")
	rootCmd.Flags().StringVar(&options.Current.Repo, FlagRepo, options.Current.Repo, "azure devops repo-(name/id)")
	rootCmd.Flags().StringVar(&options.Current.ApiVersion, FlagApiVersion, options.Current.ApiVersion, "azure devops api version")
}
