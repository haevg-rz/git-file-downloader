package options

type Options struct {
	ProjectId int
}

func NewOptions() *Options {
	return &Options{}
}

var Current *Options = &Options{
	ProjectId: -1,
}
