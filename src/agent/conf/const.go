package conf

type constConf struct {
	ServiceName string
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName: "owl-agent",
	}
}
