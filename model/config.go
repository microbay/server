package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type API struct {
	Name      string     `json:"name"`
	Portal    string     `json:"portal"`
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Auth   int     `json:"auth"`
	Path   string  `json:"path"`
	Micros []Micro `json:"micros"`
}

type Micro struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

func Load() API {
	var api API
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Loading api.json ", err)
	}
	if err := viper.Marshal(&api); err != nil {
		log.Fatal("Malformed api.json ", err)
	}
	log.Debug("Config loaded")
	return api
}

func init() {
	setConfigLocations()
	setupEnvVars()
	initLogLevel()
}

func setConfigLocations() {
	viper.SetConfigName("api")         // api.json file name
	viper.AddConfigPath("/etc/apigo/") // deb package config location
	viper.AddConfigPath("config/")     //	local dev config location
}

func setupEnvVars() {
	viper.SetEnvPrefix("apigo")
	viper.BindEnv("env")
	viper.SetDefault("env", "development")
	viper.BindEnv("loglevel")
	viper.SetDefault("loglevel", "debug")
	viper.BindEnv("host")
	viper.SetDefault("host", "localhost:7777")
}

func initLogLevel() {
	logLevel, err := log.ParseLevel(viper.GetString("loglevel"))
	if err != nil {
		log.Warn("Unsupported log level ", viper.Get("loglevel"))
	}
	log.SetLevel(logLevel)
}
