package options

type Options struct {
	OutPath        string
	RemotePath     string
	Branch         string
	IncludePattern string
	ExcludePattern string
	GitProvider    string
	Owner          string
	Repo           string
	ProjectNumber  int
	LogLevel       int
	Api            *ApiOptions
}

type ApiOptions struct {
	UserAgent string
	Auth      string
	BaseUrl   string
}

func NewOptions() *Options {
	return &Options{}
}

var Current *Options = &Options{
	OutPath:        "",
	RemotePath:     "",
	Branch:         "main",
	IncludePattern: "",
	ExcludePattern: "",
	Owner:          "",
	Repo:           "",
	GitProvider:    "",
	ProjectNumber:  -1,
	LogLevel:       3,
	Api: &ApiOptions{
		UserAgent: "Go-http-client/1.1",
		Auth:      "",
		BaseUrl:   "",
	},
}
