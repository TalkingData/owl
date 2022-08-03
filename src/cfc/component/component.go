package component

import (
	"context"
	"owl/cfc/conf"
	"owl/common/logger"
	"owl/dao"
)

// Component 组件接口
type Component interface {
	// Start 启动Component服务
	Start() error
	// Stop 关闭Component服务
	Stop()
}

// NewCfcComponent 创建Cfc组件
func NewCfcComponent(ctx context.Context, dao *dao.Dao, conf *conf.Conf, lg *logger.Logger) Component {
	return newCfc(ctx, dao, conf, lg)
}
