package types

import (
	"encoding/json"
	"time"
)

type Host struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	IP           string    `json:"ip"`
	SN           string    `json:"sn"`
	Hostname     string    `json:"hostname"`
	AgentVersion string    `json:"agent_version"`
	Status       string    `json:"status"`
	CreateAt     time.Time `json:"-"`
	UpdateAt     time.Time `json:"-"`
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

func (Host) TableName() string {
	return "host"
}
