# Issue: Implementasi Tabel Soal dan API CRUD Soal

## 📋 Deskripsi
Implementasi lengkap module soal (questions) dengan migration database, soft delete, dan 6 endpoint CRUD lengkap dengan dokumentasi API.

## 🎯 Objektif
Membuat module soal yang terintegrasi dengan bank_soal dan mendukung:
- Soal dengan opsi jawaban (A-E)
- Gambar untuk soal dan opsi jawaban
- Kunci jawaban
- Soft delete
- Dokumentasi API lengkap

---

## 📋 Tahapan Implementasi

### Phase 1: Database Migration & Model
**Durasi: ~30 menit**
**Tujuan**: Setup database structure untuk tabel soal

#### 1.1 Buat Migration File
**File**: `internal/database/migrations/XXXX_create_soal_table.go` (jika menggunakan SQL files di folder)
**Atau langsung modifikasi**: `internal/database/migrate.go`

**Struktur Tabel Soal:**
```sql
CREATE TABLE soal (
  -- Primary Key & Foreign Key
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  id_bank_soal UUID NOT NULL REFERENCES bank_soal(id) ON DELETE CASCADE,
  
  -- Konten Soal
  soal TEXT NOT NULL,
  gambar_soal VARCHAR(500),
  
  -- Opsi Jawaban
  opsi_a TEXT NOT NULL,
  opsi_b TEXT NOT NULL,
  opsi_c TEXT NOT NULL,
  opsi_d TEXT,
  opsi_e TEXT,
  
  -- Gambar Opsi (opsional - untuk soal berbentuk gambar)
  gambar_a VARCHAR(500),
  gambar_b VARCHAR(500),
  gambar_c VARCHAR(500),
  gambar_d VARCHAR(500),
  gambar_e VARCHAR(500),
  
  -- Jawaban Benar (A, B, C, D, atau E)
  kunci VARCHAR(1) NOT NULL,
  
  -- Waktu
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  -- Tracking (opsional)
  created_by UUID,
  updated_by UUID,
  
  -- Index
  INDEX idx_id_bank_soal (id_bank_soal),
  INDEX idx_deleted_at (deleted_at)
);
```

**Penjelasan:**
- `id_bank_soal`: Foreign key ke bank_soal, cascade delete jika bank_soal dihapus
- `soal`: Pertanyaan utama (text panjang)
- `gambar_soal`: URL/path gambar pertanyaan (opsional)
- `opsi_a, opsi_b, opsi_c`: Mandatory options (3 opsi wajib)
- `opsi_d, opsi_e`: Optional (bisa NULL untuk soal dengan 3 opsi)
- `gambar_*`: URL/path gambar untuk setiap opsi (opsional)
- `kunci`: Jawaban benar (constraint: hanya A-E, dan minimal 3 opsi harus ada)
- Timestamps & soft delete: Konsisten dengan module lain

#### 1.2 Buat Model Struct
**File**: `internal/modules/soal/model/soal_model.go`

```go
package model

import (
	"database/sql"
	"time"
)

type Soal struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IdBankSoal   string         `gorm:"type:uuid;index" json:"id_bank_soal"`
	Soal         string         `gorm:"type:text" json:"soal"`
	GambarSoal   string         `gorm:"type:varchar(500)" json:"gambar_soal"`
	OpsiA        string         `gorm:"type:text" json:"opsi_a"`
	OpsiB        string         `gorm:"type:text" json:"opsi_b"`
	OpsiC        string         `gorm:"type:text" json:"opsi_c"`
	OpsiD        string         `gorm:"type:text" json:"opsi_d"`
	OpsiE        string         `gorm:"type:text" json:"opsi_e"`
	GambarA      string         `gorm:"type:varchar(500)" json:"gambar_a"`
	GambarB      string         `gorm:"type:varchar(500)" json:"gambar_b"`
	GambarC      string         `gorm:"type:varchar(500)" json:"gambar_c"`
	GambarD      string         `gorm:"type:varchar(500)" json:"gambar_d"`
	GambarE      string         `gorm:"type:varchar(500)" json:"gambar_e"`
	Kunci        string         `gorm:"type:varchar(1)" json:"kunci"` // A, B, C, D, E
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time     `gorm:"index" json:"deleted_at"`
	CreatedBy    sql.NullString `gorm:"type:uuid" json:"created_by"`
	UpdatedBy    sql.NullString `gorm:"type:uuid" json:"updated_by"`
}

func (Soal) TableName() string {
	return "soal"
}
```

