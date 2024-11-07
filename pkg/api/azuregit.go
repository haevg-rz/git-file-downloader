package api

type AzureGitApi struct {
	Base *Config
}

var _ IGitApi = &AzureGitApi{}

func NewAzureGitApi(auth, userAgent, url string) *AzureGitApi {
	return &AzureGitApi{
		Base: &Config{
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
