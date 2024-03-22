package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"owl/common/logger"
	commonpb "owl/common/proto"
)

func (proxySrv *OwlProxyService) ReceiveTimeSeriesData(
	ctx context.Context, req *commonpb.TsData,
) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.ReceiveTimeSeriesData called.")
	defer proxySrv.logger.Debug("proxySrv.ReceiveTimeSeriesData end.")

	ret, err := proxySrv.repCli.ReceiveTimeSeriesData(ctx, req)
	if err != nil {
		proxySrv.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling proxySrv.repCli.ReceiveTimeSeriesData, Skipped.")
		return ret, err
	}

	return ret, nil
}

func (proxySrv *OwlProxyService) ReceiveTimeSeriesDataArray(
	ctx context.Context, req *commonpb.TsDataArray,
) (*emptypb.Empty, error) {
	proxySrv.logger.Debug("proxySrv.ReceiveTimeSeriesDataArray called.")
	defer proxySrv.logger.Debug("proxySrv.ReceiveTimeSeriesDataArray end.")

	ret, err := proxySrv.repCli.ReceiveTimeSeriesDataArray(ctx, req)
	if err != nil {
		proxySrv.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while calling proxySrv.repCli.ReceiveTimeSeriesDataArray, Skipped.")
		return ret, err
	}

	return ret, nil
}