**Penjelasan:**
- Mengikuti pattern yang sama dengan BankSoal model
- Semua field string untuk opsi (TEXT type untuk mendukung konten panjang)
- Kunci sebagai VARCHAR(1) dengan validasi di service layer

#### 1.3 Register Migration di Database Init
**File**: `internal/database/migrate.go`

Tambahkan ke `RunMigrations()`:
```go
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		// ... existing models
		&soalmodel.Soal{},
	)
}
```

**Penjelasan:**
- GORM akan auto-create tabel berdasarkan struct
- Foreign key relationship auto-handled oleh GORM

---

### Phase 2: DTO (Data Transfer Objects)
**Durasi: ~20 menit**
**Tujuan**: Definisi request/response format

#### 2.1 Buat DTO File
**File**: `internal/modules/soal/dto/soal_dto.go`

```go
package dto

type CreateSoalRequest struct {
	IdBankSoal string `json:"id_bank_soal" validate:"required"`
	Soal       string `json:"soal" validate:"required"`
	GambarSoal string `json:"gambar_soal"`
	OpsiA      string `json:"opsi_a" validate:"required"`
	OpsiB      string `json:"opsi_b" validate:"required"`
	OpsiC      string `json:"opsi_c" validate:"required"`
	OpsiD      string `json:"opsi_d"`
	OpsiE      string `json:"opsi_e"`
	GambarA    string `json:"gambar_a"`
	GambarB    string `json:"gambar_b"`
	GambarC    string `json:"gambar_c"`
	GambarD    string `json:"gambar_d"`
	GambarE    string `json:"gambar_e"`
	Kunci      string `json:"kunci" validate:"required,len=1"` // Hanya A-E
}

type UpdateSoalRequest struct {
	Soal       string `json:"soal" validate:"required"`
	GambarSoal string `json:"gambar_soal"`
	OpsiA      string `json:"opsi_a" validate:"required"`
	OpsiB      string `json:"opsi_b" validate:"required"`
	OpsiC      string `json:"opsi_c" validate:"required"`
	OpsiD      string `json:"opsi_d"`
	OpsiE      string `json:"opsi_e"`
	GambarA    string `json:"gambar_a"`
	GambarB    string `json:"gambar_b"`
	GambarC    string `json:"gambar_c"`
	GambarD    string `json:"gambar_d"`
	GambarE    string `json:"gambar_e"`
	Kunci      string `json:"kunci" validate:"required,len=1"`
}

type SoalResponse struct {
	ID         string `json:"id"`
	IdBankSoal string `json:"id_bank_soal"`
	Soal       string `json:"soal"`
	GambarSoal string `json:"gambar_soal"`
	OpsiA      string `json:"opsi_a"`
	OpsiB      string `json:"opsi_b"`
	OpsiC      string `json:"opsi_c"`
	OpsiD      string `json:"opsi_d"`
	OpsiE      string `json:"opsi_e"`
	GambarA    string `json:"gambar_a"`
	GambarB    string `json:"gambar_b"`
	GambarC    string `json:"gambar_c"`
	GambarD    string `json:"gambar_d"`
	GambarE    string `json:"gambar_e"`
	Kunci      string `json:"kunci"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type SoalListResponse struct {
	Data      []SoalResponse `json:"data"`
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"page_size"`
	TotalPage int            `json:"total_page"`
}
```

**Penjelasan:**
- CreateSoalRequest: 3 opsi wajib (A, B, C), 2 opsi opsional (D, E)
- Validasi di struct tags untuk basic validation
- Business logic validation di service layer (e.g., kunci harus valid)

---

### Phase 3: Repository Layer
**Durasi: ~45 menit**
**Tujuan**: Data access layer dengan soft delete support

#### 3.1 Buat Repository Interface & Implementation
**File**: `internal/modules/soal/repository/soal_repository.go`

**Interface:**
```go
package repository

import (
	"backend/internal/modules/soal/model"
	"gorm.io/gorm"
)

type SoalRepository interface {
	Create(soal *model.Soal) error
	GetByID(id string) (*model.Soal, error)
	GetAll(page, pageSize int) ([]model.Soal, int64, error)
	GetByBankSoalID(bankSoalID string, page, pageSize int) ([]model.Soal, int64, error)
	Update(soal *model.Soal) error
	Delete(id string) error
	Restore(id string) error
	HardDelete(id string) error
}
```

**Implementation:**
```go
type soalRepository struct {
	db *gorm.DB
}

func NewSoalRepository(db *gorm.DB) SoalRepository {
	return &soalRepository{db: db}
}

// Create - Insert soal baru
func (r *soalRepository) Create(soal *model.Soal) error {
	return r.db.Create(soal).Error
}

// GetByID - Get soal by ID (exclude soft deleted)
func (r *soalRepository) GetByID(id string) (*model.Soal, error) {
	var soal model.Soal
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&soal).Error
	if err != nil {
		return nil, err
	}
	return &soal, nil
}

// GetAll - Get all soal with pagination (exclude soft deleted)
func (r *soalRepository) GetAll(page, pageSize int) ([]model.Soal, int64, error) {
	var soals []model.Soal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.Soal{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Find(&soals).Error

	return soals, total, err
}

// GetByBankSoalID - Get soal by bank_soal_id with pagination
func (r *soalRepository) GetByBankSoalID(bankSoalID string, page, pageSize int) ([]model.Soal, int64, error) {
	var soals []model.Soal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.Soal{}).
		Where("id_bank_soal = ? AND deleted_at IS NULL", bankSoalID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("id_bank_soal = ? AND deleted_at IS NULL", bankSoalID).
		Offset(offset).
		Limit(pageSize).
		Find(&soals).Error

	return soals, total, err
}

// Update - Update soal
func (r *soalRepository) Update(soal *model.Soal) error {
	return r.db.Save(soal).Error
}

// Delete - Soft delete (set deleted_at = NOW())
func (r *soalRepository) Delete(id string) error {
	now := time.Now()
	return r.db.Model(&model.Soal{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// Restore - Set deleted_at = NULL
func (r *soalRepository) Restore(id string) error {
	return r.db.Model(&model.Soal{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
}

// HardDelete - Permanent delete
func (r *soalRepository) HardDelete(id string) error {
	return r.db.Unscoped().Delete(&model.Soal{}, "id = ?", id).Error
}
```

**Penjelasan:**
- Semua SELECT query filter `WHERE deleted_at IS NULL` (soft delete pattern)
- Delete menggunakan explicit UPDATE dengan timestamp
- GetByBankSoalID penting untuk filter soal per bank_soal
- Restore menggunakan `gorm.Expr("NULL")` untuk null handling

---

### Phase 4: Service Layer
**Durasi: ~45 menit**
**Tujuan**: Business logic & validation

#### 4.1 Buat Service Interface & Implementation
**File**: `internal/modules/soal/service/soal_service.go`

**Interface:**
```go
package service

import "backend/internal/modules/soal/dto"

type SoalService interface {
	CreateSoal(req *dto.CreateSoalRequest) (*dto.SoalResponse, error)
	GetSoalByID(id string) (*dto.SoalResponse, error)
	GetAllSoal(page, pageSize int) (*dto.SoalListResponse, error)
	GetSoalByBankSoal(bankSoalID string, page, pageSize int) (*dto.SoalListResponse, error)
	UpdateSoal(id string, req *dto.UpdateSoalRequest) (*dto.SoalResponse, error)
	DeleteSoal(id string) error
	RestoreSoal(id string) error
}
```

**Implementation Key Points:**

```go
type soalService struct {
	repo repository.SoalRepository
}

// CreateSoal - Validate & create
func (s *soalService) CreateSoal(req *dto.CreateSoalRequest) (*dto.SoalResponse, error) {
	// Validation logic
	if err := s.validateKunci(req.Kunci, req.OpsiA, req.OpsiB, req.OpsiC, req.OpsiD, req.OpsiE); err != nil {
		return nil, err
	}

	soal := &model.Soal{
		IdBankSoal: req.IdBankSoal,
		Soal:       req.Soal,
		// ... copy fields
	}

	if err := s.repo.Create(soal); err != nil {
		return nil, err
	}

	return s.modelToResponse(soal), nil
}

// Helper: validateKunci
func (s *soalService) validateKunci(kunci string, opsiA, opsiB, opsiC, opsiD, opsiE string) error {
	validKeys := map[string]bool{
		"A": true, "B": true, "C": true, "D": true, "E": true,
	}
	
	if !validKeys[kunci] {
		return errors.New("kunci harus A, B, C, D, atau E")
	}

	// Pastikan kunci merujuk ke opsi yang valid
	switch kunci {
	case "D":
		if opsiD == "" {
			return errors.New("opsi D tidak boleh kosong jika kunci D")
		}
	case "E":
		if opsiE == "" {
			return errors.New("opsi E tidak boleh kosong jika kunci E")
		}
	}
	
	return nil
}

// Helper: modelToResponse
func (s *soalService) modelToResponse(soal *model.Soal) *dto.SoalResponse {
	return &dto.SoalResponse{
		ID:         soal.ID,
		IdBankSoal: soal.IdBankSoal,
		Soal:       soal.Soal,
		// ... map all fields
		CreatedAt:  soal.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  soal.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
```

**Penjelasan:**
- `validateKunci()`: Ensure kunci valid (A-E) dan opsi yang dirujuk tidak kosong
- `modelToResponse()`: Convert model ke DTO dengan format timestamp yang konsisten
- GetAllSoal & GetSoalByBankSoal: Handle pagination response format

---

### Phase 5: Controller Layer
**Durasi: ~40 menit**
**Tujuan**: HTTP endpoint handlers

#### 5.1 Buat Controller
**File**: `internal/modules/soal/controller/soal_controller.go`

**6 Endpoints:**
```go
package controller

import (
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

// POST /api/soal - Create
func (c *SoalController) CreateSoal(ctx *fiber.Ctx) error {
	var req dto.CreateSoalRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateSoal(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create soal successfully", resp)
}

// GET /api/soal/:id - Get Detail
func (c *SoalController) GetSoalByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	resp, err := c.service.GetSoalByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}
	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get soal successfully", resp)
}

// GET /api/soal - Get All
func (c *SoalController) GetAllSoal(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")

	pageNum, _ := strconv.Atoi(page)
	pageSizeNum, _ := strconv.Atoi(pageSize)

	resp, err := c.service.GetAllSoal(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all soal successfully", resp)
}

// GET /api/soal/bank/:bank_soal_id - Get by Bank Soal
func (c *SoalController) GetSoalByBankSoal(ctx *fiber.Ctx) error {
	bankSoalID := ctx.Params("bank_soal_id")
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")

	pageNum, _ := strconv.Atoi(page)
	pageSizeNum, _ := strconv.Atoi(pageSize)

	resp, err := c.service.GetSoalByBankSoal(bankSoalID, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get soal by bank successfully", resp)
}

// PUT /api/soal/:id - Update
func (c *SoalController) UpdateSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateSoalRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateSoal(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update soal successfully", resp)
}

// DELETE /api/soal/:id - Soft Delete
func (c *SoalController) DeleteSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.service.DeleteSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}
	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete soal successfully", nil)
}

// PATCH /api/soal/:id/restore - Restore
func (c *SoalController) RestoreSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.service.RestoreSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}
	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore soal successfully", nil)
}
```

**Penjelasan:**
- 6 endpoints standard untuk CRUD + restore
- Query parameters untuk pagination (page, page_size)
- Status code konsisten: 201 Create, 200 OK, 400 Bad Request, 404 Not Found
- Error handling dengan helper function

---

### Phase 6: Routes Setup
**Durasi: ~15 menit**
**Tujuan**: Register routes di application

#### 6.1 Buat Routes File
**File**: `internal/modules/soal/routes/soal_routes.go`

```go
package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/soal/controller"
	"backend/internal/modules/soal/repository"
	"backend/internal/modules/soal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupSoalRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewSoalRepository(db)
	svc := service.NewSoalService(repo)
	ctrl := controller.NewSoalController(svc)

	api := app.Group("/api")
	soal := api.Group("/soal")

	// Public endpoints (GET)
	soal.Get("/", ctrl.GetAllSoal)
	soal.Get("/:id", ctrl.GetSoalByID)
	soal.Get("/bank/:bank_soal_id", ctrl.GetSoalByBankSoal)

	// Protected endpoints (write operations)
	soal.Post("/", middleware.JWTAuth(), ctrl.CreateSoal)
	soal.Put("/:id", middleware.JWTAuth(), ctrl.UpdateSoal)
	soal.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteSoal)
	soal.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreSoal)
}
```

**Penjelasan:**
- GET endpoints public (tidak perlu auth)
- POST, PUT, DELETE, PATCH memerlukan JWT auth
- Routes organized dalam group `/api/soal`
- Dependency injection: repo → svc → ctrl

#### 6.2 Register Routes di Main App
**File**: `cmd/server/main.go`

Tambahkan di `setupRoutes()`:
```go
import soalroutes "backend/internal/modules/soal/routes"

func setupRoutes(app *fiber.App) {
	// ... existing routes
	soalroutes.SetupSoalRoutes(app, database.DB)
}
```

---

### Phase 7: API Documentation
**Durasi: ~60 menit**
**Tujuan**: Complete API documentation dengan examples

#### 7.1 Buat API Documentation File
**File**: `docs/SOAL_API.md`

**Dokumentasi mencakup:**

1. **Base URL & Authentication**
   - Endpoint base: `http://localhost:3000/api/soal`
   - Auth: JWT token di header untuk write operations

2. **Endpoints (6 total):**
   - `POST /api/soal` - Create soal
   - `GET /api/soal` - List all soal dengan pagination
   - `GET /api/soal/:id` - Get detail soal
   - `GET /api/soal/bank/:bank_soal_id` - Filter soal by bank_soal
   - `PUT /api/soal/:id` - Update soal
   - `DELETE /api/soal/:id` - Soft delete soal
   - `PATCH /api/soal/:id/restore` - Restore soft-deleted soal

3. **Setiap endpoint dokumentasikan:**
   - Method & Path
   - Authentication requirement
   - Request body (contoh JSON)
   - Request parameters
   - Success response (200/201)
   - Error responses (400/404/500)

**Example: POST /api/soal**

```markdown
### POST - Create Soal

**Endpoint:** `POST /api/soal`

**Authentication:** ✅ Required (JWT Token)

**Request Body:**
```json
{
  "id_bank_soal": "5112e444-25d8-4ca6-859f-3d24099f45ce",
  "soal": "Berapa hasil dari 2 + 2?",
  "gambar_soal": "https://example.com/soal.jpg",
  "opsi_a": "3",
  "opsi_b": "4",
  "opsi_c": "5",
  "opsi_d": "6",
  "opsi_e": "7",
  "gambar_a": "https://example.com/a.jpg",
  "gambar_b": "https://example.com/b.jpg",
  "gambar_c": "https://example.com/c.jpg",
  "gambar_d": "https://example.com/d.jpg",
  "gambar_e": "https://example.com/e.jpg",
  "kunci": "B"
}
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Create soal successfully",
  "data": {
    "id": "abc123...",
    "id_bank_soal": "5112e444-25d8-4ca6-859f-3d24099f45ce",
    "soal": "Berapa hasil dari 2 + 2?",
    "gambar_soal": "https://example.com/soal.jpg",
    "opsi_a": "3",
    "opsi_b": "4",
    "opsi_c": "5",
    "opsi_d": "6",
    "opsi_e": "7",
    "gambar_a": "https://example.com/a.jpg",
    "gambar_b": "https://example.com/b.jpg",
    "gambar_c": "https://example.com/c.jpg",
    "gambar_d": "https://example.com/d.jpg",
    "gambar_e": "https://example.com/e.jpg",
    "kunci": "B",
    "created_at": "2026-05-18 14:00:00",
    "updated_at": "2026-05-18 14:00:00"
  },
  "errors": null
}
```

**Field Validation:**
- `id_bank_soal`: Required, must exist in bank_soal table
- `soal`: Required, min length 10
- `opsi_a`, `opsi_b`, `opsi_c`: Required
- `opsi_d`, `opsi_e`: Optional
- `kunci`: Required, must be A/B/C/D/E and match valid opsi
- Gambar fields: Optional, URL format
```

