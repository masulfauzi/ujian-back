package main

import (
	"log"

	"backend/configs"
	"backend/internal/database"
	"backend/internal/database/seeders"
)

func main() {
	if err := configs.LoadEnv(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := seeders.SeedAdmin(database.DB); err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	log.Println("✅ Seeding admin completed!")
	_ = database.Close()
}
