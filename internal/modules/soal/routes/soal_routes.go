package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/soal/controller"
	"backend/internal/modules/soal/repository"
	"backend/internal/modules/soal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupSoalRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewSoalRepository(db)
	svc := service.NewSoalService(repo)
	ctrl := controller.NewSoalController(svc)

	api := app.Group("/api")
	soal := api.Group("/soal")

	// Public endpoints (GET)
	soal.Get("/", ctrl.GetAllSoal)
	soal.Get("/bank/:bank_soal_id", ctrl.GetSoalByBankSoal)

	// Protected endpoints (write operations)
	soal.Post("/", middleware.JWTAuth(), ctrl.CreateSoal)
	soal.Post("/import", middleware.JWTAuth(), ctrl.ImportSoalFromExcel)
	soal.Put("/:id", middleware.JWTAuth(), ctrl.UpdateSoal)
	soal.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteSoal)
	soal.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreSoal)

	// Public endpoint (GET) - dynamic route at the end
	soal.Get("/:id", ctrl.GetSoalByID)
}
