package model

type HostPlugin struct {
	Id       uint   `json:"id"`
	HostId   string `json:"host_id"`
	PluginId uint   `json:"plugin_id"`
	Args     string `json:"args"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
	Comment  string `json:"comment"`
}

func (*HostPlugin) TableName() string {
	return "host_plugin"
}
