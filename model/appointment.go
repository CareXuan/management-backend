package model

type Appointment struct {
	Id        int `json:"id" xorm:"pk autoincr INT(11)"`
	DeviceId  int `json:"device_id" xorm:"INT(11) not null"`
	MemberId  int `json:"member_id" xorm:"INT(11) not null"`
	StartTime int `json:"start_time" xorm:"INT(11) not null"`
	EndTime   int `json:"end_time" xorm:"INT(11) not null"`
	Status    int `json:"status" xorm:"INT(3) not null"`
	CreatedAt int `json:"created_at" xorm:"INT(11) not null default 0"`
}
