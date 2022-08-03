package component

import (
	"owl/common/logger"
	"owl/repeater/conf"
)

// Component 组件接口
type Component interface {
	// Start 启动Component服务
	Start() error
	// Stop 关闭Component服务
	Stop()
}

// NewRepeaterComponent 创建Repeater组件
func NewRepeaterComponent(conf *conf.Conf, lg *logger.Logger) Component {
	return newRepeater(conf, lg)
}
