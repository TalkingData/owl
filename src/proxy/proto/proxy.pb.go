// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/proxy.proto

package proxy_proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	math "math"
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

type TsData struct {
	Metric               string            `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
	DataType             string            `protobuf:"bytes,2,opt,name=data_type,json=dataType,proto3" json:"data_type,omitempty"`
	Value                float64           `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`
	Timestamp            int64             `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Cycle                int32             `protobuf:"varint,5,opt,name=cycle,proto3" json:"cycle,omitempty"`
	Tags                 map[string]string `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *TsData) Reset()         { *m = TsData{} }
func (m *TsData) String() string { return proto.CompactTextString(m) }
func (*TsData) ProtoMessage()    {}
func (*TsData) Descriptor() ([]byte, []int) {
	return fileDescriptor_58b4a54be18c47e6, []int{0}
}

func (m *TsData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TsData.Unmarshal(m, b)
}
func (m *TsData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TsData.Marshal(b, m, deterministic)
}
func (m *TsData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TsData.Merge(m, src)
}
func (m *TsData) XXX_Size() int {
	return xxx_messageInfo_TsData.Size(m)
}
func (m *TsData) XXX_DiscardUnknown() {
	xxx_messageInfo_TsData.DiscardUnknown(m)
}

var xxx_messageInfo_TsData proto.InternalMessageInfo

func (m *TsData) GetMetric() string {
	if m != nil {
		return m.Metric
	}
	return ""
}

func (m *TsData) GetDataType() string {
	if m != nil {
		return m.DataType
	}
	return ""
}

func (m *TsData) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *TsData) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *TsData) GetCycle() int32 {
	if m != nil {
		return m.Cycle
	}
	return 0
}

func (m *TsData) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

type HostIdReq struct {
	HostId               string   `protobuf:"bytes,1,opt,name=host_id,json=hostId,proto3" json:"host_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HostIdReq) Reset()         { *m = HostIdReq{} }
func (m *HostIdReq) String() string { return proto.CompactTextString(m) }
func (*HostIdReq) ProtoMessage()    {}
func (*HostIdReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_58b4a54be18c47e6, []int{1}
}

func (m *HostIdReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HostIdReq.Unmarshal(m, b)
}
func (m *HostIdReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HostIdReq.Marshal(b, m, deterministic)
}
func (m *HostIdReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HostIdReq.Merge(m, src)
}
func (m *HostIdReq) XXX_Size() int {
	return xxx_messageInfo_HostIdReq.Size(m)
}
func (m *HostIdReq) XXX_DiscardUnknown() {
	xxx_messageInfo_HostIdReq.DiscardUnknown(m)
}

var xxx_messageInfo_HostIdReq proto.InternalMessageInfo

func (m *HostIdReq) GetHostId() string {
	if m != nil {
		return m.HostId
	}
	return ""
}

type DownloadPluginReq struct {
	RelPath              string   `protobuf:"bytes,1,opt,name=rel_path,json=relPath,proto3" json:"rel_path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DownloadPluginReq) Reset()         { *m = DownloadPluginReq{} }
func (m *DownloadPluginReq) String() string { return proto.CompactTextString(m) }
func (*DownloadPluginReq) ProtoMessage()    {}
func (*DownloadPluginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_58b4a54be18c47e6, []int{2}
}

func (m *DownloadPluginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DownloadPluginReq.Unmarshal(m, b)
}
func (m *DownloadPluginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DownloadPluginReq.Marshal(b, m, deterministic)
}
func (m *DownloadPluginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DownloadPluginReq.Merge(m, src)
}
func (m *DownloadPluginReq) XXX_Size() int {
	return xxx_messageInfo_DownloadPluginReq.Size(m)
}
func (m *DownloadPluginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_DownloadPluginReq.DiscardUnknown(m)
}

var xxx_messageInfo_DownloadPluginReq proto.InternalMessageInfo

func (m *DownloadPluginReq) GetRelPath() string {
	if m != nil {
		return m.RelPath
	}
	return ""
}

type AgentInfo struct {
	HostId               string            `protobuf:"bytes,1,opt,name=host_id,json=hostId,proto3" json:"host_id,omitempty"`
	Ip                   string            `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	Hostname             string            `protobuf:"bytes,3,opt,name=hostname,proto3" json:"hostname,omitempty"`
	AgentVersion         string            `protobuf:"bytes,4,opt,name=agent_version,json=agentVersion,proto3" json:"agent_version,omitempty"`
	AgentOs              string            `protobuf:"bytes,5,opt,name=agent_os,json=agentOs,proto3" json:"agent_os,omitempty"`
	AgentArch            string            `protobuf:"bytes,6,opt,name=agent_arch,json=agentArch,proto3" json:"agent_arch,omitempty"`
	Uptime               float64           `protobuf:"fixed64,7,opt,name=uptime,proto3" json:"uptime,omitempty"`
	IdlePct              float64           `protobuf:"fixed64,8,opt,name=idle_pct,json=idlePct,proto3" json:"idle_pct,omitempty"`
	Metadata             map[string]string `protobuf:"bytes,9,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AgentInfo) Reset()         { *m = AgentInfo{} }
func (m *AgentInfo) String() string { return proto.CompactTextString(m) }
func (*AgentInfo) ProtoMessage()    {}
func (*AgentInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_58b4a54be18c47e6, []int{3}
}

func (m *AgentInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AgentInfo.Unmarshal(m, b)
}
func (m *AgentInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AgentInfo.Marshal(b, m, deterministic)
}
func (m *AgentInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AgentInfo.Merge(m, src)
}
func (m *AgentInfo) XXX_Size() int {
	return xxx_messageInfo_AgentInfo.Size(m)
}
func (m *AgentInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_AgentInfo.DiscardUnknown(m)
}

var xxx_messageInfo_AgentInfo proto.InternalMessageInfo

func (m *AgentInfo) GetHostId() string {
	if m != nil {
		return m.HostId
	}
	return ""
}

func (m *AgentInfo) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *AgentInfo) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *AgentInfo) GetAgentVersion() string {
	if m != nil {
		return m.AgentVersion
	}
	return ""
}

func (m *AgentInfo) GetAgentOs() string {
	if m != nil {
		return m.AgentOs
	}
	return ""
}

func (m *AgentInfo) GetAgentArch() string {
	if m != nil {
		return m.AgentArch
	}
	return ""
}

func (m *AgentInfo) GetUptime() float64 {
	if m != nil {
		return m.Uptime
	}
	return 0
}

func (m *AgentInfo) GetIdlePct() float64 {
	if m != nil {
		return m.IdlePct
	}
	return 0
}

func (m *AgentInfo) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Plugins struct {
	Plugins              []*Plugin `protobuf:"bytes,1,rep,name=plugins,proto3" json:"plugins,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Plugins) Reset()         { *m = Plugins{} }
func (m *Plugins) String() string { return proto.CompactTextString(m) }
func (*Plugins) ProtoMessage()    {}
func (*Plugins) Descriptor() ([]byte, []int) {
	return fileDescriptor_58b4a54be18c47e6, []int{4}
}

func (m *Plugins) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Plugins.Unmarshal(m, b)
}
func (m *Plugins) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Plugins.Marshal(b, m, deterministic)
}
func (m *Plugins) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Plugins.Merge(m, src)
}
func (m *Plugins) XXX_Size() int {
	return xxx_messageInfo_Plugins.Size(m)
}
func (m *Plugins) XXX_DiscardUnknown() {
	xxx_messageInfo_Plugins.DiscardUnknown(m)
}

