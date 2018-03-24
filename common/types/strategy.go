package types

const (
	PRIORITY_HIGH_LEVEL = iota + 1
	PRIORITY_MIDDLE_LEVEL
	PRIORITY_LOW_LEVEL
)

type Strategy struct {
	ID          int    `json:"id"`
	ProductID   int    `json:"product_id" db:"product_id"`
	Name        string `json:"name"`
	Priority    int    `json:"priority"`
	AlarmCount  int    `json:"alarm_count" db:"alarm_count"`
	Cycle       int    `json:"cycle" `
	Expression  string `json:"expression"`
	Description string `json:"description"`
	UserID      int    `json:"user_id" db:"user_id"`
	Enable      bool   `json:"enable"`
}
