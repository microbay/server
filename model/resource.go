package model

import (
	"github.com/SHMEDIALIMITED/apigo/plugin"
	"github.com/SHMEDIALIMITED/apigo/server/backends"
)

type Resource struct {
	Auth       string   `json:"auth"`
	Path       string   `json:"path"`
	Methods    []string `json:"methods"`
	Micros     []Micro  `json:"micros"`
	Plugins    []string `json:"plugins"`
	Middleware []plugin.Plugin
	Backends   backends.Backends // Load balancer
}
