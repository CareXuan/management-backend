package model

type AddAppointmentReq struct {
	MemberId  int    `json:"member_id"`
	DeviceId  int    `json:"device_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Times     int    `json:"times"`
	Remark    string `json:"remark"`
	Pic       string `json:"pic"`
}

type VerifyAppointmentReq struct {
	AppointmentId int `json:"appointment_id"`
	Status        int `json:"status"`
}
