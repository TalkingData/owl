package biz

import (
	"owl/cfc/conf"
	"owl/common/logger"
	"owl/dao"
)

type Biz struct {
	dao    *dao.Dao
	conf   *conf.Conf
	logger *logger.Logger
}

func NewBiz(dao *dao.Dao, conf *conf.Conf, lg *logger.Logger) *Biz {
	return &Biz{
		dao:    dao,
		conf:   conf,
		logger: lg,
	}
}
