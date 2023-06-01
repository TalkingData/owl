package conf

import "owl/common/global"

type constConf struct {
	ServiceName    string
	RpcRegisterKey string
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName:    global.OwlRepeaterServiceName,
		RpcRegisterKey: global.OwlRepeaterRpcRegisterKey,
	}
}
