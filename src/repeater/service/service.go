package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"owl/common/logger"
	"owl/dto"
	"owl/repeater/backend"
	repProto "owl/repeater/proto"
)

// OwlRepeaterService struct
type OwlRepeaterService struct {
	backend backend.Backend
	logger  *logger.Logger
}

// NewOwlRepeaterService 新建Repeater服务
func NewOwlRepeaterService(backend backend.Backend, logger *logger.Logger) *OwlRepeaterService {
	return &OwlRepeaterService{
		backend: backend,
		logger:  logger,
	}
}

// ReceiveTimeSeriesData 中继器接收数据
func (repSrv *OwlRepeaterService) ReceiveTimeSeriesData(_ context.Context, tsData *repProto.TsData) (*emptypb.Empty, error) {
	repSrv.logger.Debug("repSrv.ReceiveTimeSeriesData called.")
	defer repSrv.logger.Debug("repSrv.ReceiveTimeSeriesData end.")

	empty := new(emptypb.Empty)

	err := repSrv.backend.Write(dto.TransRepeaterTsData2Dto(tsData))
	if err != nil {
		// 汇报类操作，错误不需要返给Agent，记录日志即可
		repSrv.logger.ErrorWithFields(logger.Fields{
			"time_series_data": tsData,
			"error":            err,
		}, "An error occurred while repSrv.backend.Write in repSrv.ReceiveTimeSeriesData.")
	}

	return empty, nil
}