var xxx_messageInfo_Plugins proto.InternalMessageInfo

func (m *Plugins) GetPlugins() []*Plugin {
	if m != nil {
		return m.Plugins
	}
	return nil
}

type Plugin struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Path                 string   `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
	Checksum             string   `protobuf:"bytes,4,opt,name=checksum,proto3" json:"checksum,omitempty"`
	Args                 string   `protobuf:"bytes,5,opt,name=args,proto3" json:"args,omitempty"`
	Interval             int32    `protobuf:"varint,6,opt,name=interval,proto3" json:"interval,omitempty"`
	Timeout              int32    `protobuf:"varint,7,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Comment              string   `protobuf:"bytes,8,opt,name=comment,proto3" json:"comment,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Plugin) Reset()         { *m = Plugin{} }
func (m *Plugin) String() string { return proto.CompactTextString(m) }
func (*Plugin) ProtoMessage()    {}
func (*Plugin) Descriptor() ([]byte, []int) {
	return fileDescriptor_58b4a54be18c47e6, []int{5}
}

func (m *Plugin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Plugin.Unmarshal(m, b)
}
func (m *Plugin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Plugin.Marshal(b, m, deterministic)
}
func (m *Plugin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Plugin.Merge(m, src)
}
func (m *Plugin) XXX_Size() int {
	return xxx_messageInfo_Plugin.Size(m)
}
func (m *Plugin) XXX_DiscardUnknown() {
	xxx_messageInfo_Plugin.DiscardUnknown(m)
}

