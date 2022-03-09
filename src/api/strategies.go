package main

import (
	"fmt"
	"net/http"
	"owl/common/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Strategy 报警策略结构体
type Strategy struct {
	ID          int64  `json:"id"`
	ProductID   int    `json:"product_id" db:"product_id"`
	Name        string `json:"name"`
	Priority    int    `json:"priority"`
	AlarmCount  int    `json:"alarm_count" db:"alarm_count"`
	Cycle       int    `json:"cycle"`
	Expression  string `json:"expression"`
	Description string `json:"description"`
	UserID      int    `json:"user_id" db:"user_id"`
	Enable      int    `json:"enable"`
}

// Trigger 报警策略逻辑表达式结构体
type Trigger struct {
	ID          int64   `form:"id" json:"-" `
	StrategyID  int     `form:"strategy_id" json:"strategy_id" db:"strategy_id"`
	Metric      string  `form:"metric" json:"metric"`
	Tags        string  `form:"tags" json:"tags"`
	Number      int     `form:"number" json:"number"`
	Index       string  `form:"index" json:"index" `
	Name        string  `form:"name" json:"name"`
	Method      string  `form:"method" json:"method" `
	Symbol      string  `form:"symbol" json:"symbol" `
	Threshold   float64 `form:"threshold" json:"threshold" `
	Description string  `form:"description" json:"description"`
}

// StrategyGroup 报警策略与主机组结构体
type StrategyGroup struct {
	ID         int64 `json:"-"`
	StrategyID int   `json:"-"`
	GroupID    int   `json:"group_id"`
}

//StrategySimple 简单的策略结构体
type StrategySimple struct {
	Strategy
	UserName string `json:"user_name" db:"user_name"`
}

//StrategySummary 查询策略的结构体
type StrategySummary struct {
	StrategySimple
	AlertCount  int `json:"alert_count" db:"alert_count"`
	NodataCount int `json:"nodata_count" db:"nodata_count"`
	UnknowCount int `json:"unknow_count" db:"unknow_count"`
}

//StrategyDetail 创建策略的结构体
type StrategyDetail struct {
	Strategy
	ExcludeHosts []*AlarmHost     `form:"exclude_hosts" json:"exclude_hosts"`
	Groups       []*StrategyGroup `form:"groups" json:"groups"`
	Triggers     []*Trigger       `form:"triggers" json:"triggers" binding:"required"`
	Actions      []*ActionDetail  `form:"actions" json:"actions" binding:"required"`
}

//StrategyInfo 单个策略详细信息的结构体
type StrategyInfo struct {
	Strategy
	ExcludeHosts []*AlarmHost   `json:"exclude_hosts"`
	Groups       []*types.Group `json:"groups"`
	Triggers     []*Trigger     `json:"triggers"`
	Actions      []*ActionInfo  `json:"actions"`
}

// Action 报警动作结构体
type Action struct {
	ID              int
	StrategyID      int    `form:"strategy_id" db:"strategy_id" json:"strategy_id" binding:"required"`
	Type            int    `form:"type" json:"type" binding:"required"`
	Kind            int    `form:"kind" json:"kind" binding:"required"`
	ScriptID        int    `form:"script_id" db:"script_id" json:"script_id"`
	AlarmSubject    string `form:"alarm_subject" json:"alarm_subject" db:"alarm_subject"`
	AlarmTemplate   string `form:"alarm_template" json:"alarm_template" db:"alarm_template"`
	RestoreSubject  string `form:"restore_subject" json:"restore_subject" db:"restore_subject"`
	RestoreTemplate string `form:"restore_template" json:"restore_template" db:"restore_template"`
	BeginTime       string `form:"begin_time" db:"begin_time" json:"begin_time" binding:"require"`
	EndTime         string `form:"end_time" db:"end_time" json:"end_time" binding:"require"`
	TimePeriod      int    `form:"time_period" db:"time_period" json:"time_period" binding:"require"`
}

// ActionUserGroup 报警执行动作的通知人员组
type ActionUserGroup struct {
	ID          int `json:"-"`
	ActionID    int `json:"-"`
	UserGroupID int `form:"user_group_id" json:"user_group_id"`
}

// ActionDetail 执行动作详细信息
type ActionDetail struct {
	Action
	UserGroups []ActionUserGroup `form:"user_groups" json:"user_groups" binding:"required"`
}

// ActionInfo 执行动作信息
type ActionInfo struct {
	Action
	Script     *Script            `json:"script"`
	UserGroups []*types.UserGroup `json:"user_groups"`
}

// AlarmHost 需要排除的报警主机
type AlarmHost struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

func strategyList(c *gin.Context) {
	productID := c.GetInt("product_id")
	query := c.GetString("query")
	my := c.DefaultQuery("my", "false")
	status := c.DefaultQuery("status", "all")
	where := fmt.Sprintf("s.product_id=%d", productID)
	order := c.GetString("order")
	if query != "" {
		where += fmt.Sprintf(" and s.name LIKE '%%%s%%'", query)
	}
	if my == "true" {
		if user := mydb.getUserProfile(c.GetString("username")); user != nil {
			where += fmt.Sprintf(" and s.user_id = %d", user.ID)
		}
	}
	if status == "ok" {
		where += " and alert_count = 0 and nodata_count = 0 and unknow_count = 0"
	} else if status == "problem" {
		where += " and alert_count != 0 or nodata_count != 0 or unknow_count != 0"
	}
	total := mydb.GetStrategiesCount(where)
	limit := fmt.Sprintf("%d, %d", c.GetInt("offset"), c.GetInt("limit"))
	strategies := mydb.GetStrategies(where, order, limit)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "strategies": strategies, "total": total})
}

