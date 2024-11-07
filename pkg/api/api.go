package api

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

var (
	HttpGetFunc = httpGetInternal
)

// IGitApi Describes the expected behaviour of the gitLabApi.
type IGitApi interface {
	GetAvailableBranches() ([]string, error)
	GetRemoteFile(filePath, branch string) (*GitRepoFile, error)
	GetFilesFromFolder(folderPath, branch string) ([]GitRepoNode, error)
}

// Config base struct for all implementations of IGitApi.
type Config struct {
	url           string
	defaultHeader map[string]string
}

// GitRepoFile Describes a single git file, independent of the git-platform
type GitRepoFile struct {
	Name    string
	Sha256  string
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

func NewConfig() *Config {
	return &Config{}
}

func CreateUrl(templateUrl string, args ...string) string {
	escapedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		// TODO FIX
		//escapedArgs[i] = url.PathEscape(arg)
		escapedArgs[i] = arg
	}
	return fmt.Sprintf(templateUrl, escapedArgs...)
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
