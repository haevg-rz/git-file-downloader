package logic

import (
	"encoding/base64"
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	globalOptions "github.com/haevg-rz/git-file-downloader/pkg/cli/options"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"os"
	"path"
	"regexp"
)

type IGitFileDownloader interface {
	HandleFile(string, string, string) (bool, error)
	HandleFolder(string, string, string, string, string) (bool, error)
	Handle(string, string, string, string) (bool, error)
}

type RegexRules struct {
	Include, Exclude string
}

type Context struct {
	OutPath    string
	RemotePath string
	Branch     string
	Patterns   *RegexRules
}

type GitFileDownloader struct {
	gitApi api.IGitApi
}

func NewRegexRules() *RegexRules {
	return &RegexRules{
		Include: "*",
		Exclude: "",
	}
}

func NewGitFileDownloader(gitApi api.IGitApi) *GitFileDownloader {
	return &GitFileDownloader{gitApi: gitApi}
}

// todo include exclude

func (g *GitFileDownloader) Handle(ctx *Context, modeArg string) error {
	if ctx.Patterns == nil {
		ctx.Patterns = NewRegexRules()
	}

	exists, err := api.ValidateBranch(g.gitApi, globalOptions.Current.Branch)
	if err != nil {
		exit.Code = exit.BranchOrRepoNotFound
		return err
	}
	if !exists {
		exit.Code = exit.BranchOrRepoNotFound
		return fmt.Errorf("branch %s does not exist", globalOptions.Current.Branch)
	}

	var modified bool
	switch modeArg {
	case "file":
		modified, err = g.HandleFile(
			ctx.OutPath,
			ctx.RemotePath,
			ctx.Branch)
	case "folder":
		modified, err = g.HandleFolder(
			ctx.OutPath,
			ctx.RemotePath,
			ctx.Branch,
			ctx.Patterns.Include,
			ctx.Patterns.Exclude)
	default:
		return fmt.Errorf("unsupported mode")
	}

	if !modified {
		log.V(1).Println("everything is up to date.")
		return nil
	}

	log.V(1).Println("synced file(s) successfully")
	return nil
}

func (g *GitFileDownloader) HandleFile(outFile, repoFilePath, branch string) (bool, error) {
	validPath, dir := GetDirFromFilepath(outFile)
	if !validPath {
		exit.Code = exit.InvalidOutPath
		return false, fmt.Errorf("GetDirFromFilepath: folder '%v' doesn't exist", dir)
	}

	gitFile, err := g.gitApi.GetRemoteFile(repoFilePath, branch)
	if err != nil {
		exit.Code = exit.FailedToRetrieveRemoteFile
		return false, fmt.Errorf("API Call exit: %v", err)
	}

	if !FileExists(outFile) {
		_, err = os.Create(outFile)
		if err != nil {
			exit.Code = exit.FailedToCreateFile
			return false, fmt.Errorf("CreateFile: '%v'", err)
		}
		log.V(2).Printf("Created File: '%s' because it didn't exist\n", outFile)
	} else {
		isEqual, err := IsHashEqual(outFile, gitFile.Sha, g.gitApi.GetHash())
		if err != nil {
			exit.Code = exit.FailedToOpenFile
			return false, fmt.Errorf("IsHashEqual: %v", err)
		}

		if isEqual {
			return false, nil
		}
	}

	fileData, err := base64.StdEncoding.DecodeString(gitFile.Content)
	if err != nil {
		exit.Code = exit.FailedToDecodeRemoteFileContent
		return false, fmt.Errorf("DecodeString: %v", err)
	}

	err = os.WriteFile(outFile, fileData, 0644)
	if err != nil {
		exit.Code = exit.FailedToWriteToFile
		return false, fmt.Errorf("WriteFile: %v", err)
	}
	return true, nil
}

func (g *GitFileDownloader) HandleFolder(outFolder, repoFolderPath, branch, include, exclude string) (bool, error) {
	updated := false

	if !IsValidPath(outFolder) {
		err := os.Mkdir(outFolder, 0755)
		if err != nil {
			exit.Code = exit.FailedToCreateFolder
			return false, err
		}
	}

	// rename files to nodes
	files, err := g.gitApi.GetFilesFromFolder(repoFolderPath, branch)
	if err != nil {
		exit.Code = exit.FailedToGetFilesFromRemoteFolder
		return false, err
	}

	log.V(2).Println("Sync", len(files), "files, from remote folder", repoFolderPath)

	for _, file := range files {
		if file.Type == "tree" || file.Type == "dir" {
			if updated, err = g.HandleFolder(path.Join(outFolder, file.Name), file.Path, branch, include, exclude); err != nil {
				return updated, err
			}
			continue
		}

		if include != "" {
			matched, err := regexp.MatchString(include, file.Name)
			if err == nil {
				if !matched {
					log.V(2).Printf("Skip: '%s' because of include rule: '%s'\n", file.Name, include)
					continue
				}
			}
		}

		if exclude != "" {
			matched, err := regexp.MatchString(exclude, file.Name)
			if err == nil {
				if matched {
					log.V(2).Printf("Skip: '%s' because of exclude rule: '%s'\n", file.Name, exclude)
					continue
				}
			}
		}

		modifiedOrCreated, err := g.HandleFile(path.Join(outFolder, path.Base(file.Path)), file.Path, branch)
		if err != nil {
			return updated, err
		}

		if !modifiedOrCreated {
			log.V(2).Printf("Skip: %s because content is equal\n", file.Path)
			continue
		}

		updated = true
		log.V(2).Printf("Wrote file: %s because is new or updated\n", file.Path)
	}

	return updated, nil
}
