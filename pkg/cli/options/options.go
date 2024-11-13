package options

type Options struct {
	OutPath        string
	RemotePath     string
	Branch         string
	IncludePattern string
	ExcludePattern string
	LogLevel       int
	LogToFile      bool
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
	LogLevel:       3,
	LogToFile:      false,
	Api: &ApiOptions{
		UserAgent: "Go-http-client/1.1",
		Auth:      "",
	},
}
