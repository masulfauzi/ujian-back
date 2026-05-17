package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/configs"
	"backend/internal/database"
	"backend/internal/middleware"
	authroutes "backend/internal/modules/auth/routes"
	userroutes "backend/internal/modules/user/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := configs.LoadEnv(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	appConfig := configs.GetAppConfig()

	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.RunMigrations(database.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: appConfig.Name,
	})

	app.Use(middleware.CORS())
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())

	setupRoutes(app)

	go func() {
		addr := fmt.Sprintf(":%d", appConfig.Port)
		log.Printf("Starting server on %s\n", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	_ = app.Shutdown()
	_ = database.Close()
	log.Println("Server shut down successfully")
}

func setupRoutes(app *fiber.App) {
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status": "ok",
			"service": "Fiber Backend API",
		})
	})

	authroutes.SetupAuthRoutes(app, database.DB)
	userroutes.SetupUserRoutes(app, database.DB)
}
