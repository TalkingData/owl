package types

import (
	"time"
)

type Chart struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	UserID     int       `json:"-"`
	Size       int       `json:"size"`
	ReferCount int       `json:"refer_count"`
	Thumbnail  string    `json:"thumbnail"`
	CreateAt   time.Time `json:"create_at"`
	UpdateAt   time.Time `json:"update_at"`

	Elements []*ChartElement `json:"elements"`
}

type ChartElement struct {
	ID      int    `json:"id"`
	ChartID int    `json:"-"`
	Name    string `json:"name"`
	Metric  string `json:"metric"`
	Tags    string `json:"tags"`
}

func (Chart) TableName() string {
	return "chart"
}

func (ChartElement) TableName() string {
	return "chart_element"
}
