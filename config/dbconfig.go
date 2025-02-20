package config

import (
	"embed"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConfig DBConf
var DB *gorm.DB

type DBConf struct {
	URL string `yaml:"dburl"`
}

//go:embed dbConfig.yaml
var f2 embed.FS

func init() {
	file, err := f2.ReadFile("dbConfig.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &DBConfig)

	if err != nil {
		panic(err)
	}
	DBInit()
}

func DBInit() {
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		})
	db, err := gorm.Open(mysql.Open(DBConfig.URL), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panicf("go orm open mysql connection failure %s \n", err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Panicf("connect to mysql failure %s \n", err.Error())
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	DB = db
}
