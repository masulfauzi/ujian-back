# Issue: Implementasi Fitur Bank Soal (Question Bank Management)

## 📋 Deskripsi
Implementasi fitur manajemen bank soal (kumpulan soal per mata pelajaran) termasuk:
1. Database migration untuk tabel bank_soal
2. Seeder untuk data awal bank_soal
3. CRUD API endpoints untuk bank_soal
4. Model dan Repository pattern
5. Relasi dengan tabel mapel

---

## 🎯 Objectives

✅ Membuat tabel `bank_soal` di database dengan soft delete support
✅ Membuat relasi foreign key dengan tabel `mapel`
✅ Membuat seeder untuk populate data bank_soal awal
✅ Setup CRUD endpoints untuk bank_soal
✅ Implementasi Repository pattern
✅ Dokumentasi API

---

## 📊 Database Schema

### Tabel: `bank_soal`

```sql
CREATE TABLE bank_soal (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_bank_soal VARCHAR(255) NOT NULL UNIQUE,
    id_mapel UUID NOT NULL,
    jml_soal INTEGER NOT NULL DEFAULT 0,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by UUID,
    updated_by UUID,
    CONSTRAINT fk_bank_soal_mapel FOREIGN KEY (id_mapel) REFERENCES mapel(id)
);

CREATE INDEX idx_bank_soal_nama ON bank_soal(nama_bank_soal);
CREATE INDEX idx_bank_soal_mapel ON bank_soal(id_mapel);
CREATE INDEX idx_bank_soal_deleted ON bank_soal(deleted_at);
```

### Field Details:
- **id** (UUID): Primary key
- **nama_bank_soal** (VARCHAR 255): Nama bank soal, UNIQUE
- **id_mapel** (UUID): Foreign key ke tabel mapel, NOT NULL
- **jml_soal** (INTEGER): Jumlah soal dalam bank soal
- **deskripsi** (TEXT): Deskripsi bank soal (optional)
- **created_at** (TIMESTAMP): Waktu pembuatan record
- **updated_at** (TIMESTAMP): Waktu update terakhir
- **deleted_at** (TIMESTAMP NULL): Soft delete timestamp
- **created_by** (UUID): ID user yang membuat
- **updated_by** (UUID): ID user yang update terakhir

### Relasi:
- **Foreign Key**: `id_mapel` → `mapel.id` (One Mapel has Many Bank Soal)

---

## 🛠️ Tahapan Implementasi

### **FASE 1: Database Setup**

#### Step 1.1: Buat Model File
**File**: `internal/modules/bank_soal/model/bank_soal_model.go`

```go
package model

import (
	"database/sql"
	"time"
)

type BankSoal struct {
	ID            string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NamaBankSoal  string         `gorm:"type:varchar(255);uniqueIndex" json:"nama_bank_soal"`
	IdMapel       string         `gorm:"type:uuid;index" json:"id_mapel"`
	JmlSoal       int            `gorm:"type:integer;default:0" json:"jml_soal"`
	Deskripsi     string         `gorm:"type:text" json:"deskripsi"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time     `gorm:"index" json:"deleted_at"`
	CreatedBy     sql.NullString `gorm:"type:uuid" json:"created_by"`
	UpdatedBy     sql.NullString `gorm:"type:uuid" json:"updated_by"`
}

func (BankSoal) TableName() string {
	return "bank_soal"
}
```

#### Step 1.2: Update Migration
**File**: `internal/database/migrate.go`

Tambahkan `&banksoalmodel.BankSoal{}` ke dalam `db.AutoMigrate()`:

```go
package database

import (
	banksoalmodel "backend/internal/modules/bank_soal/model"
	mapelmodel "backend/internal/modules/mapel/model"
	usermodel "backend/internal/modules/user/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&usermodel.User{},
		&mapelmodel.Mapel{},
		&banksoalmodel.BankSoal{},
	)
}
```

#### Step 1.3: Jalankan Migration

```bash
# Jika sudah ada server running, cukup restart
go run ./cmd/server/main.go

# Jika ingin jalankan manual
go run ./cmd/seed/main.go  # Ini akan otomatis run migrations
```

**Expected Output**:
```
✓ Migration applied: User table
✓ Migration applied: Mapel table
✓ Migration applied: BankSoal table
```

---

### **FASE 2: DTO (Data Transfer Objects)**

#### Step 2.1: Buat DTO File
**File**: `internal/modules/bank_soal/dto/bank_soal_dto.go`

```go
package dto

