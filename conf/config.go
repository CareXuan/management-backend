package conf

import (
	"github.com/go-xorm/xorm"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"my-gpt-server/model"
	"time"
	"xorm.io/core"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Mysql MysqlConfig `yaml:"mysql"`
	Gpt   GptConfig   `yaml:"gpt"`
}

type MysqlConfig struct {
	M string `yaml:"m"`
	S string `yaml:"s"`
}

type GptConfig struct {
	Key string `yaml:"key"`
}

var (
	Mysql *xorm.EngineGroup
	Conf  *Config
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
		new(model.GptQuestion),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}
