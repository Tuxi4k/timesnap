package main

import (
	"log"

	"github.com/Tuxi4k/timesnap/internal/config"
	"github.com/Tuxi4k/timesnap/internal/database"
	"github.com/Tuxi4k/timesnap/internal/modules/deadline"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	deadlinesRoutes := app.Group("deadlines/")
	deadlineRepo := deadline.NewRepository(db)
	deadlineService := deadline.NewService(deadlineRepo)
	deadlineHandler := deadline.NewHandler(deadlineService)

	deadlineHandler.RegisterRoutes(deadlinesRoutes)

	log.Fatalf("Server error: %v", app.Listen(":"+cfg.Server.Port))
}
