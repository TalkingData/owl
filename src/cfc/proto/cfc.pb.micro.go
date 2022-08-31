// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/cfc.proto

package cfc_proto

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

// Api Endpoints for OwlCfcService service

func NewOwlCfcServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for OwlCfcService service

type OwlCfcService interface {
	// RegisterAgent 客户端注册
	RegisterAgent(ctx context.Context, in *AgentInfo, opts ...client.CallOption) (*emptypb.Empty, error)
	// ListAgentPlugins 列出客户端需要执行的插件
	ListAgentPlugins(ctx context.Context, in *HostIdReq, opts ...client.CallOption) (*Plugins, error)
	// ReceiveAgentHeartbeat 接收客户端上报的心跳
	ReceiveAgentHeartbeat(ctx context.Context, in *AgentInfo, opts ...client.CallOption) (*emptypb.Empty, error)
	// ReceiveAgentMetric 接收客户端上报的Metric
	ReceiveAgentMetric(ctx context.Context, in *Metric, opts ...client.CallOption) (*emptypb.Empty, error)
}

type owlCfcService struct {
	c    client.Client
	name string
}

func NewOwlCfcService(name string, c client.Client) OwlCfcService {
	return &owlCfcService{
		c:    c,
		name: name,
	}
}

func (c *owlCfcService) RegisterAgent(ctx context.Context, in *AgentInfo, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "OwlCfcService.RegisterAgent", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *owlCfcService) ListAgentPlugins(ctx context.Context, in *HostIdReq, opts ...client.CallOption) (*Plugins, error) {
	req := c.c.NewRequest(c.name, "OwlCfcService.ListAgentPlugins", in)
	out := new(Plugins)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *owlCfcService) ReceiveAgentHeartbeat(ctx context.Context, in *AgentInfo, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "OwlCfcService.ReceiveAgentHeartbeat", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *owlCfcService) ReceiveAgentMetric(ctx context.Context, in *Metric, opts ...client.CallOption) (*emptypb.Empty, error) {
	req := c.c.NewRequest(c.name, "OwlCfcService.ReceiveAgentMetric", in)
	out := new(emptypb.Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OwlCfcService service

type OwlCfcServiceHandler interface {
	// RegisterAgent 客户端注册
	RegisterAgent(context.Context, *AgentInfo, *emptypb.Empty) error
	// ListAgentPlugins 列出客户端需要执行的插件
	ListAgentPlugins(context.Context, *HostIdReq, *Plugins) error
	// ReceiveAgentHeartbeat 接收客户端上报的心跳
	ReceiveAgentHeartbeat(context.Context, *AgentInfo, *emptypb.Empty) error
	// ReceiveAgentMetric 接收客户端上报的Metric
	ReceiveAgentMetric(context.Context, *Metric, *emptypb.Empty) error
}

func RegisterOwlCfcServiceHandler(s server.Server, hdlr OwlCfcServiceHandler, opts ...server.HandlerOption) error {
	type owlCfcService interface {
		RegisterAgent(ctx context.Context, in *AgentInfo, out *emptypb.Empty) error
		ListAgentPlugins(ctx context.Context, in *HostIdReq, out *Plugins) error
		ReceiveAgentHeartbeat(ctx context.Context, in *AgentInfo, out *emptypb.Empty) error
		ReceiveAgentMetric(ctx context.Context, in *Metric, out *emptypb.Empty) error
	}
	type OwlCfcService struct {
		owlCfcService
	}
	h := &owlCfcServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&OwlCfcService{h}, opts...))
}

type owlCfcServiceHandler struct {
	OwlCfcServiceHandler
}

func (h *owlCfcServiceHandler) RegisterAgent(ctx context.Context, in *AgentInfo, out *emptypb.Empty) error {
	return h.OwlCfcServiceHandler.RegisterAgent(ctx, in, out)
}

func (h *owlCfcServiceHandler) ListAgentPlugins(ctx context.Context, in *HostIdReq, out *Plugins) error {
	return h.OwlCfcServiceHandler.ListAgentPlugins(ctx, in, out)
}

func (h *owlCfcServiceHandler) ReceiveAgentHeartbeat(ctx context.Context, in *AgentInfo, out *emptypb.Empty) error {
	return h.OwlCfcServiceHandler.ReceiveAgentHeartbeat(ctx, in, out)
}

func (h *owlCfcServiceHandler) ReceiveAgentMetric(ctx context.Context, in *Metric, out *emptypb.Empty) error {
	return h.OwlCfcServiceHandler.ReceiveAgentMetric(ctx, in, out)
}
