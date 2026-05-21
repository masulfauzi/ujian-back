# Issue: Implementasi Modul Jadwal (Exam Schedule Management)

## Deskripsi

Buat modul pengelolaan jadwal ujian. Jadwal mengatur kapan sebuah bank soal akan diujikan, dengan mencatat waktu mulai dan waktu selesai. Setiap jadwal berelasi ke satu bank soal.

---

## Konteks Proyek

Proyek ini menggunakan:
- **Bahasa:** Go (Golang)
- **Framework HTTP:** [Fiber v2](https://gofiber.io/)
- **ORM:** GORM dengan driver PostgreSQL
- **Struktur Modul:** `internal/modules/<nama_modul>/` berisi `model/`, `dto/`, `repository/`, `service/`, `controller/`, `routes/`
- **Module Go:** `backend` (lihat `go.mod`)

Sebelum mulai, pelajari contoh modul yang sudah ada: `internal/modules/bank_soal/` — pola implementasi jadwal harus identik.

---

## Skema Tabel

Nama tabel: `jadwal`

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID | Primary key, auto-generate via `gen_random_uuid()` |
| `id_bank_soal` | UUID | Foreign key ke tabel `bank_soal`, NOT NULL |
| `wkt_mulai` | timestamp | Waktu mulai ujian, NOT NULL |
| `wkt_selesai` | timestamp | Waktu selesai ujian, NOT NULL |
| `created_at` | timestamp | Auto-fill saat create |
| `updated_at` | timestamp | Auto-update saat update |
| `deleted_at` | timestamp (nullable) | Soft delete — null = aktif, terisi = terhapus |

---

## Tahapan Implementasi

### Tahap 1 — Buat Model

**File:** `internal/modules/jadwal/model/jadwal_model.go`

```go
package model

import (
	"time"
)

type Jadwal struct {
	ID          string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDBankSoal  string     `gorm:"type:uuid;not null;index" json:"id_bank_soal"`
	WktMulai    time.Time  `gorm:"type:timestamp;not null" json:"wkt_mulai"`
	WktSelesai  time.Time  `gorm:"type:timestamp;not null" json:"wkt_selesai"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

func (Jadwal) TableName() string {
	return "jadwal"
}
```

> **Catatan:**
> - `DeletedAt` bertipe `*time.Time` (pointer) bukan `gorm.DeletedAt` — ini mengikuti pola yang sama dengan modul `bank_soal`.
> - `WktMulai` dan `WktSelesai` bertipe `time.Time` karena menyimpan tanggal + waktu (datetime).

---

### Tahap 2 — Daftarkan Model ke Migration

**File:** `internal/database/migrate.go`

Tambahkan import model jadwal dan daftarkan di dalam fungsi `RunMigrations`.

```go
package database

import (
	banksoalmodel "backend/internal/modules/bank_soal/model"
	jadwalmodel   "backend/internal/modules/jadwal/model"   // <-- tambahkan ini
	jurusanmodel  "backend/internal/modules/jurusan/model"
	kelasmodel    "backend/internal/modules/kelas/model"
	mapelmodel    "backend/internal/modules/mapel/model"
	soalmodel     "backend/internal/modules/soal/model"
	usermodel     "backend/internal/modules/user/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	db.Exec("DROP INDEX IF EXISTS idx_jurusans_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS uni_jurusan_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS idx_jurusan_nama_jurusan")

	return db.AutoMigrate(
		&usermodel.User{},
		&mapelmodel.Mapel{},
		&banksoalmodel.BankSoal{},   // <-- jadwal harus didaftarkan SETELAH bank_soal
		&soalmodel.Soal{},
		&jurusanmodel.Jurusan{},
		&kelasmodel.Kelas{},
		&jadwalmodel.Jadwal{},        // <-- tambahkan di sini
	)
}
```

> **Penting:** `Jadwal` harus didaftarkan **setelah** `BankSoal` karena tabel `jadwal` memiliki foreign key ke tabel `bank_soal`.

---

### Tahap 3 — Buat DTO

**File:** `internal/modules/jadwal/dto/jadwal_dto.go`

DTO mendefinisikan struktur data untuk request (input) dan response (output) API.

```go
package dto

type CreateJadwalRequest struct {
	IDBankSoal string `json:"id_bank_soal" validate:"required"`
	WktMulai   string `json:"wkt_mulai" validate:"required"`   // format: "2006-01-02 15:04:05"
	WktSelesai string `json:"wkt_selesai" validate:"required"` // format: "2006-01-02 15:04:05"
}

type UpdateJadwalRequest struct {
	IDBankSoal string `json:"id_bank_soal" validate:"required"`
	WktMulai   string `json:"wkt_mulai" validate:"required"`
	WktSelesai string `json:"wkt_selesai" validate:"required"`
}

type JadwalResponse struct {
	ID           string `json:"id"`
	IDBankSoal   string `json:"id_bank_soal"`
	NamaBankSoal string `json:"nama_bank_soal"`
	WktMulai     string `json:"wkt_mulai"`
	WktSelesai   string `json:"wkt_selesai"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type JadwalListResponse struct {
	Data      []JadwalResponse `json:"data"`
	Total     int64            `json:"total"`
	Page      int              `json:"page"`
	PageSize  int              `json:"page_size"`
	TotalPage int              `json:"total_page"`
}
```

> **Catatan:** `WktMulai` dan `WktSelesai` di request dikirim sebagai `string` lalu di-parse menjadi `time.Time` di service. Di response, dikembalikan sebagai `string` dengan format `"2006-01-02 15:04:05"`.

---

### Tahap 4 — Buat Repository

**File:** `internal/modules/jadwal/repository/jadwal_repository.go`

Repository menangani semua interaksi langsung dengan database. Query menggunakan JOIN ke tabel `bank_soal` agar `nama_bank_soal` ikut tampil di response.

```go
package repository

import (
	"backend/internal/modules/jadwal/model"
	"time"

	"gorm.io/gorm"
)

type JadwalWithBankSoal struct {
	ID           string  `gorm:"column:id"`
	IDBankSoal   string  `gorm:"column:id_bank_soal"`
	NamaBankSoal string  `gorm:"column:nama_bank_soal"`
	WktMulai     string  `gorm:"column:wkt_mulai"`
	WktSelesai   string  `gorm:"column:wkt_selesai"`
	CreatedAt    string  `gorm:"column:created_at"`
	UpdatedAt    string  `gorm:"column:updated_at"`
}

type JadwalRepository interface {
	Create(jadwal *model.Jadwal) error
	GetByID(id string) (*model.Jadwal, error)
	GetByIDWithBankSoal(id string) (*JadwalWithBankSoal, error)
	GetAllWithBankSoal(page, pageSize int) ([]JadwalWithBankSoal, int64, error)
	GetByBankSoalID(bankSoalID string, page, pageSize int) ([]JadwalWithBankSoal, int64, error)
	Update(jadwal *model.Jadwal) error
	Delete(id string) error
	Restore(id string) error
}

type jadwalRepository struct {
	db *gorm.DB
}

func NewJadwalRepository(db *gorm.DB) JadwalRepository {
	return &jadwalRepository{db: db}
}

func (r *jadwalRepository) Create(jadwal *model.Jadwal) error {
	return r.db.Create(jadwal).Error
}

func (r *jadwalRepository) GetByID(id string) (*model.Jadwal, error) {
	var jadwal model.Jadwal
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&jadwal).Error
	if err != nil {
		return nil, err
	}
	return &jadwal, nil
}

func (r *jadwalRepository) GetByIDWithBankSoal(id string) (*JadwalWithBankSoal, error) {
	var jadwal JadwalWithBankSoal
	err := r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.wkt_mulai, jadwal.wkt_selesai, jadwal.created_at, jadwal.updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id = ? AND jadwal.deleted_at IS NULL", id).
		First(&jadwal).Error
	if err != nil {
		return nil, err
	}
	return &jadwal, nil
}

func (r *jadwalRepository) GetAllWithBankSoal(page, pageSize int) ([]JadwalWithBankSoal, int64, error) {
	var jadwalList []JadwalWithBankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Table("jadwal").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.wkt_mulai, jadwal.wkt_selesai, jadwal.created_at, jadwal.updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Scan(&jadwalList).Error

	return jadwalList, total, err
}

func (r *jadwalRepository) GetByBankSoalID(bankSoalID string, page, pageSize int) ([]JadwalWithBankSoal, int64, error) {
	var jadwalList []JadwalWithBankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Table("jadwal").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id_bank_soal = ? AND jadwal.deleted_at IS NULL", bankSoalID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.wkt_mulai, jadwal.wkt_selesai, jadwal.created_at, jadwal.updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id_bank_soal = ? AND jadwal.deleted_at IS NULL", bankSoalID).
		Offset(offset).
		Limit(pageSize).
		Scan(&jadwalList).Error

	return jadwalList, total, err
}

func (r *jadwalRepository) Update(jadwal *model.Jadwal) error {
	return r.db.Save(jadwal).Error
}

func (r *jadwalRepository) Delete(id string) error {
	now := time.Now()
	return r.db.Model(&model.Jadwal{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *jadwalRepository) Restore(id string) error {
	return r.db.Model(&model.Jadwal{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
}
```

---

### Tahap 5 — Buat Service

**File:** `internal/modules/jadwal/service/jadwal_service.go`

Service berisi business logic, termasuk parsing string datetime menjadi `time.Time`.

```go
package service

import (
	"errors"
	"math"
	"time"

	"backend/internal/constants"
	"backend/internal/modules/jadwal/dto"
	"backend/internal/modules/jadwal/model"
	"backend/internal/modules/jadwal/repository"

	"gorm.io/gorm"
)

const timeLayout = "2006-01-02 15:04:05"

type JadwalService interface {
	CreateJadwal(req *dto.CreateJadwalRequest) (*dto.JadwalResponse, error)
	GetJadwalByID(id string) (*dto.JadwalResponse, error)
	GetAllJadwal(page, pageSize int) (*dto.JadwalListResponse, error)
	GetJadwalByBankSoal(bankSoalID string, page, pageSize int) (*dto.JadwalListResponse, error)
	UpdateJadwal(id string, req *dto.UpdateJadwalRequest) (*dto.JadwalResponse, error)
	DeleteJadwal(id string) error
	RestoreJadwal(id string) error
}

type jadwalService struct {
	repo repository.JadwalRepository
}

func NewJadwalService(repo repository.JadwalRepository) JadwalService {
	return &jadwalService{repo: repo}
}

func (s *jadwalService) CreateJadwal(req *dto.CreateJadwalRequest) (*dto.JadwalResponse, error) {
	wktMulai, err := time.Parse(timeLayout, req.WktMulai)
	if err != nil {
		return nil, errors.New("format wkt_mulai tidak valid, gunakan: 2006-01-02 15:04:05")
	}

	wktSelesai, err := time.Parse(timeLayout, req.WktSelesai)
	if err != nil {
		return nil, errors.New("format wkt_selesai tidak valid, gunakan: 2006-01-02 15:04:05")
	}

	if !wktSelesai.After(wktMulai) {
		return nil, errors.New("wkt_selesai harus setelah wkt_mulai")
	}

	jadwal := &model.Jadwal{
		IDBankSoal: req.IDBankSoal,
		WktMulai:   wktMulai,
		WktSelesai: wktSelesai,
	}

	if err := s.repo.Create(jadwal); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByIDWithBankSoal(jadwal.ID)
	if err != nil {
		return nil, err
	}

	return joinedToResponse(created), nil
}

func (s *jadwalService) GetJadwalByID(id string) (*dto.JadwalResponse, error) {
	jadwal, err := s.repo.GetByIDWithBankSoal(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return joinedToResponse(jadwal), nil
}

func (s *jadwalService) GetAllJadwal(page, pageSize int) (*dto.JadwalListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	jadwalList, total, err := s.repo.GetAllWithBankSoal(page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := []dto.JadwalResponse{}
	for _, j := range jadwalList {
		responses = append(responses, *joinedToResponse(&j))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.JadwalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *jadwalService) GetJadwalByBankSoal(bankSoalID string, page, pageSize int) (*dto.JadwalListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	jadwalList, total, err := s.repo.GetByBankSoalID(bankSoalID, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := []dto.JadwalResponse{}
	for _, j := range jadwalList {
		responses = append(responses, *joinedToResponse(&j))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.JadwalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *jadwalService) UpdateJadwal(id string, req *dto.UpdateJadwalRequest) (*dto.JadwalResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	wktMulai, err := time.Parse(timeLayout, req.WktMulai)
	if err != nil {
		return nil, errors.New("format wkt_mulai tidak valid, gunakan: 2006-01-02 15:04:05")
	}

	wktSelesai, err := time.Parse(timeLayout, req.WktSelesai)
	if err != nil {
		return nil, errors.New("format wkt_selesai tidak valid, gunakan: 2006-01-02 15:04:05")
	}

	if !wktSelesai.After(wktMulai) {
		return nil, errors.New("wkt_selesai harus setelah wkt_mulai")
	}

	existing.IDBankSoal = req.IDBankSoal
	existing.WktMulai   = wktMulai
	existing.WktSelesai = wktSelesai

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDWithBankSoal(id)
	if err != nil {
		return nil, err
	}

	return joinedToResponse(updated), nil
}

func (s *jadwalService) DeleteJadwal(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(id)
}

func (s *jadwalService) RestoreJadwal(id string) error {
	return s.repo.Restore(id)
}

func joinedToResponse(j *repository.JadwalWithBankSoal) *dto.JadwalResponse {
	return &dto.JadwalResponse{
		ID:           j.ID,
		IDBankSoal:   j.IDBankSoal,
		NamaBankSoal: j.NamaBankSoal,
		WktMulai:     j.WktMulai,
		WktSelesai:   j.WktSelesai,
		CreatedAt:    j.CreatedAt,
		UpdatedAt:    j.UpdatedAt,
	}
}
```

> **Aturan bisnis penting:**
> - `wkt_selesai` harus selalu **setelah** `wkt_mulai` — validasi dilakukan di service.
> - Format datetime yang diterima dari client: `"2006-01-02 15:04:05"` (contoh: `"2025-08-01 08:00:00"`).

---

### Tahap 6 — Buat Controller

**File:** `internal/modules/jadwal/controller/jadwal_controller.go`

```go
package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jadwal/dto"
	"backend/internal/modules/jadwal/service"

	"github.com/gofiber/fiber/v2"
)

type JadwalController struct {
	service service.JadwalService
}

func NewJadwalController(service service.JadwalService) *JadwalController {
	return &JadwalController{service: service}
}

func (c *JadwalController) CreateJadwal(ctx *fiber.Ctx) error {
	var req dto.CreateJadwalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJadwal(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jadwal successfully", resp)
}

func (c *JadwalController) GetAllJadwal(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetAllJadwal(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jadwal successfully", resp)
}

func (c *JadwalController) GetJadwalByBankSoal(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetJadwalByBankSoal(bankSoalID, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jadwal by bank soal successfully", resp)
}

func (c *JadwalController) GetJadwalByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJadwalByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jadwal successfully", resp)
}

func (c *JadwalController) UpdateJadwal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJadwalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJadwal(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jadwal successfully", resp)
}

func (c *JadwalController) DeleteJadwal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJadwal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jadwal successfully", nil)
}

func (c *JadwalController) RestoreJadwal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreJadwal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore jadwal successfully", nil)
}
```

---

### Tahap 7 — Buat Routes

**File:** `internal/modules/jadwal/routes/jadwal_routes.go`

```go
package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/jadwal/controller"
	"backend/internal/modules/jadwal/repository"
	"backend/internal/modules/jadwal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupJadwalRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewJadwalRepository(db)
	svc := service.NewJadwalService(repo)
	ctrl := controller.NewJadwalController(svc)

	api := app.Group("/api")
	jadwal := api.Group("/jadwal")

	jadwal.Post("/", middleware.JWTAuth(), ctrl.CreateJadwal)
	jadwal.Get("/", ctrl.GetAllJadwal)
	jadwal.Get("/bank-soal/:bank_soal_id", ctrl.GetJadwalByBankSoal)
	jadwal.Get("/:id", ctrl.GetJadwalByID)
	jadwal.Put("/:id", middleware.JWTAuth(), ctrl.UpdateJadwal)
	jadwal.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteJadwal)
	jadwal.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreJadwal)
}
```

---

### Tahap 8 — Daftarkan Routes di Server

**File:** `cmd/server/main.go`

Tambahkan import dan panggilan `SetupJadwalRoutes` di fungsi `setupRoutes`.

```go
// Tambahkan di bagian import:
jadwalroutes "backend/internal/modules/jadwal/routes"

// Tambahkan di dalam fungsi setupRoutes, setelah kelasroutes:
jadwalroutes.SetupJadwalRoutes(app, database.DB)
```

Setelah perubahan, blok import dan fungsi `setupRoutes` akan terlihat seperti ini:

```go
import (
	// ... import yang sudah ada ...
	jadwalroutes  "backend/internal/modules/jadwal/routes"  // <-- tambahkan
	kelasroutes   "backend/internal/modules/kelas/routes"
)

func setupRoutes(app *fiber.App) {
	// ... routes yang sudah ada ...
	kelasroutes.SetupKelasRoutes(app, database.DB)
	jadwalroutes.SetupJadwalRoutes(app, database.DB)  // <-- tambahkan
}
```

---

### Tahap 9 — Build dan Test

```bash
# Pastikan tidak ada error compile
go build ./...

# Jalankan server (migration otomatis berjalan saat startup)
go run cmd/server/main.go

# Test endpoint
curl http://localhost:3000/api/jadwal
```

---

## Ringkasan File yang Perlu Dibuat / Diubah

### File Baru (buat dari awal):
```
internal/modules/jadwal/model/jadwal_model.go
internal/modules/jadwal/dto/jadwal_dto.go
internal/modules/jadwal/repository/jadwal_repository.go
internal/modules/jadwal/service/jadwal_service.go
internal/modules/jadwal/controller/jadwal_controller.go
internal/modules/jadwal/routes/jadwal_routes.go
```

### File yang Dimodifikasi (tambahkan beberapa baris):
```
internal/database/migrate.go   — tambahkan import & daftarkan &jadwalmodel.Jadwal{}
cmd/server/main.go             — tambahkan import & panggilan SetupJadwalRoutes
```

---

## Endpoint API yang Dihasilkan

| Method | URL | Auth | Deskripsi |
|--------|-----|------|-----------|
| `GET` | `/api/jadwal` | Tidak | Daftar semua jadwal (pagination) |
| `GET` | `/api/jadwal/:id` | Tidak | Detail satu jadwal berdasarkan ID |
| `GET` | `/api/jadwal/bank-soal/:bank_soal_id` | Tidak | Daftar jadwal berdasarkan bank soal |
| `POST` | `/api/jadwal` | JWT | Buat jadwal baru |
| `PUT` | `/api/jadwal/:id` | JWT | Update data jadwal |
| `DELETE` | `/api/jadwal/:id` | JWT | Soft delete jadwal |
| `PATCH` | `/api/jadwal/:id/restore` | JWT | Restore jadwal yang dihapus |

### Query Parameter untuk GET /api/jadwal:
| Parameter | Tipe | Default | Keterangan |
|-----------|------|---------|------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |

---

## Aturan Bisnis Penting

1. **Format datetime:** `wkt_mulai` dan `wkt_selesai` dikirim dalam format `"2006-01-02 15:04:05"` (contoh: `"2025-08-01 08:00:00"`).

2. **Validasi urutan waktu:** `wkt_selesai` harus selalu **setelah** `wkt_mulai`. Jika tidak, API mengembalikan error 400.

3. **Soft delete:** Data tidak benar-benar dihapus dari database. Field `deleted_at` diisi timestamp saat dihapus. Data yang sudah dihapus tidak muncul di GET.

4. **Relasi ke bank_soal:** Setiap jadwal wajib memiliki `id_bank_soal` yang valid. Response API menyertakan `nama_bank_soal` hasil JOIN dari tabel `bank_soal`.

5. **Endpoint filter by bank soal:** Gunakan `GET /api/jadwal/bank-soal/:bank_soal_id` untuk melihat semua jadwal milik satu bank soal tertentu.

---

## Referensi

Semua pola implementasi mengacu pada modul `bank_soal` yang sudah ada:
- Model: `internal/modules/bank_soal/model/bank_soal_model.go`
- DTO: `internal/modules/bank_soal/dto/bank_soal_dto.go`
- Repository: `internal/modules/bank_soal/repository/bank_soal_repository.go`
- Service: `internal/modules/bank_soal/service/bank_soal_service.go`
- Controller: `internal/modules/bank_soal/controller/bank_soal_controller.go`
- Routes: `internal/modules/bank_soal/routes/bank_soal_routes.go`
