package options

type Options struct {
	OutFolder      string
	RepoFolderPath string

	// not used, as recursive search is currently the default behaviour
	Recursive bool
}

func NewOptions() *Options {
	return &Options{}
}

var Current *Options = &Options{
	OutFolder:      "",
	RepoFolderPath: "",
	Recursive:      false,
}
