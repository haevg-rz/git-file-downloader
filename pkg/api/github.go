package api

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"net/url"
)

// GitHubApi is a specific implementation of GitApi. Defines owner and repoName for remote repository.
type GitHubApi struct {
	base  *SharedConfig
	owner string
	repo  string
}

// GitHubConfig defines all fields needed for the api.
type GitHubConfig struct {
	Owner      string
	Repo       string
	ApiVersion string
}

// GitHubRepoNode defines a node returned from the githubApi.
type GitHubRepoNode struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
}

// We ignore the given hash and calculate our own using the base64-decoded file content. It is what it is.

// GitHubRepoFile defines a file returned by the githubApi.
type GitHubRepoFile struct {
	Name    string `json:"name"`
	Sha     string `json:"omitempty"`
	Content string `json:"content"`
}

const (
	// ContentNodeUrlTemplate OWNER, REPO, PATH, BRANCH
	githubNodeTemplate   = "%s/repos/%s/%s/contents/%s?ref=%s"
	githubBranchTemplate = "%s/repos/%s/%s/branches"
)

var _ GitApi = &GitHubApi{}

// todo

// NewGitHubApi creates a new instance of the GitHubApi.
func NewGitHubApi(baseConfig *BaseConfig, githubConfig *GitHubConfig) *GitHubApi {
	return &GitHubApi{
		base: &SharedConfig{
			url: baseConfig.Url,
			defaultHeader: map[string]string{
				"Authorization":        fmt.Sprintf("Bearer %s", baseConfig.Auth),
				"User-Agent":           baseConfig.UserAgent,
				"X-GitHub-Api-Version": githubConfig.ApiVersion,
				"Accept":               "application/vnd.github+json",
			},
		},
		owner: githubConfig.Owner,
		repo:  githubConfig.Repo,
	}
}

// GetHash returns the hash-method for the githubApi.
func (g *GitHubApi) GetHash() hash.Hash {
	return sha1.New()
}

// GetAvailableBranches returns a list of all available branches from the defined repository.
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

// GetRemoteFile returns a GitRepoFile from the given remotePath and branch.
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

	decodedContent, err := base64.StdEncoding.DecodeString(githubFile.Content)
	if err != nil {
		return nil, err
	}

	h := g.GetHash()
	if _, err = h.Write(decodedContent); err != nil {
		return nil, err
	}

	return &GitRepoFile{
		Name:    githubFile.Name,
		Sha:     hex.EncodeToString(h.Sum(nil)),
		Content: githubFile.Content,
	}, nil
}

// GetFilesFromFolder returns a slice of GitRepoNode from a given folderPath and branch.
func (g *GitHubApi) GetFilesFromFolder(folderPath, branch string) ([]GitRepoNode, error) {
	gitHubNodes := make([]GitHubRepoNode, 0)
	gitNodes := make([]GitRepoNode, 0)
	fullUrl := fmt.Sprintf(
		githubNodeTemplate,
		g.base.url,
		url.PathEscape(g.owner),
		url.PathEscape(g.repo),
		folderPath,
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
