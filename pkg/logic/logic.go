package logic

import (
	"encoding/base64"
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"log"
	"os"
	"path"
	"regexp"
)

type IGitFileDownloader interface {
	HandleFile(string, string, string) (bool, error)
}

type GitFileDownloader struct {
	gitLabApi api.IGitLabApi
}

func NewGitFileDownloader(gitLabApi api.IGitLabApi) *GitFileDownloader {
	return &GitFileDownloader{gitLabApi: gitLabApi}
}

func (g *GitFileDownloader) HandleFile(outFile, repoFilePath, branch string) (bool, error) {
	validPath, dir := GetDirFromFilepath(outFile)
	if !validPath {
		return false, fmt.Errorf("GetDirFromFilepath: folder %v doesn't exist", dir)
	}

	fileExists := FileExists(outFile)
	if !fileExists {
		_, err := os.Create(outFile)
		if err != nil {
			return false, fmt.Errorf("CreateFile: %v", err)
		}
		log.Printf("Created File: %s because it didn't exist\n", outFile)
	}

	gitFile, err := g.gitLabApi.GetFile(repoFilePath, branch)
	if err != nil {
		return false, fmt.Errorf("API Call error: %v", err)
	}

	fileData, err := base64.StdEncoding.DecodeString(gitFile.Content)
	if err != nil {
		return false, fmt.Errorf("DecodeString: %v", err)
	}

	isEqual, err := IsHashEqual(outFile, gitFile.ContentSha256)
	if err != nil {
		return false, fmt.Errorf("IsHashEqual: %v", err)
	}

	if isEqual {
		return false, nil
	}

	err = os.WriteFile(outFile, fileData, 0644)
	if err != nil {
		return false, fmt.Errorf("WriteFile: %v", err)
	}
	return true, nil
}

func (g *GitFileDownloader) HandleFolder(outFolder, repoFolderPath, branch, include, exclude string) (bool, error) {
	updated := false

	if !IsValidPath(outFolder) {
		err := os.Mkdir(outFolder, 0755)
		if err != nil {
			return false, err
		}
	}

	files, err := g.gitLabApi.GetFilesFromFolder(repoFolderPath, branch)
	if err != nil {
		return false, err
	}

	log.Println("Sync", len(files), "files, from remote folder", repoFolderPath)

	for _, file := range files {
		if include != "" {
			matched, err := regexp.MatchString(include, file.Name)
			if err == nil {
				if !matched {
					log.Println("Skip:", file.Name, "because of include rule:", include)
					continue
				}
			}
		}
		if exclude != "" {
			matched, err := regexp.MatchString(exclude, file.Name)
			if err == nil {
				if matched {
					log.Println("Skip:", file.Name, "because of exclude rule:", exclude)
					continue
				}
			}
		}

		if file.Type == "tree" {
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
			log.Printf("Skip: %s because content is equal\n", file.Path)
			continue
		}

		updated = true
		log.Printf("Wrote file: %s because is new or updated\n", file.Path)
	}

	return updated, err
}
