package api

import (
	"crypto/tls"
	"fmt"
	"hash"
	"io"
	"net/http"
	"slices"
)

var (
	HttpGetFunc = httpGetInternal
)

// IGitApi Describes the expected behaviour of the gitLabApi.
type IGitApi interface {
	GetAvailableBranches() ([]string, error)
	GetRemoteFile(filePath, branch string) (*GitRepoFile, error)
	GetFilesFromFolder(folderPath, branch string) ([]GitRepoNode, error)
	GetHash() hash.Hash
}

// SharedConfig base struct for all implementations of IGitApi.
type SharedConfig struct {
	url           string
	defaultHeader map[string]string
}

// GitRepoFile Describes a single git file, independent of the git-platform
type GitRepoFile struct {
	Name    string
	Sha     string
	Content string
}

type GitRepoNode struct {
	Name string
	Type string
	Path string
}

type GitBranch struct {
	Name string `json:"name"`
}

func NewSharedConfig() *SharedConfig {
	return &SharedConfig{}
}

func ValidateBranch(api IGitApi, branch string) (bool, error) {
	branches, err := api.GetAvailableBranches()
	if err != nil {
		return false, err
	}

	return (branches != nil) && slices.Contains(branches, branch), nil
}

// httpGetInternal sends GET-Request with given fullUrl, privateToken (for api) and userAgent. Returns the response body.
func httpGetInternal(fullUrl string, header map[string]string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range header {
		req.Header.Add(key, val)
	}

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

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
