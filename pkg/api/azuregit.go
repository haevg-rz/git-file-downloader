package api

type AzureGitApi struct {
	Base *GitApi
}

var _ IGitApi = &AzureGitApi{}

func NewAzureGitApi(auth, userAgent, url string) *AzureGitApi {
	return &AzureGitApi{
		Base: &GitApi{
			AuthToken: auth,
			UserAgent: userAgent,
			Url:       url,
		},
	}
}

func (a *AzureGitApi) GetAvailableBranches() ([]GitLabBranch, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AzureGitApi) BranchExists(s string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AzureGitApi) GetFile(s string, s2 string) (*GitLabFile, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AzureGitApi) GetFilesFromFolder(s string, s2 string) ([]GitLabRepoNode, error) {
	//TODO implement me
	panic("implement me")
}
