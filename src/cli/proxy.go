package cli

import (
	"google.golang.org/grpc"
	proxypb "owl/proxy/proto"
)

// NewProxyClient 创建Proxy客户端
func NewProxyClient(addr string) (proxypb.OwlProxyServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return proxypb.NewOwlProxyServiceClient(conn), nil
}
