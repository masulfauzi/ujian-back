package controller

import (
	"backend/internal/helpers"
	"backend/internal/modules/user/dto"
	"backend/internal/modules/user/service"
	"backend/internal/modules/user/validator"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

// isSelf: apakah id yang diubah adalah akun yang sedang login (ctx.Locals("user_id") dari JWTAuth).
func isSelf(ctx *fiber.Ctx, id string) bool {
	userID, _ := ctx.Locals("user_id").(string)
	return userID != "" && userID == id
}

// Create godoc
// @Summary Daftarkan akun admin/guru
// @Description Hanya admin. Mendaftarkan username yang boleh masuk sebagai admin/guru. Jika password dikosongkan (disarankan untuk guru), login memakai username & password sianida; jika diisi, akun memakai password lokal (mis. akun cadangan). Email boleh kosong (diisi otomatis <username>@sianida.local). Role default admin.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "Data user"
// @Success 201 {object} helpers.Response{data=dto.UserResponse} "User created successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal atau username/email sudah terdaftar"
// @Failure 401 {object} helpers.Response "Token tidak ada/tidak valid"
// @Failure 403 {object} helpers.Response "Akses hanya untuk admin"
// @Router /users [post]
func (c *UserController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateUserRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	if err := validator.ValidateCreateUser(req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Validation error", nil)
	}

	resp, err := c.service.Create(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "User created successfully", resp)
}

// GetAll godoc
// @Summary List semua user
// @Description Hanya admin. Mengambil daftar seluruh user (admin/guru). Field login_sianida = true berarti akun login memakai password sianida.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helpers.Response{data=[]dto.UserResponse} "Get all users successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Failure 401 {object} helpers.Response "Token tidak ada/tidak valid"
// @Failure 403 {object} helpers.Response "Akses hanya untuk admin"
// @Router /users [get]
func (c *UserController) GetAll(ctx *fiber.Ctx) error {
	users, err := c.service.GetAll()
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get users", nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all users successfully", users)
}

// GetByID godoc
// @Summary Detail user
// @Description Hanya admin. Mengambil satu user berdasarkan ID.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID user (uuid)"
// @Success 200 {object} helpers.Response{data=dto.UserResponse} "Get user successfully"
// @Failure 404 {object} helpers.Response "User tidak ditemukan"
// @Failure 401 {object} helpers.Response "Token tidak ada/tidak valid"
// @Failure 403 {object} helpers.Response "Akses hanya untuk admin"
// @Router /users/{id} [get]
func (c *UserController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	user, err := c.service.GetByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get user successfully", user)
}

// Update godoc
// @Summary Update user
// @Description Hanya admin. Field yang tidak dikirim tidak diubah. status: active/inactive (inactive = tidak bisa login). password: tidak dikirim = tetap, "" = pindah ke login sianida, diisi = set password lokal (min 6). Admin tidak bisa menonaktifkan akunnya sendiri.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID user (uuid)"
// @Param request body dto.UpdateUserRequest true "Field yang diubah"
// @Success 200 {object} helpers.Response{data=dto.UserResponse} "User updated successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Failure 404 {object} helpers.Response "User tidak ditemukan"
// @Failure 401 {object} helpers.Response "Token tidak ada/tidak valid"
// @Failure 403 {object} helpers.Response "Akses hanya untuk admin"
// @Router /users/{id} [put]
func (c *UserController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req dto.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	if err := validator.ValidateUpdateUser(req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Validation error", nil)
	}

	if req.Status == "inactive" && isSelf(ctx, id) {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Tidak bisa menonaktifkan akun sendiri", nil)
	}

	user, err := c.service.Update(id, &req)
	if err != nil {
		if err.Error() == "Resource not found" {
			return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
		}
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "User updated successfully", user)
}

// Delete godoc
// @Summary Hapus user
// @Description Hanya admin. Menghapus user berdasarkan ID. Admin tidak bisa menghapus akunnya sendiri.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID user (uuid)"
// @Success 200 {object} helpers.Response "User deleted successfully"
// @Failure 404 {object} helpers.Response "User tidak ditemukan"
// @Failure 500 {object} helpers.Response "Gagal menghapus user"
// @Failure 401 {object} helpers.Response "Token tidak ada/tidak valid"
// @Failure 403 {object} helpers.Response "Akses hanya untuk admin"
// @Router /users/{id} [delete]
func (c *UserController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if isSelf(ctx, id) {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Tidak bisa menghapus akun sendiri", nil)
	}

	if err := c.service.Delete(id); err != nil {
		if err.Error() == "Resource not found" {
			return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
		}
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete user", nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "User deleted successfully", nil)
}
