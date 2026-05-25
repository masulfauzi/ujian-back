package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jawaban/dto"
	"backend/internal/modules/jawaban/service"

	"github.com/gofiber/fiber/v2"
)

type JawabanController struct {
	service service.JawabanService
}

func NewJawabanController(service service.JawabanService) *JawabanController {
	return &JawabanController{service: service}
}

func (c *JawabanController) CreateJawaban(ctx *fiber.Ctx) error {
	var req dto.CreateJawabanRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJawaban(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jawaban successfully", resp)
}

func (c *JawabanController) GetAllJawaban(ctx *fiber.Ctx) error {
	page      := ctx.Query("page", "1")
	pageSize  := ctx.Query("page_size", "10")
	idNilai   := ctx.Query("id_nilai", "")
	idPeserta := ctx.Query("id_peserta", "")
	idSoal    := ctx.Query("id_soal", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}
	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllJawaban(pageNum, pageSizeNum, idNilai, idPeserta, idSoal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jawaban successfully", resp)
}

func (c *JawabanController) GetJawabanByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJawabanByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban successfully", resp)
}

func (c *JawabanController) GetJawabanByNilai(ctx *fiber.Ctx) error {
	idNilai := ctx.Params("id_nilai")

	resp, err := c.service.GetJawabanByNilai(idNilai)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban by nilai successfully", resp)
}

func (c *JawabanController) GetJawabanByPeserta(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetJawabanByPeserta(idPeserta, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban by peserta successfully", resp)
}

func (c *JawabanController) UpdateJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJawabanRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJawaban(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jawaban successfully", resp)
}

func (c *JawabanController) DeleteJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJawaban(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jawaban successfully", nil)
}

func (c *JawabanController) RestoreJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreJawaban(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore jawaban successfully", nil)
}