func strategyHostGroup(c *gin.Context) {
	hostGroupID, err := strconv.Atoi(c.Param("host_group_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "message": "params is invalid"})
		return
	}
	query := c.GetString("query")
	where := fmt.Sprintf("sg.group_id = %d", hostGroupID)
	if query != "" {
		where += fmt.Sprintf(" AND (s.name LIKE '%%%s%%' OR u.username LIKE '%%%s%%')", query, query)

	}
	limit := fmt.Sprintf("%d, %d", c.GetInt("offset"), c.GetInt("limit"))
	strategies := mydb.GetStrategiesByHostGroupID(where, limit)
	total := mydb.GetStrategiesByHostGroupIDCount(where)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "strategies": strategies, "total": total})
}

func strategyCreate(c *gin.Context) {
	productID := c.GetInt("product_id")
	var strategy StrategyDetail
	if err := c.BindJSON(&strategy); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	user := mydb.getUserProfile(c.GetString("username"))
	if user == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "user not found"})
		return
	}
	strategy.UserID = user.ID
	strategy.ProductID = productID
	if err := mydb.CreateStrategy(&strategy); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

func strategyInfo(c *gin.Context) {
	productID := c.GetInt("product_id")
	strategyID, err := strconv.ParseInt(c.Param("strategy_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "params value is invalid"})
		return
	}
	strategy := mydb.GetStrategy(strategyID, productID)
	if strategy == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "strategy not found"})
		return
	}
	hostGroups := mydb.GetHostGroupsByStrategyID(strategyID)
	triggers := mydb.GetTriggersByStrategyID(strategyID)
	actions := mydb.GetActionsByStrategyID(strategyID)
	hosts := mydb.GetHostsExByStrategyID(strategyID)
	for _, action := range actions {
		action.Script = mydb.GetScriptByScriptID(action.ScriptID)
		action.UserGroups = mydb.GetUserGroupsByActionID(action.ID)
	}
	strategyInfo := &StrategyInfo{*strategy, hosts, hostGroups, triggers, actions}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "strategy": &strategyInfo})
}

func strategyUpdate(c *gin.Context) {
	productID := c.GetInt("product_id")
	var strategy StrategyDetail
	if err := c.BindJSON(&strategy); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if s := mydb.GetStrategy(strategy.ID, productID); s == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "strategy not found"})
		return
	}
	if err := mydb.UpdateStrategy(&strategy); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

func strategySwitch(c *gin.Context) {
	productID := c.GetInt("product_id")
	strategyIDs := c.QueryArray("strategy_id")
	enable := c.Query("enable")
	if len(strategyIDs) == 0 || enable == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "params not be applied"})
		return
	}
	if err := mydb.UpdateStrategiesStatus(strategyIDs, productID, enable); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}

func strategyDelete(c *gin.Context) {
	productID := c.GetInt("product_id")
	strategyIDs := c.QueryArray("strategy_id")
	if len(strategyIDs) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "params not be applied"})
		return
	}
	if err := mydb.DeleteStrategies(strategyIDs, productID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}
