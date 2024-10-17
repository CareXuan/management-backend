package ammeter

type Board struct {
	BoardId         int    `json:"board_id"`
	Iccid           string `json:"iccid"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	Cmd             int    `json:"cmd"`
	Comment         string `json:"comment"`
	FactoryId       int    `json:"factory_id"`
	Pwd1            int    `json:"pwd1"`
	Pwd1Base        int    `json:"pwd1_base"`
	Pwd2            int    `json:"pwd2"`
	Pwd2Base        int    `json:"pwd2_base"`
	Pwd3            int    `json:"pwd3"`
	Pwd3Base        int    `json:"pwd3_base"`
	Pwd4            int    `json:"pwd4"`
	Pwd4Base        int    `json:"pwd4_base"`
	Ts              string `json:"ts"`
	Ip              string `json:"ip"`
	Ip2             string `json:"ip2"`
	TimeData        int    `json:"time_data"`
	RetryTimes      int    `json:"retry_times"`
	Hxqd            int    `json:"hxqd"`
	SignalIntensity int    `json:"signal_intensity"`
}
