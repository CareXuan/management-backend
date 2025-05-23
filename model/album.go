package model

type Album struct {
	Id        int          `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string       `json:"name" xorm:"VARCHAR(24) not null default '' comment('礼物组名称')"`
	HasCnt    int          `json:"has_cnt" xorm:"-"`
	AllCnt    int          `json:"all_cnt" xorm:"-"`
	CreateAt  int          `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt  int          `json:"delete_at" xorm:"INT(10) not null default 0"`
	AlbumGift []*AlbumGift `json:"album_gift" xorm:"-"`
}

type AlbumGift struct {
	Id       int `json:"id" xorm:"pk autoincr INT(11)"`
	AlbumId  int `json:"album_id" xorm:"INT(10) not null default 0 comment('相册ID')"`
	GiftId   int `json:"gift_id" xorm:"INT(10) not null default 0 comment('礼物ID')"`
	CreateAt int `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt int `json:"delete_at" xorm:"INT(10) not null default 0"`
}