4. **GET /api/soal**
   - Query params: page, page_size
   - Response dengan pagination info
   - Exclude soft-deleted items

5. **GET /api/soal/bank/:bank_soal_id**
   - Filter soal berdasarkan bank_soal
   - Dengan pagination
   - Contoh: GET /api/soal/bank/5112e444-25d8-4ca6-859f-3d24099f45ce?page=1&page_size=20

6. **Error Handling**
   - 400: Invalid request format, validasi gagal
   - 404: Soal atau bank_soal tidak ditemukan
   - 401: Unauthorized (missing/invalid token)
   - 500: Server error

7. **cURL Examples**
   - Untuk setiap endpoint
   - Dengan sample data

---

### Phase 8: Testing & Validation
**Durasi: ~30 menit**
**Tujuan**: Ensure semua endpoint bekerja

#### 8.1 Manual Testing Checklist

**Create Soal:**
- [ ] Create dengan 3 opsi (A, B, C)
- [ ] Create dengan 5 opsi (A-E)
- [ ] Validate kunci must match valid opsi
- [ ] Validate id_bank_soal exists
- [ ] Test auth requirement

**Read Soal:**
- [ ] Get all soal (pagination)
- [ ] Get soal by ID
- [ ] Get soal by bank_soal_id
- [ ] Soft-deleted items tidak appear

