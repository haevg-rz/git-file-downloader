package options

type Options struct {
	Owner      string
	Repo       string
	ApiVersion string
}

func NewOptions() *Options {
	return &Options{}
}

var Current *Options = &Options{
	Owner:      "",
	Repo:       "",
	ApiVersion: "2022-11-28",
}
