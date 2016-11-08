package types

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

//消息类型可读映射
var MessageTypeText map[MessageType]string = map[MessageType]string{
	MESS_POST_HOST_CONFIG:          "post host config",
	MESS_POST_METRIC:               "post metric",
	MESS_POST_TSD:                  "post time series data",
	MESS_GET_HOST_PLUGIN_LIST:      "get host plugin list",
	MESS_GET_HOST_PLUGIN_LIST_RESP: "get host plugin list response",
	MESS_GET_ALL_PLUGIN_MD5:        "get all plugin md5 list",
	MESS_GET_ALL_PLUGIN_MD5_RESP:   "get all plugin md5 list response",
	MESS_PULL_PLUGIN:               "pull plugin file",
	MESS_POST_HOST_ALIVE:           "post host alive",
}

//消息接口
type Message interface {
	Encode() []byte
}

func Pack(t MessageType, m Message) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, t)
	binary.Write(&buf, binary.BigEndian, m.Encode())
	return buf.Bytes()
}

type PostMetric struct {
	HostID  string           `json:"host_id"`
	Metrics []TimeSeriesData `json:"metrics"` // cpu.idle/ip=10.10.32.10,cpu=1
}

func (this *PostMetric) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

type GetPluginResp struct {
	HostID  string   `json:"host_id"`
	Plugins []Plugin `json:"plugins"`
}

func (this *GetPluginResp) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

func (this *GetPluginResp) Decode(data []byte) error {
	return json.Unmarshal(data, this)
}
