package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"owl/common/types"
)

func eventsList(c *gin.Context) {
	var events []types.StrategyEvent
	var total int
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	page_size, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	priority, err := strconv.Atoi(c.DefaultQuery("priority", "0"))
	status, err := strconv.Atoi(c.DefaultQuery("status", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	key := c.DefaultQuery("key", "")
	start := c.DefaultQuery("start", time.Now().AddDate(0, 0, -7).Format("2006-01-02 00:00:00"))
	end := c.DefaultQuery("end", time.Now().Format("2006-01-02 15:04:05"))

	offset := (page - 1) * page_size
	sort := "status asc,update_time desc,priority asc"
	where := fmt.Sprintf("`update_time` BETWEEN '%s' AND '%s'", start, end)

	if priority != 0 {
		where += fmt.Sprintf(" AND `priority` = %d", priority)
	}
	if key != "" {
		key = strings.TrimSpace(key)
		where += fmt.Sprintf(" AND (`strategy_name` LIKE '%%%s%%' OR"+
			"`host_cname` LIKE '%%%s%%' OR"+
			"`ip` LIKE '%%%s%%' OR"+
			"`process_user` LIKE '%%%s%%')", key, key, key, key)
	}
	if status != 0 {
		where += fmt.Sprintf(" AND `status` = %d", status)
	}

	mydb.Table("strategy_event").Where(where).Count(&total)
	mydb.Where(where).Order(sort).Offset(offset).Limit(page_size).Find(&events)

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "events": &events, "total": total})
}

type EventsCount struct {
	Total    int       `json:"total"`
	High     int       `json:"high"`
	Middle   int       `json:"middle"`
	Low      int       `json:"low"`
	Unclosed int       `json:"unclosed"`
	Active   int       `json:"active"`
	LastTime time.Time `json:"last_time"`
}

func (e EventsCount) MarshalJSON() ([]byte, error) {
	type Alias EventsCount
	return json.Marshal(&struct {
		LastTime string `json:"last_time"`
		Alias
	}{
		LastTime: e.LastTime.Format("2006-01-02 15:04:05"),
		Alias:    (Alias)(e),
	})
}

func eventsCount(c *gin.Context) {
	var events []types.StrategyEvent
	var events_count EventsCount
	start := c.DefaultQuery("start", time.Now().AddDate(0, 0, -7).Format("2006-01-02 00:00:00"))
	end := c.DefaultQuery("end", time.Now().Format("2006-01-02 15:04:05"))
	where := fmt.Sprintf("`update_time` BETWEEN '%s' AND '%s'", start, end)
	mydb.Where(where).Order("`update_time` desc").Find(&events)

	if len(events) != 0 {
		events_count.LastTime = events[0].UpdateTime
		events_count.Total = len(events)
		for _, event := range events {
			switch event.Priority {
			case types.PRIORITY_HIGH_LEVEL:
				events_count.High += 1
			case types.PRIORITY_MIDDLE_LEVEL:
				events_count.Middle += 1
			case types.PRIORITY_LOW_LEVEL:
				events_count.Low += 1
			}

			if event.Status == types.EVENT_NEW {
				events_count.Active += 1
			}
			if event.Status != types.EVENT_CLOSED {
				events_count.Unclosed += 1
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "count": events_count})
}

func eventsStatus(c *gin.Context) {
	var events []types.StrategyEvent
	start := c.DefaultQuery("start", time.Now().AddDate(0, 0, -7).Format("2006-01-02 00:00:00"))
	end := c.DefaultQuery("end", time.Now().Format("2006-01-02 15:04:05"))
	where := fmt.Sprintf("`update_time` BETWEEN '%s' AND '%s'", start, end)
	mydb.Where(where).Find(&events)

	var active, aware, closed int
	if len(events) != 0 {
		for _, event := range events {
			switch event.Status {
			case types.EVENT_NEW:
				active += 1
			case types.EVENT_AWARE:
				aware += 1
			case types.EVENT_CLOSED:
				closed += 1
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "active": active, "awared": aware, "closed": closed})
}

func eventInform(c *gin.Context) {
	user := GetUser(c)
	event_id := c.Param("id")
	var event = types.StrategyEvent{}
	if len(event_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "event_id should be applied"})
		return
	}

	if err := mydb.Where("id = ?", event_id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	if event.Status == types.EVENT_NEW {
		event.Status = types.EVENT_AWARE
		event.ProcessUser = user.Username
		mydb.Save(&event)
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "status": event.Status})
}

