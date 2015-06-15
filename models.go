package server

import (
	"errors"
	//log "github.com/Sirupsen/logrus"
	"github.com/microbay/server/backends"
	"github.com/microbay/server/plugin"
	"net/http"
	"regexp"
)

type API struct {
	Name      string      `json:"name"`
	Portal    string      `json:"portal"`
	Resources []*Resource `json:"resources"`
	plugins   map[string]map[string]interface{}
}

func (a *API) FindResourceByRequest(req *http.Request) (*Resource, error) {
	for _, resource := range a.Resources {
		if resource.Regex.MatchString(req.URL.Path) == true {
			if stringInSlice(req.Method, resource.Methods) {
				return resource, nil
			} else {
				return resource, errors.New("Method")
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
	Micros     map[string][]Micro       `json:"micros"`
	Plugins    []map[string]interface{} `json:"plugins"`
	Middleware []plugin.Interface
	Backends   map[string]backends.Backends // Load balancer
	Regex      *regexp.Regexp
	Keys       []string
}

type Micro struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}
