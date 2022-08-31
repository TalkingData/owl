package cli

import (
	"google.golang.org/grpc"
	proxyProto "owl/proxy/proto"
)

// NewProxyClient 创建Proxy客户端
func NewProxyClient(addr string) (proxyProto.OwlProxyServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return proxyProto.NewOwlProxyServiceClient(conn), nil
}
