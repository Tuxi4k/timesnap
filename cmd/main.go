package main

import (
	"log"

	"github.com/Tuxi4k/timesnap/internal/config"
	"github.com/Tuxi4k/timesnap/internal/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	_, err = database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Fatalf("Server error: %v", app.Listen(":"+cfg.Server.Port))
}
