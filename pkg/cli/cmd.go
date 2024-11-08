package cli

import (
	"errors"
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/version"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"github.com/haevg-rz/git-file-downloader/pkg/logic"
	"github.com/haevg-rz/git-file-downloader/pkg/provider"
	"github.com/spf13/cobra"
	"slices"
	"strings"
)

const (
	// SOURCE, DEST

	FlagOutPath    = "out-path"
	FlagRemotePath = "remote-path"

	FlagBranch         = "branch"
	FlagIncludePattern = "include"
	FlagExcludePattern = "exclude"
	FlagLogLevel       = "log-level"

	// API

	FlagAuthToken = "token"
	FlagUrl       = "url"
	FlagUserAgent = "user-agent"

	// GITLAB

	FlagProjectNumber = "project-number"

	// GITHUB

	FlagOwner = "owner"
	FlagRepo  = "repo"
)

var (
	ProviderToUrl = map[string]string{
		provider.Github: "https://api.github.com",
		provider.Gitlab: "https://gitlab.com/api/v4",
	}
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:       "gdown",
	Short:     "gdown",
	Long:      "git-file-downloader",
	Args:      cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	ValidArgs: []string{"file", "folder", "github", "gitlab"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// set global log level
		log.Level = options.Current.LogLevel

		// no flags needed for version subcommand
		if cmd.CalledAs() == "version" || cmd.CalledAs() == "help" {
			return nil
		}

		return validate.Flags(map[string]interface{}{
			FlagAuthToken:  options.Current.Api.Auth,
			FlagOutPath:    options.Current.OutPath,
			FlagRemotePath: options.Current.RemotePath,
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			exit.Code = exit.MissingFlags
			return fmt.Errorf("please specify mode[file|folder] and git-provider[github|gitlab|azure]")
		}

		mode := strings.ToLower(args[0])
		gitProvider := strings.ToLower(args[1])

		var gitApi api.IGitApi

		switch gitProvider {
		case provider.Github:
			gitApi = api.NewGitHubApi(
				options.Current.Api.Auth,
				options.Current.Api.UserAgent,
				ProviderToUrl[gitProvider],
				options.Current.Owner,
				options.Current.Repo)
		case provider.Gitlab:
			gitApi = api.NewGitLabApi(
				options.Current.Api.Auth,
				options.Current.Api.UserAgent,
				ProviderToUrl[gitProvider],
				options.Current.ProjectNumber)
		default:
			exit.Code = exit.UnknownGitProvider
			return fmt.Errorf("unsupported git provider: %s", options.Current.GitProvider)
		}

		// check branch

		branches, err := gitApi.GetAvailableBranches()
		if err != nil {
			exit.Code = exit.BranchNotFound
			return err
		}

		if (branches != nil) && slices.Contains(branches, options.Current.Branch) == false {
			exit.Code = exit.BranchNotFound
			return errors.New(fmt.Sprintf("branch '%s' does not exist\n", options.Current.Branch))
		}

		// begin

		fileDownloader := logic.NewGitFileDownloader(gitApi)

		var created bool

		switch mode {
		case "file":
			created, err = fileDownloader.HandleFile(
				options.Current.OutPath,
				options.Current.RemotePath,
				options.Current.Branch)
		case "folder":
			created, err = fileDownloader.HandleFolder(
				options.Current.OutPath,
				options.Current.RemotePath,
				options.Current.Branch,
				options.Current.IncludePattern,
				options.Current.ExcludePattern)
		}

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
	rootCmd.AddCommand(version.Command())
	return rootCmd
}

func init() {
	// source, dest
	rootCmd.PersistentFlags().StringVar(&options.Current.OutPath, FlagOutPath, options.Current.OutPath, "Path to write file to disk")
	rootCmd.PersistentFlags().StringVar(&options.Current.RemotePath, FlagRemotePath, options.Current.RemotePath, "Path to file/folder from remote source")

	rootCmd.PersistentFlags().StringVar(&options.Current.Branch, FlagBranch, options.Current.Branch, "Branch name")
	rootCmd.PersistentFlags().StringVar(&options.Current.IncludePattern, FlagIncludePattern, options.Current.IncludePattern, "Include this regex pattern")
	rootCmd.PersistentFlags().StringVar(&options.Current.ExcludePattern, FlagExcludePattern, options.Current.ExcludePattern, "Exclude this regex pattern")
	rootCmd.PersistentFlags().IntVar(&options.Current.LogLevel, FlagLogLevel, options.Current.LogLevel, "Higher loglevel leads to more verbose logging. Set log level to 0 if you dont want any logging.")

	// gitlab def
	rootCmd.PersistentFlags().IntVar(&options.Current.ProjectNumber, FlagProjectNumber, options.Current.ProjectNumber, "The Project ID from your project")

	// github def
	rootCmd.PersistentFlags().StringVar(&options.Current.Owner, FlagOwner, options.Current.Owner, "github repo owner")
	rootCmd.PersistentFlags().StringVar(&options.Current.Repo, FlagRepo, options.Current.Repo, "repo name")

	// api flag def
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.Auth, FlagAuthToken, options.Current.Api.Auth, "Private-Token with access right for \"api\" and \"read_repository\", role must be minimum \"Reporter\"")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.UserAgent, FlagUserAgent, options.Current.Api.UserAgent, "User agent")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.BaseUrl, FlagUrl, options.Current.Api.BaseUrl, "url to Api v4, like https://my-git-lab-server.local/api/v4/")
}
