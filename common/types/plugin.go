package types

import (
	"encoding/json"
	"fmt"
)

type Plugin struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Args     string `json:"args"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

func (this *Plugin) Encode() []byte {
	data, _ := json.Marshal(this)
	return data
}

func (this *Plugin) Decode(data []byte) error {
	return json.Unmarshal(data, &this)
}

func (this Plugin) String() string {
	return fmt.Sprintf("{id:%d, name:%s, args:%s, interval:%d, timeout:%d}",
		this.ID,
		this.Name,
		this.Args,
		this.Interval,
		this.Timeout,
	)
}

func (Plugin) TableName() string {
	return "plugin"
}
