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
	GetAvailableBranches() ([]GitLabBranch, error)
	BranchExists(string) (bool, error)
	GetFile(string, string) (*GitLabFile, error)
	GetFilesFromFolder(string, string) ([]GitLabRepoNode, error)
}

// GitApi Base struct for all implementations of IGitApi.
type GitApi struct {
	AuthToken string
	UserAgent string
	Url       string
}

// GitFile Describes a single git file, independent of the git-platform
type GitFile struct {
	Name          string
	ContentSha256 string
	Content       string
}

// httpGetInternal sends GET-Request with given fullUrl, privateToken (for api) and userAgent. Returns the response body.
func httpGetInternal(fullUrl, privateToken, userAgent string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Private-Token", privateToken)
	req.Header.Add("User-Agent", userAgent)

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

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
