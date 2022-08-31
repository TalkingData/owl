// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/repeater.proto

package repeater_proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for OwlRepeater service

func NewOwlRepeaterEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for OwlRepeater service

type OwlRepeaterService interface {
	// ReceiveTimeSeriesData 中继器接收数据
	ReceiveTimeSeriesData(ctx context.Context, in *TsData, opts ...client.CallOption) (*emptypb.Empty, error)
}

type owlRepeaterService struct {
	c    client.Client
	name string
}

func NewOwlRepeaterService(name string, c client.Client) OwlRepeaterService {
	return &owlRepeaterService{
		c:    c,
		name: name,
	}
}

func (c *owlRepeaterService) ReceiveTimeSeriesData(ctx context.Context, in *TsData, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "OwlRepeater.ReceiveTimeSeriesData", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OwlRepeater service

type OwlRepeaterHandler interface {
	// ReceiveTimeSeriesData 中继器接收数据
	ReceiveTimeSeriesData(context.Context, *TsData, *emptypb.Empty) error
}

func RegisterOwlRepeaterHandler(s server.Server, hdlr OwlRepeaterHandler, opts ...server.HandlerOption) error {
	type owlRepeater interface {
		ReceiveTimeSeriesData(ctx context.Context, in *TsData, out *emptypb.Empty) error
	}
	type OwlRepeater struct {
		owlRepeater
	}
	h := &owlRepeaterHandler{hdlr}
	return s.Handle(s.NewHandler(&OwlRepeater{h}, opts...))
}

type owlRepeaterHandler struct {
	OwlRepeaterHandler
}

func (h *owlRepeaterHandler) ReceiveTimeSeriesData(ctx context.Context, in *TsData, out *emptypb.Empty) error {
	return h.OwlRepeaterHandler.ReceiveTimeSeriesData(ctx, in, out)
}
