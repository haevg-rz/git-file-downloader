package logic

import (
	"encoding/base64"
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"os"
	"path"
	"regexp"
)

type IGitFileDownloader interface {
	HandleFile(string, string, string) (bool, error)
	HandleFolder(string, string, string, string, string) (bool, error)
}

type GitFileDownloader struct {
	gitApi api.IGitApi
}

func NewGitFileDownloader(gitApi api.IGitApi) *GitFileDownloader {
	return &GitFileDownloader{gitApi: gitApi}
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

	fileExists := FileExists(outFile)
	if !fileExists {
		_, err = os.Create(outFile)
		if err != nil {
			exit.Code = exit.FailedToCreateFile
			return false, fmt.Errorf("CreateFile: '%v'", err)
		}
		log.V(2).Printf("Created File: '%s' because it didn't exist\n", outFile)
	}

	fileData, err := base64.StdEncoding.DecodeString(gitFile.Content)
	if err != nil {
		exit.Code = exit.FailedToDecodeRemoteFileContent
		return false, fmt.Errorf("DecodeString: %v", err)
	}

	isEqual, err := IsHashEqual(outFile, gitFile.Sha, g.gitApi.GetHash())
	if err != nil {
		exit.Code = exit.FailedToOpenFile
		return false, fmt.Errorf("IsHashEqual: %v", err)
	}

	if isEqual {
		return false, nil
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

	files, err := g.gitApi.GetFilesFromFolder(repoFolderPath, branch)
	if err != nil {
		exit.Code = exit.FailedToGetFilesFromRemoteFolder
		return false, err
	}

	log.V(2).Println("Sync", len(files), "files, from remote folder", repoFolderPath)

	for _, file := range files {
		if include != "" {
			matched, err := regexp.MatchString(include, file.Name)
			if err == nil {
				if !matched {
					log.V(3).Printf("Skip: '%s' because of include rule: '%s'\n", file.Name, include)
					continue
				}
			}
		}
		if exclude != "" {
			matched, err := regexp.MatchString(exclude, file.Name)
			if err == nil {
				if matched {
					log.V(3).Printf("Skip: '%s' because of exclude rule: '%s'\n", file.Name, exclude)
					continue
				}
			}
		}

		if file.Type == "tree" || file.Type == "dir" {
			updated, err = g.HandleFolder(path.Join(outFolder, file.Name), file.Path, branch, include, exclude)
			if err != nil {
				return updated, err
			}
			continue
		}

		modifiedOrCreated, err := g.HandleFile(path.Join(outFolder, path.Base(file.Path)), file.Path, branch)
		if err != nil {
			return updated, err
		}

		if !modifiedOrCreated {
			log.V(3).Printf("Skip: %s because content is equal\n", file.Path)
			continue
		}

		updated = true
		log.V(3).Printf("Wrote file: %s because is new or updated\n", file.Path)
	}

	return updated, nil
}
