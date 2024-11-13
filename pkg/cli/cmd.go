package cli

import (
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/azure"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/github"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/gitlab"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/validate"
	"github.com/haevg-rz/git-file-downloader/pkg/cli/version"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"github.com/spf13/cobra"
	goLog "log"
	"os"
	"time"
)

const (
	LogOutputPath = "./logs"
	LogFileFormat = "YYYY-MM-DD-HH-MM-SS"
	// SOURCE, DEST

	FlagOutPath    = "out"
	FlagRemotePath = "remote-path"

	FlagBranch         = "branch"
	FlagIncludePattern = "include"
	FlagExcludePattern = "exclude"
	FlagLogLevel       = "log-level"

	// API

	FlagAuthToken = "token"
	FlagUrl       = "url"
	FlagUserAgent = "user-agent"
	FlagLogToFile = "logfile"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "gdown <github|gitlab|azure|help|version>",
	Short: "gdown",
	Long:  "git-file-downloader",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.Level = options.Current.LogLevel

		if options.Current.LogToFile {
			if err := initLog(); err != nil {
				return err
			}
		}

		if cmd.CalledAs() == "version" || cmd.CalledAs() == "help" || cmd.CalledAs() == "gdown" {
			return nil
		}

		return validate.Flags(map[string]interface{}{
			FlagAuthToken:  options.Current.Api.Auth,
			FlagOutPath:    options.Current.OutPath,
			FlagRemotePath: options.Current.RemotePath,
			FlagLogToFile:  options.Current.LogToFile,
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand..")
		return nil
	},
}

func Command() *cobra.Command {
	rootCmd.AddCommand(version.Command())
	rootCmd.AddCommand(github.Command())
	rootCmd.AddCommand(gitlab.Command())
	rootCmd.AddCommand(azure.Command())
	return rootCmd
}

func initLog() error {
	var err error

	defer func() {
		if err != nil {
			exit.Code = exit.InternalError
		}
	}()

	fileInfo, err := os.Stat(LogOutputPath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() == false {
		exit.Code = exit.InternalError
		err = fmt.Errorf("%s is not a directory", LogOutputPath)
		return err
	}

	logFile, err := os.Open(fmt.Sprintf("%s/%s-log.txt", LogOutputPath, time.Now().Format(LogFileFormat)))
	if err != nil {
		exit.Code = exit.InternalError
		return err
	}

	goLog.SetOutput(logFile)
	return nil
}

func init() {
	// source, dest
	rootCmd.PersistentFlags().StringVar(&options.Current.OutPath, FlagOutPath, options.Current.OutPath, "Path to write file to disk")
	rootCmd.PersistentFlags().StringVar(&options.Current.RemotePath, FlagRemotePath, options.Current.RemotePath, "Path to file/folder from remote source")

	rootCmd.PersistentFlags().StringVar(&options.Current.Branch, FlagBranch, options.Current.Branch, "Branch name")
	rootCmd.PersistentFlags().StringVar(&options.Current.IncludePattern, FlagIncludePattern, options.Current.IncludePattern, "Include this regex pattern")
	rootCmd.PersistentFlags().StringVar(&options.Current.ExcludePattern, FlagExcludePattern, options.Current.ExcludePattern, "Exclude this regex pattern")
	rootCmd.PersistentFlags().IntVar(&options.Current.LogLevel, FlagLogLevel, options.Current.LogLevel, "Higher loglevel leads to more verbose logging. Set log level to 0 if you dont want any logging.")
	rootCmd.PersistentFlags().BoolVar(&options.Current.LogToFile, FlagLogToFile, options.Current.LogToFile, "Write to file instead of stdout")

	// api flag def
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.Auth, FlagAuthToken, options.Current.Api.Auth, "Private-Token with access right for \"api\" and \"read_repository\", role must be minimum \"Reporter\"")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.UserAgent, FlagUserAgent, options.Current.Api.UserAgent, "User agent")
	rootCmd.PersistentFlags().StringVar(&options.Current.Api.BaseUrl, FlagUrl, options.Current.Api.BaseUrl, "url to Api v4, like https://my-git-lab-server.local/api/v4/")
}
