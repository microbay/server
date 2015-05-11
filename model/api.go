package model

import (
	"errors"
	"net/http"
)

type API struct {
	Name      string      `json:"name"`
	Portal    string      `json:"portal"`
	Resources []*Resource `json:"resources"`
	Key       []byte
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
