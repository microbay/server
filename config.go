package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() API {
	var api API
	var plugins map[string]map[string]interface{}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed loading server.json ", err)
	}
	if err := viper.MarshalKey("api", &api); err != nil {
		log.Fatal("Malformed api section in server.json ", err)
	}
	if err := viper.MarshalKey("plugins", &plugins); err != nil {
		log.Fatal("Malformed plugins section in server.json ", err)
	}
	api.plugins = plugins
	return api
}

func init() {
	setConfigLocations()
	setupEnvVars()
	initLogLevel()
}

func setConfigLocations() {
	viper.SetConfigName("server")         // server.json file name
	viper.AddConfigPath("/etc/apigo/")    // package config location
	viper.AddConfigPath("server/config/") //	local dev config location
}

func setupEnvVars() {
	viper.SetEnvPrefix("mb")
	viper.BindEnv("env")
	viper.SetDefault("env", "development")
	viper.BindEnv("loglevel")
	viper.SetDefault("loglevel", "debug")
	viper.BindEnv("host")
	viper.SetDefault("host", "localhost:7777")
	viper.BindEnv("redis_host")
	viper.SetDefault("redis_host", "localhost:6379")
	viper.BindEnv("redis_password")
	viper.SetDefault("redis_password", "")
}

func initLogLevel() {
	logLevel, err := log.ParseLevel(viper.GetString("loglevel"))
	if err != nil {
		log.Warn("Unsupported log level ", viper.Get("loglevel"))
	}
	log.SetLevel(logLevel)
}