type CreateBankSoalRequest struct {
	NamaBankSoal string `json:"nama_bank_soal" validate:"required"`
	IdMapel      string `json:"id_mapel" validate:"required"`
	JmlSoal      int    `json:"jml_soal" validate:"required,min=0"`
	Deskripsi    string `json:"deskripsi"`
}

type UpdateBankSoalRequest struct {
	NamaBankSoal string `json:"nama_bank_soal" validate:"required"`
	IdMapel      string `json:"id_mapel" validate:"required"`
	JmlSoal      int    `json:"jml_soal" validate:"required,min=0"`
	Deskripsi    string `json:"deskripsi"`
}

type BankSoalResponse struct {
	ID           string `json:"id"`
	NamaBankSoal string `json:"nama_bank_soal"`
	IdMapel      string `json:"id_mapel"`
	JmlSoal      int    `json:"jml_soal"`
	Deskripsi    string `json:"deskripsi"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type BankSoalListResponse struct {
	Data      []BankSoalResponse `json:"data"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	TotalPage int                `json:"total_page"`
}
```

---

### **FASE 3: Repository & Service**

#### Step 3.1: Buat Repository
**File**: `internal/modules/bank_soal/repository/bank_soal_repository.go`

```go
package repository

import (
	"backend/internal/modules/bank_soal/model"

	"gorm.io/gorm"
)

type BankSoalRepository interface {
	Create(bankSoal *model.BankSoal) error
	GetByID(id string) (*model.BankSoal, error)
	GetAll(page, pageSize int) ([]model.BankSoal, int64, error)
	GetByMapelID(mapelID string, page, pageSize int) ([]model.BankSoal, int64, error)
	Update(bankSoal *model.BankSoal) error
	Delete(id string) error // Soft delete
	Restore(id string) error
	HardDelete(id string) error
}

type bankSoalRepository struct {
	db *gorm.DB
}

func NewBankSoalRepository(db *gorm.DB) BankSoalRepository {
	return &bankSoalRepository{db: db}
}

func (r *bankSoalRepository) Create(bankSoal *model.BankSoal) error {
	return r.db.Create(bankSoal).Error
}

func (r *bankSoalRepository) GetByID(id string) (*model.BankSoal, error) {
	var bankSoal model.BankSoal
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&bankSoal).Error
	if err != nil {
		return nil, err
	}
	return &bankSoal, nil
}

func (r *bankSoalRepository) GetAll(page, pageSize int) ([]model.BankSoal, int64, error) {
	var bankSoals []model.BankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.BankSoal{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Find(&bankSoals).Error

	return bankSoals, total, err
}

func (r *bankSoalRepository) GetByMapelID(mapelID string, page, pageSize int) ([]model.BankSoal, int64, error) {
	var bankSoals []model.BankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.BankSoal{}).
		Where("id_mapel = ? AND deleted_at IS NULL", mapelID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("id_mapel = ? AND deleted_at IS NULL", mapelID).
		Offset(offset).
		Limit(pageSize).
		Find(&bankSoals).Error

	return bankSoals, total, err
}

func (r *bankSoalRepository) Update(bankSoal *model.BankSoal) error {
	return r.db.Save(bankSoal).Error
}

func (r *bankSoalRepository) Delete(id string) error {
	// Soft delete
	return r.db.Delete(&model.BankSoal{}, "id = ?", id).Error
}

func (r *bankSoalRepository) Restore(id string) error {
	// Restore - clear deleted_at
	return r.db.Table("bank_soal").Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *bankSoalRepository) HardDelete(id string) error {
	return r.db.Unscoped().Delete(&model.BankSoal{}, "id = ?", id).Error
}
```

#### Step 3.2: Buat Service
**File**: `internal/modules/bank_soal/service/bank_soal_service.go`

```go
package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/bank_soal/dto"
	"backend/internal/modules/bank_soal/model"
	"backend/internal/modules/bank_soal/repository"

	"gorm.io/gorm"
)

