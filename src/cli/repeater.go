package cli

import (
	"google.golang.org/grpc"
	repProto "owl/repeater/proto"
)

// NewRepeaterClient 创建Repeater客户端
func NewRepeaterClient(addr string) (repProto.OwlRepeaterServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return repProto.NewOwlRepeaterServiceClient(conn), nil
}
