package model

import (
	//"container/ring"
	"github.com/SHMEDIALIMITED/apigo/server/backends"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	REDIS string = "redis"
)

type API struct {
	Name      string      `json:"name"`
	Portal    string      `json:"portal"`
	Resources []*Resource `json:"resources"`
}

func (a *API) FindResourceByPath(p string) *Resource {
	for _, v := range a.Resources {
		if v.Path == p {
			return v
		}
	}
	return nil
}

type Resource struct {
	Auth     string            `json:"auth"`
	Path     string            `json:"path"`
	Micros   []Micro           `json:"micros"`
	Backends backends.Backends // Load balancer
}

type Micro struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

func Load() API {
	var api API
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed loading server.json ", err)
	}
	if err := viper.MarshalKey("api", &api); err != nil {
		log.Fatal("Malformed server.json ", err)
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
	viper.SetConfigName("server")      // server.json file name
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