var xxx_messageInfo_Plugin proto.InternalMessageInfo

func (m *Plugin) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Plugin) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Plugin) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Plugin) GetChecksum() string {
	if m != nil {
		return m.Checksum
	}
	return ""
}

func (m *Plugin) GetArgs() string {
	if m != nil {
		return m.Args
	}
	return ""
}

func (m *Plugin) GetInterval() int32 {
	if m != nil {
		return m.Interval
	}
	return 0
}

func (m *Plugin) GetTimeout() int32 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *Plugin) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

type Metric struct {
	HostId               string            `protobuf:"bytes,1,opt,name=host_id,json=hostId,proto3" json:"host_id,omitempty"`
	Metric               string            `protobuf:"bytes,2,opt,name=metric,proto3" json:"metric,omitempty"`
	DataType             string            `protobuf:"bytes,3,opt,name=data_type,json=dataType,proto3" json:"data_type,omitempty"`
	Cycle                int32             `protobuf:"varint,4,opt,name=cycle,proto3" json:"cycle,omitempty"`
	Tags                 map[string]string `protobuf:"bytes,5,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Metric) Reset()         { *m = Metric{} }
func (m *Metric) String() string { return proto.CompactTextString(m) }
func (*Metric) ProtoMessage()    {}
func (*Metric) Descriptor() ([]byte, []int) {
	return fileDescriptor_58b4a54be18c47e6, []int{6}
}

func (m *Metric) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metric.Unmarshal(m, b)
}
func (m *Metric) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metric.Marshal(b, m, deterministic)
}
func (m *Metric) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metric.Merge(m, src)
}
func (m *Metric) XXX_Size() int {
	return xxx_messageInfo_Metric.Size(m)
}
func (m *Metric) XXX_DiscardUnknown() {
	xxx_messageInfo_Metric.DiscardUnknown(m)
}

var xxx_messageInfo_Metric proto.InternalMessageInfo

func (m *Metric) GetHostId() string {
	if m != nil {
		return m.HostId
	}
	return ""
}

func (m *Metric) GetMetric() string {
	if m != nil {
		return m.Metric
	}
	return ""
}

func (m *Metric) GetDataType() string {
	if m != nil {
		return m.DataType
	}
	return ""
}

func (m *Metric) GetCycle() int32 {
	if m != nil {
		return m.Cycle
	}
	return 0
}

func (m *Metric) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

type PluginFile struct {
	Buffer               []byte   `protobuf:"bytes,1,opt,name=buffer,proto3" json:"buffer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PluginFile) Reset()         { *m = PluginFile{} }
func (m *PluginFile) String() string { return proto.CompactTextString(m) }
func (*PluginFile) ProtoMessage()    {}
func (*PluginFile) Descriptor() ([]byte, []int) {
	return fileDescriptor_58b4a54be18c47e6, []int{7}
}

func (m *PluginFile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PluginFile.Unmarshal(m, b)
}
func (m *PluginFile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PluginFile.Marshal(b, m, deterministic)
}
func (m *PluginFile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PluginFile.Merge(m, src)
}
func (m *PluginFile) XXX_Size() int {
	return xxx_messageInfo_PluginFile.Size(m)
}
func (m *PluginFile) XXX_DiscardUnknown() {
	xxx_messageInfo_PluginFile.DiscardUnknown(m)
}

var xxx_messageInfo_PluginFile proto.InternalMessageInfo

func (m *PluginFile) GetBuffer() []byte {
	if m != nil {
		return m.Buffer
	}
	return nil
}

func init() {
	proto.RegisterType((*TsData)(nil), "proxy_proto.TsData")
	proto.RegisterMapType((map[string]string)(nil), "proxy_proto.TsData.TagsEntry")
	proto.RegisterType((*HostIdReq)(nil), "proxy_proto.HostIdReq")
	proto.RegisterType((*DownloadPluginReq)(nil), "proxy_proto.DownloadPluginReq")
	proto.RegisterType((*AgentInfo)(nil), "proxy_proto.AgentInfo")
	proto.RegisterMapType((map[string]string)(nil), "proxy_proto.AgentInfo.MetadataEntry")
	proto.RegisterType((*Plugins)(nil), "proxy_proto.Plugins")
	proto.RegisterType((*Plugin)(nil), "proxy_proto.Plugin")
	proto.RegisterType((*Metric)(nil), "proxy_proto.Metric")
	proto.RegisterMapType((map[string]string)(nil), "proxy_proto.Metric.TagsEntry")
	proto.RegisterType((*PluginFile)(nil), "proxy_proto.PluginFile")
}

func init() { proto.RegisterFile("proto/proxy.proto", fileDescriptor_58b4a54be18c47e6) }

var fileDescriptor_58b4a54be18c47e6 = []byte{
	// 743 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xcb, 0x4e, 0x1b, 0x4b,
	0x10, 0xd5, 0xf8, 0x35, 0x9e, 0x02, 0xdf, 0x0b, 0x0d, 0x17, 0xe6, 0x9a, 0x10, 0x59, 0x0e, 0x0b,
	0x6f, 0x32, 0x24, 0x64, 0x11, 0x94, 0x2c, 0x02, 0x12, 0x24, 0x20, 0x05, 0x61, 0x35, 0x56, 0xb6,
	0x56, 0x7b, 0x5c, 0x1e, 0xb7, 0x98, 0x57, 0x7a, 0xda, 0x26, 0xfe, 0xb2, 0x6c, 0xf3, 0x19, 0xd9,
	0xe6, 0x1b, 0xf2, 0x03, 0x51, 0x77, 0x8f, 0x07, 0x2c, 0xec, 0x05, 0x52, 0x76, 0x7d, 0xaa, 0x4e,
	0xd7, 0x9c, 0xa9, 0xd3, 0x55, 0xb0, 0x99, 0x8a, 0x44, 0x26, 0x87, 0xa9, 0x48, 0xbe, 0xcd, 0x3c,
	0x7d, 0x26, 0x6b, 0x1a, 0xf4, 0x35, 0x68, 0xee, 0x05, 0x49, 0x12, 0x84, 0x78, 0xa8, 0xd1, 0x60,
	0x32, 0x3a, 0xc4, 0x28, 0x95, 0x39, 0xb3, 0xfd, 0xdb, 0x82, 0x5a, 0x2f, 0x3b, 0x63, 0x92, 0x91,
	0x1d, 0xa8, 0x45, 0x28, 0x05, 0xf7, 0x5d, 0xab, 0x65, 0x75, 0x1c, 0x9a, 0x23, 0xb2, 0x07, 0xce,
	0x90, 0x49, 0xd6, 0x97, 0xb3, 0x14, 0xdd, 0x92, 0x4e, 0xd5, 0x55, 0xa0, 0x37, 0x4b, 0x91, 0x6c,
	0x43, 0x75, 0xca, 0xc2, 0x09, 0xba, 0xe5, 0x96, 0xd5, 0xb1, 0xa8, 0x01, 0xe4, 0x19, 0x38, 0x92,
	0x47, 0x98, 0x49, 0x16, 0xa5, 0x6e, 0xa5, 0x65, 0x75, 0xca, 0xf4, 0x3e, 0xa0, 0xee, 0xf8, 0x33,
	0x3f, 0x44, 0xb7, 0xda, 0xb2, 0x3a, 0x55, 0x6a, 0x00, 0x79, 0x0d, 0x15, 0xc9, 0x82, 0xcc, 0xad,
	0xb5, 0xca, 0x9d, 0xb5, 0xa3, 0x7d, 0xef, 0xc1, 0x2f, 0x78, 0x46, 0xa1, 0xd7, 0x63, 0x41, 0x76,
	0x1e, 0x4b, 0x31, 0xa3, 0x9a, 0xda, 0x7c, 0x0b, 0x4e, 0x11, 0x22, 0x1b, 0x50, 0xbe, 0xc5, 0x59,
	0xae, 0x5d, 0x1d, 0xef, 0xb5, 0x19, 0xd1, 0x06, 0xbc, 0x2b, 0x1d, 0x5b, 0xed, 0x03, 0x70, 0x2e,
	0x92, 0x4c, 0x5e, 0x0e, 0x29, 0x7e, 0x25, 0xbb, 0x60, 0x8f, 0x93, 0x4c, 0xf6, 0xf9, 0x70, 0xfe,
	0xe3, 0x63, 0x9d, 0x6b, 0x7b, 0xb0, 0x79, 0x96, 0xdc, 0xc5, 0x61, 0xc2, 0x86, 0xdd, 0x70, 0x12,
	0xf0, 0x58, 0xb1, 0xff, 0x87, 0xba, 0xc0, 0xb0, 0x9f, 0x32, 0x39, 0xce, 0xe9, 0xb6, 0xc0, 0xb0,
	0xcb, 0xe4, 0xb8, 0xfd, 0xab, 0x04, 0xce, 0x69, 0x80, 0xb1, 0xbc, 0x8c, 0x47, 0xc9, 0xca, 0xb2,
	0xe4, 0x1f, 0x28, 0xf1, 0x34, 0xd7, 0x54, 0xe2, 0x29, 0x69, 0x42, 0x5d, 0x65, 0x62, 0x16, 0x99,
	0x2e, 0x3a, 0xb4, 0xc0, 0xe4, 0x05, 0x34, 0x98, 0xaa, 0xd8, 0x9f, 0xa2, 0xc8, 0x78, 0x12, 0xeb,
	0x66, 0x3a, 0x74, 0x5d, 0x07, 0xbf, 0x98, 0x98, 0x92, 0x64, 0x48, 0x49, 0xa6, 0x5b, 0xea, 0x50,
	0x5b, 0xe3, 0xeb, 0x8c, 0xec, 0x03, 0x98, 0x14, 0x13, 0xfe, 0xd8, 0xad, 0xe9, 0xa4, 0xa3, 0x23,
	0xa7, 0xc2, 0x1f, 0x2b, 0xcb, 0x27, 0xa9, 0x32, 0xc6, 0xb5, 0xb5, 0x7d, 0x39, 0x52, 0x15, 0xf9,
	0x30, 0xc4, 0x7e, 0xea, 0x4b, 0xb7, 0xae, 0x33, 0xb6, 0xc2, 0x5d, 0x5f, 0x92, 0x13, 0xa8, 0x47,
	0x28, 0x99, 0x7a, 0x00, 0xae, 0xa3, 0xad, 0x3a, 0x58, 0xb0, 0xaa, 0x68, 0x80, 0x77, 0x95, 0xd3,
	0x8c, 0x63, 0xc5, 0xad, 0xe6, 0x7b, 0x68, 0x2c, 0xa4, 0x9e, 0xe4, 0xdc, 0x31, 0xd8, 0xc6, 0x8b,
	0x8c, 0xbc, 0x04, 0x3b, 0x35, 0x47, 0xd7, 0xd2, 0x42, 0xb6, 0x16, 0x84, 0xe4, 0x96, 0xcd, 0x39,
	0xed, 0x1f, 0x16, 0xd4, 0x4c, 0x4c, 0x3b, 0x60, 0x5c, 0x69, 0xd0, 0x12, 0x1f, 0x12, 0x02, 0x15,
	0xdd, 0x7d, 0xf3, 0x35, 0x7d, 0x56, 0x31, 0xed, 0xb1, 0x71, 0x44, 0x9f, 0x95, 0x53, 0xfe, 0x18,
	0xfd, 0xdb, 0x6c, 0x12, 0xe5, 0x46, 0x14, 0x58, 0xf1, 0x99, 0x08, 0xe6, 0x06, 0xe8, 0xb3, 0xe2,
	0xf3, 0x58, 0xa2, 0x98, 0xb2, 0x50, 0xf7, 0xbe, 0x4a, 0x0b, 0x4c, 0x5c, 0xb0, 0x55, 0xab, 0x93,
	0x89, 0xd4, 0xbd, 0xaf, 0xd2, 0x39, 0x54, 0x19, 0x3f, 0x89, 0x22, 0x8c, 0x4d, 0xef, 0x1d, 0x3a,
	0x87, 0xed, 0x9f, 0x16, 0xd4, 0xae, 0xcc, 0x50, 0xae, 0x7c, 0x5d, 0xf7, 0x53, 0x5c, 0x5a, 0x3d,
	0xc5, 0xe5, 0xc7, 0x53, 0x6c, 0x26, 0xb2, 0xb2, 0x6c, 0x22, 0xab, 0x4b, 0x26, 0xd2, 0xc8, 0xf8,
	0x9b, 0x13, 0x09, 0xc6, 0x9c, 0x8f, 0x3c, 0x44, 0xf5, 0x13, 0x83, 0xc9, 0x68, 0x84, 0x42, 0x5f,
	0x5e, 0xa7, 0x39, 0x3a, 0xfa, 0x5e, 0x86, 0x7f, 0xaf, 0xef, 0xc2, 0xae, 0x12, 0x72, 0x83, 0x62,
	0xca, 0x7d, 0x24, 0x67, 0xf0, 0x1f, 0x45, 0x1f, 0xf9, 0x14, 0x7b, 0x3c, 0xc2, 0x1b, 0x14, 0x1c,
	0xcd, 0x3e, 0xdb, 0x5a, 0xb2, 0x42, 0x9a, 0x3b, 0x9e, 0xd9, 0x86, 0xde, 0x7c, 0x1b, 0x7a, 0xe7,
	0x6a, 0x1b, 0x92, 0x4f, 0x45, 0x15, 0xfd, 0x80, 0x2f, 0x90, 0x09, 0x39, 0x40, 0x26, 0xc9, 0xce,
	0xf2, 0xd7, 0xbd, 0xb2, 0xd0, 0x29, 0x90, 0x87, 0x85, 0x72, 0xbb, 0xb6, 0x96, 0x34, 0x6f, 0x65,
	0x89, 0x0f, 0xd0, 0xa0, 0x18, 0xf0, 0x4c, 0xa2, 0xd0, 0x35, 0x9e, 0xac, 0xe1, 0x04, 0x36, 0x3e,
	0xf3, 0x4c, 0x6a, 0xe2, 0x7c, 0x5a, 0x16, 0x6b, 0x14, 0xdb, 0xaf, 0xb9, 0xbd, 0x64, 0x68, 0x32,
	0x72, 0x0d, 0x64, 0x71, 0xf5, 0x69, 0x5b, 0x9e, 0x2f, 0x70, 0x1f, 0xed, 0xc6, 0xe6, 0xee, 0x92,
	0x5a, 0xea, 0xe2, 0x2b, 0x6b, 0x50, 0xd3, 0xb1, 0x37, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x83,
	0xd8, 0xd2, 0xee, 0xad, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OwlProxyServiceClient is the client API for OwlProxyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OwlProxyServiceClient interface {
	// ReceiveTimeSeriesData 中继器接收数据
	ReceiveTimeSeriesData(ctx context.Context, in *TsData, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// ReceiveAgentHeartbeat 接收客户端上报的心跳
	ReceiveAgentHeartbeat(ctx context.Context, in *AgentInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// ReceiveAgentMetric 接收客户端上报的Metric
	ReceiveAgentMetric(ctx context.Context, in *Metric, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// RegisterAgent 客户端注册
	RegisterAgent(ctx context.Context, in *AgentInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// ListAgentPlugins 列出客户端需要执行的插件
	ListAgentPlugins(ctx context.Context, in *HostIdReq, opts ...grpc.CallOption) (*Plugins, error)
	// DownloadPluginFile 下载插件文件
	DownloadPluginFile(ctx context.Context, in *DownloadPluginReq, opts ...grpc.CallOption) (OwlProxyService_DownloadPluginFileClient, error)
}

type owlProxyServiceClient struct {
	cc *grpc.ClientConn
}

func NewOwlProxyServiceClient(cc *grpc.ClientConn) OwlProxyServiceClient {
	return &owlProxyServiceClient{cc}
}

func (c *owlProxyServiceClient) ReceiveTimeSeriesData(ctx context.Context, in *TsData, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proxy_proto.OwlProxyService/ReceiveTimeSeriesData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *owlProxyServiceClient) ReceiveAgentHeartbeat(ctx context.Context, in *AgentInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proxy_proto.OwlProxyService/ReceiveAgentHeartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *owlProxyServiceClient) ReceiveAgentMetric(ctx context.Context, in *Metric, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proxy_proto.OwlProxyService/ReceiveAgentMetric", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *owlProxyServiceClient) RegisterAgent(ctx context.Context, in *AgentInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proxy_proto.OwlProxyService/RegisterAgent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *owlProxyServiceClient) ListAgentPlugins(ctx context.Context, in *HostIdReq, opts ...grpc.CallOption) (*Plugins, error) {
	out := new(Plugins)
	err := c.cc.Invoke(ctx, "/proxy_proto.OwlProxyService/ListAgentPlugins", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *owlProxyServiceClient) DownloadPluginFile(ctx context.Context, in *DownloadPluginReq, opts ...grpc.CallOption) (OwlProxyService_DownloadPluginFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &_OwlProxyService_serviceDesc.Streams[0], "/proxy_proto.OwlProxyService/DownloadPluginFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &owlProxyServiceDownloadPluginFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OwlProxyService_DownloadPluginFileClient interface {
	Recv() (*PluginFile, error)
	grpc.ClientStream
}

type owlProxyServiceDownloadPluginFileClient struct {
	grpc.ClientStream
}

func (x *owlProxyServiceDownloadPluginFileClient) Recv() (*PluginFile, error) {
	m := new(PluginFile)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OwlProxyServiceServer is the server API for OwlProxyService service.
type OwlProxyServiceServer interface {
	// ReceiveTimeSeriesData 中继器接收数据
	ReceiveTimeSeriesData(context.Context, *TsData) (*emptypb.Empty, error)
	// ReceiveAgentHeartbeat 接收客户端上报的心跳
	ReceiveAgentHeartbeat(context.Context, *AgentInfo) (*emptypb.Empty, error)
	// ReceiveAgentMetric 接收客户端上报的Metric
	ReceiveAgentMetric(context.Context, *Metric) (*emptypb.Empty, error)
	// RegisterAgent 客户端注册
	RegisterAgent(context.Context, *AgentInfo) (*emptypb.Empty, error)
	// ListAgentPlugins 列出客户端需要执行的插件
	ListAgentPlugins(context.Context, *HostIdReq) (*Plugins, error)
	// DownloadPluginFile 下载插件文件
	DownloadPluginFile(*DownloadPluginReq, OwlProxyService_DownloadPluginFileServer) error
}

func RegisterOwlProxyServiceServer(s *grpc.Server, srv OwlProxyServiceServer) {
	s.RegisterService(&_OwlProxyService_serviceDesc, srv)
}

func _OwlProxyService_ReceiveTimeSeriesData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TsData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OwlProxyServiceServer).ReceiveTimeSeriesData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy_proto.OwlProxyService/ReceiveTimeSeriesData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OwlProxyServiceServer).ReceiveTimeSeriesData(ctx, req.(*TsData))
	}
	return interceptor(ctx, in, info, handler)
}

func _OwlProxyService_ReceiveAgentHeartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OwlProxyServiceServer).ReceiveAgentHeartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy_proto.OwlProxyService/ReceiveAgentHeartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OwlProxyServiceServer).ReceiveAgentHeartbeat(ctx, req.(*AgentInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _OwlProxyService_ReceiveAgentMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Metric)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OwlProxyServiceServer).ReceiveAgentMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy_proto.OwlProxyService/ReceiveAgentMetric",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OwlProxyServiceServer).ReceiveAgentMetric(ctx, req.(*Metric))
	}
	return interceptor(ctx, in, info, handler)
}

