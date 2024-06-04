package rbac

type UserRole struct {
	Id        int `json:"id" xorm:"pk autoincr INT(11)"`
	UserId    int `json:"user_id" xorm:"INT(11) not null"`
	RoleId    int `json:"role_id" xorm:"INT(11) not null"`
	CreatedAt int `json:"-" xorm:"INT(11) not null default 0"`
	UpdatedAt int `json:"-" xorm:"INT(11) not null default 0"`
	DeletedAt int `json:"-" xorm:"INT(11) not null default 0"`
}
