package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"owl/common/utils"
)

const (
	EVENT_NEW = iota + 1
	EVENT_AWARED
	EVENT_CLOSED
)

var STRATEGY_STATUS_MAPPING = map[int]string{1: "活跃报警", 2: "已知悉报警", 3: "已关闭报警"}

type StrategyEvent struct {
	ID              int64     `json:"id"`
	StrategyID      int       `json:"strategy_id"`
	StrategyName    string    `json:"strategy_name"`
	StrategyType    int       `json:"strategy_type"`
	Priority        int       `json:"priority"`
	Cycle           int       `json:"cycle"`
	AlarmCount      int       `json:"alarm_count"`
	Expression      string    `json:"expression"`
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
	Count           int       `json:"count"`
	Status          int       `json:"status"`
	HostID          string    `json:"host_id"`
	HostCname       string    `json:"host_cname"`
	HostName        string    `json:"host_name"`
	IP              string    `json:"ip"`
	SN              string    `json:"sn"`
	ProcessUser     string    `json:"process_user"`
	ProcessComments string    `json:"process_comments"`
	ProcessTime     time.Time `json:"process_time"`
}

func (s StrategyEvent) MarshalJSON() ([]byte, error) {
	type Alias StrategyEvent
	return json.Marshal(&struct {
		Alias
		ProcessTime string `json:"process_time"`
		CreateTime  string `json:"create_time"`
		UpdateTime  string `json:"update_time"`
	}{
		Alias:       (Alias)(s),
		ProcessTime: s.ProcessTime.Format("2006-01-02 15:04:05"),
		CreateTime:  s.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:  s.UpdateTime.Format("2006-01-02 15:04:05"),
	})
}

func (StrategyEvent) TableName() string {
	return "strategy_event"
}

func NewStrategyEvent(strategy_id int,
	strategy_name string,
	strategy_type int,
	priority int,
	cycle int,
	alarm_count int,
	expression string,
	create_time time.Time,
	host_id string,
	host_cname string,
	host_name string,
	ip string,
	sn string) *StrategyEvent {
	return &StrategyEvent{
		StrategyID:      strategy_id,
		StrategyName:    strategy_name,
		StrategyType:    strategy_type,
		Priority:        priority,
		Cycle:           cycle,
		AlarmCount:      alarm_count,
		Expression:      expression,
		CreateTime:      create_time,
		UpdateTime:      create_time,
		Count:           1,
		Status:          1,
		HostID:          host_id,
		HostCname:       host_cname,
		HostName:        host_name,
		IP:              ip,
		SN:              sn,
		ProcessUser:     "",
		ProcessComments: ""}
}

type TriggerEvent struct {
	StrategyEventID  int64   `json:"strategy_event_id"`
	Index            string  `json:"index"`
	Metric           string  `json:"metric"`
	Tags             string  `json:"tags"`
	Number           int     `json:"number"`
	AggregateTags    string  `json:"aggregate_tags"`
	CurrentThreshold float64 `json:"current_threshold"`
	Method           string  `json:"method"`
	Symbol           string  `json:"symbol"`
	Threshold        float64 `json:"threshold"`
	Triggered        bool    `json:"triggered"`
}

func NewTriggerEvent(strategy_event_id int64, index, metric, tags, aggregate_tags, symbol, method string, number int, threshold, current_threshold float64, triggered bool) *TriggerEvent {
	return &TriggerEvent{
		StrategyEventID:  strategy_event_id,
		Index:            index,
		Metric:           metric,
		Tags:             tags,
		Number:           number,
		AggregateTags:    aggregate_tags,
		CurrentThreshold: current_threshold,
		Method:           method,
		Symbol:           symbol,
		Threshold:        threshold,
		Triggered:        triggered}
}

func (this *TriggerEvent) String() string {
	number := ""
	if this.Number != 0 {
		number = strconv.Itoa(this.Number)
	}
	return fmt.Sprintf("\n%v: %v %v %v %v %v %v %v", this.Index, this.Metric, this.Tags, this.Method, number, utils.Bytes2Human(this.CurrentThreshold), this.Symbol, utils.Bytes2Human(this.Threshold))
}

type TemplateStrategy struct {
	ID                int64
	NAME              string
	TYPE              string
	CYCLE             int
	PRIORITY          string
	STATUS            string
	ALARM_COUNT       int
	COUNT             int
	UPDATE_TIME       string
	EXPRESSION        string
	EXPRESSION_DETAIL string
}

type TemplateHost struct {
	CNAME  string
	NAME   string
	IP     string
	STATUS string
	SN     string
}

type Template struct {
	STRATEGY TemplateStrategy
	HOST     TemplateHost
}
