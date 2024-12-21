package main

import (
	"flag"
	"log"

	"github.com/vanshavenger/goproxynginx/server"
	"github.com/vanshavenger/goproxynginx/utils"
)

func main() {
	configPtr := flag.String("config", "config.yaml", "Configuration file")
	flag.StringVar(configPtr, "c", "config.yaml", "Configuration file (shorthand)")
	flag.Parse()

	fileContents, err := utils.ParseYAMLConfig(*configPtr)
	if err != nil {
		log.Fatal(err)
	}

	validateConfig, err := utils.ValidateConfig(fileContents)

	if err != nil {
		log.Fatal(err)
	}

	server.CreateServer(validateConfig)

}