func _OwlProxyService_RegisterAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OwlProxyServiceServer).RegisterAgent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy_proto.OwlProxyService/RegisterAgent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OwlProxyServiceServer).RegisterAgent(ctx, req.(*AgentInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _OwlProxyService_ListAgentPlugins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OwlProxyServiceServer).ListAgentPlugins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy_proto.OwlProxyService/ListAgentPlugins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OwlProxyServiceServer).ListAgentPlugins(ctx, req.(*HostIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OwlProxyService_DownloadPluginFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadPluginReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OwlProxyServiceServer).DownloadPluginFile(m, &owlProxyServiceDownloadPluginFileServer{stream})
}

type OwlProxyService_DownloadPluginFileServer interface {
	Send(*PluginFile) error
	grpc.ServerStream
}

type owlProxyServiceDownloadPluginFileServer struct {
	grpc.ServerStream
}

func (x *owlProxyServiceDownloadPluginFileServer) Send(m *PluginFile) error {
	return x.ServerStream.SendMsg(m)
}

var _OwlProxyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proxy_proto.OwlProxyService",
	HandlerType: (*OwlProxyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReceiveTimeSeriesData",
			Handler:    _OwlProxyService_ReceiveTimeSeriesData_Handler,
		},
		{
			MethodName: "ReceiveAgentHeartbeat",
			Handler:    _OwlProxyService_ReceiveAgentHeartbeat_Handler,
		},
		{
			MethodName: "ReceiveAgentMetric",
			Handler:    _OwlProxyService_ReceiveAgentMetric_Handler,
		},
		{
			MethodName: "RegisterAgent",
			Handler:    _OwlProxyService_RegisterAgent_Handler,
		},
		{
			MethodName: "ListAgentPlugins",
			Handler:    _OwlProxyService_ListAgentPlugins_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadPluginFile",
			Handler:       _OwlProxyService_DownloadPluginFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/proxy.proto",
}