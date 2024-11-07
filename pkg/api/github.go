package api

import (
	"encoding/json"
	"fmt"
)

type GitHubApi struct {
	base  *Config
	owner string
	repo  string
}

// GitHub type defs

type GitHubRepoNode struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
}

type GitHubRepoFile struct {
	Name    string `json:"name"`
	Sha     string `json:"sha"`
	Content string `json:"content"`
}

const (
	// ContentNodeUrlTemplate OWNER, REPO, PATH, BRANCH
	githubNodeTemplate   = "%s/repos/%s/%s/contents/%s?ref=%s"
	githubBranchTemplate = "%s/repos/%s/%s/branches"
)

var _ IGitApi = &GitHubApi{}

func NewGitHubApi(bearerToken, userAgent, url, owner, repo string) *GitHubApi {
	return &GitHubApi{
		base: &Config{
			url: url,
			defaultHeader: map[string]string{
				"Authorization": fmt.Sprintf("Bearer: %s", bearerToken),
				"User-Agent":    userAgent,
			},
		},
		owner: owner,
		repo:  repo,
	}
}

func (g *GitHubApi) GetAvailableBranches() ([]string, error) {
	var branches []GitBranch
	var branchesStr []string
	url := CreateUrl(githubBranchTemplate, g.base.url, g.owner, g.repo)

	body, err := HttpGetFunc(url, g.base.defaultHeader)
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

func (g *GitHubApi) GetRemoteFile(path, branch string) (*GitRepoFile, error) {
	var githubFile *GitHubRepoFile
	fullUrl := CreateUrl(githubNodeTemplate, g.base.url, g.owner, g.repo, path, branch)

	body, err := HttpGetFunc(fullUrl, g.base.defaultHeader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &githubFile)
	if err != nil {
		return nil, err
	}

	return &GitRepoFile{
		Name:    githubFile.Name,
		Sha256:  githubFile.Sha,
		Content: githubFile.Content,
	}, nil
}

func (g *GitHubApi) GetFilesFromFolder(path, branch string) ([]GitRepoNode, error) {
	gitHubNodes := make([]GitHubRepoNode, 0)
	gitNodes := make([]GitRepoNode, 0)
	url := CreateUrl(githubNodeTemplate, g.base.url, g.owner, g.repo, path, branch)

	body, err := HttpGetFunc(url, g.base.defaultHeader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &gitHubNodes)
	if err != nil {
		return nil, err
	}

	for _, node := range gitHubNodes {
		gitNodes = append(gitNodes, GitRepoNode{
			Name: node.Name,
			Type: node.Type,
			Path: node.Path,
		})
	}

	return gitNodes, nil
}
