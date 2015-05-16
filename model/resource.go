package model

import (
	"github.com/SHMEDIALIMITED/apigo/server/backends"
	"github.com/gocraft/web"
)

type Plugin interface {
	Inbound(req *web.Request)
}

type Resource struct {
	Auth       string   `json:"auth"`
	Path       string   `json:"path"`
	Methods    []string `json:"methods"`
	Micros     []Micro  `json:"micros"`
	Plugins    []string `json:"plugins"`
	Middleware []Plugin
	Backends   backends.Backends // Load balancer
}
