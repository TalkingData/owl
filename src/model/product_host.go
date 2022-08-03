package model

type ProductHost struct {
	Id        uint   `json:"id"`
	ProductId uint   `json:"product_id"`
	HostId    string `json:"host_id"`
}

func (*ProductHost) TableName() string {
	return "product_host"
}
