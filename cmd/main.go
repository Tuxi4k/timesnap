package main

import (
	"context"
	"log"
	"time"

	"github.com/Tuxi4k/swaggen"
	"github.com/Tuxi4k/timesnap/internal/config"
	"github.com/Tuxi4k/timesnap/internal/database"
	"github.com/Tuxi4k/timesnap/internal/modules/deadline"
	"github.com/Tuxi4k/timesnap/internal/pkg/worker"
	"github.com/Tuxi4k/timesnap/pkg/utils/ptr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title TimeSnap API
// @version 1.0
// @description API for managing deadlines
func main() {
	var swaggerJSON *string = ptr.To(`{"swagger":"2.0","info":{"title":"Loading..."}}`)
	go generateSwagger(swaggerJSON)

	app := fiber.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	deadlineWorker := worker.New(db)
	go deadlineWorker.Start(context.Background(), time.Minute)

	registerSwagger(app, swaggerJSON)

	deadlinesRoutes := app.Group("deadlines/")
	deadlineRepo := deadline.NewRepository(db)
	deadlineService := deadline.NewService(deadlineRepo)
	deadlineHandler := deadline.NewHandler(deadlineService)

	deadlineHandler.RegisterRoutes(deadlinesRoutes)

	log.Fatalf("Server error: %v", app.Listen(":"+cfg.Server.Port))
}

func generateSwagger(swaggerJSON *string) {
	*swaggerJSON, _ = swaggen.Generate(swaggen.WithMainAPIFile("cmd/main.go"))
}

func registerSwagger(app *fiber.App, swaggerJSON *string) {
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString(*swaggerJSON)
	})

	app.Get("/swagger/*", swagger.HandlerDefault)
}
