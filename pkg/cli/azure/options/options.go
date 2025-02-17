package options

type Options struct {
	Organization string
	Project      string
	Repo         string
	ApiVersion   string
}

func NewOptions() *Options {
	return &Options{}
}

var Current *Options = &Options{
	Organization: "",
	Project:      "",
	Repo:         "",
	ApiVersion:   "7.1",
}
