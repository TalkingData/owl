package types

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (Group) TableName() string {
	return "group"
}
