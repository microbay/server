package server

import (
	"errors"
	"github.com/microbay/server/backends"
	"github.com/microbay/server/plugin"
	"net/http"
)

type API struct {
	Name      string      `json:"name"`
	Portal    string      `json:"portal"`
	Resources []*Resource `json:"resources"`
	plugins   map[string]map[string]interface{}
}

func (a *API) FindResourceByRequest(req *http.Request) (*Resource, error) {
	for _, resources := range a.Resources {
		if resources.Path == req.URL.Path {
			if stringInSlice(req.Method, resources.Methods) {
				return resources, nil
			} else {
				return nil, errors.New("Method")
			}
		}
	}
	return nil, errors.New("Resource")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type Resource struct {
	Auth       string                   `json:"auth"`
	Path       string                   `json:"path"`
	Methods    []string                 `json:"methods"`
	Micros     []Micro                  `json:"micros"`
	Plugins    []map[string]interface{} `json:"plugins"`
	Middleware []plugin.Interface
	Backends   backends.Backends // Load balancer
}

type Micro struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}
