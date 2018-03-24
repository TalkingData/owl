package types

const (
	USER = iota
	ADMIN
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     int    `json:"role"`
	Phone    string `json:"phone"`
	Mail     string `json:"mail"`
	Wechat   string `json:"wechat"`
	Status   int    `json:"status"`
}

func (this *User) IsAdmin() bool {
	return this.Role == ADMIN
}
