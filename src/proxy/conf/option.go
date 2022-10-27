package conf

const (
	defaultConfFilePathname = "../conf/owl_proxy.conf"
)

type Option func(o *Options)

// Option struct
type Options struct {
	ConfFilePathname string
}

func newOptions(opts ...Option) Options {
	opt := Options{
		ConfFilePathname: defaultConfFilePathname,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// ConfFilePathname 设置ConfFilePathname
func ConfFilePathname(in string) Option {
	return func(o *Options) {
		o.ConfFilePathname = in
	}
}
