package api

import (
	"hash"
)

type AzureGitApi struct {
	Base *SharedConfig
}

func (a AzureGitApi) GetHash() hash.Hash {
	//TODO implement me
	panic("implement me")
}

var _ IGitApi = &AzureGitApi{}

func NewAzureGitApi(auth, userAgent, url string) *AzureGitApi {
	return &AzureGitApi{
		Base: &SharedConfig{
			url:           url,
			defaultHeader: nil,
		},
	}
}

func (a AzureGitApi) GetAvailableBranches() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (a AzureGitApi) GetRemoteFile(filePath, branch string) (*GitRepoFile, error) {
	//TODO implement me
	panic("implement me")
}

func (a AzureGitApi) GetFilesFromFolder(folderPath, branch string) ([]GitRepoNode, error) {
	//TODO implement me
	panic("implement me")
}
