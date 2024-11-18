package options

type Options struct {
	Organization string
	Project      string
	Repo         string
}

func NewOptions() *Options {
	return &Options{}
}

var Current *Options = &Options{
	Organization: "",
	Project:      "",
	Repo:         "",
}
