package model

type HostGroupHost struct {
	Id          uint32 `json:"id"`
	HostGroupId uint32 `json:"host_group_id"`
	HostId      string `json:"host_id"`
}

func (*HostGroupHost) TableName() string {
	return "host_group_host"
}
