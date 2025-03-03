package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 定义了应用程序配置的结构体。
// 它包含了应用程序和数据库相关的配置信息。
type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
	}
}

// AppConfig 保存了应用程序的全局配置实例。
var AppConfig *Config

// InitConfig 初始化应用程序的配置。
// 该函数负责读取配置文件并将其解析到 AppConfig 全局变量中。
// 同时，它还会初始化数据库和 Redis 的连接。
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error reading config file", err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalln("Error unmarshaling config file", err)
	}

	initDB()
	initRedis()
}
