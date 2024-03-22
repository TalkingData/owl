package model

type ProductHost struct {
	Id        uint32 `json:"id"`
	ProductId uint32 `json:"product_id"`
	HostId    string `json:"host_id"`
}

func (*ProductHost) TableName() string {
	return "product_host"
}
