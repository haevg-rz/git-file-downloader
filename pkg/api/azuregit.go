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

type AzureGitApi struct {
	Base         *SharedConfig
	Organization string
	Project      string
	Repo         string
}

type AzureGitRepoNode struct {
	Name string `json:"omitempty"`
	Path string `json:"path"`
	Type string `json:"gitObjectType"`
}

type AzureGitRepoNodes struct {
	Value []AzureGitRepoNode `json:"value"`
}

type AzureGitBranch struct {
	Name string `json:"name"`
}

type AzureGitBranches struct {
	Value []AzureGitBranch `json:"value"`
}

const (
	// ENDPOINT, ORGANIZATION, PROJECT, REPO

	azureBranchTemplate = "%s/%s/%s/_apis/git/repositories/%s/refs?api-version=7.1"

	// i really want to know why the azure dev ops api is so confusing compared to github & gitlab
	azureItemsTemplate = "%s/%s/%s/_apis/git/repositories/%s/items?scopePath=%s&versionDescriptor.version=%s&api-version=7.1&recursionLevel=oneLevel"

	azureFileTemplate = "%s/%s/%s/_apis/git/repositories/%s/items?scopePath=%s&versionDescriptor.version=%s&api-version=7.1"
)

var _ IGitApi = &AzureGitApi{}

func NewAzureGitApi(auth, userAgent, url, organization, project, repo string) *AzureGitApi {
	return &AzureGitApi{
		Base: &SharedConfig{
			url: url,
			defaultHeader: map[string]string{
				"User-Agent":    userAgent,
				"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(":"+auth))),
			},
		},
		Organization: organization,
		Project:      project,
		Repo:         repo,
	}
}

func (a *AzureGitApi) GetHash() hash.Hash {
	return sha1.New()
}

func (a *AzureGitApi) GetAvailableBranches() ([]string, error) {
	var branches *AzureGitBranches
	var branchesStr []string
	fullUrl := fmt.Sprintf(
		azureBranchTemplate,
		a.Base.url,
		url.PathEscape(a.Organization),
		url.PathEscape(a.Project),
		url.PathEscape(a.Repo))

	body, err := HttpGetFunc(fullUrl, a.Base.defaultHeader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &branches)
	if err != nil {
		return nil, err
	}

	for _, branch := range branches.Value {
		_, branchName := path.Split(branch.Name)
		branchesStr = append(branchesStr, branchName)
	}

	return branchesStr, nil
}

func (a *AzureGitApi) GetRemoteFile(filePath, branch string) (*GitRepoFile, error) {
	fullUrl := fmt.Sprintf(
		azureFileTemplate,
		a.Base.url,
		url.PathEscape(a.Organization),
		url.PathEscape(a.Project),
		url.PathEscape(a.Repo),
		filePath,
		branch)

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
		url.PathEscape(branch))

	body, err := HttpGetFunc(fullUrl, a.Base.defaultHeader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &azureNodes)
	if err != nil {
		return nil, err
	}

	// set node name manually
	// i couldn't figure out how to do it with the api. please help.
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
