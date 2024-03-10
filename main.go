package main

import (
	"log"

	"sentinel_bot/cmd"
	"sentinel_bot/internal/config"
)

func main() {
	conf := config.Load()

	err := cmd.Run(conf)
	if err != nil {
		log.Fatal(err)
	}
}
