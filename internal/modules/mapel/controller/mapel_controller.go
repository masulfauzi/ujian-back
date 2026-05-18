package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/mapel/dto"
	"backend/internal/modules/mapel/service"

	"github.com/gofiber/fiber/v2"
)

type MapelController struct {
	service service.MapelService
}

func NewMapelController(service service.MapelService) *MapelController {
	return &MapelController{service: service}
}

func (c *MapelController) CreateMapel(ctx *fiber.Ctx) error {
	var req dto.CreateMapelRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateMapel(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create mapel successfully", resp)
}

func (c *MapelController) GetAllMapel(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetAllMapel(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all mapel successfully", resp)
}

func (c *MapelController) GetMapelByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetMapelByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get mapel successfully", resp)
}

func (c *MapelController) UpdateMapel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateMapelRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateMapel(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update mapel successfully", resp)
}

func (c *MapelController) DeleteMapel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteMapel(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete mapel successfully", nil)
}

func (c *MapelController) RestoreMapel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreMapel(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore mapel successfully", nil)
}
