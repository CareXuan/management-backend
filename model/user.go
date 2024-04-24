package model

type User struct {
	Id        int    `json:"-"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	Token     string `json:"token"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}
