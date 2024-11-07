package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

var _ IGitApi = &GitLabApi{}

// GitLabFile describes a file returned from the gitLabApi
type GitLabFile struct {
	FileName      string `json:"file_name"`
	ContentSha256 string `json:"content_sha256"`
	Content       string
}

// GitLabRepoNode describes either a file or directory ("tree") returned from the gitLabApi.
// Contains metadata about path, type, id, etc.
type GitLabRepoNode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
	Mode string `json:"mode"`
}

// GitLabBranch contains the name of a branch returned from the gitLabApi.
type GitLabBranch struct {
	Name string `json:"name"`
}

// GitLabApi is used for communication to the gitLabApi. Instance fields are used as base-configuration for every request.
// Implements IGitApi.
type GitLabApi struct {
	Base          *GitApi
	ProjectNumber int
}

// NewGitLabApi creates a new instance of the git lab api
func NewGitLabApi(userAgent, apiBaseUrl, privateToken string, projectNumber int) *GitLabApi {
	return &GitLabApi{
		Base: &GitApi{
			AuthToken: privateToken,
			UserAgent: userAgent,
			Url:       apiBaseUrl,
		},
		ProjectNumber: projectNumber,
	}
}

// GetFile retrieves remote file from a given path and branch
func (g *GitLabApi) GetFile(repoFilePath, branchName string) (*GitLabFile, error) {
	var gitFile *GitLabFile

	path := url.QueryEscape(repoFilePath)
	branch := url.QueryEscape(branchName)

	fullUrl := fmt.Sprintf("%vprojects/%v/repository/files/%v?ref=%v", g.Base.Url, g.ProjectNumber, path, branch)

	body, err := HttpGetFunc(fullUrl, g.Base.AuthToken, g.Base.UserAgent)
	if err != nil {
		return &GitLabFile{}, err
	}

	err = json.Unmarshal(body, &gitFile)
	if err != nil {
		return &GitLabFile{}, err
	}

	return gitFile, nil
}

// GetFilesFromFolder retrieves All remote files/folders from a given path and branch
func (g *GitLabApi) GetFilesFromFolder(repoFolderPath, branch string) ([]GitLabRepoNode, error) {
	path := url.QueryEscape(repoFolderPath)
	branchEsc := url.QueryEscape(branch)
	fullUrl := fmt.Sprintf("%vprojects/%v/repository/tree/?ref=%v&path=%v", g.Base.Url, g.ProjectNumber, branchEsc, path)

	body, err := HttpGetFunc(fullUrl, g.Base.AuthToken, g.Base.UserAgent)
	if err != nil {
		return nil, err
	}

	var result []GitLabRepoNode
	err = json.Unmarshal(body, &result)

	return result, err
}

// GetAvailableBranches retrieves all available branches.
func (g *GitLabApi) GetAvailableBranches() ([]GitLabBranch, error) {
	body, err := HttpGetFunc(fmt.Sprintf("%vprojects/%v/repository/branches", g.Base.Url, g.ProjectNumber), g.Base.AuthToken, g.Base.UserAgent)
	if err != nil {
		return nil, err
	}

	var availableBranches []GitLabBranch
	err = json.Unmarshal(body, &availableBranches)

	return availableBranches, err
}

// BranchExists checks whether a certain branch exists. Calls GetAvailableBranches internally.
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
