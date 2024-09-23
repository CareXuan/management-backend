package utils

const AMMETER_SUPERVISOR = 1
const AMMETER_MANAGER = 2

const AMMETER_DATA_TYPE_LEAKAGE = 1
const AMMETER_DATA_TYPE_CURRENT = 2
const AMMETER_DATA_TYPE_VOLTAGE = 3
const AMMETER_DATA_TYPE_TEMPERATURE = 4
const AMMETER_DATA_TYPE_POWER = 5
const AMMETER_DATA_TYPE_CONSUMPTION = 6

const AMMETER_STATUS_ONLINE = 1
const AMMETER_STATUS_OFFLINE = 2

const AMMETER_STATUS_SWITCH_OPEN = 1
const AMMETER_STATUS_SWITCH_CLOSE = 2

var UPDATE_AMMETER_CONFIG_PARAMS = map[string]string{
	"overload_cd":                  "0220",
	"overload_action":              "0221",
	"overload_warning":             "0222",
	"overload_trip":                "0223",
	"overload_deal":                "0224",
	"lack_voltage_action":          "0225",
	"lack_voltage_warning":         "0226",
	"lack_voltage_trip":            "0227",
	"lack_voltage_deal":            "0228",
	"over_voltage_action":          "0229",
	"over_voltage_warning":         "022A",
	"over_voltage_trip":            "022B",
	"over_voltage_deal":            "022C",
	"leakage_action":               "022D",
	"leakage_warning":              "022E",
	"leakage_type":                 "022F",
	"leakage_deal":                 "0230",
	"temperature_action":           "0231",
	"temperature_warning":          "0232",
	"temperature_trip":             "0233",
	"temperature_deal":             "0234",
	"instant_action":               "0235",
	"instant_deal":                 "0236",
	"lack_phase_action":            "0237",
	"lack_phase_warning":           "0238",
	"lack_phase_trip":              "0239",
	"lack_phase_deal":              "023A",
	"power_action":                 "023B",
	"power_warning":                "023C",
	"power_trip":                   "023D",
	"power_deal":                   "023E",
	"strike_cycle_number":          "023F",
	"strike_cycle_warning":         "0240",
	"strike_cycle_time":            "0241",
	"strike_cycle_deal":            "0242",
	"three_unbalance_action":       "0243",
	"three_unbalance_warning":      "0244",
	"three_unbalance_trip":         "0245",
	"three_unbalance_deal":         "0246",
	"malignant_action":             "0247",
	"malignant_warning":            "0248",
	"malignant_time":               "0249",
	"malignant_deal":               "024A",
	"anti_isolated_island_action":  "024B",
	"anti_isolated_island_warning": "024C",
	"anti_isolated_island_time":    "024D",
	"anti_isolated_island_deal":    "024E",
}
