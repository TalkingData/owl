package types

import (
	"fmt"
	"strconv"
	"time"
)

var DEFAULT_TIME, _ = time.Parse(time.RFC3339, "1980-01-01T00:00:00+00:00")

const (
	EVENT_NEW = iota + 1
	EVENT_AWARE
	EVENT_CLOSED
	EVENT_NODATA
	EVENT_UNKNOW
)

type StrategyEvent struct {
	ID           int64     `json:"id" db:"id"`
	ProductID    int       `json:"product_id" db:"product_id"`
	StrategyID   int       `json:"strategy_id" db:"strategy_id"`
	StrategyName string    `json:"strategy_name" db:"strategy_name"`
	Priority     int       `json:"priority"`
	Cycle        int       `json:"cycle"`
	AlarmCount   int       `json:"alarm_count" db:"alarm_count"`
	Expression   string    `json:"expression"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
	UpdateTime   time.Time `json:"update_time" db:"update_time"`
	AwareEndTime time.Time `json:"aware_end_time" db:"aware_end_time"`
	Count        int       `json:"count"`
	Status       int       `json:"status"`
	HostID       string    `json:"host_id" db:"host_id"`
	HostName     string    `json:"host_name" db:"host_name"`
	IP           string    `json:"ip"`
	ErrorMessage string    `json:"error_message"`
}

func NewStrategyEvent(
	product_id int,
	strategy_id int,
	strategy_name string,
	priority int,
	cycle int,
	alarm_count int,
	expression string,
	create_time time.Time,
	host_id string,
	host_name string,
	ip string,
	error_message string) *StrategyEvent {
	return &StrategyEvent{
		ProductID:    product_id,
		StrategyID:   strategy_id,
		StrategyName: strategy_name,
		Priority:     priority,
		Cycle:        cycle,
		AlarmCount:   alarm_count,
		Expression:   expression,
		CreateTime:   create_time,
		UpdateTime:   create_time,
		AwareEndTime: DEFAULT_TIME,
		Count:        1,
		Status:       1,
		HostID:       host_id,
		HostName:     host_name,
		IP:           ip,
		ErrorMessage: error_message}
}

type TriggerEvent struct {
	StrategyEventID  int64   `json:"strategy_event_id" db:"strategy_event_id"`
	Index            string  `json:"index" db:"index"`
	Metric           string  `json:"metric" db:"metric"`
	Tags             string  `json:"tags"`
	Number           int     `json:"number"`
	AggregateTags    string  `json:"aggregate_tags" db:"aggregate_tags"`
	CurrentThreshold float64 `json:"current_threshold" db:"current_threshold"`
	Method           string  `json:"method"`
	Symbol           string  `json:"symbol"`
	Threshold        float64 `json:"threshold"`
	Triggered        bool    `json:"triggered"`
	TriggerChanged   bool    `json:"-"`
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
	return fmt.Sprintf("\n%v: %v %v %v %v %v %v 当前值:%v 结果:%v", this.Index, this.Metric, this.Tags, this.Method, number, this.Symbol, bytes2Human(this.Threshold), bytes2Human(this.CurrentThreshold), this.Triggered)
}

func bytes2Human(num float64) string {
	if num < 1000.00 {
		return fmt.Sprintf("%3.2f", num)
	}
	suffix := []string{"", "K", "M", "G", "T", "P", "E", "Z"}
	for _, unit := range suffix {
		if num < 1000.00 {
			return fmt.Sprintf("%3.2f%s%s", num, unit, "B")
		}
		num /= 1000.00
	}
	return fmt.Sprintf("%.2f%s%s", num, "Y", "B")
}
