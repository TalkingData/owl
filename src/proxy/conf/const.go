package conf

import "owl/common/global"

type constConf struct {
	ServiceName string

	CfcServiceName      string
	RepeaterServiceName string
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName: global.OwlProxyServiceName,

		CfcServiceName:      global.OwlCfcServiceName,
		RepeaterServiceName: global.OwlRepeaterServiceName,
	}
}
