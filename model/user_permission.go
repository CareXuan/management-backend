package model

type UserPermission struct {
	Id           int `json:"id" xorm:"pk autoincr INT(11)"`
	UserId       int `json:"user_id" xorm:"INT(11) not null"`
	PermissionId int `json:"permission_id" xorm:"INT(11) not null"`
	CreatedAt    int `json:"created_at" xorm:"INT(11) not null default 0"`
	UpdatedAt    int `json:"updated_at" xorm:"INT(11) not null default 0"`
	DeletedAt    int `json:"deleted_at" xorm:"INT(11) not null default 0"`
}
