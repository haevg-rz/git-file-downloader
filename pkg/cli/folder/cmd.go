package folder

import (
	"errors"
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/folder/options"
	globalOptions "github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"github.com/haevg-rz/git-file-downloader/pkg/logic"
	"github.com/spf13/cobra"
)

const (
	FlagOutFolder      = "out-folder"
	FlagRepoFolderPath = "repo-folder-path"
	FlagRecursive      = "recursive"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "folder",
	Short: "runs in folder mode",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validate.Flags(map[string]interface{}{
			FlagOutFolder:      options.Current.OutFolder,
			FlagRepoFolderPath: options.Current.RepoFolderPath,
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.V(1).Println("FOLDER MODE")

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

		created, err := fileDownloader.HandleFolder(
			options.Current.OutFolder,
			options.Current.RepoFolderPath,
			globalOptions.Current.Branch,
			globalOptions.Current.IncludePattern,
			globalOptions.Current.ExcludePattern)

		if err != nil {
			return err
		}

		if !created {
			log.V(1).Println("folder is up to date.")
			return nil
		}

		log.V(1).Println("synced folder successfully")

		return nil
	},
}

func Command() *cobra.Command {
	return rootCmd
}

func init() {
	// flag definitions
	rootCmd.Flags().StringVar(&options.Current.OutFolder, FlagOutFolder, options.Current.OutFolder, "Folder to write file to disk")
	rootCmd.Flags().StringVar(&options.Current.RepoFolderPath, FlagRepoFolderPath, options.Current.RepoFolderPath, "Folder to write file to disk")
	rootCmd.Flags().BoolVar(&options.Current.Recursive, FlagRecursive, options.Current.Recursive, "recursive mode")
}
