package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/nilai/dto"
	"backend/internal/modules/nilai/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type NilaiController struct {
	service service.NilaiService
}

func NewNilaiController(service service.NilaiService) *NilaiController {
	return &NilaiController{service: service}
}

func (c *NilaiController) CreateNilai(ctx *fiber.Ctx) error {
	var req dto.CreateNilaiRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateNilai(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create nilai successfully", resp)
}

func (c *NilaiController) GetAllNilai(ctx *fiber.Ctx) error {
	page      := ctx.Query("page", "1")
	pageSize  := ctx.Query("page_size", "10")
	idPeserta := ctx.Query("id_peserta", "")
	idJadwal  := ctx.Query("id_jadwal", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllNilai(pageNum, pageSizeNum, idPeserta, idJadwal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all nilai successfully", resp)
}

func (c *NilaiController) GetNilaiByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetNilaiByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get nilai successfully", resp)
}

func (c *NilaiController) GetNilaiByPeserta(ctx *fiber.Ctx) error {
	idPeserta := ctx.Params("id_peserta")
	page      := ctx.Query("page", "1")
	pageSize  := ctx.Query("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetNilaiByPeserta(idPeserta, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get nilai by peserta successfully", resp)
}

func (c *NilaiController) GetNilaiByJadwal(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
	page     := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetNilaiByJadwal(idJadwal, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get nilai by jadwal successfully", resp)
}

func (c *NilaiController) UpdateNilai(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateNilaiRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateNilai(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update nilai successfully", resp)
}

func (c *NilaiController) DeleteNilai(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteNilai(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete nilai successfully", nil)
}

func (c *NilaiController) RestoreNilai(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreNilai(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore nilai successfully", nil)
}

func (c *NilaiController) MulaiUjian(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
	if idJadwal == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_jadwal tidak boleh kosong", nil)
	}

	// Ambil user_id (= id_peserta) dari JWT
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "Invalid token claims", nil)
	}
	idPeserta, ok := claims["user_id"].(string)
	if !ok || idPeserta == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "user_id tidak ditemukan di token", nil)
	}

	resp, isNew, err := c.service.MulaiUjian(idPeserta, idJadwal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if isNew {
		return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Mulai ujian successfully", resp)
	}
	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Lanjutkan ujian successfully", resp)
}
