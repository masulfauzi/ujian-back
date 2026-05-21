package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/jadwal/controller"
	"backend/internal/modules/jadwal/repository"
	"backend/internal/modules/jadwal/service"
	jadwalkelasrepo "backend/internal/modules/jadwal_kelas/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupJadwalRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewJadwalRepository(db)
	jadwalKelasRepo := jadwalkelasrepo.NewJadwalKelasRepository(db)
	svc := service.NewJadwalService(repo, jadwalKelasRepo)
	ctrl := controller.NewJadwalController(svc)

	api := app.Group("/api")
	jadwal := api.Group("/jadwal")

	jadwal.Post("/", middleware.JWTAuth(), ctrl.CreateJadwal)
	jadwal.Get("/", ctrl.GetAllJadwal)
	jadwal.Get("/bank-soal/:bank_soal_id", ctrl.GetJadwalByBankSoal)
	jadwal.Get("/:id", ctrl.GetJadwalByID)
	jadwal.Put("/:id", middleware.JWTAuth(), ctrl.UpdateJadwal)
	jadwal.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteJadwal)
	jadwal.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreJadwal)
}
