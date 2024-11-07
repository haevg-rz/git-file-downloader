package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var _ IGitApi = &GitLabApi{}

// GitLabRepoFile describes a file returned from the gitLabApi
type GitLabRepoFile struct {
	Name    string `json:"file_name"`
	Sha256  string `json:"content_sha256"`
	Content string `json:"content"`
}

// GitLabRepoNode describes either a file or directory ("tree") returned from the gitLabApi.
// Contains metadata about path, type, id, etc.
type GitLabRepoNode struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
}

// GitLabApi is used for communication to the gitLabApi. Instance fields are used as base-configuration for every request.
// Implements IGitApi.
type GitLabApi struct {
	base          *Config
	projectNumber int
}

const (
	gitlabNodeTemplate   = "%s/projects/%s/repository/tree/?ref=%s&path=%s"
	filePath             = "%s/projects/%s/repository/files/%s?ref=%s"
	gitlabBranchTemplate = "%s/projects/%s/repository/branches"
)

// NewGitLabApi creates a new instance of the git lab api
func NewGitLabApi(privateToken, userAgent, apiBaseUrl string, projectNumber int) *GitLabApi {
	return &GitLabApi{
		base: &Config{
			url: apiBaseUrl,
			defaultHeader: map[string]string{
				"Private-Token": privateToken,
				"User-Agent":    userAgent,
			},
		},
		projectNumber: projectNumber,
	}
}

// GetRemoteFile retrieves remote file from a given path and branch
func (g *GitLabApi) GetRemoteFile(path, branch string) (*GitRepoFile, error) {
	var gitLabFile *GitRepoFile
	fullUrl := CreateUrl(filePath, g.base.url, strconv.Itoa(g.projectNumber), path, branch)
	fmt.Println(fullUrl)

	body, err := HttpGetFunc(fullUrl, g.base.defaultHeader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &gitLabFile)
	if err != nil {
		return nil, err
	}

	return &GitRepoFile{
		Name:    gitLabFile.Name,
		Sha256:  gitLabFile.Sha256,
		Content: gitLabFile.Content,
	}, nil
}

// GetFilesFromFolder retrieves All remote files/folders from a given path and branch
func (g *GitLabApi) GetFilesFromFolder(path, branch string) ([]GitRepoNode, error) {
	gitLabNodes := make([]GitLabRepoNode, 0)
	gitNodes := make([]GitRepoNode, 0)
	fullUrl := CreateUrl(gitlabNodeTemplate, g.base.url, strconv.Itoa(g.projectNumber), branch, path)

	body, err := HttpGetFunc(fullUrl, g.base.defaultHeader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &gitLabNodes)
	if err != nil {
		return nil, err
	}

	for _, node := range gitLabNodes {
		gitNodes = append(gitNodes, GitRepoNode{
			Name: node.Name,
			Type: node.Type,
			Path: node.Path,
		})
	}

	return gitNodes, nil
}

// GetAvailableBranches retrieves all available branches.
func (g *GitLabApi) GetAvailableBranches() ([]string, error) {
	var branches []GitBranch
	var branchesStr []string
	fullUrl := CreateUrl(gitlabBranchTemplate, g.base.url, strconv.Itoa(g.projectNumber))
	fmt.Println(fullUrl)

	body, err := HttpGetFunc(fullUrl, g.base.defaultHeader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &branches)
	if err != nil {
		return nil, err
	}

	for _, branch := range branches {
		branchesStr = append(branchesStr, branch.Name)
	}

	return branchesStr, nil
}
