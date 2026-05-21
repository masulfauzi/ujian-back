# Issue: Implementasi Modul Peserta (Student Management)

## Deskripsi

Buat modul pengelolaan data peserta ujian. Peserta adalah siswa yang terdaftar dalam sebuah kelas dan akan mengikuti ujian. Setiap peserta memiliki akun login dengan `username` dan `password`.

---

## Konteks Proyek

Proyek ini menggunakan:
- **Bahasa:** Go (Golang)
- **Framework HTTP:** [Fiber v2](https://gofiber.io/)
- **ORM:** GORM dengan driver PostgreSQL
- **Struktur Modul:** `internal/modules/<nama_modul>/` berisi `model/`, `dto/`, `repository/`, `service/`, `controller/`, `routes/`
- **Password hashing:** `golang.org/x/crypto/bcrypt` — sudah ada helper di `internal/utils/password.go`
- **Module Go:** `backend` (lihat `go.mod`)

Sebelum mulai, pelajari contoh modul yang sudah ada: `internal/modules/kelas/` — pola implementasi peserta harus identik.

---

## Skema Tabel

Nama tabel: `peserta`

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID | Primary key, auto-generate via `gen_random_uuid()` |
| `nama` | varchar(255) | Nama lengkap peserta, NOT NULL |
| `id_kelas` | UUID | Foreign key ke tabel `kelas`, NOT NULL |
| `username` | varchar(100) | Username untuk login, NOT NULL, UNIQUE (partial: where deleted_at IS NULL) |
| `password` | varchar(255) | Password ter-hash bcrypt, NOT NULL |
| `created_at` | timestamp | Auto-fill saat create |
| `updated_at` | timestamp | Auto-update saat update |
| `deleted_at` | timestamp (nullable) | Soft delete — null = aktif, terisi = terhapus |

---

## Tahapan Implementasi

### Tahap 1 — Buat Model

**File:** `internal/modules/peserta/model/peserta_model.go`

Model adalah representasi struct Go dari tabel database. Ikuti pola yang persis sama seperti `internal/modules/kelas/model/kelas_model.go`.

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

type Peserta struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Nama      string         `gorm:"type:varchar(255);not null" json:"nama"`
	IDKelas   string         `gorm:"type:uuid;not null;index" json:"id_kelas"`
	Username  string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_peserta_username,where:deleted_at IS NULL" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Peserta) TableName() string {
	return "peserta"
}
```

> **Catatan penting:**
> - `json:"-"` pada field `Password` berarti password TIDAK pernah dikembalikan di response API.
> - `uniqueIndex` dengan `where:deleted_at IS NULL` memastikan username hanya unik di antara peserta yang aktif (belum dihapus). Peserta yang sudah di-soft delete boleh memiliki username yang sama dengan peserta baru.

---

### Tahap 2 — Daftarkan Model ke Migration

**File:** `internal/database/migrate.go`

Tambahkan import model peserta dan daftarkan di dalam fungsi `RunMigrations`.

```go
package database

import (
	banksoalmodel "backend/internal/modules/bank_soal/model"
	jurusanmodel  "backend/internal/modules/jurusan/model"
	kelasmodel    "backend/internal/modules/kelas/model"
	mapelmodel    "backend/internal/modules/mapel/model"
	pesertamodel  "backend/internal/modules/peserta/model"  // <-- tambahkan ini
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
		&banksoalmodel.BankSoal{},
		&soalmodel.Soal{},
		&jurusanmodel.Jurusan{},
		&kelasmodel.Kelas{},
		&pesertamodel.Peserta{},  // <-- tambahkan ini (setelah Kelas karena peserta FK ke kelas)
	)
}
```

> **Penting:** `Peserta` harus didaftarkan **setelah** `Kelas` karena tabel `peserta` memiliki foreign key ke tabel `kelas`.

---

### Tahap 3 — Buat Seeder

**File:** `internal/database/seeders/peserta_seeder.go`

Seeder bertugas mengisi data awal ke database. Seeder peserta akan:
1. Mengambil semua data kelas yang aktif dari database
2. Untuk setiap kelas, membuat 5 peserta contoh
3. Password di-hash menggunakan bcrypt sebelum disimpan

```go
package seeders

