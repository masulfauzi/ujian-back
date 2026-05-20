package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/kelas/dto"
	"backend/internal/modules/kelas/service"

	"github.com/gofiber/fiber/v2"
)

type KelasController struct {
	service service.KelasService
}

func NewKelasController(service service.KelasService) *KelasController {
	return &KelasController{service: service}
}

func (c *KelasController) CreateKelas(ctx *fiber.Ctx) error {
	var req dto.CreateKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateKelas(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create kelas successfully", resp)
}

func (c *KelasController) GetAllKelas(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idJurusan := ctx.Query("id_jurusan", "")
	tingkat := ctx.Query("tingkat", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllKelas(pageNum, pageSizeNum, idJurusan, tingkat)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all kelas successfully", resp)
}

func (c *KelasController) GetKelasByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetKelasByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get kelas successfully", resp)
}

func (c *KelasController) UpdateKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateKelas(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update kelas successfully", resp)
}

func (c *KelasController) DeleteKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteKelas(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete kelas successfully", nil)
}

func (c *KelasController) RestoreKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreKelas(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore kelas successfully", nil)
}
