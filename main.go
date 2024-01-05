package main

import (
	"fmt"
	"log"

	"github.com/AkifhanIlgaz/foody-api/cfg"
)

func main() {
	config, err := cfg.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not read environment variables", err)
	}

	fmt.Println(config)
}
