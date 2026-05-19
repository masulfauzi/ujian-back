package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jurusan/dto"
	"backend/internal/modules/jurusan/service"

	"github.com/gofiber/fiber/v2"
)

type JurusanController struct {
	service service.JurusanService
}

func NewJurusanController(service service.JurusanService) *JurusanController {
	return &JurusanController{service: service}
}

func (c *JurusanController) CreateJurusan(ctx *fiber.Ctx) error {
	var req dto.CreateJurusanRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJurusan(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jurusan successfully", resp)
}

func (c *JurusanController) GetAllJurusan(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllJurusan(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jurusan successfully", resp)
}

func (c *JurusanController) GetJurusanByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJurusanByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jurusan successfully", resp)
}

func (c *JurusanController) UpdateJurusan(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJurusanRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJurusan(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jurusan successfully", resp)
}

func (c *JurusanController) DeleteJurusan(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJurusan(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jurusan successfully", nil)
}

func (c *JurusanController) RestoreJurusan(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreJurusan(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore jurusan successfully", nil)
}
