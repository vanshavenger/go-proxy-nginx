package main

import (
	"flag"
	"log"

	"github.com/vanshavenger/goproxynginx/server"
	"github.com/vanshavenger/goproxynginx/utils"
)

const (
	// DefaultConfigPath is the default configuration file path
	DefaultConfigPath = "config.yaml"
)

func main() {
	configPtr := flag.String("config", DefaultConfigPath, "Configuration file")
	flag.StringVar(configPtr, "c", DefaultConfigPath, "Configuration file (shorthand)")
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
