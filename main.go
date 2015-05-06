package main

import(
	"github.com/SHMEDIALIMITED/apigo/server"
	"github.com/SHMEDIALIMITED/apigo/config"
)

func main() {
	config := config.New();
	server.Bootstrap(config)
}