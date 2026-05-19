package controller

import (
	"path/filepath"
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/soal/dto"
	"backend/internal/modules/soal/service"

	"github.com/gofiber/fiber/v2"
)

type SoalController struct {
	service service.SoalService
}

func NewSoalController(service service.SoalService) *SoalController {
	return &SoalController{service: service}
}

func (c *SoalController) CreateSoal(ctx *fiber.Ctx) error {
	req := new(dto.CreateSoalRequest)

	req.IdBankSoal = ctx.FormValue("id_bank_soal")
	req.Soal = ctx.FormValue("soal")
	req.OpsiA = ctx.FormValue("opsi_a")
	req.OpsiB = ctx.FormValue("opsi_b")
	req.OpsiC = ctx.FormValue("opsi_c")
	req.OpsiD = ctx.FormValue("opsi_d")
	req.OpsiE = ctx.FormValue("opsi_e")
	req.Kunci = ctx.FormValue("kunci")

	noSoalStr := ctx.FormValue("no_soal")
	if noSoal, err := strconv.Atoi(noSoalStr); err == nil {
		req.NoSoal = noSoal
	}

	if file, err := ctx.FormFile("gambar_soal"); err == nil {
		req.GambarSoal = file
	}
	if file, err := ctx.FormFile("gambar_a"); err == nil {
		req.GambarA = file
	}
	if file, err := ctx.FormFile("gambar_b"); err == nil {
		req.GambarB = file
	}
	if file, err := ctx.FormFile("gambar_c"); err == nil {
		req.GambarC = file
	}
	if file, err := ctx.FormFile("gambar_d"); err == nil {
		req.GambarD = file
	}
	if file, err := ctx.FormFile("gambar_e"); err == nil {
		req.GambarE = file
	}

	resp, err := c.service.CreateSoal(req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create soal successfully", resp)
}

func (c *SoalController) GetSoalByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetSoalByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get soal successfully", resp)
}

func (c *SoalController) GetAllSoal(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetAllSoal(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all soal successfully", resp)
}

func (c *SoalController) GetSoalByBankSoal(ctx *fiber.Ctx) error {
	bankSoalID := ctx.Params("bank_soal_id")
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

	resp, err := c.service.GetSoalByBankSoal(bankSoalID, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get soal by bank successfully", resp)
}

func (c *SoalController) UpdateSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	req := new(dto.UpdateSoalRequest)

	req.Soal = ctx.FormValue("soal")
	req.OpsiA = ctx.FormValue("opsi_a")
	req.OpsiB = ctx.FormValue("opsi_b")
	req.OpsiC = ctx.FormValue("opsi_c")
	req.OpsiD = ctx.FormValue("opsi_d")
	req.OpsiE = ctx.FormValue("opsi_e")
	req.Kunci = ctx.FormValue("kunci")

	noSoalStr := ctx.FormValue("no_soal")
	if noSoal, err := strconv.Atoi(noSoalStr); err == nil {
		req.NoSoal = noSoal
	}

	if file, err := ctx.FormFile("gambar_soal"); err == nil {
		req.GambarSoal = file
	}
	if file, err := ctx.FormFile("gambar_a"); err == nil {
		req.GambarA = file
	}
	if file, err := ctx.FormFile("gambar_b"); err == nil {
		req.GambarB = file
	}
	if file, err := ctx.FormFile("gambar_c"); err == nil {
		req.GambarC = file
	}
	if file, err := ctx.FormFile("gambar_d"); err == nil {
		req.GambarD = file
	}
	if file, err := ctx.FormFile("gambar_e"); err == nil {
		req.GambarE = file
	}

	resp, err := c.service.UpdateSoal(id, req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update soal successfully", resp)
}

func (c *SoalController) DeleteSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete soal successfully", nil)
}

func (c *SoalController) RestoreSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore soal successfully", nil)
}

func (c *SoalController) ImportSoalFromExcel(ctx *fiber.Ctx) error {
	// 1. Parse multipart form
	file, err := ctx.FormFile("file")
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "File tidak ditemukan", map[string]string{
			"error": "Silakan upload file excel",
		})
	}

	// 2. Validate file size (max 10MB)
	const maxFileSize = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "File terlalu besar", map[string]string{
			"error": "Max file size adalah 10MB",
		})
	}

	// 3. Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".xls" && ext != ".xlsx" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Format file tidak valid", map[string]string{
			"error": "File harus berupa .xls atau .xlsx",
		})
	}

	// 4. Build request
	req := &dto.ImportSoalRequest{
		IdBankSoal: ctx.FormValue("id_bank_soal"),
		File:       file,
	}

	// 5. Validate required fields
	if req.IdBankSoal == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_bank_soal tidak ditemukan", nil)
	}

	// 6. Call service
	resp, err := c.service.ImportSoalFromExcel(ctx.Context(), req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Import soal gagal", map[string]string{
			"error": err.Error(),
		})
	}

	// 7. Return success response
	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Import soal berhasil", resp)
}
