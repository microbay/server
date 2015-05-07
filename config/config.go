package config 

import(
	"github.com/gocraft/web"
	"fmt"
)


type Context struct {

}

type API struct {
	*Context
	Name string
	Resources []Resource
}


type Resource struct {
	*API
	Type int
	Name string
	Path string
	Micros []Micro  
}	


type Micro struct {
	URL string 
	Weight int
}



func (api *API) Root(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "hello, %q", api.Name)
}


func New() (config API) {
	c := API{}
	c.Name = "httpbin"
	return c
}