type BankSoalService interface {
	CreateBankSoal(req *dto.CreateBankSoalRequest) (*dto.BankSoalResponse, error)
	GetBankSoalByID(id string) (*dto.BankSoalResponse, error)
	GetAllBankSoal(page, pageSize int) (*dto.BankSoalListResponse, error)
	GetBankSoalByMapel(mapelID string, page, pageSize int) (*dto.BankSoalListResponse, error)
	UpdateBankSoal(id string, req *dto.UpdateBankSoalRequest) (*dto.BankSoalResponse, error)
	DeleteBankSoal(id string) error
	RestoreBankSoal(id string) error
}

type bankSoalService struct {
	repo repository.BankSoalRepository
}

func NewBankSoalService(repo repository.BankSoalRepository) BankSoalService {
	return &bankSoalService{repo: repo}
}

func (s *bankSoalService) CreateBankSoal(req *dto.CreateBankSoalRequest) (*dto.BankSoalResponse, error) {
	bankSoal := &model.BankSoal{
		NamaBankSoal: req.NamaBankSoal,
		IdMapel:      req.IdMapel,
		JmlSoal:      req.JmlSoal,
		Deskripsi:    req.Deskripsi,
	}

	if err := s.repo.Create(bankSoal); err != nil {
		return nil, err
	}

	return s.modelToResponse(bankSoal), nil
}

func (s *bankSoalService) GetBankSoalByID(id string) (*dto.BankSoalResponse, error) {
	bankSoal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.modelToResponse(bankSoal), nil
}

