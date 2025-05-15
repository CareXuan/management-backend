package conf

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/go-xorm/xorm"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"management-backend/model"
	"time"
	"xorm.io/core"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Mysql  MysqlConfig  `yaml:"mysql"`
	Gpt    GptConfig    `yaml:"gpt"`
	Upload UploadConfig `yaml:"upload"`
	Tcp    TcpConfig    `yaml:"tcp"`
	Wechat WechatConfig `yaml:"wechat"`
}

type MysqlConfig struct {
	M string `yaml:"m"`
	S string `yaml:"s"`
}

type GptConfig struct {
	Key string `yaml:"key"`
}

type UploadConfig struct {
	Url string `yaml:"url"`
}

type TcpConfig struct {
	Host1 string `yaml:"host1"`
	Port1 string `yaml:"port1"`
	Host2 string `yaml:"host2"`
	Port2 string `yaml:"port2"`
	Host3 string `yaml:"host3"`
	Port3 string `yaml:"port3"`
}

type WechatConfig struct {
	AppId     string        `yaml:"app_id"`
	AppSecret string        `yaml:"app_secret"`
	Token     string        `yaml:"token"`
	Warning   wechatWarning `yaml:"warning"`
}

type wechatWarning struct {
	Common      string `yaml:"common"`
	TestWarning string `yaml:"test_warning"`
}

var (
	Mysql     *xorm.EngineGroup
	WechatApp *officialAccount.OfficialAccount
	Conf      *Config
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
			Stdout: false, // 是否打印在终端
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
		new(model.User),
		new(model.Role),
		new(model.UserRole),
		new(model.Permission),
		new(model.RolePermission),
		new(model.Organization),
		new(model.OrganizationUser),
		new(model.Device),
		new(model.DeviceCommonData),
		new(model.DeviceLocationHistory),
		new(model.DeviceServiceData),
		new(model.DeviceNewServiceData),
		new(model.DeviceChangeLog),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}
