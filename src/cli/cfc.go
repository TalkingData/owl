package cli

import (
	"google.golang.org/grpc"
	cfcProto "owl/cfc/proto"
)

// NewCfcClient 创建Cfc客户端
func NewCfcClient(addr string) (cfcProto.OwlCfcServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return cfcProto.NewOwlCfcServiceClient(conn), nil
}
