package model

type AmmeterData struct {
	Id         int `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId  int `json:"ammeterId" xorm:"INT(11) not null"`
	Type       int `json:"type" xorm:"INT(11) not null"`
	Value      int `json:"value" xorm:"INT(11) not null"`
	CreateTime int `json:"-" xorm:"int(11) not null default 0"`
}

type AmmeterStatisticRes struct {
	Data                            []StatisticForm `json:"data"`
	TodayElectricityConsumption     int             `json:"today_electricity_consumption"`
	YesterdayElectricityConsumption int             `json:"yesterday_electricity_consumption"`
	MonthElectricityConsumption     int             `json:"month_electricity_consumption"`
	LastMonthElectricityConsumption int             `json:"last_month_electricity_consumption"`
	YearElectricityConsumption      int             `json:"year_electricity_consumption"`
	LastYearElectricityConsumption  int             `json:"last_year_electricity_consumption"`
}

type StatisticForm struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}
