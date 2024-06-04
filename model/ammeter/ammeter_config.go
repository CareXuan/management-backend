package ammeter

type AmmeterConfig struct {
	Id                 int    `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId          int    `json:"ammeter_id" xorm:"int(11) not null"`
	OverloadAction     int    `json:"overload_action" xorm:"int(11) not null default 0"`
	OverloadWarning    int    `json:"overload_warning" xorm:"int(11) not null default 0"`
	LeakageAction      int    `json:"leakage_action" xorm:"int(11) not null default 0"`
	LeakageWarning     int    `json:"leakage_warning" xorm:"int(11) not null default 0"`
	TemperatureAction  int    `json:"temperature_action" xorm:"int(11) not null default 0"`
	TemperatureWarning int    `json:"temperature_warning" xorm:"int(11) not null default 0"`
	ElectricCurrent    int    `json:"electric_current" xorm:"int(11) not null default 0"`
	LeakageCurrent     int    `json:"leakage_current" xorm:"int(11) not null default 0"`
	TimingCloseSwitch  int    `json:"timing_close_switch" xorm:"int(11) not null default 0"`
	TimingCloseTime    int    `json:"timing_close_time" xorm:"int(11) not null default 0"`
	TimingCloseTimeStr string `json:"timing_close_time_str" xorm:"-"`
	TimingOpenSwitch   int    `json:"timing_open_switch" xorm:"int(11) not null default 0"`
	TimingOpenTime     int    `json:"timing_open_time" xorm:"int(11) not null default 0"`
	TimingOpenTimeStr  string `json:"timing_open_time_str" xorm:"-"`
}