type Process struct {
	ProcessUser     string    `json:"process_user"`
	ProcessComments string    `form:"process_comments" json:"process_comments" binding:"required"`
	ProcessTime     time.Time `json:"process_time"`
}

func (p Process) MarshalJSON() ([]byte, error) {
	type Alias Process
	return json.Marshal(&struct {
		ProcessTime string `json:"process_time"`
		Alias
	}{
		ProcessTime: p.ProcessTime.Format("2006-01-02 15:04:05"),
		Alias:       (Alias)(p),
	})
}

func eventClose(c *gin.Context) {
	user := GetUser(c)
	event_id := c.Param("id")
	var event = types.StrategyEvent{}
	var process Process

	if err := c.BindJSON(&process); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	if len(event_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "event_id should be applied"})
		return
	}

	if err := mydb.Table("strategy_event").Where("id = ?", event_id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	if event.Status == types.EVENT_AWARE {
		event.Status = types.EVENT_CLOSED
		event.ProcessUser = user.Username
		event.ProcessComments = process.ProcessComments
		event.ProcessTime = time.Now()
		mydb.Save(&event)
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "status": event.Status})
}

func processInfo(c *gin.Context) {
	event_id := c.Param("id")
	var process Process
	var event types.StrategyEvent

	if len(event_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "event_id should be applied"})
		return
	}

	if err := mydb.Table("strategy_event").Where("id = ?", event_id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	process.ProcessUser = event.ProcessUser
	process.ProcessComments = event.ProcessComments
	process.ProcessTime = event.ProcessTime

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "process": &process})
}

type EventDetail struct {
	StrategyEvent types.StrategyEvent  `json:"strategy_event"`
	TriggerEvents []types.TriggerEvent `json:"trigger_events"`
	ActionResults []types.ActionResult `json:"action_results"`
}

func eventDetail(c *gin.Context) {
	event_id := c.Param("id")
	var detail EventDetail
	var strategy_event types.StrategyEvent
	var trigger_events []types.TriggerEvent
	var action_results []types.ActionResult
	if len(event_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "event_id should be applied"})
		return
	}

	if err := mydb.Table("strategy_event").Select("strategy_event.*").Where("id = ?", event_id).First(&strategy_event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	mydb.Table("trigger_event").
		Where("strategy_event_id = ?", event_id).
		Order("`index` asc").
		Find(&trigger_events)

	mydb.Table("action_result").
		Where("strategy_event_id = ?", event_id).
		Find(&action_results)

	detail.StrategyEvent = strategy_event
	detail.TriggerEvents = trigger_events
	detail.ActionResults = action_results

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "detail": detail})
}

type EventPanel struct {
	Days        string `json:"days"`
	HighCount   int    `json:"high_count"`
	MiddleCount int    `json:"middle_count"`
	LowCount    int    `json:"low_count"`
}

func eventsPanel(c *gin.Context) {
	var panels []EventPanel
	start := c.DefaultQuery("start", time.Now().AddDate(0, 0, -7).Format("2006-01-02 00:00:00"))
	end := c.DefaultQuery("end", time.Now().Format("2006-01-02 15:04:05"))
	where := fmt.Sprintf("`create_time` BETWEEN '%s' AND '%s'", start, end)
	mydb.Table("strategy_event").Select("DATE_FORMAT(create_time,'%Y-%m-%d') as days, sum(case when priority=1 then 1 else 0 end) high_count, sum(case when priority=2 then 1 else 0 end) middle_count,  sum(case when priority=3 then 1 else 0 end) low_count").Where(where).Group("days").Find(&panels)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "panel": panels})
}
