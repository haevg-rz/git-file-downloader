package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/haevg-rz/git-file-downloader/internal"
)

var (
	HttpGetFunc func(apiUrl string, settings internal.Settings) ([]byte, error) = httpGetInternal
)

func httpGetInternal(apiUrl string, settings internal.Settings) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Private-Token", settings.PrivateToken)
	req.Header.Add("User-Agent", settings.UserAgent)

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP GET failed with status code %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetBranches(settings internal.Settings) ([]GitLabBranch, error) {
	apiUrl := fmt.Sprintf("%vprojects/%v/repository/branches", settings.ApiUrl, settings.ProjectNumber)
	body, err := HttpGetFunc(apiUrl, settings)
	if err != nil {
		return nil, err
	}

	var responseStruct []GitLabBranch
	err = json.Unmarshal(body, &responseStruct)

	return responseStruct, err
}

type GitLabBranch struct {
	Name string `json:"name"`
}

func GetFilesFromFolder(settings internal.Settings) ([]GitLabRepoFile, error) {
	path := url.QueryEscape(settings.RepoFolderPath)
	branch := url.QueryEscape(settings.Branch)
	apiUrl := fmt.Sprintf("%vprojects/%v/repository/tree/?ref=%v&path=%v", settings.ApiUrl, settings.ProjectNumber, branch, path)

	body, err := HttpGetFunc(apiUrl, settings)
	if err != nil {
		return nil, err
	}

	var responseStruct []GitLabRepoFile
	err = json.Unmarshal(body, &responseStruct)

	return responseStruct, err
}

type GitLabRepoFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
	Mode string `json:"mode"`
}

func GetFile(settings internal.Settings) (GitLapFile, error) {
	path := url.QueryEscape(settings.RepoFilePath)
	branch := url.QueryEscape(settings.Branch)
	apiUrl := fmt.Sprintf("%vprojects/%v/repository/files/%v?ref=%v", settings.ApiUrl, settings.ProjectNumber, path, branch)

	body, err := HttpGetFunc(apiUrl, settings)
	if err != nil {
		return GitLapFile{}, err
	}

	file, err := createGitLapFile(body)
	if err != nil {
		return GitLapFile{}, err
	}
	return file, nil
}

func createGitLapFile(data []byte) (GitLapFile, error) {
	var gitLapFile GitLapFile
	err := json.Unmarshal(data, &gitLapFile)
	return gitLapFile, err
}

type GitLapFile struct {
	FileName      string `json:"file_name"`
	ContentSha256 string `json:"content_sha256"`
	Content       string
}
