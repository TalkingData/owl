package model

import "owl/common/utils"

const (
	HostStatusDown     = "0"
	HostStatusOk       = "1"
	HostStatusDisabled = "2"
	HostStatusNew      = "3"
)

var (
	HostStatusCodenameMap = map[string]string{
		HostStatusDown:     "Down",
		HostStatusOk:       "Ok",
		HostStatusDisabled: "Disabled",
		HostStatusNew:      "New",
	}
)

type Host struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"` // 经确认，这个字段在v5版本中没有用到，永远赋空字符串
	Ip           string           `json:"ip"`
	Hostname     string           `json:"hostname"`
	AgentVersion string           `json:"agent_version"`
	Status       string           `json:"status" gorm:"default:3"`
	CreateAt     *utils.LocalTime `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt     *utils.LocalTime `json:"update_at" gorm:"autoUpdateTime"`
	MuteTime     *utils.LocalTime `json:"mute_time"`
	Uptime       float64          `json:"uptime"`
	IdlePct      float64          `json:"idle_pct"`
}

func (h *Host) GetStatusCodename() string {
	return HostStatusCodenameMap[h.Status]
}

func (*Host) TableName() string {
	return "host"
}
