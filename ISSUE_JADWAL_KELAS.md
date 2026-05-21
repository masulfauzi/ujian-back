# Issue: Implementasi Modul Jadwal Kelas (Schedule-Class Assignment)

## Deskripsi

Buat modul pengelolaan penugasan kelas ke jadwal ujian. Tabel `jadwal_kelas` adalah tabel pivot yang menghubungkan `jadwal` dengan `kelas` — artinya menentukan kelas mana saja yang mengikuti jadwal ujian tertentu.

**Perbedaan penting dibanding modul lain:**
1. **Hard delete** — data benar-benar dihapus permanen dari database, tidak ada soft delete.
2. **Tidak ada kolom `deleted_at`** — hanya ada `created_at` dan `updated_at`.
3. **Tidak ada endpoint restore** — data yang dihapus tidak bisa dikembalikan.
4. **Unique constraint** — kombinasi `(id_jadwal, id_kelas)` harus unik — satu kelas tidak boleh didaftarkan dua kali ke jadwal yang sama.

---

## Konteks Proyek

Proyek ini menggunakan:
- **Bahasa:** Go (Golang)
- **Framework HTTP:** [Fiber v2](https://gofiber.io/)
- **ORM:** GORM dengan driver PostgreSQL
- **Struktur Modul:** `internal/modules/<nama_modul>/` berisi `model/`, `dto/`, `repository/`, `service/`, `controller/`, `routes/`
- **Module Go:** `backend` (lihat `go.mod`)

Sebelum mulai, pelajari contoh modul yang sudah ada: `internal/modules/kelas/` — pola umum implementasi harus mengikuti modul tersebut, **kecuali** bagian yang berkaitan dengan soft delete.

---

## Skema Tabel

Nama tabel: `jadwal_kelas`

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID | Primary key, auto-generate via `gen_random_uuid()` |
| `id_jadwal` | UUID | Foreign key ke tabel `jadwal`, NOT NULL |
| `id_kelas` | UUID | Foreign key ke tabel `kelas`, NOT NULL |
| `created_at` | timestamp | Auto-fill saat create |
| `updated_at` | timestamp | Auto-update saat update |

**Constraint tambahan:**
- `UNIQUE(id_jadwal, id_kelas)` — satu kelas tidak boleh didaftarkan lebih dari satu kali ke jadwal yang sama.

---

## Tahapan Implementasi

### Tahap 1 — Buat Model

**File:** `internal/modules/jadwal_kelas/model/jadwal_kelas_model.go`

Model adalah representasi struct Go dari tabel database. Perhatikan: **tidak ada field `DeletedAt`** karena menggunakan hard delete.

```go
package model

import (
	"time"
)

type JadwalKelas struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDJadwal  string    `gorm:"type:uuid;not null;index" json:"id_jadwal"`
	IDKelas   string    `gorm:"type:uuid;not null;index" json:"id_kelas"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (JadwalKelas) TableName() string {
	return "jadwal_kelas"
}
```

> **Catatan penting:**
> - Tidak ada field `DeletedAt` sama sekali — ini yang membuat hard delete bekerja secara otomatis dengan GORM.
> - GORM hanya melakukan soft delete jika model memiliki field `DeletedAt`. Karena tidak ada, `db.Delete()` akan langsung menghapus baris dari database.
> - Unique constraint `(id_jadwal, id_kelas)` **tidak** dideklarasikan di sini, melainkan di migration (lihat Tahap 2).

---

### Tahap 2 — Daftarkan Model ke Migration

**File:** `internal/database/migrate.go`

Tambahkan import model dan daftarkan di dalam fungsi `RunMigrations`. Juga tambahkan kode untuk membuat unique constraint setelah AutoMigrate.

```go
package database

import (
	banksoalmodel    "backend/internal/modules/bank_soal/model"
	jadwalmodel      "backend/internal/modules/jadwal/model"
	jadwalkelasmodel "backend/internal/modules/jadwal_kelas/model"  // <-- tambahkan
	jurusanmodel     "backend/internal/modules/jurusan/model"
	kelasmodel       "backend/internal/modules/kelas/model"
	mapelmodel       "backend/internal/modules/mapel/model"
	soalmodel        "backend/internal/modules/soal/model"
	usermodel        "backend/internal/modules/user/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	db.Exec("DROP INDEX IF EXISTS idx_jurusans_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS uni_jurusan_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS idx_jurusan_nama_jurusan")

	if err := db.AutoMigrate(
		&usermodel.User{},
		&mapelmodel.Mapel{},
		&banksoalmodel.BankSoal{},
		&soalmodel.Soal{},
		&jurusanmodel.Jurusan{},
		&kelasmodel.Kelas{},
		&jadwalmodel.Jadwal{},
		&jadwalkelasmodel.JadwalKelas{},  // <-- tambahkan (setelah jadwal dan kelas)
	); err != nil {
		return err
	}

	// Buat unique constraint untuk mencegah duplikasi assignment kelas ke jadwal
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_jadwal_kelas_unique ON jadwal_kelas(id_jadwal, id_kelas)")

	return nil
}
```

> **Penting:**
> - `JadwalKelas` harus didaftarkan **setelah** `Jadwal` dan `Kelas` karena memiliki FK ke kedua tabel tersebut.
> - Unique index dibuat manual via `db.Exec` karena GORM AutoMigrate tidak mendukung composite unique index secara deklaratif di struct.
> - `IF NOT EXISTS` memastikan perintah aman dijalankan berulang kali saat server restart.

---

### Tahap 3 — Buat DTO

**File:** `internal/modules/jadwal_kelas/dto/jadwal_kelas_dto.go`

DTO mendefinisikan struktur data untuk request (input) dan response (output) API. Response menyertakan data dari tabel `jadwal` dan `kelas` via JOIN.

```go
package dto

type CreateJadwalKelasRequest struct {
	IDJadwal string `json:"id_jadwal" validate:"required"`
	IDKelas  string `json:"id_kelas" validate:"required"`
}

type UpdateJadwalKelasRequest struct {
	IDJadwal string `json:"id_jadwal" validate:"required"`
	IDKelas  string `json:"id_kelas" validate:"required"`
}

type JadwalKelasResponse struct {
	ID           string `json:"id"`
	IDJadwal     string `json:"id_jadwal"`
	IDKelas      string `json:"id_kelas"`
	NamaKelas    string `json:"nama_kelas"`
	NamaBankSoal string `json:"nama_bank_soal"`
	WktMulai     string `json:"wkt_mulai"`
	WktSelesai   string `json:"wkt_selesai"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type JadwalKelasListResponse struct {
	Data      []JadwalKelasResponse `json:"data"`
	Total     int64                 `json:"total"`
	Page      int                   `json:"page"`
	PageSize  int                   `json:"page_size"`
	TotalPage int                   `json:"total_page"`
}
```

> **Catatan:**
> - `NamaKelas` berasal dari JOIN ke tabel `kelas`.
> - `NamaBankSoal`, `WktMulai`, `WktSelesai` berasal dari JOIN ke tabel `jadwal` → `bank_soal`.
> - Tidak ada field restore di DTO karena tidak ada endpoint restore.

---

### Tahap 4 — Buat Repository

**File:** `internal/modules/jadwal_kelas/repository/jadwal_kelas_repository.go`

Repository menangani semua interaksi langsung dengan database. Query menggunakan JOIN ke tabel `jadwal`, `bank_soal`, dan `kelas` agar data relasi ikut tampil di response.

```go
package repository

import (
	"backend/internal/modules/jadwal_kelas/model"

	"gorm.io/gorm"
)

type JadwalKelasWithDetail struct {
	ID           string `gorm:"column:id"`
	IDJadwal     string `gorm:"column:id_jadwal"`
	IDKelas      string `gorm:"column:id_kelas"`
	NamaKelas    string `gorm:"column:nama_kelas"`
	NamaBankSoal string `gorm:"column:nama_bank_soal"`
	WktMulai     string `gorm:"column:wkt_mulai"`
	WktSelesai   string `gorm:"column:wkt_selesai"`
	CreatedAt    string `gorm:"column:created_at"`
	UpdatedAt    string `gorm:"column:updated_at"`
}

type JadwalKelasRepository interface {
	Create(jadwalKelas *model.JadwalKelas) error
	GetByID(id string) (*model.JadwalKelas, error)
	GetByIDWithDetail(id string) (*JadwalKelasWithDetail, error)
	GetAllWithDetail(page, pageSize int, idJadwal string, idKelas string) ([]JadwalKelasWithDetail, int64, error)
	CheckDuplicate(idJadwal, idKelas string) (bool, error)
	Update(jadwalKelas *model.JadwalKelas) error
	Delete(id string) error
}

type jadwalKelasRepository struct {
	db *gorm.DB
}

func NewJadwalKelasRepository(db *gorm.DB) JadwalKelasRepository {
	return &jadwalKelasRepository{db: db}
}

func (r *jadwalKelasRepository) Create(jadwalKelas *model.JadwalKelas) error {
	return r.db.Create(jadwalKelas).Error
}

func (r *jadwalKelasRepository) GetByID(id string) (*model.JadwalKelas, error) {
	var jadwalKelas model.JadwalKelas
	err := r.db.Where("id = ?", id).First(&jadwalKelas).Error
	if err != nil {
		return nil, err
	}
	return &jadwalKelas, nil
}

func (r *jadwalKelasRepository) GetByIDWithDetail(id string) (*JadwalKelasWithDetail, error) {
	var result JadwalKelasWithDetail
	err := r.db.
		Table("jadwal_kelas").
		Select(`
			jadwal_kelas.id,
			jadwal_kelas.id_jadwal,
			jadwal_kelas.id_kelas,
			kelas.nama_kelas,
			bank_soal.nama_bank_soal,
			jadwal.wkt_mulai,
			jadwal.wkt_selesai,
			jadwal_kelas.created_at,
			jadwal_kelas.updated_at
		`).
		Joins("INNER JOIN jadwal ON jadwal_kelas.id_jadwal = jadwal.id").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal_kelas.id = ?", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *jadwalKelasRepository) GetAllWithDetail(page, pageSize int, idJadwal string, idKelas string) ([]JadwalKelasWithDetail, int64, error) {
	var results []JadwalKelasWithDetail
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	countQuery := r.db.Table("jadwal_kelas").
		Joins("INNER JOIN jadwal ON jadwal_kelas.id_jadwal = jadwal.id").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id")

	if idJadwal != "" {
		countQuery = countQuery.Where("jadwal_kelas.id_jadwal = ?", idJadwal)
	}
	if idKelas != "" {
		countQuery = countQuery.Where("jadwal_kelas.id_kelas = ?", idKelas)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.
		Table("jadwal_kelas").
		Select(`
			jadwal_kelas.id,
			jadwal_kelas.id_jadwal,
			jadwal_kelas.id_kelas,
			kelas.nama_kelas,
			bank_soal.nama_bank_soal,
			jadwal.wkt_mulai,
			jadwal.wkt_selesai,
			jadwal_kelas.created_at,
			jadwal_kelas.updated_at
		`).
		Joins("INNER JOIN jadwal ON jadwal_kelas.id_jadwal = jadwal.id").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id")

	if idJadwal != "" {
		query = query.Where("jadwal_kelas.id_jadwal = ?", idJadwal)
	}
	if idKelas != "" {
		query = query.Where("jadwal_kelas.id_kelas = ?", idKelas)
	}

	err := query.Offset(offset).Limit(pageSize).Scan(&results).Error
	return results, total, err
}

func (r *jadwalKelasRepository) CheckDuplicate(idJadwal, idKelas string) (bool, error) {
	var count int64
	err := r.db.Model(&model.JadwalKelas{}).
		Where("id_jadwal = ? AND id_kelas = ?", idJadwal, idKelas).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *jadwalKelasRepository) Update(jadwalKelas *model.JadwalKelas) error {
	return r.db.Save(jadwalKelas).Error
}

// Delete melakukan hard delete — baris dihapus permanen dari database
func (r *jadwalKelasRepository) Delete(id string) error {
	return r.db.Delete(&model.JadwalKelas{}, "id = ?", id).Error
}
```

> **Catatan penting tentang hard delete:**
> - `r.db.Delete(&model.JadwalKelas{}, "id = ?", id)` melakukan **hard delete** karena model `JadwalKelas` tidak memiliki field `DeletedAt`.
> - GORM otomatis menentukan jenis delete berdasarkan ada/tidaknya field `DeletedAt` di model. Tidak perlu `Unscoped()`.
> - Method `CheckDuplicate` digunakan oleh service untuk mencegah assignment duplikat sebelum create maupun update.

---

### Tahap 5 — Buat Service

**File:** `internal/modules/jadwal_kelas/service/jadwal_kelas_service.go`

Service berisi business logic, termasuk pengecekan duplikat assignment.

```go
package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/jadwal_kelas/dto"
	"backend/internal/modules/jadwal_kelas/model"
	"backend/internal/modules/jadwal_kelas/repository"

	"gorm.io/gorm"
)

type JadwalKelasService interface {
	CreateJadwalKelas(req *dto.CreateJadwalKelasRequest) (*dto.JadwalKelasResponse, error)
	GetJadwalKelasByID(id string) (*dto.JadwalKelasResponse, error)
	GetAllJadwalKelas(page, pageSize int, idJadwal string, idKelas string) (*dto.JadwalKelasListResponse, error)
	UpdateJadwalKelas(id string, req *dto.UpdateJadwalKelasRequest) (*dto.JadwalKelasResponse, error)
	DeleteJadwalKelas(id string) error
}

type jadwalKelasService struct {
	repo repository.JadwalKelasRepository
}

func NewJadwalKelasService(repo repository.JadwalKelasRepository) JadwalKelasService {
	return &jadwalKelasService{repo: repo}
}

func (s *jadwalKelasService) CreateJadwalKelas(req *dto.CreateJadwalKelasRequest) (*dto.JadwalKelasResponse, error) {
	// Cek duplikat: satu kelas tidak boleh didaftarkan dua kali ke jadwal yang sama
	exists, err := s.repo.CheckDuplicate(req.IDJadwal, req.IDKelas)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("kelas ini sudah terdaftar di jadwal tersebut")
	}

	jadwalKelas := &model.JadwalKelas{
		IDJadwal: req.IDJadwal,
		IDKelas:  req.IDKelas,
	}

	if err := s.repo.Create(jadwalKelas); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByIDWithDetail(jadwalKelas.ID)
	if err != nil {
		return nil, err
	}

	return detailToResponse(created), nil
}

func (s *jadwalKelasService) GetJadwalKelasByID(id string) (*dto.JadwalKelasResponse, error) {
	result, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return detailToResponse(result), nil
}

func (s *jadwalKelasService) GetAllJadwalKelas(page, pageSize int, idJadwal string, idKelas string) (*dto.JadwalKelasListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	results, total, err := s.repo.GetAllWithDetail(page, pageSize, idJadwal, idKelas)
	if err != nil {
		return nil, err
	}

	responses := []dto.JadwalKelasResponse{}
	for _, r := range results {
		responses = append(responses, *detailToResponse(&r))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.JadwalKelasListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *jadwalKelasService) UpdateJadwalKelas(id string, req *dto.UpdateJadwalKelasRequest) (*dto.JadwalKelasResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	// Cek duplikat hanya jika ada perubahan pada id_jadwal atau id_kelas
	if req.IDJadwal != existing.IDJadwal || req.IDKelas != existing.IDKelas {
		exists, err := s.repo.CheckDuplicate(req.IDJadwal, req.IDKelas)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("kelas ini sudah terdaftar di jadwal tersebut")
		}
	}

	existing.IDJadwal = req.IDJadwal
	existing.IDKelas  = req.IDKelas

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		return nil, err
	}

	return detailToResponse(updated), nil
}

func (s *jadwalKelasService) DeleteJadwalKelas(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(id)
}

func detailToResponse(r *repository.JadwalKelasWithDetail) *dto.JadwalKelasResponse {
	return &dto.JadwalKelasResponse{
		ID:           r.ID,
		IDJadwal:     r.IDJadwal,
		IDKelas:      r.IDKelas,
		NamaKelas:    r.NamaKelas,
		NamaBankSoal: r.NamaBankSoal,
		WktMulai:     r.WktMulai,
		WktSelesai:   r.WktSelesai,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
```

> **Aturan bisnis penting:**
> - Pada saat create, cek duplikat dilakukan terhadap semua record (tidak ada kondisi `deleted_at IS NULL` karena tidak ada soft delete).
> - Pada saat update, cek duplikat dilakukan **hanya jika** ada perubahan nilai — jika `id_jadwal` dan `id_kelas` tidak berubah, tidak perlu cek duplikat.
> - Tidak ada method `RestoreJadwalKelas` karena hard delete.

---

### Tahap 6 — Buat Controller

**File:** `internal/modules/jadwal_kelas/controller/jadwal_kelas_controller.go`

Perhatikan: **tidak ada handler `RestoreJadwalKelas`**.

```go
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
```

---

### Tahap 7 — Buat Routes

**File:** `internal/modules/jadwal_kelas/routes/jadwal_kelas_routes.go`

Perhatikan: **tidak ada route restore**.

```go
package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/jadwal_kelas/controller"
	"backend/internal/modules/jadwal_kelas/repository"
	"backend/internal/modules/jadwal_kelas/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupJadwalKelasRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewJadwalKelasRepository(db)
	svc  := service.NewJadwalKelasService(repo)
	ctrl := controller.NewJadwalKelasController(svc)

	api       := app.Group("/api")
	jadwalKelas := api.Group("/jadwal-kelas")

	jadwalKelas.Post("/", middleware.JWTAuth(), ctrl.CreateJadwalKelas)
	jadwalKelas.Get("/", ctrl.GetAllJadwalKelas)
	jadwalKelas.Get("/:id", ctrl.GetJadwalKelasByID)
	jadwalKelas.Put("/:id", middleware.JWTAuth(), ctrl.UpdateJadwalKelas)
	jadwalKelas.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteJadwalKelas)
}
```

---

### Tahap 8 — Daftarkan Routes di Server

**File:** `cmd/server/main.go`

Tambahkan import dan panggilan `SetupJadwalKelasRoutes` di fungsi `setupRoutes`.

```go
// Tambahkan di bagian import:
jadwalkelasroutes "backend/internal/modules/jadwal_kelas/routes"

// Tambahkan di dalam fungsi setupRoutes, setelah jadwalroutes:
jadwalkelasroutes.SetupJadwalKelasRoutes(app, database.DB)
```

Setelah perubahan, blok import dan fungsi `setupRoutes` akan terlihat seperti ini:

```go
import (
	// ... import yang sudah ada ...
	jadwalroutes      "backend/internal/modules/jadwal/routes"
	jadwalkelasroutes "backend/internal/modules/jadwal_kelas/routes"  // <-- tambahkan
	kelasroutes       "backend/internal/modules/kelas/routes"
)

func setupRoutes(app *fiber.App) {
	// ... routes yang sudah ada ...
	jadwalroutes.SetupJadwalRoutes(app, database.DB)
	jadwalkelasroutes.SetupJadwalKelasRoutes(app, database.DB)  // <-- tambahkan
	kelasroutes.SetupKelasRoutes(app, database.DB)
}
```

---

### Tahap 9 — Build dan Test

```bash
# Pastikan tidak ada error compile
go build ./...

# Jalankan server (migration dan unique index otomatis dibuat saat startup)
go run cmd/server/main.go

# Test endpoint
curl http://localhost:3000/api/jadwal-kelas
```

---

## Ringkasan File yang Perlu Dibuat / Diubah

### File Baru (buat dari awal):
```
internal/modules/jadwal_kelas/model/jadwal_kelas_model.go
internal/modules/jadwal_kelas/dto/jadwal_kelas_dto.go
internal/modules/jadwal_kelas/repository/jadwal_kelas_repository.go
internal/modules/jadwal_kelas/service/jadwal_kelas_service.go
internal/modules/jadwal_kelas/controller/jadwal_kelas_controller.go
internal/modules/jadwal_kelas/routes/jadwal_kelas_routes.go
```

### File yang Dimodifikasi (tambahkan beberapa baris):
```
internal/database/migrate.go   — ubah return menjadi if err, tambahkan import & register model, tambahkan CREATE UNIQUE INDEX
cmd/server/main.go             — tambahkan import & panggilan SetupJadwalKelasRoutes
```

---

## Endpoint API yang Dihasilkan

| Method | URL | Auth | Deskripsi |
|--------|-----|------|-----------|
| `GET` | `/api/jadwal-kelas` | Tidak | Daftar semua assignment (pagination + filter) |
| `GET` | `/api/jadwal-kelas/:id` | Tidak | Detail satu assignment berdasarkan ID |
| `POST` | `/api/jadwal-kelas` | JWT | Daftarkan kelas ke jadwal ujian |
| `PUT` | `/api/jadwal-kelas/:id` | JWT | Update assignment (ganti jadwal atau kelas) |
| `DELETE` | `/api/jadwal-kelas/:id` | JWT | Hapus assignment secara permanen |

> **Tidak ada endpoint restore** — delete bersifat permanent (hard delete).

### Query Parameter untuk GET /api/jadwal-kelas:

| Parameter | Tipe | Default | Keterangan |
|-----------|------|---------|------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |
| `id_jadwal` | string (UUID) | _(opsional)_ | Filter berdasarkan jadwal tertentu |
| `id_kelas` | string (UUID) | _(opsional)_ | Filter berdasarkan kelas tertentu |

---

## Aturan Bisnis Penting

1. **Hard delete** — `DELETE /api/jadwal-kelas/:id` menghapus baris secara permanen dari database. Tidak ada cara untuk mengembalikan data yang sudah dihapus. Tidak ada endpoint restore.

2. **Unique constraint** — Kombinasi `(id_jadwal, id_kelas)` harus unik. Satu kelas tidak boleh didaftarkan lebih dari satu kali ke jadwal yang sama. Pelanggaran akan menghasilkan error 400.

3. **Validasi duplikat di service** — Pengecekan duplikat dilakukan di level service (bukan hanya database constraint) agar error message lebih informatif.

4. **Tidak ada `deleted_at`** — Karena hard delete, semua query tidak perlu menambahkan kondisi `WHERE deleted_at IS NULL`.

5. **Response menyertakan data relasi** — Response API menyertakan `nama_kelas`, `nama_bank_soal`, `wkt_mulai`, `wkt_selesai` yang diperoleh melalui JOIN:
   - `jadwal_kelas` → `kelas` (untuk `nama_kelas`)
   - `jadwal_kelas` → `jadwal` → `bank_soal` (untuk `nama_bank_soal`, `wkt_mulai`, `wkt_selesai`)

---

## Perbedaan dengan Modul Soft Delete (Jadwal, Kelas, dll.)

| Aspek | Modul Soft Delete | Modul Ini (Hard Delete) |
|-------|-------------------|------------------------|
| Kolom `deleted_at` | Ada | **Tidak ada** |
| Cara delete di repo | `Update deleted_at` | `db.Delete()` langsung |
| Endpoint restore | Ada | **Tidak ada** |
| Kondisi query | `WHERE deleted_at IS NULL` | Tidak diperlukan |
| Data setelah delete | Masih di DB, tersembunyi | **Terhapus permanen** |

---

## Referensi

- Modul jadwal (FK pertama): `internal/modules/jadwal/`
- Modul kelas (FK kedua): `internal/modules/kelas/`
- Pola repository + service + controller: `internal/modules/kelas/`
- Cara registrasi model: `internal/database/migrate.go`
- Cara registrasi routes: `cmd/server/main.go`
