package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"owl/common/types"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//StrategyEvent 报警事件结构体
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
}

//MarshalJSON 当json输出时转换时间格式
func (s StrategyEvent) MarshalJSON() ([]byte, error) {
	type Alias StrategyEvent
	return json.Marshal(&struct {
		Alias
		CreateTime   string `json:"create_time"`
		UpdateTime   string `json:"update_time"`
		AwareEndTime string `json:"aware_end_time"`
	}{
		Alias:        (Alias)(s),
		CreateTime:   s.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:   s.UpdateTime.Format("2006-01-02 15:04:05"),
		AwareEndTime: s.AwareEndTime.Format("2006-01-02 15:04:05"),
	})
}

//StrategyEventFailed 报警事件失败信息的结构体
type StrategyEventFailed struct {
	Status     int       `json:"status"`
	HostName   string    `json:"hostname" db:"hostname"`
	IP         string    `json:"ip"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
	Message    string    `json:"message"`
}

//MarshalJSON 当json输出时转换时间格式
func (s StrategyEventFailed) MarshalJSON() ([]byte, error) {
	type Alias StrategyEventFailed
	return json.Marshal(&struct {
		Alias
		UpdateTime string `json:"update_time"`
	}{
		Alias:      (Alias)(s),
		UpdateTime: s.UpdateTime.Format("2006-01-02 15:04:05"),
	})
}

//StrategyEventProcess 报警事件处理记录
type StrategyEventProcess struct {
	ProcessStatus   int       `json:"process_status" db:"process_status"`
	ProcessUser     string    `json:"process_user" db:"process_user"`
	ProcessComments string    `json:"process_comments" db:"process_comments"`
	ProcessTime     time.Time `json:"process_time" db:"process_time"`
}

//MarshalJSON 当json输出时转换时间格式
func (s StrategyEventProcess) MarshalJSON() ([]byte, error) {
	type Alias StrategyEventProcess
	return json.Marshal(&struct {
		Alias
		ProcessTime string `json:"process_time"`
	}{
		Alias:       (Alias)(s),
		ProcessTime: s.ProcessTime.Format("2006-01-02 15:04:05"),
	})
}

//ActionResult 报警执行动作结果
type ActionResult struct {
	StrategyEventID int64  `json:"strategy_event_id" db:"strategy_event_id"`
	Count           int    `json:"count" db:"count"`
	ActionID        int    `json:"action_id" db:"action_id"`
	ActionType      int    `json:"action_type" db:"action_type"`
	ActionKind      int    `json:"action_kind" db:"action_kind"`
	ScriptID        int    `json:"script_id" db:"script_id"`
	ScriptName      string `json:"script_name" db:"script_name"`
	FilePath        string `json:"file_path" db:"file_path"`
	UserID          int    `json:"user_id" db:"user_id"`
	Username        string `json:"username" db:"username"`
	Phone           string `json:"phone"`
	Mail            string `json:"mail"`
	Weixin          string `json:"weixin"`
	Subject         string `json:"subject"`
	Content         string `json:"content"`
	Success         bool   `json:"success"`
	Response        string `json:"response"`
}

//AlarmRecord 报警记录
type AlarmRecord struct {
	StrategyEvent *StrategyEventRecord  `json:"strategy_event"`
	TriggerEvents []*TriggerEventRecord `json:"trigger_events"`
	ActionResults []*ActionResult       `json:"action_results"`
}

//StrategyEventRecord 报警事件结构体
type StrategyEventRecord struct {
	StrategyEventID int64     `json:"strategy_event_id" db:"strategy_event_id"`
	Count           int       `json:"count"`
	StrategyID      int       `json:"strategy_id" db:"strategy_id"`
	StrategyName    string    `json:"strategy_name" db:"strategy_name"`
	Priority        int       `json:"priority"`
	Cycle           int       `json:"cycle"`
	AlarmCount      int       `json:"alarm_count" db:"alarm_count"`
	Expression      string    `json:"expression"`
	CreateTime      time.Time `json:"create_time" db:"create_time"`
	UpdateTime      time.Time `json:"update_time" db:"update_time"`
	AwareEndTime    time.Time `json:"aware_end_time" db:"aware_end_time"`
	Status          int       `json:"status"`
	HostID          string    `json:"host_id" db:"host_id"`
	HostName        string    `json:"host_name" db:"host_name"`
	IP              string    `json:"ip"`
}

//TriggerEventRecord 报警的表达式
type TriggerEventRecord struct {
	StrategyEventID  int64   `json:"strategy_event_id" db:"strategy_event_id"`
	Count            int     `json:"count" db:"count"`
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

//MarshalJSON 当json输出时转换时间格式
func (s StrategyEventRecord) MarshalJSON() ([]byte, error) {
	type Alias StrategyEventRecord
	return json.Marshal(&struct {
		Alias
		CreateTime   string `json:"create_time"`
		UpdateTime   string `json:"update_time"`
		AwareEndTime string `json:"aware_end_time"`
	}{
		Alias:        (Alias)(s),
		CreateTime:   s.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:   s.UpdateTime.Format("2006-01-02 15:04:05"),
		AwareEndTime: s.AwareEndTime.Format("2006-01-02 15:04:05"),
	})
}

func eventsList(c *gin.Context) {
	productID := c.GetInt("product_id")
	strategyID, _ := strconv.Atoi(c.DefaultQuery("strategy_id", "0"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	query := c.GetString("query")
	inputOrder := c.GetString("order")
	order := "status ASC, create_time DESC, priority ASC"
	if len(inputOrder) > 0 {
		order = inputOrder
	}
	where := fmt.Sprintf("product_id=%d", productID)
	if strategyID != 0 {
		where += fmt.Sprintf(" AND strategy_id = %d", strategyID)
	}
	if query != "" {
		query = strings.TrimSpace(query)
		where += fmt.Sprintf(" AND (`strategy_name` LIKE '%%%s%%' OR"+
			"`host_name` LIKE '%%%s%%' OR"+
			"`ip` LIKE '%%%s%%')", query, query, query)
	}
	if status != 0 {
		where += fmt.Sprintf(" AND `status` = %d", status)
	}
	limit := fmt.Sprintf("%d, %d", c.GetInt("offset"), c.GetInt("limit"))
	events := mydb.GetStrategyEvents(where, order, limit)
	total := mydb.GetStrategiesEventsCount(where)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "events": &events, "total": total})
}

func eventAware(c *gin.Context) {
	productID := c.GetInt("product_id")
	eventIDs := c.QueryArray("event_id")
	awareEndTime := c.Query("aware_end_time")
	if len(eventIDs) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "params value is invalid"})
		return
	}
	name := c.GetString("username")
	mydb.UpdateStrategyEventsStatus(eventIDs, awareEndTime, productID, types.EVENT_AWARE)
	mydb.CreateStrategyEventProcesses(eventIDs, types.EVENT_AWARE, name, "知悉")
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

func eventsFailed(c *gin.Context) {
	strategyID := c.DefaultQuery("strategy_id", "")
	status := c.DefaultQuery("status", "")
	query := c.GetString("query")
	order := "update_time DESC"
	if strategyID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "params value is invalid"})
		return
	}
	where := fmt.Sprintf("strategy_id = %s", strategyID)
	if status != "" {
		where += fmt.Sprintf(" AND sef.status = %s", status)
	}
	if query != "" {
		query = strings.TrimSpace(query)
		where += fmt.Sprintf(" AND (`hostname` LIKE '%%%s%%' OR `ip` LIKE '%%%s%%')", query, query)
	}
	limit := fmt.Sprintf("%d, %d", c.GetInt("offset"), c.GetInt("limit"))
	eventsFailed := mydb.GetStrategyEventsFailed(where, order, limit)
	total := mydb.GetStrategyEventsFailedCount(where)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "events_failed": &eventsFailed, "total": total})
}

func eventProcessRecord(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("event_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "params value is invalid"})
		return
	}
	records := mydb.GetStrategyEventProcessRecord(eventID)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "events_process": &records})
}

func eventDetail(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("event_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "params value is invalid"})
		return
	}
	order := "update_time DESC"
	limit := fmt.Sprintf("%d, %d", c.GetInt("offset"), c.GetInt("limit"))
	records, total := mydb.GetAlarmRecords(eventID, order, limit)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "records": &records, "total": total})
}
