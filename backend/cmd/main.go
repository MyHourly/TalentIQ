package main

import (
	"log"

	"github.com/MyHourly/TalentIQ/backend/api"
	"github.com/MyHourly/TalentIQ/backend/config"
)

func main() {
	cfg := config.Load()
	router := api.NewRouter()

	log.Printf("TalentIQ server starting on port %s in %s environment", cfg.ServerPort, cfg.Environment)

	if err := router.Run("127.0.0.1:" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
