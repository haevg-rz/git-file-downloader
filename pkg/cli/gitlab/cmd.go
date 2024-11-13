package gitlab

import (
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/gitlab/options"
	globalOptions "github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/logic"
	"github.com/spf13/cobra"
	"strings"
)

const (
	FlagProjectId = "project"
	Endpoint      = "https://gitlab.com/api/v4"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "gitlab <file|folder> <flags>",
	Short: "retrieves data from gitlab",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validate.Flags(map[string]interface{}{
			FlagProjectId: options.Current.ProjectId,
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var gitApi api.IGitApi = api.NewGitLabApi(
			globalOptions.Current.Api.Auth,
			globalOptions.Current.Api.UserAgent,
			Endpoint,
			options.Current.ProjectId)

		return logic.NewGitFileDownloader(gitApi).Handle(
			globalOptions.Current.OutPath,
			globalOptions.Current.RemotePath,
			globalOptions.Current.Branch,
			strings.ToLower(args[0]))
	},
}

func Command() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.Flags().IntVar(&options.Current.ProjectId, FlagProjectId, options.Current.ProjectId, "project id")
}
