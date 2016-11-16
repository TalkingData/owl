package types

const (
	SEND_MAIL = iota + 1
	SEND_SMS
	SEND_WECHAT
)

const (
	ACTION_ALARM = iota + 1
	ACTION_RESTORE
	ACTION_CUSTOM
)

type Action struct {
	ID              int    `json:"-"`
	StrategyID      int    `form:"strategy_id" json:"strategy_id" binding:"required"`
	Type            int    `form:"type" json:"type" binding:"required"`
	FilePath        string `form:"file_path" json:"file_path"`
	AlarmSubject    string `form:"alarm_subject" json:"alarm_subject"`
	RestoreSubject  string `form:"restore_subject" json:"restore_subject"`
	AlarmTemplate   string `form:"alarm_template" json:"alarm_template"`
	RestoreTemplate string `form:"restore_template" json:"restore_template"`
	SendType        int    `form:"send_type" json:"send_type" binding:"require"`
}

func (Action) TableName() string {
	return "action"
}

type ActionResult struct {
	StrategyEventID int64  `json:"strategy_event_id"`
	ActionID        int    `json:"action_id"`
	ActionType      int    `json:"action_type"`
	ActionSendType  int    `json:"action_send_type"`
	UserID          int    `json:"user_id"`
	Username        string `json:"username"`
	Phone           string `json:"phone"`
	Mail            string `json:"mail"`
	Weixin          string `json:"weixin"`
	Subject         string `json:"subject"`
	Content         string `json:"content"`
	Success         bool   `json:"success"`
	Response        string `json:"response"`
}

type ActionUser struct {
	ID       int `json:"-"`
	ActionID int `json:"-"`
	UserID   int `json:"user_id"`
}

func (ActionUser) TableName() string {
	return "action_user"
}

type ActionUserGroup struct {
	ID          int `json:"-"`
	ActionID    int `json:"-"`
	UserGroupID int `form:"user_group_id" json:"user_group_id"`
}

func (ActionUserGroup) TableName() string {
	return "action_user_group"
}
