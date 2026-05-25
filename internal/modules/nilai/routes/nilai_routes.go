package routes

import (
	"backend/internal/middleware"
	jawabanrepo "backend/internal/modules/jawaban/repository"
	"backend/internal/modules/nilai/controller"
	"backend/internal/modules/nilai/repository"
	"backend/internal/modules/nilai/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupNilaiRoutes(app *fiber.App, db *gorm.DB) {
	repo              := repository.NewNilaiRepository(db)
	jawabanRepository := jawabanrepo.NewJawabanRepository(db)
	svc               := service.NewNilaiService(repo, jawabanRepository, db)
	ctrl              := controller.NewNilaiController(svc)

	api   := app.Group("/api")
	nilai := api.Group("/nilai")

	// Dua-segmen path: aman dari konflik dengan /:id (beda kedalaman)
	nilai.Get("/peserta/:id_peserta", ctrl.GetNilaiByPeserta)
	nilai.Get("/jadwal/:id_jadwal", ctrl.GetNilaiByJadwal)
	nilai.Post("/mulai-ujian/:id_jadwal", middleware.JWTAuth(), ctrl.MulaiUjian)

	nilai.Post("/", middleware.JWTAuth(), ctrl.CreateNilai)
	nilai.Get("/", ctrl.GetAllNilai)
	nilai.Get("/:id", ctrl.GetNilaiByID)
	nilai.Put("/:id", middleware.JWTAuth(), ctrl.UpdateNilai)
	nilai.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteNilai)
	nilai.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreNilai)
}
