package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MeteorsLiu/Light/interfaces"
)

func readConf(configPath string) *interfaces.Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	var conf interfaces.Config
	if err = json.Unmarshal(data, &conf); err != nil {
		log.Fatal(err)
	}
	return &conf
}
