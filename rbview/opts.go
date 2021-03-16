package rbview

type Options struct {
	TemplateName string
}

type Option func(*Options)

func Template(name string) Option {
	return func(opts *Options) {
		opts.TemplateName = name
	}
}
