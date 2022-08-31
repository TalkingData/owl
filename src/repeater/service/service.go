package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"owl/common/logger"
	"owl/dto"
	"owl/repeater/backend"
	"owl/repeater/conf"
	repProto "owl/repeater/proto"
)

// OwlRepeaterService struct
type OwlRepeaterService struct {
	backend backend.Backend
	logger  *logger.Logger
}

// NewOwlRepeaterService 新建Repeater服务
func NewOwlRepeaterService(conf *conf.Conf, lg *logger.Logger) *OwlRepeaterService {
	bk, err := backend.NewBackend(conf)
	if err != nil {
		lg.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while NewOwlRepeaterService.")
		panic(err)
	}

	return &OwlRepeaterService{
		backend: bk,
		logger:  lg,
	}
}

// ReceiveTimeSeriesData 中继器接收数据
func (repSrv *OwlRepeaterService) ReceiveTimeSeriesData(_ context.Context, tsData *repProto.TsData, _ *emptypb.Empty) error {
	repSrv.logger.Debug("repSrv.ReceiveTimeSeriesData called.")
	defer repSrv.logger.Debug("repSrv.ReceiveTimeSeriesData end.")

	if err := repSrv.backend.Write(dto.TransRepeaterTsData2Dto(tsData)); err != nil {
		// 汇报类操作，错误不需要返给Agent，记录日志即可
		repSrv.logger.ErrorWithFields(logger.Fields{
			"time_series_data": tsData,
			"error":            err,
		}, "An error occurred while repSrv.backend.Write in repSrv.ReceiveTimeSeriesData.")
	}

	return nil
}
