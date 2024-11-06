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
	"github.com/spf13/cobra"
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

		var gitApi api.IGitLabApi = api.NewGitLabApi(
			globalOptions.Current.Api.UserAgent,
			globalOptions.Current.Api.ApiBaseUrl,
			globalOptions.Current.Api.PrivateToken,
			globalOptions.Current.Api.ProjectNumber)

		exists, err := gitApi.BranchExists(globalOptions.Current.Branch)
		if !exists {
			exit.Code = exit.BranchNotFound

			if err != nil {
				return err
			}
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
			log.V(1).Println("Skip:", options.Current.RepoFilePath, ", because content is equal")
			return nil
		}

		log.V(1).Println("Wrote file:", options.Current.RepoFilePath, ", because is new or changed")

		return nil
	},
}

func Command() *cobra.Command {
	return rootCmd
}

func init() {
	// flag definitions
	rootCmd.Flags().StringVar(&options.Current.OutFile, FlagOutFile, options.Current.OutFile, "output file")
	rootCmd.Flags().StringVar(&options.Current.RepoFilePath, FlagRepoFilePath, options.Current.RepoFilePath, "File path in repo, like src/main.go")
}
