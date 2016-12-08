package types

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type AlarmMessageType byte

//报警服务消息类型定义
const (
	ALAR_MESS_INSPECTOR_HEARTBEAT AlarmMessageType = iota + 1
	ALAR_MESS_INSPECTOR_TASK_REQUEST
	ALAR_MESS_INSPECTOR_TASKS
	ALAR_MESS_INSPECTOR_RESULT
)

//报警服务消息类型可读映射
var AlarmMessageTypeText map[AlarmMessageType]string = map[AlarmMessageType]string{
	ALAR_MESS_INSPECTOR_HEARTBEAT:    "inspector heartbeat",
	ALAR_MESS_INSPECTOR_TASK_REQUEST: "inspector task request",
	ALAR_MESS_INSPECTOR_TASKS:        "inspector tasks",
	ALAR_MESS_INSPECTOR_RESULT:       "inspector result",
}

//报警服务消息接口
type AlarmMessage interface {
	Encode() []byte
}

func AlarmPack(t AlarmMessageType, m AlarmMessage) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, t)
	binary.Write(&buf, binary.BigEndian, m.Encode())
	return buf.Bytes()
}

type AlarmTask struct {
	ID       string
	Host     *Host
	Strategy *Strategy
	Triggers map[string]*Trigger
}

func NewAlarmTask(host *Host, strategy *Strategy, triggers map[string]*Trigger) *AlarmTask {
	id := fmt.Sprintf("%v@%v", strategy.ID, host.ID)
	return &AlarmTask{id, host, strategy, triggers}
}

type GetTasksResp struct {
	AlarmTasks []*AlarmTask
}

func (this *GetTasksResp) Encode() []byte {
	data, err := json.Marshal(this)
	if err != nil {
		fmt.Println(err.Error())
	}
	return data
}

func (this *GetTasksResp) Decode(data []byte) error {
	return json.Unmarshal(data, this)
}

type NodePool struct {
	Nodes map[string]*Node
	Lock  *sync.Mutex
}

func NewNodePool() *NodePool {
	return &NodePool{make(map[string]*Node), &sync.Mutex{}}
}

type Node struct {
	IP       string
	Hostname string
	Update   time.Time
}

func (this *Node) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

type HeartBeat struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

func NewHeartBeat(ip, hostname string) *HeartBeat {
	return &HeartBeat{ip, hostname}
}

func (this *HeartBeat) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

func (this *HeartBeat) Decode(data []byte) error {
	return json.Unmarshal(data, &this)
}

type StrategyResult struct {
	TaskID            string
	Priority          int
	TriggerResultSets map[string]*TriggerResultSet
	ErrorMessage      string
	Triggered         bool
	CreateTime        time.Time
}

func (this *StrategyResult) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

func (this *StrategyResult) Decode(data []byte) error {
	return json.Unmarshal(data, &this)
}

type TriggerResultSet struct {
	TriggerResults []*TriggerResult
	Triggered      bool
}

type TriggerResult struct {
	Index            string
	Tags             string
	AggregateTags    string
	CurrentThreshold float64
	Triggered        bool
}

func NewStrategyResult(task_id string, priority int, trigger_results map[string]*TriggerResultSet, error_message string, triggered bool, create_time time.Time) *StrategyResult {
	return &StrategyResult{task_id, priority, trigger_results, error_message, triggered, create_time}
}

func NewTriggerResult(index string, tags map[string]string, aggregate_tags []string, current_threshold float64, triggered bool) *TriggerResult {
	tags_string := ""
	aggregate_tags_string := ""
	tags_limiter := ""
	aggregate_tags_limiter := ""

	if _, ok := tags["hostname"]; ok {
		delete(tags, "hostname")
	}

	if _, ok := tags["uuid"]; ok {
		delete(tags, "uuid")
	}

	if len(tags) > 1 {
		tags_limiter = ","
	}

	if len(aggregate_tags) > 1 {
		aggregate_tags_limiter = ","
	}

	for tagk, tagv := range tags {
		tags_string += fmt.Sprintf("%s=%s%s", tagk, tagv, tags_limiter)
	}

	for _, aggregate_tag := range aggregate_tags {
		aggregate_tags_string += fmt.Sprintf("%s%s", aggregate_tag, aggregate_tags_limiter)
	}

	if strings.HasSuffix(tags_string, ",") {
		tags_string = strings.TrimSuffix(tags_string, ",")
	}

	if strings.HasSuffix(aggregate_tags_string, ",") {
		aggregate_tags_string = strings.TrimSuffix(aggregate_tags_string, ",")
	}

	return &TriggerResult{index, tags_string, aggregate_tags_string, current_threshold, triggered}
}
