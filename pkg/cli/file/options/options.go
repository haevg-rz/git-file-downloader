package options

type Options struct {
	OutFile      string
	RepoFilePath string
}

func NewOptions() *Options {
	return &Options{}
}

var Current *Options = &Options{
	OutFile:      "",
	RepoFilePath: "",
}
