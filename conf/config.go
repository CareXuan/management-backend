package conf

import (
	"encoding/json"
	"fmt"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"management-backend/common"
	"management-backend/model/ammeter"
	"management-backend/model/rbac"
	"management-backend/utils"
	"strconv"
	"time"
	"xorm.io/core"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Mysql   MysqlConfig   `yaml:"mysql"`
	Rabbit  RabbitConfig  `yaml:"rabbit"`
	Redis   RedisConfig   `yaml:"redis"`
	Vehicle VehicleConfig `yaml:"vehicle_config"`
	Tcp     TcpConfig     `yaml:"tcp"`
	Ammeter AmmeterConfig `yaml:"ammeter"`
	Admin   int           `yaml:"admin"`
	Gpt     GptConfig     `yaml:"gpt"`
	Upload  UploadConfig  `yaml:"upload"`
	Wechat  WechatConfig  `yaml:"wechat"`
}

type MysqlConfig struct {
	M string `yaml:"m"`
	S string `yaml:"s"`
}

type RabbitConfig struct {
	Host string `yaml:"host"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
}

type VehicleConfig struct {
	Host string `yaml:"host"`
}

type TcpConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AmmeterConfig struct {
	Supervisor int `yaml:"supervisor"`
	Manager    int `yaml:"manager"`
}

type GptConfig struct {
	Key string `yaml:"key"`
}

type UploadConfig struct {
	Url string `yaml:"url"`
}

type WechatConfig struct {
	AppId     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

var (
	Mysql     *xorm.EngineGroup
	Rabbit    *amqp.Connection
	Redis     redis.Conn
	Conf      *Config
	WechatApp *officialAccount.OfficialAccount
)

func NewConfig(configPath string) {

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = yaml.Unmarshal(content, &Conf)
	if err != nil {
		log.Fatal(err)
		return
	}
	connectMysql()
	initWechatApp()
	//connectRabbitMQ()
	//connectRedis()
	//GetVehicleConfig(Conf.Vehicle.Host)
}

func connectMysql() {

	var err error
	Mysql, err = xorm.NewEngineGroup("mysql", []string{Conf.Mysql.M, Conf.Mysql.S})
	if err != nil {
		log.Fatal(err)
	}
	// 连接池参数:空闲数、最大连接数、连接最大生存时间
	Mysql.SetMaxIdleConns(10)
	Mysql.SetMaxOpenConns(100)
	Mysql.SetConnMaxLifetime(3 * time.Hour)

	Mysql.SetLogLevel(core.LOG_DEBUG)
	Mysql.ShowSQL(true)
	Mysql.ShowExecTime(true)

	if err := Mysql.Ping(); err != nil {
		log.Fatal(err)

	}

	// 同步表
	syncTables()
}

func initWechatApp() {
	OfficialAccountApp, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:  Conf.Wechat.AppId,     // 公众号、小程序的appid
		Secret: Conf.Wechat.AppSecret, //

		Log: officialAccount.Log{
			Level:  "debug",
			File:   "./wechat.log",
			Stdout: false, //  是否打印在终端
		},

		HttpDebug: true,
		Debug:     false,
	})
	if err != nil {
		log.Fatal(err)
	}
	WechatApp = OfficialAccountApp
}

func syncTables() {
	err := Mysql.Sync2(
		new(rbac.User),
		new(rbac.Role),
		new(rbac.UserRole),
		new(rbac.Permission),
		new(rbac.RolePermission),
		new(ammeter.Ammeter),
		new(ammeter.AmmeterConfig),
		new(ammeter.AmmeterData),
		new(ammeter.AmmeterManage),
		new(ammeter.AmmeterManageConfig),
		new(ammeter.AmmeterManageLog),
		new(ammeter.AmmeterWarning),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func connectRabbitMQ() {
	var err error
	// 连接到RabbitMQ服务器
	Rabbit, err = amqp.Dial(Conf.Rabbit.Host)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
}

func connectRedis() {
	var err error
	Redis, err = redis.Dial("tcp", Conf.Redis.Host)
	if err != nil {
		log.Fatal(err)
		return
	}
}
func GetVehicleConfig(configUrl string) {
	var reqParams = make(map[string]string)
	reqParams["current"] = "0"
	reqParams["page_size"] = "9999"

	res, err := common.DoGet(configUrl, reqParams)
	if err != nil {
		log.Fatal(err)
		return
	}
	configData := res.Body.(map[string]any)["list"]
	for _, v := range configData.([]interface{}) {
		data := v.(map[string]interface{})
		deviceId, _ := strconv.Atoi(data["deviceId"].(string))
		deviceConfigJson, _ := json.Marshal(data)
		_, err := Redis.Do("SET", fmt.Sprintf(utils.REDIS_KEY_VEHICLE_CONFIG, deviceId), string(deviceConfigJson))
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
