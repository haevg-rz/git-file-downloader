package api

import (
	"fmt"
	"hash"
	"io"
	"net/http"
	"slices"
)

var (
	HttpGetFunc = httpGetInternal
)

// GitApi Describes the expected behaviour of a gitApiHandler.
type GitApi interface {
	GetAvailableBranches() ([]string, error)
	GetRemoteFile(filePath, branch string) (*GitRepoFile, error)
	GetFilesFromFolder(folderPath, branch string) ([]GitRepoNode, error)
	GetHash() hash.Hash
}

// SharedConfig Defines shared fields for all implementations of GitApi.
type SharedConfig struct {
	url           string
	defaultHeader map[string]string
}

// BaseConfig Describes the base config shared by all implementations.
type BaseConfig struct {
	Url       string
	Auth      string
	UserAgent string
}

// GitRepoFile Describes a single git file, independent of the git-platform
type GitRepoFile struct {
	Name    string
	Sha     string
	Content string
}

// GitRepoNode Describes a generic git-repo-file independent of the specific provider.
type GitRepoNode struct {
	Name string
	Type string
	Path string
}

// GitBranch Describes a singular branch from a remote repository.
type GitBranch struct {
	Name string `json:"name"`
}

// NewSharedConfig Creates a new instance of SharedConfig
func NewSharedConfig() *SharedConfig {
	return &SharedConfig{}
}

// ValidateBranch retrieves the available branches from a remote repository and returns true if the given branch is available.
func ValidateBranch(api GitApi, branch string) (bool, error) {
	branches, err := api.GetAvailableBranches()
	if err != nil {
		return false, err
	}

	return (branches != nil) && slices.Contains(branches, branch), nil
}

// httpGetInternal sends GET-Request with given fullUrl, privateToken (for api) and userAgent. Returns the response body.
func httpGetInternal(fullUrl string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range header {
		req.Header.Set(key, val)
	}

	client := &http.Client{}
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

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
