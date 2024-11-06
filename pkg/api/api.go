package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	HttpGetFunc = httpGetInternal
)

type IGitLabApi interface {
	GetAvailableBranches() ([]GitLabBranch, error)
	BranchExists(string) (bool, error)
	GetFile(string, string) (*GitLabFile, error)
	GetFilesFromFolder(string, string) ([]GitLabRepoFile, error)
}

type GitLabFile struct {
	FileName      string `json:"file_name"`
	ContentSha256 string `json:"content_sha256"`
	Content       string
}

type GitLabRepoFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
	Mode string `json:"mode"`
}

type GitLabBranch struct {
	Name string `json:"name"`
}

type GitLabApi struct {
	UserAgent     string
	ApiBaseUrl    string
	PrivateToken  string
	ProjectNumber int
}

func NewGitLabApi(userAgent, apiBaseUrl, privateToken string, projectNumber int) *GitLabApi {
	return &GitLabApi{
		UserAgent:     userAgent,
		ApiBaseUrl:    apiBaseUrl,
		PrivateToken:  privateToken,
		ProjectNumber: projectNumber,
	}
}

func (g *GitLabApi) GetFile(repoFilePath, branchName string) (*GitLabFile, error) {
	var gitFile *GitLabFile

	path := url.QueryEscape(repoFilePath)
	branch := url.QueryEscape(branchName)

	fullUrl := fmt.Sprintf("%vprojects/%v/repository/files/%v?ref=%v", g.ApiBaseUrl, g.ProjectNumber, path, branch)

	body, err := HttpGetFunc(fullUrl, g.PrivateToken, g.UserAgent)
	if err != nil {
		return &GitLabFile{}, err
	}

	err = json.Unmarshal(body, &gitFile)
	if err != nil {
		return &GitLabFile{}, err
	}

	return gitFile, nil
}

func (g *GitLabApi) GetFilesFromFolder(repoFolderPath, branch string) ([]GitLabRepoFile, error) {
	path := url.QueryEscape(repoFolderPath)
	branchEsc := url.QueryEscape(branch)
	fullUrl := fmt.Sprintf("%vprojects/%v/repository/tree/?ref=%v&path=%v", g.ApiBaseUrl, g.ProjectNumber, branchEsc, path)

	body, err := HttpGetFunc(fullUrl, g.PrivateToken, g.UserAgent)
	if err != nil {
		return nil, err
	}

	var result []GitLabRepoFile
	err = json.Unmarshal(body, &result)

	return result, err
}

func (g *GitLabApi) GetAvailableBranches() ([]GitLabBranch, error) {
	body, err := HttpGetFunc(fmt.Sprintf("%vprojects/%v/repository/branches", g.ApiBaseUrl, g.ProjectNumber), g.PrivateToken, g.UserAgent)
	if err != nil {
		return nil, err
	}

	var availableBranches []GitLabBranch
	err = json.Unmarshal(body, &availableBranches)

	return availableBranches, err
}

func (g *GitLabApi) BranchExists(branch string) (bool, error) {
	branches, err := g.GetAvailableBranches()
	if err != nil {
		return false, err
	}

	for _, val := range branches {
		if val.Name == branch {
			return true, nil
		}
	}
	return false, nil
}

func httpGetInternal(fullUrl, privateToken, userAgent string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Private-Token", privateToken)
	req.Header.Add("User-Agent", userAgent)

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

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
