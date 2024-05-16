package conf

import (
	"encoding/json"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"management-backend/common"
	"management-backend/model"
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

var (
	Mysql  *xorm.EngineGroup
	Rabbit *amqp.Connection
	Redis  redis.Conn
	Conf   *Config
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
	connectRabbitMQ()
	connectRedis()
	GetVehicleConfig(Conf.Vehicle.Host)
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

func syncTables() {
	err := Mysql.Sync2(
		new(model.User),
		new(model.Role),
		new(model.UserRole),
		new(model.Permission),
		new(model.RolePermission),
		new(model.Ammeter),
		new(model.AmmeterConfig),
		new(model.AmmeterData),
		new(model.AmmeterManage),
		new(model.AmmeterWarning),
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
