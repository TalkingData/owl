package types

import (
	"encoding/json"
	"time"
)

type Host struct {
	ID           string    `json:"id"`
	IP           string    `json:"ip"`
	Hostname     string    `json:"hostname"`
	AgentVersion string    `json:"agent_version" db:"agent_version"`
	Status       string    `json:"status"`
	CreateAt     time.Time `json:"create_at" db:"create_at"`
	UpdateAt     time.Time `json:"update_at" db:"update_at"`
	Uptime       float64   `json:"uptime" db:"uptime"`
	IdlePct      float64   `json:"idle_pct" db:"idle_pct"`
	MuteTime     time.Time `json:"-" db:"mute_time"`
}

func (this *Host) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

func (this *Host) Decode(data []byte) error {
	return json.Unmarshal(data, &this)
}

func (this *Host) IsAlive() bool {
	return this.Status == "1"
}

func (this *Host) IsMute() bool {
	return time.Now().Before(this.MuteTime)
}
