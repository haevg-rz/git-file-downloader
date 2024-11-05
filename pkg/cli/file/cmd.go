package file

import (
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/file/options"
	globalOptions "github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/logic"
	"github.com/spf13/cobra"
	"log"
)

const (
	FlagOutFile      = "out-file"
	FlagRepoFilePath = "repo-file-path"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "file",
	Short: "runs in file mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("running file mode.")

		fileDownloader := logic.NewGitFileDownloader(api.NewGitLabApi(
			globalOptions.Current.Api.UserAgent,
			globalOptions.Current.Api.ApiBaseUrl,
			globalOptions.Current.Api.PrivateToken,
			globalOptions.Current.Api.ProjectNumber))

		created, err := fileDownloader.HandleFile(
			options.Current.OutFile,
			options.Current.RepoFilePath,
			globalOptions.Current.Branch)

		if err != nil {
			return err
		}

		if !created {
			log.Println("Skip:", options.Current.RepoFilePath, ", because content is equal")
			return nil
		}

		log.Println("Wrote file:", options.Current.RepoFilePath, ", because is new or changed")

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
