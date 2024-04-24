package model

type Permission struct {
	Id        int           `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string        `json:"name" xorm:"varchar(64) not null"`
	Path      string        `json:"path" xorm:"varchar(256) not null"`
	Icon      string        `json:"icon" xorm:"varchar(512) not null"`
	Sort      int           `json:"sort" xorm:"int(11) not null default 100"`
	Children  []*Permission `json:"-" xorm:"-"`
	ParentId  int           `json:"parent_id" xorm:"int(11) not null default 0"`
	CreatedAt int           `json:"created_at" xorm:"int(11) not null default 0"`
	UpdatedAt int           `json:"updated_at" xorm:"int(11) not null default 0"`
	DeletedAt int           `json:"deleted_at" xorm:"int(11) not null default 0"`
}
