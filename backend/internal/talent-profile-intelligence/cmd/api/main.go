package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"talentiq/talent-profile-intelligence/internal/config"
	"talentiq/talent-profile-intelligence/internal/handler"
)

func main() {
	cfg := config.Load()

	router := gin.Default()

	router.GET("/health", handler.Health)

	log.Printf(
		"%s starting on port %s",
		cfg.AppName,
		cfg.ServerPort,
	)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}