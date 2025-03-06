package conf

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/go-xorm/xorm"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	cron2 "prize-draw/cron"
	"prize-draw/model"
	"time"
	"xorm.io/core"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Mysql  MysqlConfig  `yaml:"mysql"`
	Gpt    GptConfig    `yaml:"gpt"`
	Upload UploadConfig `yaml:"upload"`
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

type WechatConfig struct {
	AppId     string        `yaml:"app_id"`
	AppSecret string        `yaml:"app_secret"`
	Wechat    wechatWarning `yaml:"wechat"`
}

type wechatWarning struct {
	Common string `yaml:"common"`
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
	initCronTask()
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

func initCronTask() {
	c := cron.New()
	_, err := c.AddFunc("02 00 */1 * *", cron2.AllotTask)
	if err != nil {
		log.Fatal("添加定时任务失败")
	}
	c.Start()
}

func syncTables() {
	err := Mysql.Sync2(
		new(model.User),
		new(model.Role),
		new(model.UserRole),
		new(model.Permission),
		new(model.RolePermission),
		new(model.Gift),
		new(model.GiftGroup),
		new(model.GiftGroupGift),
		new(model.GiftExtract),
		new(model.GiftPackage),
		new(model.Config),
		new(model.Achievement),
		new(model.AchievementTask),
		new(model.AchievementGift),
		new(model.Task),
		new(model.TaskGift),
		new(model.TaskDo),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}
