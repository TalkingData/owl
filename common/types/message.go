package types

import (
	"encoding/json"

	"github.com/wuyingsong/tcp"
)

type MessageType byte

//消息类型定义
const (
	MESS_POST_HOST_CONFIG MessageType = iota
	MESS_POST_METRIC
	MESS_POST_TSD
	MESS_GET_HOST_PLUGIN_LIST
	MESS_GET_HOST_PLUGIN_LIST_RESP
	MESS_GET_ALL_PLUGIN_MD5
	MESS_GET_ALL_PLUGIN_MD5_RESP
	MESS_PULL_PLUGIN
	MESS_POST_HOST_ALIVE
)

const (
	// Agent
	MsgAgentRegister tcp.PacketType = iota
	MsgAgentSendMetricInfo
	MsgAgentSendTimeSeriesData
	MsgAgentGetPluginsList
	MsgAgentRequestSyncPlugins
	MsgAgentSendHeartbeat

	// CFC
	MsgCFCSendPluginsList
	MsgCFCSendPlugin
	MsgCFCSendReconnect

	// Repeater
	MsgRepeaterPostTimeSeriesData
)

var MsgTextMap map[tcp.PacketType]string = map[tcp.PacketType]string{
	MsgAgentRegister:              "MsgAgentRegister",
	MsgAgentSendMetricInfo:        "MsgAgentSendMetricInfo",
	MsgAgentSendTimeSeriesData:    "MsgAgentSendTimeSeriesData",
	MsgAgentGetPluginsList:        "MsgAgentGetPluginList",
	MsgAgentRequestSyncPlugins:    "MsgAgentRequestSyncPlugins",
	MsgAgentSendHeartbeat:         "MsgAgentSendHeartbeat",
	MsgCFCSendPluginsList:         "MsgCFCSendPluginsList",
	MsgCFCSendPlugin:              "MsgCFCSendPlugin",
	MsgCFCSendReconnect:           "MsgCFCSendReconnect",
	MsgRepeaterPostTimeSeriesData: "MsgRepeaterPostTimeSeriesData",
}

type AgentPostMetricRequest struct {
	HostID  string           `json:"host_id"`
	Metrics []TimeSeriesData `json:"metrics"` // cpu.idle/ip=10.10.32.10,cpu=1
}

func (this *AgentPostMetricRequest) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

type GetPluginResp struct {
	HostID  string   `json:"host_id"` // 当agent通过proxy连接，需要通过hostid来查找映射表
	Plugins []Plugin `json:"plugins"`
}

func (this *GetPluginResp) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

func (this *GetPluginResp) Decode(data []byte) error {
	return json.Unmarshal(data, this)
}

type SyncPluginResponse struct {
	HostID string
	Path   string
	Body   []byte
}

func (sp *SyncPluginResponse) Encode() []byte {
	data, _ := json.Marshal(sp)
	return data
}

func (sp *SyncPluginResponse) Decode(data []byte) error {
	return json.Unmarshal(data, sp)
}

type SyncPluginRequest struct {
	HostID string
	Plugin
}

func (spr *SyncPluginRequest) Encode() []byte {
	data, _ := json.Marshal(spr)
	return data
}

func (spr *SyncPluginRequest) Decode(data []byte) error {
	return json.Unmarshal(data, spr)
}

type MetricConfig struct {
	HostID     string         `json:"host_id"`
	SeriesData TimeSeriesData `json:"time_series_data"`
}

func (tsd *MetricConfig) Encode() []byte {
	data, _ := json.Marshal(tsd)
	return data
}

func (tsd *MetricConfig) Decode(data []byte) error {
	return json.Unmarshal(data, &tsd)
}
