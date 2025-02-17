package api

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"net/url"
	"path"
)

// AzureGitApi defines a specific implementation of GitApi.
type AzureGitApi struct {
	Base         *SharedConfig
	Organization string
	Project      string
	Repo         string
	ApiVersion   string
}

// AzureGitConfig defines all fields needed for the api.
type AzureGitConfig struct {
	ApiVersion   string
	Organization string
	Project      string
	Repo         string
}

// AzureGitRepoNode defines a single node returned from the azureGitApi.
type AzureGitRepoNode struct {
	Name string `json:"omitempty"`
	Path string `json:"path"`
	Type string `json:"gitObjectType"`
}

// AzureGitRepoNodes defines a slice of AzureGitRepoNode returned from the azureGitApi.
type AzureGitRepoNodes struct {
	Value []AzureGitRepoNode `json:"value"`
}

// AzureGitBranch defines an individual branch from a remote azure repository. Is returned from the azureGitApi.
type AzureGitBranch struct {
	Name string `json:"name"`
}

// AzureGitBranches defines the available branches of a remote azure repository. Is returned from the azureGitApi.
type AzureGitBranches struct {
	Value []AzureGitBranch `json:"value"`
}

const (
	// ENDPOINT, ORGANIZATION, PROJECT, REPO

	azureBranchTemplate = "%s/%s/%s/_apis/git/repositories/%s/refs?api-version=%s"

	// I really want to know why the azure dev ops api is so confusing compared to gitHub & gitlab
	azureItemsTemplate = "%s/%s/%s/_apis/git/repositories/%s/items?scopePath=%s&versionDescriptor.version=%s&api-version=%s&recursionLevel=oneLevel"

	azureFileTemplate = "%s/%s/%s/_apis/git/repositories/%s/items?scopePath=%s&versionDescriptor.version=%s&api-version=%s"
)

var _ GitApi = &AzureGitApi{}

// NewAzureGitApi creates a new instance of AzureGitApi.
func NewAzureGitApi(baseConfig *BaseConfig, azureConfig *AzureGitConfig) *AzureGitApi {
	return &AzureGitApi{
		Base: &SharedConfig{
			url: baseConfig.Url,
			defaultHeader: map[string]string{
				"User-Agent":    baseConfig.UserAgent,
				"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(":"+baseConfig.Auth))),
			},
		},
		Organization: azureConfig.Organization,
		Project:      azureConfig.Project,
		Repo:         azureConfig.Repo,
		ApiVersion:   azureConfig.ApiVersion,
	}
}

// GetHash returns the hash-method for the azureGitApi.
func (a *AzureGitApi) GetHash() hash.Hash {
	return sha1.New()
}

// GetAvailableBranches returns a list of all available branches from the defined repository.
func (a *AzureGitApi) GetAvailableBranches() ([]string, error) {
	var branchesStr []string
	branches := &AzureGitBranches{}
	fullUrl := fmt.Sprintf(
		azureBranchTemplate,
		a.Base.url,
		url.PathEscape(a.Organization),
		url.PathEscape(a.Project),
		url.PathEscape(a.Repo),
		a.ApiVersion)

	body, err := HttpGetFunc(fullUrl, a.Base.defaultHeader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, branches)
	if err != nil {
		return nil, err
	}

	for _, branch := range branches.Value {
		_, branchName := path.Split(branch.Name)
		branchesStr = append(branchesStr, branchName)
	}

	return branchesStr, nil
}

// GetRemoteFile returns a GitRepoFile from the given remotePath and branch.
func (a *AzureGitApi) GetRemoteFile(filePath, branch string) (*GitRepoFile, error) {
	fullUrl := fmt.Sprintf(
		azureFileTemplate,
		a.Base.url,
		url.PathEscape(a.Organization),
		url.PathEscape(a.Project),
		url.PathEscape(a.Repo),
		filePath,
		branch,
		a.ApiVersion)

	body, err := HttpGetFunc(fullUrl, a.Base.defaultHeader)
	if err != nil {
		return nil, err
	}

	_, fileName := path.Split(filePath)
	contentEncoded := base64.StdEncoding.EncodeToString(body)

	h := a.GetHash()
	if _, err = h.Write(body); err != nil {
		return nil, err
	}

	return &GitRepoFile{
		Name:    fileName,
		Sha:     hex.EncodeToString(h.Sum(nil)),
		Content: contentEncoded,
	}, nil
}

// GetFilesFromFolder returns a slice of GitRepoNode from a given folderPath and branch.
func (a *AzureGitApi) GetFilesFromFolder(folderPath, branch string) ([]GitRepoNode, error) {
	var gitNodes []GitRepoNode
	var azureNodes *AzureGitRepoNodes

	fullUrl := fmt.Sprintf(
		azureItemsTemplate,
		a.Base.url,
		url.PathEscape(a.Organization),
		url.PathEscape(a.Project),
		url.PathEscape(a.Repo),
		folderPath,
		url.PathEscape(branch),
		a.ApiVersion)

	body, err := HttpGetFunc(fullUrl, a.Base.defaultHeader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &azureNodes)
	if err != nil {
		return nil, err
	}

	// set node name manually
	// I couldn't figure out how to do it with the api. please help.
	_, folderName := path.Split(folderPath)
	for _, node := range azureNodes.Value {
		_, file := path.Split(node.Path)
		if file == "" || file == folderName {
			continue
		}
		gitNodes = append(gitNodes, GitRepoNode{
			Name: file,
			Type: node.Type,
			Path: node.Path,
		})
	}

	return gitNodes, nil
}
