package options

type Options struct {
	OutPath        string
	Branch         string
	IncludePattern string
	ExcludePattern string
	GitProvider    string
	LogLevel       int
	Api            *ApiOptions
}

type ApiOptions struct {
	UserAgent     string
	PrivateToken  string
	BaseUrl       string
	ProjectNumber int
}

func NewOptions() *Options {
	return &Options{}
}

var Current *Options = &Options{
	OutPath:        "",
	Branch:         "main",
	IncludePattern: "",
	ExcludePattern: "",
	GitProvider:    "",
	LogLevel:       1,
	Api: &ApiOptions{
		UserAgent:     "Go-http-client/1.1",
		PrivateToken:  "",
		BaseUrl:       "https://gitlab.com/api/v4",
		ProjectNumber: -1,
	},
}