**Update Soal:**
- [ ] Update soal content
- [ ] Update kunci
- [ ] Validate updated data
- [ ] Test auth requirement

**Delete & Restore:**
- [ ] Soft delete (item masih di DB, deleted_at set)
- [ ] Item hilang dari GET queries
- [ ] Restore soft-deleted item
- [ ] Item muncul lagi di GET queries
- [ ] Test auth requirement

**Pagination:**
- [ ] GetAll pagination works
- [ ] GetByBankSoal pagination works
- [ ] Total count accurate
- [ ] Page calculation correct

---

## 📝 Template yang Bisa Direferensikan

Untuk konsistensi, gunakan module yang sudah ada sebagai template:

1. **Model Template**: `internal/modules/bank_soal/model/bank_soal_model.go`
2. **DTO Template**: `internal/modules/bank_soal/dto/bank_soal_dto.go`
3. **Repository Template**: `internal/modules/bank_soal/repository/bank_soal_repository.go`
4. **Service Template**: `internal/modules/bank_soal/service/bank_soal_service.go`
5. **Controller Template**: `internal/modules/bank_soal/controller/bank_soal_controller.go`
6. **Routes Template**: `internal/modules/bank_soal/routes/bank_soal_routes.go`
7. **API Doc Template**: `docs/BANK_SOAL_API.md`

