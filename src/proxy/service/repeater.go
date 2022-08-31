package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"owl/common/logger"
	"owl/dto"
	proxyProto "owl/proxy/proto"
)

func (proxySrv *OwlProxyService) ReceiveTimeSeriesData(ctx context.Context, req *proxyProto.TsData) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.ReceiveTimeSeriesData called.")
	defer proxySrv.logger.Debug("proxySrv.ReceiveTimeSeriesData end.")

	empty := new(emptypb.Empty)

	_, err := proxySrv.repCli.ReceiveTimeSeriesData(ctx, dto.TransProxyTsData2Repeater(req))
	if err != nil {
		proxySrv.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while proxySrv.repCli.ReceiveTimeSeriesData, Skipped.")
		return empty, err
	}

	return empty, nil
}
