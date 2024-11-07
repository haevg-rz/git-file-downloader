package api

type GitHubApi struct {
	Base *GitApi
}

type GitHubApiPathNode struct {
	name     string
	nodeType string
}

const (
	// ContentPath OWNER, REPO, PATH, BRANCH
	ContentPath = "repos/%s/%s/contents/%s?ref=%s"
)

var _ IGitApi = &GitHubApi{}

func NewGitHubApi(auth, userAgent, url string) *GitHubApi {
	return &GitHubApi{
		Base: &GitApi{
			AuthToken: auth,
			UserAgent: userAgent,
			Url:       url,
		},
	}
}

func (g *GitHubApi) GetAvailableBranches() ([]GitLabBranch, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GitHubApi) BranchExists(s string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GitHubApi) GetFile(s string, s2 string) (*GitLabFile, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GitHubApi) GetFilesFromFolder(s string, s2 string) ([]GitLabRepoNode, error) {
	//TODO implement me
	panic("implement me")
}