func (s *bankSoalService) GetAllBankSoal(page, pageSize int) (*dto.BankSoalListResponse, error) {
	bankSoals, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var responses []dto.BankSoalResponse
	for _, bs := range bankSoals {
		responses = append(responses, *s.modelToResponse(&bs))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.BankSoalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *bankSoalService) GetBankSoalByMapel(mapelID string, page, pageSize int) (*dto.BankSoalListResponse, error) {
	bankSoals, total, err := s.repo.GetByMapelID(mapelID, page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var responses []dto.BankSoalResponse
	for _, bs := range bankSoals {
		responses = append(responses, *s.modelToResponse(&bs))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.BankSoalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *bankSoalService) UpdateBankSoal(id string, req *dto.UpdateBankSoalRequest) (*dto.BankSoalResponse, error) {
	bankSoal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	bankSoal.NamaBankSoal = req.NamaBankSoal
	bankSoal.IdMapel = req.IdMapel
	bankSoal.JmlSoal = req.JmlSoal
	bankSoal.Deskripsi = req.Deskripsi

	if err := s.repo.Update(bankSoal); err != nil {
		return nil, err
	}

	return s.modelToResponse(bankSoal), nil
}

func (s *bankSoalService) DeleteBankSoal(id string) error {
	bankSoal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(bankSoal.ID)
}

func (s *bankSoalService) RestoreBankSoal(id string) error {
	return s.repo.Restore(id)
}

func (s *bankSoalService) modelToResponse(bankSoal *model.BankSoal) *dto.BankSoalResponse {
	return &dto.BankSoalResponse{
		ID:           bankSoal.ID,
		NamaBankSoal: bankSoal.NamaBankSoal,
		IdMapel:      bankSoal.IdMapel,
		JmlSoal:      bankSoal.JmlSoal,
		Deskripsi:    bankSoal.Deskripsi,
		CreatedAt:    bankSoal.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    bankSoal.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
```

---

### **FASE 4: Controller & Routes**

#### Step 4.1: Buat Controller
**File**: `internal/modules/bank_soal/controller/bank_soal_controller.go`

```go
package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/bank_soal/dto"
	"backend/internal/modules/bank_soal/service"

	"github.com/gofiber/fiber/v2"
)

type BankSoalController struct {
	service service.BankSoalService
}

func NewBankSoalController(service service.BankSoalService) *BankSoalController {
	return &BankSoalController{service: service}
}

func (c *BankSoalController) CreateBankSoal(ctx *fiber.Ctx) error {
	var req dto.CreateBankSoalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateBankSoal(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create bank soal successfully", resp)
}

func (c *BankSoalController) GetAllBankSoal(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetAllBankSoal(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all bank soal successfully", resp)
}

func (c *BankSoalController) GetBankSoalByMapel(ctx *fiber.Ctx) error {
	mapelID := ctx.Params("mapel_id")
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

	resp, err := c.service.GetBankSoalByMapel(mapelID, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get bank soal by mapel successfully", resp)
}

func (c *BankSoalController) GetBankSoalByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetBankSoalByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get bank soal successfully", resp)
}

func (c *BankSoalController) UpdateBankSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateBankSoalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateBankSoal(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update bank soal successfully", resp)
}

func (c *BankSoalController) DeleteBankSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteBankSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete bank soal successfully", nil)
}

func (c *BankSoalController) RestoreBankSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreBankSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore bank soal successfully", nil)
}
```

#### Step 4.2: Setup Routes
**File**: `internal/modules/bank_soal/routes/bank_soal_routes.go`

```go
package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/bank_soal/controller"
	"backend/internal/modules/bank_soal/repository"
	"backend/internal/modules/bank_soal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupBankSoalRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewBankSoalRepository(db)
	svc := service.NewBankSoalService(repo)
	ctrl := controller.NewBankSoalController(svc)

	api := app.Group("/api")
	bankSoal := api.Group("/bank-soal")

	bankSoal.Post("/", middleware.JWTAuth(), ctrl.CreateBankSoal)
	bankSoal.Get("/", ctrl.GetAllBankSoal)
	bankSoal.Get("/mapel/:mapel_id", ctrl.GetBankSoalByMapel)
	bankSoal.Get("/:id", ctrl.GetBankSoalByID)
	bankSoal.Put("/:id", middleware.JWTAuth(), ctrl.UpdateBankSoal)
	bankSoal.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteBankSoal)
	bankSoal.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreBankSoal)
}
```

#### Step 4.3: Register Routes di Main
**File**: `cmd/server/main.go`

Tambahkan import dan setup route:

```go
import (
	banksoalroutes "backend/internal/modules/bank_soal/routes"
	// ... routes lainnya
)

func setupRoutes(app *fiber.App) {
	// ... routes lainnya
	banksoalroutes.SetupBankSoalRoutes(app, database.DB)
}
```

---

### **FASE 5: Seeder**

#### Step 5.1: Buat Seeder File
**File**: `internal/database/seeders/bank_soal_seeder.go`

```go
package seeders

import (
	"backend/internal/modules/bank_soal/model"
	"time"

	"gorm.io/gorm"
)

func SeedBankSoal(db *gorm.DB) error {
	// First, get sample mapel IDs (assuming mapel seeder already ran)
	var mapelIDs []string
	if err := db.Model(&struct{}{}).
		Table("mapel").
		Where("deleted_at IS NULL").
		Limit(5).
		Pluck("id", &mapelIDs).Error; err != nil {
		return err
	}

	if len(mapelIDs) == 0 {
		return nil // Skip if no mapel data
	}

	bankSoals := []model.BankSoal{
		{
			NamaBankSoal: "Bank Soal Matematika Dasar",
			IdMapel:      mapelIDs[0],
			JmlSoal:      50,
			Deskripsi:    "Kumpulan soal matematika level dasar",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			NamaBankSoal: "Bank Soal Matematika Lanjutan",
			IdMapel:      mapelIDs[0],
			JmlSoal:      75,
			Deskripsi:    "Kumpulan soal matematika level lanjutan",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			NamaBankSoal: "Bank Soal Bahasa Indonesia Umum",
			IdMapel:      mapelIDs[1],
			JmlSoal:      40,
			Deskripsi:    "Soal umum bahasa Indonesia",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			NamaBankSoal: "Bank Soal Grammar Bahasa Inggris",
			IdMapel:      mapelIDs[2],
			JmlSoal:      60,
			Deskripsi:    "Soal grammar bahasa Inggris",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			NamaBankSoal: "Bank Soal IPA Fisika",
			IdMapel:      mapelIDs[3],
			JmlSoal:      45,
			Deskripsi:    "Soal fisika IPA",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	return db.CreateInBatches(bankSoals, 100).Error
}
```

#### Step 5.2: Register Seeder
**File**: `internal/database/seed.go`

Tambahkan ke `RunSeeders` function:

```go
func RunSeeders(db *gorm.DB) error {
	if err := seeders.SeedMapel(db); err != nil {
		return fmt.Errorf("failed to seed mapel: %w", err)
	}

	if err := seeders.SeedBankSoal(db); err != nil {
		return fmt.Errorf("failed to seed bank_soal: %w", err)
	}

	return nil
}
```

#### Step 5.3: Jalankan Seeder

```bash
go run ./cmd/seed/main.go
```

---

### **FASE 6: Testing**

#### Step 6.1: Test Migration
```bash
# Verify tabel dibuat
# Jalankan server dan check database
go run ./cmd/server/main.go
```

#### Step 6.2: Test Seeder
```bash
go run ./cmd/seed/main.go
# Expected: Seeder ran successfully!
```

#### Step 6.3: Test API Endpoints

**Create Bank Soal:**
```bash
curl -X POST http://localhost:3000/api/bank-soal \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "nama_bank_soal": "Bank Soal Seni Budaya",
    "id_mapel": "mapel-uuid-here",
    "jml_soal": 30,
    "deskripsi": "Soal seni budaya"
  }'
```

**Get All Bank Soal:**
```bash
curl http://localhost:3000/api/bank-soal?page=1&page_size=10
```

**Get Bank Soal by Mapel:**
```bash
curl http://localhost:3000/api/bank-soal/mapel/{mapel_id}?page=1&page_size=10
```

**Get By ID:**
```bash
curl http://localhost:3000/api/bank-soal/{id}
```

**Update Bank Soal:**
```bash
curl -X PUT http://localhost:3000/api/bank-soal/{id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "nama_bank_soal": "Updated Name",
    "id_mapel": "mapel-uuid",
    "jml_soal": 35,
    "deskripsi": "Updated deskripsi"
  }'
```

**Soft Delete:**
```bash
curl -X DELETE http://localhost:3000/api/bank-soal/{id} \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

**Restore:**
```bash
curl -X PATCH http://localhost:3000/api/bank-soal/{id}/restore \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

## 📁 Struktur File yang Dibuat

```
project-root/
├── internal/
│   └── modules/
│       └── bank_soal/
│           ├── model/
│           │   └── bank_soal_model.go
│           ├── dto/
│           │   └── bank_soal_dto.go
│           ├── repository/
│           │   └── bank_soal_repository.go
│           ├── service/
│           │   └── bank_soal_service.go
│           ├── controller/
│           │   └── bank_soal_controller.go
│           └── routes/
│               └── bank_soal_routes.go
│   └── database/
│       └── seeders/
│           └── bank_soal_seeder.go
└── cmd/server/main.go (update untuk register routes)
```

---

## ✅ Checklist Implementasi

- [ ] Directory structure dibuat
- [ ] Model BankSoal dibuat dengan soft delete support
- [ ] DTO untuk request/response dibuat
- [ ] Repository interface dan implementasi dibuat
- [ ] Service layer dibuat dengan business logic
- [ ] Controller dengan semua endpoints dibuat
- [ ] Routes di-register di main.go
- [ ] Migration file updated (model sudah handle via AutoMigrate)
- [ ] Seeder dibuat dan dijalankan
- [ ] API tested dengan curl/Postman
- [ ] Soft delete berfungsi dengan baik
- [ ] Relasi dengan mapel working
- [ ] GetByMapel filtering working
- [ ] Documentation (API docs) ditambahkan

---

## 🔧 Catatan Teknis

### Soft Delete Implementation
- Gunakan `*time.Time` untuk DeletedAt field
- Query otomatis exclude record dengan `deleted_at IS NOT NULL`
- Gunakan `WHERE deleted_at IS NULL` di semua SELECT queries

### Foreign Key
- IdMapel harus selalu valid reference ke mapel table
- Harus ada validasi saat create/update bahwa mapel_id exists

### Pagination
- Default page: 1, Default page_size: 10
- Gunakan formula: `offset = (page - 1) * pageSize`
- Hitung totalPage: `ceil(total / pageSize)`

### Additional Endpoint
- **GET /api/bank-soal/mapel/:mapel_id** - Filter by mapel ID
- Endpoint ini penting untuk filter soal per mata pelajaran

---

## 🚀 Next Steps

Setelah semua tahapan selesai:

1. **Validation Enhancement** - Add validasi foreign key existence
2. **Error Handling** - Improve error messages
3. **API Documentation** - Create BANK_SOAL_API.md
4. **Frontend Integration** - Connect dengan Vue.js frontend
5. **Query Optimization** - Add eager loading untuk mapel relation

---

## 📞 Support

Jika ada pertanyaan saat implementasi:
- Cek structure project yang sudah ada (mapel module sebagai reference)
- Follow naming convention & pattern yang sudah ada
- Refer ke `internal/modules/mapel` sebagai reference implementation
- Semua kode sudah dijelaskan step-by-step

