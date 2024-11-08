package options

type Options struct {
	Owner string
	Repo  string
}

func NewOptions() *Options {
	return &Options{}
}

var Current *Options = &Options{
	Owner: "",
	Repo:  "",
}
