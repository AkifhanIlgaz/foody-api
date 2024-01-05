package main

import (
	"log"

	"github.com/AkifhanIlgaz/foody-api/cfg"
	"github.com/AkifhanIlgaz/foody-api/database"
)

func main() {
	config, err := cfg.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not read environment variables", err)
	}

	databases, err := database.Connect(config)
	if err != nil {
		log.Fatal("Could not connect to databases: ", err)
	}

	defer databases.Postgres.Close()
	defer databases.Redis.Close()
}
