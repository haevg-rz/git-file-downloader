package options

type Options struct {
	OutPath        string
	Branch         string
	IncludePattern string
	ExcludePattern string
	LogLevel       int
	Api            *ApiOptions
}

type ApiOptions struct {
	UserAgent     string
	PrivateToken  string
	ApiBaseUrl    string
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
	LogLevel:       1,
	Api: &ApiOptions{
		UserAgent:     "Go-http-client/1.1",
		PrivateToken:  "",
		ApiBaseUrl:    "https://gitlab.com/api/v4/",
		ProjectNumber: -1,
	},
}
