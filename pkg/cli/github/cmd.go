package github

import (
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/github/options"
	globalOptions "github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"github.com/haevg-rz/git-file-downloader/pkg/logic"
	"github.com/spf13/cobra"
	"strings"
)

const (
	FlagOwner = "owner"
	FlagRepo  = "repo"

	Endpoint = "https://api.github.com"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "github",
	Short: "github",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := validate.Flags(map[string]interface{}{
			FlagOwner: options.Current.Owner,
			FlagRepo:  options.Current.Repo,
		})
		if err != nil {
			exit.Code = exit.MissingFlags
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.V(1).Println("retrieving files from github")

		var gitApi api.IGitApi = api.NewGitHubApi(
			globalOptions.Current.Api.Auth,
			globalOptions.Current.Api.UserAgent,
			Endpoint,
			options.Current.Owner,
			options.Current.Repo)

		exists, err := api.ValidateBranch(gitApi, globalOptions.Current.Branch)
		if err != nil {
			exit.Code = exit.BranchNotFound
			return err
		}
		if !exists {
			exit.Code = exit.BranchNotFound
			return fmt.Errorf("branch %s does not exist", globalOptions.Current.Branch)
		}

		fileDownloader := logic.NewGitFileDownloader(gitApi)

		created, err := fileDownloader.Handle(
			globalOptions.Current.OutPath,
			globalOptions.Current.RemotePath,
			globalOptions.Current.Branch,
			strings.ToLower(args[0]))

		if err != nil {
			return err
		}

		if !created {
			log.V(1).Println("everything is up to date.")
			return nil
		}

		log.V(1).Println("synced file(s) successfully")
		return nil
	},
}

func Command() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.Flags().StringVar(&options.Current.Owner, FlagOwner, options.Current.Owner, "repo owner")
	rootCmd.Flags().StringVar(&options.Current.Repo, FlagRepo, options.Current.Repo, "repo name")
}
