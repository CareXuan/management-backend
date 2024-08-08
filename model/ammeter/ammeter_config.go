package ammeter

type AmmeterConfig struct {
	Id                 int    `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId          int    `json:"ammeter_id" xorm:"int(11) not null comment('设备ID')"`
	OverloadAction     int    `json:"overload_action" xorm:"int(11) not null default 0 comment('过载保护动作值')"`
	OverloadWarning    int    `json:"overload_warning" xorm:"int(11) not null default 0 comment('过载保护报警值')"`
	LeakageAction      int    `json:"leakage_action" xorm:"int(11) not null default 0 comment('漏电保护动作值')"`
	LeakageWarning     int    `json:"leakage_warning" xorm:"int(11) not null default 0 comment('漏电保护报警值')"`
	TemperatureAction  int    `json:"temperature_action" xorm:"int(11) not null default 0 comment('温度保护动作值')"`
	TemperatureWarning int    `json:"temperature_warning" xorm:"int(11) not null default 0 comment('温度保护报警值')"`
	ElectricCurrent    int    `json:"electric_current" xorm:"int(11) not null default 0 comment('额定电流')"`
	LeakageCurrent     int    `json:"leakage_current" xorm:"int(11) not null default 0 comment('额定漏电电流')"`
	TimingCloseSwitch  int    `json:"timing_close_switch" xorm:"int(11) not null default 0 comment('定时合闸开关 1：开 2：关')"`
	TimingCloseTime    int    `json:"timing_close_time" xorm:"int(11) not null default 0 comment('定时合闸时间')"`
	TimingCloseTimeStr string `json:"timing_close_time_str" xorm:"-"`
	TimingOpenSwitch   int    `json:"timing_open_switch" xorm:"int(11) not null default 0 comment('定时开闸开关 1：开 2：关')"`
	TimingOpenTime     int    `json:"timing_open_time" xorm:"int(11) not null default 0 comment('定时开闸时间')"`
	TimingOpenTimeStr  string `json:"timing_open_time_str" xorm:"-"`
}
