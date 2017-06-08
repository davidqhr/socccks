package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Daemon  bool
	LogFile os.File
	ErrFile os.File
	Address string
	Users   map[string]int
}

func loadConfig(filePath string) *Configuration {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalln(err)
	}

	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	err = decoder.Decode(&configuration)

	if err != nil {
		log.Fatalln(err)
	}

	return &configuration
}
