package main

import (
	"log"

	"github.com/begenov/backend/internal/app"
	"github.com/begenov/backend/internal/config"
)

func main() {
	cfg, err := config.Init(".")
	if err != nil {
		log.Fatalln(err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("[ERROR] Run Application: %v\n", err)
	}
}
