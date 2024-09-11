package ammeter

type AmmeterConfig struct {
	Id        int `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId int `json:"ammeter_id" xorm:"int(11) not null comment('设备ID')"`
	// 过载保护
	OverloadCd      int `json:"overload_cd" xorm:"int(11) not null default 0 comment('过载保护冷却时间')"`
	OverloadAction  int `json:"overload_action" xorm:"int(11) not null default 0 comment('过载保护动作阈值')"`
	OverloadWarning int `json:"overload_warning" xorm:"int(11) not null default 0 comment('过载保护报警阈值')"`
	OverloadTrip    int `json:"overload_trip" xorm:"int(11) not null default 0 comment('过载保护脱扣时间')"`
	OverloadDeal    int `json:"overload_deal" xorm:"int(11) not null default 0 comment('过载保护处理方式')"`
	// 欠压保护
	LackVoltageAction  int `json:"lack_voltage_action" xorm:"int(11) not null default 0 comment('欠压保护动作阈值')"`
	LackVoltageWarning int `json:"lack_voltage_warning" xorm:"int(11) not null default 0 comment('欠压保护报警阈值')"`
	LackVoltageTrip    int `json:"lack_voltage_trip" xorm:"int(11) not null default 0 comment('欠压保护脱扣时间')"`
	LackVoltageDeal    int `json:"lack_voltage_deal" xorm:"int(11) not null default 0 comment('欠压保护处理方式')"`
	// 过压保护
	OverVoltageAction  int `json:"over_voltage_action" xorm:"int(11) not null default 0 comment('过压保护动作阈值')"`
	OverVoltageWarning int `json:"over_voltage_warning" xorm:"int(11) not null default 0 comment('过压保护报警阈值')"`
	OverVoltageTrip    int `json:"over_voltage_trip" xorm:"int(11) not null default 0 comment('过压保护脱扣时间')"`
	OverVoltageDeal    int `json:"over_voltage_deal" xorm:"int(11) not null default 0 comment('过压保护处理方式')"`
	// 漏电保护
	LeakageAction  int `json:"leakage_action" xorm:"int(11) not null default 0 comment('漏电保护动作阈值')"`
	LeakageWarning int `json:"leakage_warning" xorm:"int(11) not null default 0 comment('漏电保护报警阈值')"`
	LeakageType    int `json:"leakage_type" xorm:"int(11) not null default 0 comment('漏电保护漏电类型')"`
	LeakageDeal    int `json:"leakage_deal" xorm:"int(11) not null default 0 comment('漏电保护处理方式')"`
	// 温度保护
	TemperatureAction  int `json:"temperature_action" xorm:"int(11) not null default 0 comment('温度保护动作阈值')"`
	TemperatureWarning int `json:"temperature_warning" xorm:"int(11) not null default 0 comment('温度保护报警阈值')"`
	TemperatureTrip    int `json:"temperature_trip" xorm:"int(11) not null default 0 comment('温度保护脱扣时间')"`
	TemperatureDeal    int `json:"temperature_deal" xorm:"int(11) not null default 0 comment('温度保护处理方式')"`
	// 瞬时保护
	InstantAction int `json:"instant_action" xorm:"int(11) not null default 0 comment('瞬时保护动作阈值')"`
	InstantDeal   int `json:"instant_deal" xorm:"int(11) not null default 0 comment('瞬时保护处理方式')"`
	// 缺相保护
	LackPhaseAction  int `json:"lack_phase_action" xorm:"int(11) not null default 0 comment('缺相保护动作阈值')"`
	LackPhaseWarning int `json:"lack_phase_warning" xorm:"int(11) not null default 0 comment('缺相保护报警阈值')"`
	LackPhaseTrip    int `json:"lack_phase_trip" xorm:"int(11) not null default 0 comment('缺相保护脱扣时间')"`
	LackPhaseDeal    int `json:"lack_phase_deal" xorm:"int(11) not null default 0 comment('缺相保护处理方式')"`
	// 功率限定
	PowerAction  int `json:"power_action" xorm:"int(11) not null default 0 comment('功率限定动作阈值')"`
	PowerWarning int `json:"power_warning" xorm:"int(11) not null default 0 comment('功率限定报警阈值')"`
	PowerTrip    int `json:"power_trip" xorm:"int(11) not null default 0 comment('功率限定脱扣时间')"`
	PowerDeal    int `json:"power_deal" xorm:"int(11) not null default 0 comment('功率限定处理方式')"`
	// 打火保护
	StrikeCycleNumber  int `json:"strike_cycle_number" xorm:"int(11) not null default 0 comment('打火保护通断次数')"`
	StrikeCycleWarning int `json:"strike_cycle_warning" xorm:"int(11) not null default 0 comment('打火保护报警次数')"`
	StrikeCycleTime    int `json:"strike_cycle_time" xorm:"int(11) not null default 0 comment('打火保护检测时间')"`
	StrikeCycleDeal    int `json:"strike_cycle_deal" xorm:"int(11) not null default 0 comment('打火保护处理方式')"`
	// 三相不平衡保护
	ThreeUnbalanceAction  int `json:"three_unbalance_action" xorm:"int(11) not null default 0 comment('三相不平衡保护动作阈值')"`
	ThreeUnbalanceWarning int `json:"three_unbalance_warning" xorm:"int(11) not null default 0 comment('三相不平衡保护报警阈值')"`
	ThreeUnbalanceTrip    int `json:"three_unbalance_trip" xorm:"int(11) not null default 0 comment('三相不平衡保护脱扣时间')"`
	ThreeUnbalanceDeal    int `json:"three_unbalance_deal" xorm:"int(11) not null default 0 comment('三相不平衡保护处理方式')"`
	// 恶性负载
	MalignantAction  int `json:"malignant_action" xorm:"int(11) not null default 0 comment('恶性负载动作阈值')"`
	MalignantWarning int `json:"malignant_warning" xorm:"int(11) not null default 0 comment('恶性负载报警阈值')"`
	MalignantTime    int `json:"malignant_time" xorm:"int(11) not null default 0 comment('恶性负载检测时间')"`
	MalignantDeal    int `json:"malignant_deal" xorm:"int(11) not null default 0 comment('恶性负载处理方式')"`
	// 防孤岛保护
	AntiIsolatedIslandAction  int `json:"anti_isolated_island_action" xorm:"int(11) not null default 0 comment('防孤岛保护动作阈值')"`
	AntiIsolatedIslandWarning int `json:"anti_isolated_island_warning" xorm:"int(11) not null default 0 comment('防孤岛保护报警阈值')"`
	AntiIsolatedIslandTime    int `json:"anti_isolated_island_time" xorm:"int(11) not null default 0 comment('防孤岛保护检测时间')"`
	AntiIsolatedIslandDeal    int `json:"anti_isolated_island_deal" xorm:"int(11) not null default 0 comment('防孤岛保护处理方式')"`
}

type ConfigUpdateReq struct {
	AmmeterId   int    `json:"ammeter_id"`
	UpdateType  string `json:"update_type"`
	UpdateValue int    `json:"update_value"`
}

type ConfigUpdateMsgData struct {
	Num        int    `json:"num"`
	ParamType  string `json:"param_type"`
	ParamValue string `json:"param_value"`
}
