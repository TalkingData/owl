package types

const (
	ACTION_ALARM = iota
	ACTION_RESTORE
)

const (
	ACTION_NOTIFY = iota + 1
	ACTION_RUN
)

type Action struct {
	ID              int    `json:"-"`
	StrategyID      int    `json:"strategy_id" db:"strategy_id"`
	Type            int    `json:"type"`
	Kind            int    `json:"kind"`
	ScriptID        int    `json:"script_id" db:"script_id"`
	FilePath        string `json:"file_path" db:"file_path"`
	AlarmSubject    string `json:"alarm_subject" db:"alarm_subject"`
	AlarmTemplate   string `json:"alarm_template" db:"alarm_template"`
	RestoreSubject  string `json:"restore_subject" db:"restore_subject"`
	RestoreTemplate string `json:"restore_template" db:"restore_template"`
	BeginTime       string `json:"begin_time" db:"begin_time"`
	EndTime         string `json:"end_time" db:"end_time"`
	TimePeriod      int    `json:"time_period" db:"time_period"`
}

type ActionResult struct {
	StrategyEventID int64  `json:"strategy_event_id"`
	Count           int    `json:"count"`
	ActionID        int    `json:"action_id"`
	ActionType      int    `json:"action_type"`
	ActionKind      int    `json:"action_kind"`
	ScriptID        int    `json:"script_id"`
	UserID          int    `json:"user_id"`
	Username        string `json:"username"`
	Phone           string `json:"phone"`
	Mail            string `json:"mail"`
	Weixin          string `json:"weixin"`
	Subject         string `json:"subject"`
	Content         string `json:"content"`
	Success         bool   `json:"success"`
	Response        string `json:"response"`
	FilePath        string `json:"file_path"`
}

func NewActionResult(
	strategy_event_id int64,
	count, action_id, action_type, action_kind, script_id, user_id int,
	user_name, phone, mail, weixin, subject, content, response string,
	success bool) *ActionResult {
	return &ActionResult{
		StrategyEventID: strategy_event_id,
		Count:           count,
		ActionID:        action_id,
		ActionType:      action_type,
		ActionKind:      action_kind,
		ScriptID:        script_id,
		UserID:          user_id,
		Username:        user_name,
		Phone:           phone,
		Mail:            mail,
		Weixin:          weixin,
		Subject:         subject,
		Content:         content,
		Success:         success,
		Response:        response,
	}
}

type Script struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	FilePath string `db:"file_path"`
}
