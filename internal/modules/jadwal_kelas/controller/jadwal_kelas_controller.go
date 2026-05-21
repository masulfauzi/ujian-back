package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jadwal_kelas/dto"
	"backend/internal/modules/jadwal_kelas/service"

	"github.com/gofiber/fiber/v2"
)

type JadwalKelasController struct {
	service service.JadwalKelasService
}

func NewJadwalKelasController(service service.JadwalKelasService) *JadwalKelasController {
	return &JadwalKelasController{service: service}
}

func (c *JadwalKelasController) CreateJadwalKelas(ctx *fiber.Ctx) error {
	var req dto.CreateJadwalKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJadwalKelas(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jadwal kelas successfully", resp)
}

func (c *JadwalKelasController) GetAllJadwalKelas(ctx *fiber.Ctx) error {
	page     := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idJadwal := ctx.Query("id_jadwal", "")
	idKelas  := ctx.Query("id_kelas", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllJadwalKelas(pageNum, pageSizeNum, idJadwal, idKelas)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jadwal kelas successfully", resp)
}

func (c *JadwalKelasController) GetJadwalKelasByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJadwalKelasByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jadwal kelas successfully", resp)
}

func (c *JadwalKelasController) UpdateJadwalKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJadwalKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJadwalKelas(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jadwal kelas successfully", resp)
}

func (c *JadwalKelasController) DeleteJadwalKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJadwalKelas(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jadwal kelas successfully", nil)
}
