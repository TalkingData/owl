package types

const (
	STRATEGY_GLOBAL = iota + 1
	STRATEGY_GROUP
	STRATEGY_HOST
)

const (
	PRIORITY_HIGH_LEVEL = iota + 1
	PRIORITY_MIDDLE_LEVEL
	PRIORITY_LOW_LEVEL
)

var STRATEGY_TYPE_MAPPING = map[int]string{1: "全局策略", 2: "主机组策略", 3: "主机策略"}
var STRATEGY_PRIORITY_MAPPING = map[int]string{1: "严重", 2: "较严重", 3: "注意"}

type Strategy struct {
	ID          int    `json:"id"`
	Name        string `form:"name" json:"name" binding:"required"`
	Priority    int    `form:"priority" json:"priority" binding:"required"`
	Type        int    `form:"type" json:"type" binding:"required"`
	Pid         int    `form:"pid" json:"pid"`
	AlarmCount  int    `form:"alarm_count" json:"alarm_count" binding:"required"`
	Cycle       int    `form:"cycle" json:"cycle" binding:"required"`
	Expression  string `form:"expression" json:"expression" binding:"required"`
	GroupID     int    `form:"group_id" json:"group_id"`
	HostID      string `form:"host_id" json:"host_id"`
	Description string `form:"description" json:"description"`
	Enable      bool   `form:"enable" json:"enable"`
}

func (Strategy) TableName() string {
	return "strategy"
}

type StrategyHost struct {
	ID         int64  `json:"-"`
	StrategyID int    `json:"-"`
	HostID     string `json:"host_id"`
}

func (StrategyHost) TableName() string {
	return "strategy_host"
}

type StrategyGroup struct {
	ID         int64 `json:"-"`
	StrategyID int   `json:"-"`
	GroupID    int   `json:"group_id"`
}

func (StrategyGroup) TableName() string {
	return "strategy_group"
}
