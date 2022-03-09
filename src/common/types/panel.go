package types

import (
	"time"
)

type Panel struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Thumbnail string    `json:"thumbnail"`
	UserID    int       `json:"-"`
	Favor     int       `json:"-"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
	Charts    []*Chart  `json:"charts" gorm:"many2many:panel_chart"`
}

func (Panel) TableName() string {
	return "panel"
}
