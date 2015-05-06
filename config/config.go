package config 

import(
	
)

type Config struct {
	Name string
	Resources []Resource
}


type Resource struct {
	Type int
	Name string
	Path string
	Micros []Micro  
}	


type Micro struct {
	URL string 
	Weight int
}


func New() (config Config) {
	c := Config{}
	c.Name = "httpbin"
	return c
}
