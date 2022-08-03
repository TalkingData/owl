package component

import (
	"context"
	"owl/agent/conf"
	"owl/common/logger"
)

// Component 组件接口
type Component interface {
	// Start 启动Component服务
	Start() error
	// Stop 关闭Component服务
	Stop()
}

// NewAgentComponent 创建Agent组件
func NewAgentComponent(ctx context.Context, conf *conf.Conf, lg *logger.Logger) (Component, error) {
	return newAgent(ctx, conf, lg)
}
