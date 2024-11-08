package api

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"hash"
	"net/url"
)

type GitHubApi struct {
	base  *SharedConfig
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
		base: &SharedConfig{
			url: url,
			defaultHeader: map[string]string{
				"Authorization":        fmt.Sprintf("Bearer %s", bearerToken),
				"User-Agent":           userAgent,
				"X-GitHub-Api-Version": "2022-11-28",
				"Accept":               "application/vnd.github+json",
			},
		},
		owner: owner,
		repo:  repo,
	}
}

func (g *GitHubApi) GetHash() hash.Hash {
	return sha1.New()
}

func (g *GitHubApi) GetAvailableBranches() ([]string, error) {
	var branches []GitBranch
	var branchesStr []string
	fullUrl := fmt.Sprintf(
		githubBranchTemplate,
		g.base.url,
		url.PathEscape(g.owner),
		url.PathEscape(g.repo))

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

func (g *GitHubApi) GetRemoteFile(path, branch string) (*GitRepoFile, error) {
	var githubFile *GitHubRepoFile
	fullUrl := fmt.Sprintf(
		githubNodeTemplate,
		g.base.url,
		url.PathEscape(g.owner),
		url.PathEscape(g.repo),
		path,
		url.PathEscape(branch))

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
		Sha:     githubFile.Sha,
		Content: githubFile.Content,
	}, nil
}

func (g *GitHubApi) GetFilesFromFolder(path, branch string) ([]GitRepoNode, error) {
	gitHubNodes := make([]GitHubRepoNode, 0)
	gitNodes := make([]GitRepoNode, 0)
	fullUrl := fmt.Sprintf(
		githubNodeTemplate,
		g.base.url,
		url.PathEscape(g.owner),
		url.PathEscape(g.repo),
		path,
		url.PathEscape(branch))

	body, err := HttpGetFunc(fullUrl, g.base.defaultHeader)
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
