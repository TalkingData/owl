package model

type ActionResult struct {
	StrategyEventId uint64 `json:"strategy_event_id"`
	Count           int    `json:"count"`
	ActionId        uint64 `json:"action_id"`
	ActionType      bool   `json:"action_type"`
	ActionKind      bool   `json:"action_kind"`
	ScriptId        uint   `json:"script_id"`
	UserId          uint   `json:"user_id"`
	Username        string `json:"username"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Wechat          string `json:"wechat"`
	Subject         string `json:"subject"`
	Content         string `json:"content"`
	Success         bool   `json:"success"`
	Response        string `json:"response"`
}

func (*ActionResult) TableName() string {
	return "action_result"
}
