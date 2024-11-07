package file

import (
	"errors"
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/file/options"
	globalOptions "github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"github.com/haevg-rz/git-file-downloader/pkg/logic"
	"github.com/haevg-rz/git-file-downloader/pkg/provider"
	"github.com/spf13/cobra"
	"slices"
)

const (
	FlagOutFile      = "out-file"
	FlagRepoFilePath = "repo-file-path"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "file",
	Short: "runs in file mode",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validate.Flags(map[string]interface{}{
			FlagOutFile:      options.Current.OutFile,
			FlagRepoFilePath: options.Current.RepoFilePath,
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.V(1).Println("FILE MODE")

		var gitApi api.IGitApi

		// TODO ADD AZURE DEV OPS
		// TODO CREATE SEPARATE COMMANDS
		switch globalOptions.Current.GitProvider {
		case provider.Github:
			gitApi = api.NewGitHubApi(
				globalOptions.Current.Api.PrivateToken,
				globalOptions.Current.Api.UserAgent,
				globalOptions.Current.Api.BaseUrl,
				"driema",
				"kattis")
		case provider.Gitlab:
			gitApi = api.NewGitLabApi(
				globalOptions.Current.Api.PrivateToken,
				globalOptions.Current.Api.UserAgent,
				globalOptions.Current.Api.BaseUrl,
				globalOptions.Current.Api.ProjectNumber)
		default:
			exit.Code = exit.UnknownGitProvider
			return fmt.Errorf("unsupported git provider: %s", globalOptions.Current.GitProvider)
		}

		branches, err := gitApi.GetAvailableBranches()
		if err != nil {
			exit.Code = exit.BranchNotFound
			return err
		}

		if (branches != nil) && slices.Contains(branches, globalOptions.Current.Branch) == false {
			exit.Code = exit.BranchNotFound
			return errors.New(fmt.Sprintf("branch '%s' does not exist\n", globalOptions.Current.Branch))
		}

		fileDownloader := logic.NewGitFileDownloader(gitApi)

		created, err := fileDownloader.HandleFile(
			options.Current.OutFile,
			options.Current.RepoFilePath,
			globalOptions.Current.Branch)
		if err != nil {
			return err
		}

		if !created {
			log.V(1).Printf("Skip: %s, because content is equal", options.Current.RepoFilePath)
			return nil
		}

		log.V(1).Printf("Wrote file: %s, because is new or changed", options.Current.RepoFilePath)

		return nil
	},
}

func Command() *cobra.Command {
	return rootCmd
}

func init() {
	// flag provider
	rootCmd.Flags().StringVar(&options.Current.OutFile, FlagOutFile, options.Current.OutFile, "output file")
	rootCmd.Flags().StringVar(&options.Current.RepoFilePath, FlagRepoFilePath, options.Current.RepoFilePath, "File path in repo, like src/main.go")
}
