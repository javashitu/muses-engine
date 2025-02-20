package config

import (
	"embed"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

var Server ServerConfig

type ServerConfig struct {
	HttpConfig ApplicationConfig `yaml:"application"`
	LogConfig  LogConfig         `yaml:"log"`
}

type ApplicationConfig struct {
	AppModel         string `yaml:"appMode"`
	HttpPort         string `yaml:"httpPort"`
	JwtKey           string `yaml:"jwtSecurityKey"`
	JwtExpireSeconds int    `yaml:"jwtExpireSeconds"`
}

type LogConfig struct {
	LogDir        string `yaml:"logDir"`
	ServerLogName string `yaml:"serverLogName"`
	MaxSize       int    `yaml:"maxSize"`
	MaxBackups    int    `yaml:"maxBackups"`
	MaxAge        int    `yaml:"maxAge"`
	Compress      bool   `yaml:"compress"`
}

//go:embed serverConfig.yaml
var svf embed.FS

func init() {
	file, err := svf.ReadFile("serverConfig.yaml")
	if err != nil {
		log.Panicf("read server config failure, maybe file not exits -> %s ", err)
	}
	err = yaml.Unmarshal(file, &Server)
	if err != nil {
		log.Panicf("bind server config failure, maybe file not yaml format -> %s ", err)
	}

	initLog()

}

func initLog() {
	os.MkdirAll(Server.LogConfig.LogDir, 0755)
	logFileName := Server.LogConfig.LogDir + Server.LogConfig.ServerLogName

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panicf("open log file failed,err -> %s ", err)
	}

	//lumberjack日志滚动记录器
	//lumberjack是一个日志滚动记录器。可以把日志文件根据大小、日期等分割。一般情况下，
	//lumberjack配合其他日志库，实现日志的滚动(rolling)记录。

	log.SetOutput(
		&lumberjack.Logger{
			Filename:   logFileName,                 //日志文件得位置
			MaxSize:    Server.LogConfig.MaxSize,    //切割之前日志文件得大小（单位：MB）
			MaxBackups: Server.LogConfig.MaxBackups, //保留旧文件得最大个数
			MaxAge:     Server.LogConfig.MaxAge,     //保留旧文件得最大天数
			Compress:   true,                        //是否压缩旧文件
		})
	//设置标准记录器的输出标志 ：完整文件名和行号   微秒分辨率    本地时区中的日期
	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)
	//同时写文件和屏幕
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	gin.DefaultWriter = multiWriter
}
