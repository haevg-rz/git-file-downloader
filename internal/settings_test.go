package internal

import (
	"reflect"
	"testing"
)

func TestSettings_Mode(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		want     Mode
	}{
		{
			name: "File mode",
			settings: Settings{
				OutFile:      "output.txt",
				RepoFilePath: "repo/file.txt",
			},
			want: ModeFile,
		},
		{
			name: "Folder mode",
			settings: Settings{
				OutFolder:      "output",
				RepoFolderPath: "repo/folder",
			},
			want: ModeFolder,
		},
		{
			name:     "Undefined mode",
			settings: Settings{},
			want:     ModeUndef,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.settings.Mode(); got != tt.want {
				t.Errorf("Settings.Mode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSettings_IsValid(t *testing.T) {
	tests := []struct {
		name            string
		settings        Settings
		wantValid       bool
		wantMissingArgs []string
		wantErrors      []string
	}{
		{
			name: "Valid settings",
			settings: Settings{
				PrivateToken:  "token",
				OutFile:       "output.txt",
				Branch:        "main",
				ApiUrl:        "https://api.example.com",
				RepoFilePath:  "repo/file.txt",
				ProjectNumber: "123",
			},
			wantValid:       true,
			wantMissingArgs: nil,
			wantErrors:      nil,
		},
		{
			name: "Missing required fields",
			settings: Settings{
				OutFile: "output.txt",
			},
			wantValid:       false,
			wantMissingArgs: []string{FlagNameToken, FlagNameBranch, FlagNameRepoFilePath, FlagNameUrl},
			wantErrors:      nil,
		},
		{
			name: "Conflicting output settings",
			settings: Settings{
				PrivateToken:  "token",
				OutFile:       "output.txt",
				OutFolder:     "output",
				Branch:        "main",
				ApiUrl:        "https://api.example.com",
				RepoFilePath:  "repo/file.txt",
				ProjectNumber: "123",
			},
			wantValid:       false,
			wantMissingArgs: []string{FlagNameRepoFolderPathEscaped},
			wantErrors:      []string{"You can't use both outPath and outFolder"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, gotMissingArgs, gotErrors := tt.settings.IsValid()
			if gotValid != tt.wantValid {
				t.Errorf("Settings.IsValid() gotValid = %v, want %v", gotValid, tt.wantValid)
			}
			if !reflect.DeepEqual(gotMissingArgs, tt.wantMissingArgs) {
				t.Errorf("Settings.IsValid() gotMissingArgs = %v, want %v", gotMissingArgs, tt.wantMissingArgs)
			}
			if !reflect.DeepEqual(gotErrors, tt.wantErrors) {
				t.Errorf("Settings.IsValid() gotErrors = %v, want %v", gotErrors, tt.wantErrors)
			}
		})
	}
}
