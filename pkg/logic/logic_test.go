package logic

import (
	"crypto/sha1"
	"github.com/haevg-rz/git-file-downloader/pkg/api"
	"github.com/haevg-rz/git-file-downloader/pkg/mode"
	"github.com/stretchr/testify/assert"
	"hash"
	"testing"
)

type MockGitApi struct {
}

var _ api.GitApi = &MockGitApi{}

func (m *MockGitApi) GetAvailableBranches() ([]string, error) {
	return []string{"main", "master"}, nil
}

func (m *MockGitApi) GetRemoteFile(filePath, branch string) (*api.GitRepoFile, error) {
	return &api.GitRepoFile{
		Name:    "test.txt",
		Sha:     "bar",
		Content: "foo",
	}, nil
}

func (m *MockGitApi) GetFilesFromFolder(folderPath, branch string) ([]api.GitRepoNode, error) {
	return []api.GitRepoNode{
		{
			Name: "asdf",
			Type: "asdf",
			Path: "asdf",
		},
	}, nil
}

func (m *MockGitApi) GetHash() hash.Hash {
	return sha1.New()
}

func TestLogic_Handle(t *testing.T) {
	fileDownloader := NewGitFileDownloader(&MockGitApi{})
	err := fileDownloader.Handle(&Context{
		OutPath:    "./asdf.txt",
		RemotePath: "asdf",
		Branch:     "main",
		Patterns:   nil,
	}, mode.File)

	assert.NoError(t, err)
}