**Pattern konsistensi yang harus diikuti:**
- Soft delete dengan `deleted_at IS NULL` filtering
- Repository interface dengan error handling
- Service layer untuk business logic & validation
- Controller untuk HTTP handling
- DTO untuk request/response
- Middleware JWTAuth untuk protected endpoints
- Helper functions untuk response formatting

---

## 🔍 Checklist Penyelesaian

- [ ] Phase 1: Model & Migration completed
- [ ] Phase 2: DTO defined
- [ ] Phase 3: Repository implemented
- [ ] Phase 4: Service layer completed
- [ ] Phase 5: Controller endpoints implemented
- [ ] Phase 6: Routes registered in main app
- [ ] Phase 7: API documentation written
- [ ] Phase 8: Manual testing passed
- [ ] Code compiles without errors
- [ ] All endpoints tested with cURL

---

## ⏱️ Estimasi Total Durasi
**Total: ~4-5 jam untuk full implementasi + testing**

- Phase 1: 30 min
- Phase 2: 20 min
- Phase 3: 45 min
- Phase 4: 45 min
- Phase 5: 40 min
- Phase 6: 15 min
- Phase 7: 60 min
- Phase 8: 30 min
- Buffer: 20 min

---

## 🎓 Tips untuk Junior Programmer / AI Model

1. **Ikuti struktur folder yang konsisten** - Semua module mengikuti pattern yang sama
2. **Jangan skip validation** - Validasi kunci & opsi sangat penting
3. **Test setiap phase** - Jangan lanjut ke phase berikutnya sebelum sekarang berhasil
4. **Gunakan helper functions** - ErrorResponse, SuccessResponse dari helpers package
5. **Soft delete pattern** - Selalu filter `WHERE deleted_at IS NULL` pada SELECT
6. **Database constraints** - Foreign key cascade delete setup di migration
7. **Dokumentasi lengkap** - Include curl examples untuk testing
8. **Error messages jelas** - Explain ke user apa yang salah dan cara memperbaikinya

---

**Created**: 2026-05-18
**Status**: Ready for Implementation
