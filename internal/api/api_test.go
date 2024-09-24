package api

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/dhcgn/GitLabFileDownloader/internal"
	"github.com/pkg/errors"
)

func TestGetBranches(t *testing.T) {
	mockResponse := `[{"name": "master"}, {"name": "develop"}]`

	HttpGetFunc = func(url string, s internal.Settings) ([]byte, error) {
		if strings.Contains(url, "/repository/branches") {
			return []byte(mockResponse), nil
		}
		return nil, errors.New("Unknown TESTING URL")
	}

	settings := internal.Settings{
		ApiUrl:        "https://gitlab.com/api/v4/",
		ProjectNumber: "123456",
		PrivateToken:  "test-token",
		UserAgent:     "test-agent",
	}

	branches, err := GetBranches(settings)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	if branches[0].Name != "master" {
		t.Errorf("expected branch name 'master', got %s", branches[0].Name)
	}
}

func TestGetFilesFromFolder(t *testing.T) {
	mockResponse := `[{"id": "1", "name": "file1.txt", "type": "blob", "path": "path/to/file1.txt", "mode": "100644"}]`

	HttpGetFunc = func(url string, s internal.Settings) ([]byte, error) {
		if strings.Contains(url, "/repository/tree") {
			return []byte(mockResponse), nil
		}
		return nil, errors.New("Unknown TESTING URL")
	}

	settings := internal.Settings{
		ApiUrl:         "https://gitlab.com/api/v4/",
		ProjectNumber:  "123456",
		PrivateToken:   "test-token",
		UserAgent:      "test-agent",
		RepoFolderPath: "path/to",
		Branch:         "master",
	}

	files, err := GetFilesFromFolder(settings)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	if files[0].Name != "file1.txt" {
		t.Errorf("expected file name 'file1.txt', got %s", files[0].Name)
	}
}

func TestGetFile(t *testing.T) {
	mockResponse := `{
		"file_name": "file1.txt",
		"content_sha256": "dummyhash",
		"content": "dGVzdCBjb250ZW50"
	}`

	HttpGetFunc = func(url string, s internal.Settings) ([]byte, error) {
		if strings.Contains(url, "/repository/files") {
			return []byte(mockResponse), nil
		}
		return nil, errors.New("Unknown TESTING URL")
	}

	settings := internal.Settings{
		ApiUrl:        "https://gitlab.com/api/v4/",
		ProjectNumber: "123456",
		PrivateToken:  "test-token",
		UserAgent:     "test-agent",
		RepoFilePath:  "path/to/file1.txt",
		Branch:        "master",
	}

	file, err := GetFile(settings)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if file.FileName != "file1.txt" {
		t.Errorf("expected file name 'file1.txt', got %s", file.FileName)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		t.Fatalf("failed to decode content: %v", err)
	}

	expectedContent := "test content"
	if string(decodedContent) != expectedContent {
		t.Errorf("expected content '%s', got '%s'", expectedContent, decodedContent)
	}
}
