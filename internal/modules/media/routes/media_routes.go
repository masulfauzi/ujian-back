package routes

import (
	"context"
	"fmt"

	"backend/internal/helpers"
	"backend/internal/storage"

	"github.com/gofiber/fiber/v2"
)

var allowedTypes = map[string]bool{
	"soal": true,
	"opsi": true,
}

func SetupMediaRoutes(app *fiber.App) {
	app.Get("/api/images/:type/:filename", getImage)
}

func getImage(ctx *fiber.Ctx) error {
	imgType := ctx.Params("type")
	filename := ctx.Params("filename")

	if !allowedTypes[imgType] {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "tipe gambar tidak valid, gunakan: soal atau opsi", nil)
	}

	if filename == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "filename tidak boleh kosong", nil)
	}

	obj, info, err := storage.GetFile(context.Background(), imgType, filename)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, "gambar tidak ditemukan", nil)
	}
	defer obj.Close()

	contentType := info.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	ctx.Set("Content-Type", contentType)
	ctx.Set("Content-Length", fmt.Sprintf("%d", info.Size))
	ctx.Set("Cache-Control", "public, max-age=86400")

	return ctx.SendStream(obj, int(info.Size))
}
