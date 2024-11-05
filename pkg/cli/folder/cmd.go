package folder

import (
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/folder/options"
	globalOptions "github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/logic"
	"github.com/spf13/cobra"
	"log"
)

const (
	FlagOutFolder      = "out-folder"
	FlagRepoFolderPath = "repo-folder-path"
	FlagRecursive      = "recursive"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "folder",
	Short: "runs in folder mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("running folder mode.")

		fileDownloader := logic.NewGitFileDownloader(api.NewGitLabApi(
			globalOptions.Current.Api.UserAgent,
			globalOptions.Current.Api.ApiBaseUrl,
			globalOptions.Current.Api.PrivateToken,
			globalOptions.Current.Api.ProjectNumber))

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
			log.Println("folder is up to date.")
			return nil
		}

		log.Println("synced folder successfully")

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
