package ammeter

type AirswitchData struct {
	BoardId              int    `json:"board_id"`
	AmmeterNum           int    `json:"ammeter_num"`
	Hxqd                 int    `json:"hxqd"`
	BreakerProtect       int    `json:"breaker_protect"`
	BreakerWarning1      int    `json:"breaker_warning1"`
	BreakerWarning2      int    `json:"breaker_warning2"`
	BreakerPreTrip       int    `json:"breaker_pre_trip"`
	BreakerTrip          int    `json:"breaker_trip"`
	AActivePower         int    `json:"a_active_power"`
	AVoltage             int    `json:"a_voltage"`
	AElectric            int    `json:"a_electric"`
	CombinedActiveEnergy int    `json:"combined_active_energy"`
	Leakage              int    `json:"leakage"`
	Temperature          int    `json:"temperature"`
	RunningState         int    `json:"running_state"`
	SwitchState          int    `json:"switch_state"`
	ManualAutoState      int    `json:"manual_auto_state"`
	Ts                   string `json:"ts"`
}
