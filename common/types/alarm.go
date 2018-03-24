package types

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

type AlarmMessageType byte

//报警服务消息类型定义
const (
	ALAR_MESS_INSPECTOR_HEARTBEAT AlarmMessageType = iota + 1
	ALAR_MESS_INSPECTOR_TASK_REQUEST
	ALAR_MESS_INSPECTOR_TASKS
	ALAR_MESS_INSPECTOR_RESULTS
)

//报警服务消息类型可读映射
var AlarmMessageTypeText map[AlarmMessageType]string = map[AlarmMessageType]string{
	ALAR_MESS_INSPECTOR_HEARTBEAT:    "inspector heartbeat",
	ALAR_MESS_INSPECTOR_TASK_REQUEST: "inspector task request",
	ALAR_MESS_INSPECTOR_TASKS:        "inspector tasks",
	ALAR_MESS_INSPECTOR_RESULTS:      "inspector results",
}

//报警服务消息接口
type AlarmMessage interface {
	Encode() []byte
}

func AlarmPack(t AlarmMessageType, m AlarmMessage) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, t)
	if m != nil {
		binary.Write(&buf, binary.BigEndian, m.Encode())
	}
	binary.Write(&buf, binary.BigEndian, make([]byte, 0))
	return buf.Bytes()
}

type AlarmTask struct {
	ID       string
	Host     *Host `json:"-"`
	Strategy *Strategy
	Triggers map[string]*Trigger
}

func NewAlarmTask(host *Host, strategy *Strategy, triggers map[string]*Trigger) *AlarmTask {
	id := fmt.Sprintf("%v@%v", strategy.ID, host.ID)
	return &AlarmTask{id, host, strategy, triggers}
}

type AlarmTasks struct {
	Tasks []*AlarmTask
}

func (this *AlarmTasks) Encode() []byte {
	data, err := json.Marshal(this)
	if err != nil {
		fmt.Println(err.Error())
	}
	return data
}

func (this *AlarmTasks) Decode(data []byte) error {
	return json.Unmarshal(data, this)
}

type Node struct {
	IP       string    `json:"ip"`
	Hostname string    `json:"hostname"`
	Update   time.Time `json:"update"`
}

func (this Node) MarshalJSON() ([]byte, error) {
	type Alias Node
	return json.Marshal(&struct {
		Alias
		Update string `json:"update"`
	}{
		Alias:  (Alias)(this),
		Update: this.Update.Format("2006-01-02 15:04:05"),
	})
}

func (this *Node) Encode() []byte {
	data, err := json.Marshal(this)
	if err != nil {
		fmt.Println(err.Error())
	}
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
	data, err := json.Marshal(this)
	if err != nil {
		fmt.Println(err.Error())
	}
	return data
}

func (this *HeartBeat) Decode(data []byte) error {
	return json.Unmarshal(data, &this)
}

type AlarmResults struct {
	Results []*StrategyResult
}

type StrategyResult struct {
	TaskID            string
	Priority          int
	TriggerResultSets map[string]*TriggerResultSet
	ErrorMessage      string
	Triggered         bool
	CreateTime        time.Time
}

func (this *AlarmResults) Encode() []byte {
	data, err := json.Marshal(this)
	if err != nil {
		fmt.Println(err.Error())
	}
	return data
}

func (this *AlarmResults) Decode(data []byte) error {
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
	merged_tags := make([]string, 0)
	for tagk, tagv := range tags {
		if tagk == "host" || tagk == "uuid" {
			continue
		}
		merged_tags = append(merged_tags, tagk+"="+tagv)
	}
	sort.Strings(merged_tags)
	return &TriggerResult{index, strings.Join(merged_tags, ","), strings.Join(aggregate_tags, ","), current_threshold, triggered}
}
