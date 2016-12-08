package types

const (
	PRIORITY_HIGH_LEVEL = iota + 1
	PRIORITY_MIDDLE_LEVEL
	PRIORITY_LOW_LEVEL
)

var STRATEGY_PRIORITY_MAPPING = map[int]string{1: "严重", 2: "较严重", 3: "注意"}

type Strategy struct {
	ID          int    `json:"id"`
	Name        string `form:"name" json:"name" binding:"required"`
	Priority    int    `form:"priority" json:"priority" binding:"required"`
	Pid         int    `form:"pid" json:"pid"`
	AlarmCount  int    `form:"alarm_count" json:"alarm_count" binding:"required"`
	Cycle       int    `form:"cycle" json:"cycle" binding:"required"`
	Expression  string `form:"expression" json:"expression" binding:"required"`
	Description string `form:"description" json:"description"`
	UserID      int    `json:"user_id"`
	Enable      bool   `json:"enable"`
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
