package server

import (
	"errors"
	log "github.com/Sirupsen/logrus"
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
		log.Warn(resource.Regex)
		if resource.Regex.MatchString(req.URL.Path) == true {
			if stringInSlice(req.Method, resource.Methods) {
				return resource, nil
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
	Regex      *regexp.Regexp
	Keys       []string
}

type Params map[string]string

func (r *Resource) Params(url string) Params {
	match := r.Regex.FindAllStringSubmatch(url, -1)
	//log.Warn(url)
	result := make(Params)
	for i := range match {
		if len(r.Keys) <= i {
			break
		}
		//result[r.Keys[i]] = match[i]
	}
	return result
}

type Micro struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}