import (
	"backend/internal/modules/peserta/model"
	"backend/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SeedPeserta(db *gorm.DB) error {
	type KelasRow struct {
		ID        string
		NamaKelas string
	}

	var kelasList []KelasRow
	if err := db.Table("kelas").
		Where("deleted_at IS NULL").
		Select("id, nama_kelas").
		Scan(&kelasList).Error; err != nil {
		return err
	}

	if len(kelasList) == 0 {
		return nil
	}

	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		return err
	}

	var pesertaList []model.Peserta

	for _, k := range kelasList {
		for i := 1; i <= 5; i++ {
			username := fmt.Sprintf("peserta_%s_%d", k.ID[:8], i)
			pesertaList = append(pesertaList, model.Peserta{
				Nama:      fmt.Sprintf("Peserta %d - %s", i, k.NamaKelas),
				IDKelas:   k.ID,
				Username:  username,
				Password:  hashedPassword,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	return db.CreateInBatches(pesertaList, 100).Error
}
```

---

### Tahap 4 — Daftarkan Seeder ke Runner

**File:** `internal/database/seed.go`

Tambahkan pemanggilan `SeedPeserta` di akhir fungsi `RunSeeders`.

```go
package database

import (
	"backend/internal/database/seeders"
	"fmt"

	"gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) error {
	if err := seeders.SeedMapel(db); err != nil {
		return fmt.Errorf("failed to seed mapel: %w", err)
	}

	if err := seeders.SeedBankSoal(db); err != nil {
		return fmt.Errorf("failed to seed bank_soal: %w", err)
	}

	if err := seeders.SeedJurusan(db); err != nil {
		return fmt.Errorf("failed to seed jurusan: %w", err)
	}

	if err := seeders.SeedKelas(db); err != nil {
		return fmt.Errorf("failed to seed kelas: %w", err)
	}

	if err := seeders.SeedPeserta(db); err != nil {  // <-- tambahkan ini
		return fmt.Errorf("failed to seed peserta: %w", err)
	}

	return nil
}
```

> **Penting:** `SeedPeserta` harus dipanggil **setelah** `SeedKelas` karena data peserta butuh ID kelas yang sudah ada.

---

### Tahap 5 — Jalankan Migration dan Seeder

Setelah model dan seeder siap, jalankan server agar migration otomatis dijalankan:

```bash
# Jalankan server (migration otomatis berjalan saat startup)
go run cmd/server/main.go
```

Untuk menjalankan seeder secara manual (jika ada command khusus):

```bash
# Cek apakah ada command seed
go run cmd/seed/main.go
```

Atau jalankan seeder via script yang sudah ada di `Makefile`:

```bash
make seed
```

Verifikasi tabel berhasil dibuat dengan cek di database PostgreSQL:

```sql
\d peserta
SELECT COUNT(*) FROM peserta;
```

---

### Tahap 6 — Buat DTO (Data Transfer Object)

**File:** `internal/modules/peserta/dto/peserta_dto.go`

DTO mendefinisikan struktur data untuk request (input) dan response (output) API.

```go
package dto

type CreatePesertaRequest struct {
	Nama     string `json:"nama" validate:"required"`
	IDKelas  string `json:"id_kelas" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdatePesertaRequest struct {
	Nama     string `json:"nama" validate:"required"`
	IDKelas  string `json:"id_kelas" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password"`  // opsional saat update — kosong = tidak diubah
}

type PesertaResponse struct {
	ID        string `json:"id"`
	Nama      string `json:"nama"`
	IDKelas   string `json:"id_kelas"`
	NamaKelas string `json:"nama_kelas"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PesertaListResponse struct {
	Data      []PesertaResponse `json:"data"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	TotalPage int               `json:"total_page"`
}
```

> **Perhatikan:** Tidak ada field `password` di `PesertaResponse` — password tidak pernah dikembalikan ke client.

---

### Tahap 7 — Buat Repository

**File:** `internal/modules/peserta/repository/peserta_repository.go`

Repository menangani semua interaksi langsung dengan database. Query menggunakan JOIN ke tabel `kelas` agar `nama_kelas` ikut tampil di response.

```go
package repository

import (
	"backend/internal/modules/peserta/model"

	"gorm.io/gorm"
)

type PesertaWithKelas struct {
	ID        string  `gorm:"column:id"`
	Nama      string  `gorm:"column:nama"`
	IDKelas   string  `gorm:"column:id_kelas"`
	NamaKelas string  `gorm:"column:nama_kelas"`
	Username  string  `gorm:"column:username"`
	CreatedAt string  `gorm:"column:created_at"`
	UpdatedAt string  `gorm:"column:updated_at"`
	DeletedAt *string `gorm:"column:deleted_at"`
}

func (PesertaWithKelas) TableName() string {
	return "peserta"
}

type PesertaRepository interface {
	Create(peserta *model.Peserta) error
	GetByID(id string) (*PesertaWithKelas, error)
	GetAll(page, pageSize int, idKelas string) ([]PesertaWithKelas, int64, error)
	GetRawByID(id string) (*model.Peserta, error)
	GetByUsername(username string) (*model.Peserta, error)
	Update(peserta *model.Peserta) error
	Delete(id string) error
	Restore(id string) error
}

type pesertaRepository struct {
	db *gorm.DB
}

func NewPesertaRepository(db *gorm.DB) PesertaRepository {
	return &pesertaRepository{db: db}
}

func (r *pesertaRepository) Create(peserta *model.Peserta) error {
	return r.db.Create(peserta).Error
}

func (r *pesertaRepository) GetByID(id string) (*PesertaWithKelas, error) {
	var peserta PesertaWithKelas
	err := r.db.
		Select("peserta.id, peserta.nama, peserta.id_kelas, peserta.username, peserta.created_at, peserta.updated_at, peserta.deleted_at, kelas.nama_kelas").
		Joins("LEFT JOIN kelas ON peserta.id_kelas = kelas.id").
		Where("peserta.id = ? AND peserta.deleted_at IS NULL", id).
		First(&peserta).Error
	if err != nil {
		return nil, err
	}
	return &peserta, nil
}

func (r *pesertaRepository) GetAll(page, pageSize int, idKelas string) ([]PesertaWithKelas, int64, error) {
	var pesertaList []PesertaWithKelas
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	countQuery := r.db.Table("peserta").
		Joins("LEFT JOIN kelas ON peserta.id_kelas = kelas.id").
		Where("peserta.deleted_at IS NULL")

	if idKelas != "" {
		countQuery = countQuery.Where("peserta.id_kelas = ?", idKelas)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.
		Select("peserta.id, peserta.nama, peserta.id_kelas, peserta.username, peserta.created_at, peserta.updated_at, peserta.deleted_at, kelas.nama_kelas").
		Joins("LEFT JOIN kelas ON peserta.id_kelas = kelas.id").
		Where("peserta.deleted_at IS NULL")

	if idKelas != "" {
		query = query.Where("peserta.id_kelas = ?", idKelas)
	}

	err := query.Offset(offset).Limit(pageSize).Find(&pesertaList).Error
	return pesertaList, total, err
}

func (r *pesertaRepository) GetRawByID(id string) (*model.Peserta, error) {
	var peserta model.Peserta
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&peserta).Error
	if err != nil {
		return nil, err
	}
	return &peserta, nil
}

func (r *pesertaRepository) GetByUsername(username string) (*model.Peserta, error) {
	var peserta model.Peserta
	err := r.db.Where("username = ? AND deleted_at IS NULL", username).First(&peserta).Error
	if err != nil {
		return nil, err
	}
	return &peserta, nil
}

func (r *pesertaRepository) Update(peserta *model.Peserta) error {
	return r.db.Save(peserta).Error
}

func (r *pesertaRepository) Delete(id string) error {
	return r.db.Delete(&model.Peserta{}, "id = ?", id).Error
}

func (r *pesertaRepository) Restore(id string) error {
	return r.db.Table("peserta").Where("id = ?", id).Update("deleted_at", nil).Error
}
```

> **Catatan `GetRawByID`:** Method ini mengambil model lengkap termasuk field `password` — digunakan oleh service saat update untuk mengisi ulang password lama jika password baru tidak diberikan.

---

### Tahap 8 — Buat Service

**File:** `internal/modules/peserta/service/peserta_service.go`

Service berisi business logic. Di sini password di-hash sebelum disimpan, dan username dicek keunikannya.

```go
package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/peserta/dto"
	"backend/internal/modules/peserta/model"
	"backend/internal/modules/peserta/repository"
	"backend/internal/utils"

	"gorm.io/gorm"
)

func pesertaWithKelasToResponse(p *repository.PesertaWithKelas) *dto.PesertaResponse {
	return &dto.PesertaResponse{
		ID:        p.ID,
		Nama:      p.Nama,
		IDKelas:   p.IDKelas,
		NamaKelas: p.NamaKelas,
		Username:  p.Username,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type PesertaService interface {
	CreatePeserta(req *dto.CreatePesertaRequest) (*dto.PesertaResponse, error)
	GetPesertaByID(id string) (*dto.PesertaResponse, error)
	GetAllPeserta(page, pageSize int, idKelas string) (*dto.PesertaListResponse, error)
	UpdatePeserta(id string, req *dto.UpdatePesertaRequest) (*dto.PesertaResponse, error)
	DeletePeserta(id string) error
	RestorePeserta(id string) error
}

type pesertaService struct {
	repo repository.PesertaRepository
}

func NewPesertaService(repo repository.PesertaRepository) PesertaService {
	return &pesertaService{repo: repo}
}

func (s *pesertaService) CreatePeserta(req *dto.CreatePesertaRequest) (*dto.PesertaResponse, error) {
	// Cek apakah username sudah digunakan
	existing, err := s.repo.GetByUsername(req.Username)
	if err == nil && existing != nil {
		return nil, errors.New("username sudah digunakan")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	peserta := &model.Peserta{
		Nama:     req.Nama,
		IDKelas:  req.IDKelas,
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := s.repo.Create(peserta); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByID(peserta.ID)
	if err != nil {
		return nil, err
	}

	return pesertaWithKelasToResponse(created), nil
}

func (s *pesertaService) GetPesertaByID(id string) (*dto.PesertaResponse, error) {
	peserta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return pesertaWithKelasToResponse(peserta), nil
}

func (s *pesertaService) GetAllPeserta(page, pageSize int, idKelas string) (*dto.PesertaListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	pesertaList, total, err := s.repo.GetAll(page, pageSize, idKelas)
	if err != nil {
		return nil, err
	}

	var responses []dto.PesertaResponse
	for _, p := range pesertaList {
		responses = append(responses, *pesertaWithKelasToResponse(&p))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PesertaListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *pesertaService) UpdatePeserta(id string, req *dto.UpdatePesertaRequest) (*dto.PesertaResponse, error) {
	existing, err := s.repo.GetRawByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	// Cek jika username diubah dan sudah dipakai peserta lain
	if req.Username != existing.Username {
		taken, err := s.repo.GetByUsername(req.Username)
		if err == nil && taken != nil && taken.ID != id {
			return nil, errors.New("username sudah digunakan")
		}
	}

	peserta := &model.Peserta{
		ID:       existing.ID,
		Nama:     req.Nama,
		IDKelas:  req.IDKelas,
		Username: req.Username,
		Password: existing.Password, // gunakan password lama sebagai default
	}

	// Jika password baru diberikan, hash dan ganti
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("gagal memproses password")
		}
		peserta.Password = hashedPassword
	}

	if err := s.repo.Update(peserta); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return pesertaWithKelasToResponse(updated), nil
}

func (s *pesertaService) DeletePeserta(id string) error {
	peserta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(peserta.ID)
}

func (s *pesertaService) RestorePeserta(id string) error {
	return s.repo.Restore(id)
}
```

---

### Tahap 9 — Buat Controller

**File:** `internal/modules/peserta/controller/peserta_controller.go`

Controller menangani request HTTP dan memanggil service.

```go
package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/peserta/dto"
	"backend/internal/modules/peserta/service"

	"github.com/gofiber/fiber/v2"
)

type PesertaController struct {
	service service.PesertaService
}

func NewPesertaController(service service.PesertaService) *PesertaController {
	return &PesertaController{service: service}
}

func (c *PesertaController) CreatePeserta(ctx *fiber.Ctx) error {
	var req dto.CreatePesertaRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreatePeserta(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create peserta successfully", resp)
}

func (c *PesertaController) GetAllPeserta(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idKelas := ctx.Query("id_kelas", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllPeserta(pageNum, pageSizeNum, idKelas)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all peserta successfully", resp)
}

func (c *PesertaController) GetPesertaByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetPesertaByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get peserta successfully", resp)
}

func (c *PesertaController) UpdatePeserta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdatePesertaRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdatePeserta(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update peserta successfully", resp)
}

func (c *PesertaController) DeletePeserta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeletePeserta(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete peserta successfully", nil)
}

func (c *PesertaController) RestorePeserta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestorePeserta(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore peserta successfully", nil)
}
```

---

### Tahap 10 — Buat Routes

**File:** `internal/modules/peserta/routes/peserta_routes.go`

```go
package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/peserta/controller"
	"backend/internal/modules/peserta/repository"
	"backend/internal/modules/peserta/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupPesertaRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewPesertaRepository(db)
	svc := service.NewPesertaService(repo)
	ctrl := controller.NewPesertaController(svc)

	api := app.Group("/api")
	peserta := api.Group("/peserta")

	peserta.Post("/", middleware.JWTAuth(), ctrl.CreatePeserta)
	peserta.Get("/", ctrl.GetAllPeserta)
	peserta.Get("/:id", ctrl.GetPesertaByID)
	peserta.Put("/:id", middleware.JWTAuth(), ctrl.UpdatePeserta)
	peserta.Delete("/:id", middleware.JWTAuth(), ctrl.DeletePeserta)
	peserta.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestorePeserta)
}
```

---

### Tahap 11 — Daftarkan Routes di Server

**File:** `cmd/server/main.go`

Tambahkan import dan panggilan `SetupPesertaRoutes` di fungsi `setupRoutes`.

```go
// Tambahkan import ini di bagian import
pesertaroutes "backend/internal/modules/peserta/routes"

// Tambahkan baris ini di dalam fungsi setupRoutes, setelah kelasroutes
pesertaroutes.SetupPesertaRoutes(app, database.DB)
```

Setelah perubahan, blok import dan fungsi `setupRoutes` di `cmd/server/main.go` akan terlihat seperti ini:

```go
import (
	// ... import lainnya yang sudah ada ...
	kelasroutes   "backend/internal/modules/kelas/routes"
	pesertaroutes "backend/internal/modules/peserta/routes"  // <-- tambahkan
)

func setupRoutes(app *fiber.App) {
	// ... routes lainnya yang sudah ada ...
	kelasroutes.SetupKelasRoutes(app, database.DB)
	pesertaroutes.SetupPesertaRoutes(app, database.DB)  // <-- tambahkan
}
```

---

### Tahap 12 — Build dan Test

```bash
# Pastikan tidak ada error compile
go build ./...

# Jalankan server
go run cmd/server/main.go

# Test endpoint berjalan
curl http://localhost:3000/api/peserta
```

---

### Tahap 13 — Buat Dokumentasi API

**File:** `docs/PESERTA_API.md`

Buat file dokumentasi di folder `docs/` mengikuti format yang sama dengan `docs/KELAS_API.md`. Dokumentasi harus mencakup:

1. Base URL: `http://localhost:3000/api/peserta`
2. Daftar semua endpoint beserta method HTTP-nya
3. Untuk setiap endpoint: deskripsi, query params / path params / request body, contoh request cURL, contoh response sukses, contoh response error
4. Tabel HTTP status codes
5. Catatan penting (soft delete, field password tidak dikembalikan, dll)

Lihat `docs/KELAS_API.md` sebagai referensi format yang harus diikuti.

---

## Ringkasan File yang Perlu Dibuat / Diubah

### File Baru (buat dari awal):
```
internal/modules/peserta/model/peserta_model.go
internal/modules/peserta/dto/peserta_dto.go
internal/modules/peserta/repository/peserta_repository.go
internal/modules/peserta/service/peserta_service.go
internal/modules/peserta/controller/peserta_controller.go
internal/modules/peserta/routes/peserta_routes.go
internal/database/seeders/peserta_seeder.go
docs/PESERTA_API.md
```

### File yang Dimodifikasi (tambahkan beberapa baris):
```
internal/database/migrate.go        — tambahkan import & daftarkan &pesertamodel.Peserta{}
internal/database/seed.go           — tambahkan pemanggilan seeders.SeedPeserta(db)
cmd/server/main.go                  — tambahkan import & panggilan SetupPesertaRoutes
```

---

## Endpoint API yang Dihasilkan

| Method | URL | Auth | Deskripsi |
|--------|-----|------|-----------|
| `GET` | `/api/peserta` | Tidak | Daftar semua peserta (pagination + filter kelas) |
| `GET` | `/api/peserta/:id` | Tidak | Detail satu peserta berdasarkan ID |
| `POST` | `/api/peserta` | JWT | Buat peserta baru |
| `PUT` | `/api/peserta/:id` | JWT | Update data peserta |
| `DELETE` | `/api/peserta/:id` | JWT | Soft delete peserta |
| `PATCH` | `/api/peserta/:id/restore` | JWT | Restore peserta yang dihapus |

### Query Parameter untuk GET /api/peserta:
| Parameter | Tipe | Default | Keterangan |
|-----------|------|---------|------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |
| `id_kelas` | string (UUID) | _(opsional)_ | Filter peserta berdasarkan kelas |

---

## Aturan Bisnis Penting

1. **Password tidak pernah dikembalikan** — field `password` tidak ada di response API sama sekali. Di model Go sudah ditandai `json:"-"`.

2. **Password di-hash dengan bcrypt** — selalu gunakan `utils.HashPassword()` sebelum menyimpan password ke database. Jangan pernah simpan password plain text.

3. **Username unik per peserta aktif** — dua peserta aktif tidak boleh memiliki username yang sama. Peserta yang sudah di-soft delete tidak dihitung (partial unique index).

4. **Update password bersifat opsional** — saat update, jika field `password` kosong atau tidak dikirim, password lama dipertahankan.

5. **Soft delete** — data tidak benar-benar dihapus dari database. Field `deleted_at` diisi timestamp saat dihapus. Data yang sudah dihapus tidak muncul di GET.

6. **Relasi ke kelas** — setiap peserta wajib memiliki `id_kelas` yang valid. Response API menyertakan `nama_kelas` hasil JOIN dari tabel `kelas`.

---

## Referensi

Semua pola implementasi mengacu pada modul `kelas` yang sudah ada:
- Model: `internal/modules/kelas/model/kelas_model.go`
- DTO: `internal/modules/kelas/dto/kelas_dto.go`
- Repository: `internal/modules/kelas/repository/kelas_repository.go`
- Service: `internal/modules/kelas/service/kelas_service.go`
- Controller: `internal/modules/kelas/controller/kelas_controller.go`
- Routes: `internal/modules/kelas/routes/kelas_routes.go`
- Seeder: `internal/database/seeders/kelas_seeder.go`
- Dokumentasi: `docs/KELAS_API.md`
