# Issue: Implementasi Fitur Kelas (Class Management)

## Deskripsi

Implementasi fitur manajemen kelas mencakup:
1. Model dan migration tabel `kelas` dengan relasi ke tabel `jurusan`
2. Seeder untuk data awal kelas
3. CRUD API endpoints untuk kelas
4. Dokumentasi API untuk kebutuhan frontend

---

## Struktur Direktori yang Akan Dibuat

```
internal/
└── modules/
    └── kelas/
        ├── model/
        │   └── kelas_model.go
        ├── dto/
        │   └── kelas_dto.go
        ├── repository/
        │   └── kelas_repository.go
        ├── service/
        │   └── kelas_service.go
        ├── controller/
        │   └── kelas_controller.go
        └── routes/
            └── kelas_routes.go

internal/
└── database/
    └── seeders/
        └── kelas_seeder.go   ← file baru
```

File yang perlu **dimodifikasi** (bukan dibuat baru):
- `internal/database/migrate.go`
- `internal/database/seed.go`
- `cmd/server/main.go`

---

## Database Schema

### Tabel: `kelas`

```sql
CREATE TABLE kelas (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_jurusan UUID NOT NULL REFERENCES jurusan(id),
    nama_kelas VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_kelas_deleted   ON kelas(deleted_at);
CREATE INDEX idx_kelas_id_jurusan ON kelas(id_jurusan);
```

### Keterangan Field:
| Field | Tipe | Keterangan |
|---|---|---|
| `id` | UUID | Primary key, auto-generated |
| `id_jurusan` | UUID | Foreign key ke tabel `jurusan`, wajib diisi |
| `nama_kelas` | VARCHAR(255) | Nama kelas, wajib diisi |
| `created_at` | TIMESTAMP | Waktu data dibuat, otomatis diisi |
| `updated_at` | TIMESTAMP | Waktu data terakhir diubah, otomatis diperbarui |
| `deleted_at` | TIMESTAMP (nullable) | Soft delete — jika berisi nilai, data dianggap terhapus |

> **Catatan soft delete:** Data tidak benar-benar dihapus dari database. Kolom `deleted_at` diisi dengan waktu penghapusan. Semua query SELECT harus menyertakan `WHERE kelas.deleted_at IS NULL`.

> **Catatan relasi:** `id_jurusan` harus berisi UUID yang valid dari tabel `jurusan`. Pastikan seeder jurusan sudah dijalankan sebelum seeder kelas, karena seeder kelas membutuhkan data jurusan.

---

## Tahapan Implementasi

Ikuti urutan tahapan ini secara berurutan. Jangan lewati satu pun.

---

### TAHAP 1 — Buat Model

**File:** `internal/modules/kelas/model/kelas_model.go`

Buat direktori terlebih dahulu:
```
internal/modules/kelas/model/
```

Isi file:

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

type Kelas struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDJurusan string         `gorm:"type:uuid;not null;index" json:"id_jurusan"`
	NamaKelas string         `gorm:"type:varchar(255);not null" json:"nama_kelas"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Kelas) TableName() string {
	return "kelas"
}
```

> **Penjelasan `IDJurusan`:** Disimpan sebagai `string` bertipe UUID — pola yang sama dengan `IdBankSoal` di modul `soal` dan `IdMapel` di modul `bank_soal`. GORM tidak perlu struct relasi di model untuk menyimpan foreign key.

> **Penjelasan `gorm.DeletedAt`:** GORM mengenali `gorm.DeletedAt` sebagai field soft delete secara otomatis. Saat `db.Delete()` dipanggil, GORM mengisi `deleted_at` dengan waktu saat ini, bukan benar-benar menghapus row.

---

### TAHAP 2 — Daftarkan ke Migration

**File:** `internal/database/migrate.go`

Tambahkan import model kelas dan daftarkan ke `AutoMigrate`. Kelas **harus** didaftarkan setelah `jurusanmodel.Jurusan{}` karena tabel `kelas` mempunyai foreign key ke tabel `jurusan`.

```go
package database

import (
	banksoalmodel "backend/internal/modules/bank_soal/model"
	jurusanmodel  "backend/internal/modules/jurusan/model"
	kelasmodel    "backend/internal/modules/kelas/model"
	mapelmodel    "backend/internal/modules/mapel/model"
	soalmodel     "backend/internal/modules/soal/model"
	usermodel     "backend/internal/modules/user/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Drop old non-partial unique index on jurusan.nama_jurusan before migrating
	// so AutoMigrate can create the correct partial index (where:deleted_at IS NULL)
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
	)
}
```

> **Penting:** Urutan di `AutoMigrate` tidak memengaruhi pembuatan tabel secara langsung (GORM menjalankan semua sekaligus), tetapi pastikan `jurusan` sudah ada di database sebelum `kelas` karena foreign key `id_jurusan` merujuk ke tabel `jurusan`. Urutan list di atas sudah benar.

---

### TAHAP 3 — Buat Seeder

**File:** `internal/database/seeders/kelas_seeder.go`

Seeder kelas mengambil semua ID jurusan yang ada, lalu membuat 3 kelas per jurusan (tingkat X, XI, XII).

```go
package seeders

import (
	"backend/internal/modules/kelas/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SeedKelas(db *gorm.DB) error {
	type JurusanRow struct {
		ID          string
		NamaJurusan string
	}

	var jurusanList []JurusanRow
	if err := db.Table("jurusan").
		Where("deleted_at IS NULL").
		Select("id, nama_jurusan").
		Scan(&jurusanList).Error; err != nil {
		return err
	}

	if len(jurusanList) == 0 {
		return nil
	}

	tingkatan := []string{"X", "XI", "XII"}
	var kelasList []model.Kelas

	for _, j := range jurusanList {
		for _, tingkat := range tingkatan {
			kelasList = append(kelasList, model.Kelas{
				IDJurusan: j.ID,
				NamaKelas: fmt.Sprintf("%s - %s", tingkat, j.NamaJurusan),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	return db.CreateInBatches(kelasList, 100).Error
}
```

> **Penjelasan:** Seeder mengambil data jurusan dari tabel `jurusan` menggunakan `Scan` ke struct lokal `JurusanRow`. Ini mengikuti pola yang sama dengan `SeedBankSoal` yang mengambil `mapelIDs` dari tabel `mapel`. Jika tidak ada jurusan di database, seeder langsung berhenti tanpa error.

Kemudian **modifikasi** `internal/database/seed.go` untuk mendaftarkan seeder baru **di bagian paling akhir** (karena kelas bergantung pada jurusan):

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

	return nil
}
```

> **Penting urutan:** `SeedKelas` harus dipanggil setelah `SeedJurusan` karena seeder kelas membaca data dari tabel `jurusan`.

---

### TAHAP 4 — Jalankan Migration dan Seeder

Setelah TAHAP 1–3 selesai, jalankan server sekali agar migration dan seeder berjalan otomatis:

```bash
make run
# atau
go run cmd/server/main.go
```

Cek log server. Jika berhasil, akan ada log seperti:
```
Starting server on :3000
```

Untuk verifikasi langsung ke database, jalankan query berikut di PostgreSQL:
```sql
-- Cek tabel kelas sudah ada
SELECT COUNT(*) FROM kelas;

-- Cek data seeder masuk
SELECT k.nama_kelas, j.nama_jurusan
FROM kelas k
JOIN jurusan j ON j.id = k.id_jurusan
WHERE k.deleted_at IS NULL
ORDER BY j.nama_jurusan, k.nama_kelas;
```

---

### TAHAP 5 — Buat DTO

**File:** `internal/modules/kelas/dto/kelas_dto.go`

DTO (Data Transfer Object) mendefinisikan bentuk data request dari client dan response ke client.

```go
package dto

type CreateKelasRequest struct {
	IDJurusan string `json:"id_jurusan" validate:"required"`
	NamaKelas string `json:"nama_kelas" validate:"required"`
}

type UpdateKelasRequest struct {
	IDJurusan string `json:"id_jurusan" validate:"required"`
	NamaKelas string `json:"nama_kelas" validate:"required"`
}

type KelasResponse struct {
	ID        string `json:"id"`
	IDJurusan string `json:"id_jurusan"`
	NamaKelas string `json:"nama_kelas"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type KelasListResponse struct {
	Data      []KelasResponse `json:"data"`
	Total     int64           `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalPage int             `json:"total_page"`
}
```

---

### TAHAP 6 — Buat Repository

**File:** `internal/modules/kelas/repository/kelas_repository.go`

Repository bertanggung jawab untuk semua operasi ke database. Service tidak boleh langsung akses database — harus melalui repository.

```go
package repository

import (
	"backend/internal/modules/kelas/model"

	"gorm.io/gorm"
)

type KelasRepository interface {
	Create(kelas *model.Kelas) error
	GetByID(id string) (*model.Kelas, error)
	GetAll(page, pageSize int) ([]model.Kelas, int64, error)
	GetByJurusan(idJurusan string, page, pageSize int) ([]model.Kelas, int64, error)
	Update(kelas *model.Kelas) error
	Delete(id string) error
	Restore(id string) error
}

type kelasRepository struct {
	db *gorm.DB
}

func NewKelasRepository(db *gorm.DB) KelasRepository {
	return &kelasRepository{db: db}
}

func (r *kelasRepository) Create(kelas *model.Kelas) error {
	return r.db.Create(kelas).Error
}

func (r *kelasRepository) GetByID(id string) (*model.Kelas, error) {
	var kelas model.Kelas
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&kelas).Error
	if err != nil {
		return nil, err
	}
	return &kelas, nil
}

func (r *kelasRepository) GetAll(page, pageSize int) ([]model.Kelas, int64, error) {
	var kelasList []model.Kelas
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.Kelas{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Find(&kelasList).Error

	return kelasList, total, err
}

func (r *kelasRepository) GetByJurusan(idJurusan string, page, pageSize int) ([]model.Kelas, int64, error) {
	var kelasList []model.Kelas
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.Kelas{}).
		Where("id_jurusan = ? AND deleted_at IS NULL", idJurusan).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("id_jurusan = ? AND deleted_at IS NULL", idJurusan).
		Offset(offset).
		Limit(pageSize).
		Find(&kelasList).Error

	return kelasList, total, err
}

func (r *kelasRepository) Update(kelas *model.Kelas) error {
	return r.db.Save(kelas).Error
}

func (r *kelasRepository) Delete(id string) error {
	return r.db.Delete(&model.Kelas{}, "id = ?", id).Error
}

func (r *kelasRepository) Restore(id string) error {
	return r.db.Table("kelas").Where("id = ?", id).Update("deleted_at", nil).Error
}
```

> **Penjelasan `GetByJurusan`:** Method ini mengembalikan semua kelas berdasarkan jurusan tertentu. Berguna untuk endpoint `GET /api/kelas?id_jurusan=xxx` sehingga frontend bisa memuat kelas dropdown berdasarkan jurusan yang dipilih.

---

### TAHAP 7 — Buat Service

**File:** `internal/modules/kelas/service/kelas_service.go`

Service berisi logika bisnis. Service menerima request dari controller, memanggil repository, lalu mengembalikan response.

```go
package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/kelas/dto"
	"backend/internal/modules/kelas/model"
	"backend/internal/modules/kelas/repository"

	"gorm.io/gorm"
)

type KelasService interface {
	CreateKelas(req *dto.CreateKelasRequest) (*dto.KelasResponse, error)
	GetKelasByID(id string) (*dto.KelasResponse, error)
	GetAllKelas(page, pageSize int, idJurusan string) (*dto.KelasListResponse, error)
	UpdateKelas(id string, req *dto.UpdateKelasRequest) (*dto.KelasResponse, error)
	DeleteKelas(id string) error
	RestoreKelas(id string) error
}

type kelasService struct {
	repo repository.KelasRepository
}

func NewKelasService(repo repository.KelasRepository) KelasService {
	return &kelasService{repo: repo}
}

func (s *kelasService) CreateKelas(req *dto.CreateKelasRequest) (*dto.KelasResponse, error) {
	kelas := &model.Kelas{
		IDJurusan: req.IDJurusan,
		NamaKelas: req.NamaKelas,
	}

	if err := s.repo.Create(kelas); err != nil {
		return nil, err
	}

	return s.modelToResponse(kelas), nil
}

func (s *kelasService) GetKelasByID(id string) (*dto.KelasResponse, error) {
	kelas, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.modelToResponse(kelas), nil
}

func (s *kelasService) GetAllKelas(page, pageSize int, idJurusan string) (*dto.KelasListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var kelasList []model.Kelas
	var total int64
	var err error

	if idJurusan != "" {
		kelasList, total, err = s.repo.GetByJurusan(idJurusan, page, pageSize)
	} else {
		kelasList, total, err = s.repo.GetAll(page, pageSize)
	}

	if err != nil {
		return nil, err
	}

	var responses []dto.KelasResponse
	for _, k := range kelasList {
		responses = append(responses, *s.modelToResponse(&k))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.KelasListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *kelasService) UpdateKelas(id string, req *dto.UpdateKelasRequest) (*dto.KelasResponse, error) {
	kelas, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	kelas.IDJurusan = req.IDJurusan
	kelas.NamaKelas = req.NamaKelas

	if err := s.repo.Update(kelas); err != nil {
		return nil, err
	}

	return s.modelToResponse(kelas), nil
}

func (s *kelasService) DeleteKelas(id string) error {
	kelas, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(kelas.ID)
}

func (s *kelasService) RestoreKelas(id string) error {
	return s.repo.Restore(id)
}

func (s *kelasService) modelToResponse(kelas *model.Kelas) *dto.KelasResponse {
	return &dto.KelasResponse{
		ID:        kelas.ID,
		IDJurusan: kelas.IDJurusan,
		NamaKelas: kelas.NamaKelas,
		CreatedAt: kelas.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: kelas.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
```

> **Penjelasan `GetAllKelas` dengan parameter `idJurusan`:** Jika query string `id_jurusan` dikirim dari frontend, service memanggil `GetByJurusan` untuk filter berdasarkan jurusan. Jika tidak ada, ambil semua kelas. Ini memudahkan frontend memuat kelas dropdown berdasarkan jurusan yang dipilih user.

---

### TAHAP 8 — Buat Controller

**File:** `internal/modules/kelas/controller/kelas_controller.go`

Controller menerima request HTTP, memanggil service, lalu mengembalikan response HTTP.

```go
package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/kelas/dto"
	"backend/internal/modules/kelas/service"

	"github.com/gofiber/fiber/v2"
)

type KelasController struct {
	service service.KelasService
}

func NewKelasController(service service.KelasService) *KelasController {
	return &KelasController{service: service}
}

func (c *KelasController) CreateKelas(ctx *fiber.Ctx) error {
	var req dto.CreateKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateKelas(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create kelas successfully", resp)
}

func (c *KelasController) GetAllKelas(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idJurusan := ctx.Query("id_jurusan", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllKelas(pageNum, pageSizeNum, idJurusan)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all kelas successfully", resp)
}

func (c *KelasController) GetKelasByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetKelasByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get kelas successfully", resp)
}

func (c *KelasController) UpdateKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateKelas(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update kelas successfully", resp)
}

func (c *KelasController) DeleteKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteKelas(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete kelas successfully", nil)
}

func (c *KelasController) RestoreKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreKelas(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore kelas successfully", nil)
}
```

---

### TAHAP 9 — Buat Routes

**File:** `internal/modules/kelas/routes/kelas_routes.go`

```go
package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/kelas/controller"
	"backend/internal/modules/kelas/repository"
	"backend/internal/modules/kelas/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupKelasRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewKelasRepository(db)
	svc := service.NewKelasService(repo)
	ctrl := controller.NewKelasController(svc)

	api := app.Group("/api")
	kelas := api.Group("/kelas")

	kelas.Post("/", middleware.JWTAuth(), ctrl.CreateKelas)
	kelas.Get("/", ctrl.GetAllKelas)
	kelas.Get("/:id", ctrl.GetKelasByID)
	kelas.Put("/:id", middleware.JWTAuth(), ctrl.UpdateKelas)
	kelas.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteKelas)
	kelas.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreKelas)
}
```

> **Catatan:** `GET /` dan `GET /:id` tidak memerlukan JWT karena data kelas biasanya diakses publik (misalnya untuk dropdown di halaman registrasi). Endpoint yang mengubah data (POST, PUT, DELETE, PATCH restore) memerlukan JWT.

---

### TAHAP 10 — Daftarkan Routes ke Main

**File:** `cmd/server/main.go`

Tambahkan import dan panggil `SetupKelasRoutes` di dalam fungsi `setupRoutes`. Tambahkan tepat setelah `jurusanroutes`:

```go
// Tambahkan di bagian import:
kelasroutes "backend/internal/modules/kelas/routes"

// Tambahkan di dalam fungsi setupRoutes():
kelasroutes.SetupKelasRoutes(app, database.DB)
```

Setelah modifikasi, fungsi `setupRoutes` akan terlihat seperti ini:

```go
func setupRoutes(app *fiber.App) {
    app.Get("/health", func(ctx *fiber.Ctx) error {
        return ctx.JSON(fiber.Map{
            "status":  "ok",
            "service": "Fiber Backend API",
        })
    })

    app.Static("/uploads", "./uploads")

    authroutes.SetupAuthRoutes(app, database.DB)
    userroutes.SetupUserRoutes(app, database.DB)
    mapelroutes.SetupMapelRoutes(app, database.DB)
    banksoalroutes.SetupBankSoalRoutes(app, database.DB)
    soalroutes.SetupSoalRoutes(app, database.DB)
    jurusanroutes.SetupJurusanRoutes(app, database.DB)
    kelasroutes.SetupKelasRoutes(app, database.DB)
}
```

---

### TAHAP 11 — Verifikasi Build

Setelah semua file dibuat dan dimodifikasi, jalankan perintah berikut untuk memastikan tidak ada error kompilasi:

```bash
go build ./...
```

Jika tidak ada output, artinya kode berhasil dikompilasi.

Untuk menjalankan server (sekaligus menjalankan migration dan seeder):

```bash
make run
# atau
go run cmd/server/main.go
```

---

## API Endpoints

Base URL: `http://localhost:3000`

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| `POST` | `/api/kelas` | JWT Required | Tambah kelas baru |
| `GET` | `/api/kelas` | Public | Ambil semua kelas (dengan pagination dan filter opsional) |
| `GET` | `/api/kelas/:id` | Public | Ambil detail kelas berdasarkan ID |
| `PUT` | `/api/kelas/:id` | JWT Required | Update data kelas |
| `DELETE` | `/api/kelas/:id` | JWT Required | Hapus kelas (soft delete) |
| `PATCH` | `/api/kelas/:id/restore` | JWT Required | Pulihkan kelas yang sudah dihapus |

---

## Dokumentasi API untuk Frontend

### Format Response Umum

Semua endpoint mengembalikan format JSON yang konsisten:

```json
{
  "success": true,
  "message": "pesan sukses",
  "data": { ... }
}
```

Untuk error:

```json
{
  "success": false,
  "message": "pesan error",
  "errors": null
}
```

---

### 1. Tambah Kelas

**POST** `/api/kelas`

Header yang diperlukan:
```
Authorization: Bearer <token>
Content-Type: application/json
```

Request body:
```json
{
  "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
  "nama_kelas": "X - Teknik Komputer dan Jaringan"
}
```

Response sukses (201):
```json
{
  "success": true,
  "message": "Create kelas successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551111",
    "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X - Teknik Komputer dan Jaringan",
    "created_at": "2026-05-19 10:00:00",
    "updated_at": "2026-05-19 10:00:00"
  }
}
```

Response error — body tidak valid (400):
```json
{
  "success": false,
  "message": "Invalid request format",
  "errors": null
}
```

---

### 2. Ambil Semua Kelas

**GET** `/api/kelas`

Query parameters (opsional):
| Parameter | Default | Keterangan |
|---|---|---|
| `page` | `1` | Nomor halaman |
| `page_size` | `10` | Jumlah data per halaman |
| `id_jurusan` | _(kosong)_ | Filter kelas berdasarkan ID jurusan tertentu |

Contoh request:
```
GET /api/kelas?page=1&page_size=10
GET /api/kelas?id_jurusan=550e8400-e29b-41d4-a716-446655440000
GET /api/kelas?id_jurusan=550e8400-e29b-41d4-a716-446655440000&page=1&page_size=100
```

Response sukses (200):
```json
{
  "success": true,
  "message": "Get all kelas successfully",
  "data": {
    "data": [
      {
        "id": "661e9511-f30c-52e5-b827-557766551111",
        "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
        "nama_kelas": "X - Teknik Komputer dan Jaringan",
        "created_at": "2026-05-19 10:00:00",
        "updated_at": "2026-05-19 10:00:00"
      },
      {
        "id": "661e9511-f30c-52e5-b827-557766551112",
        "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
        "nama_kelas": "XI - Teknik Komputer dan Jaringan",
        "created_at": "2026-05-19 10:00:00",
        "updated_at": "2026-05-19 10:00:00"
      }
    ],
    "total": 15,
    "page": 1,
    "page_size": 10,
    "total_page": 2
  }
}
```

> **Tips untuk frontend:** Untuk mengisi dropdown kelas berdasarkan jurusan yang dipilih user, gunakan:
> ```
> GET /api/kelas?id_jurusan=<uuid-jurusan>&page=1&page_size=100
> ```

---

### 3. Ambil Detail Kelas

**GET** `/api/kelas/:id`

Contoh: `GET /api/kelas/661e9511-f30c-52e5-b827-557766551111`

Response sukses (200):
```json
{
  "success": true,
  "message": "Get kelas successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551111",
    "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X - Teknik Komputer dan Jaringan",
    "created_at": "2026-05-19 10:00:00",
    "updated_at": "2026-05-19 10:00:00"
  }
}
```

Response tidak ditemukan (404):
```json
{
  "success": false,
  "message": "Resource not found",
  "errors": null
}
```

---

### 4. Update Kelas

**PUT** `/api/kelas/:id`

Header yang diperlukan:
```
Authorization: Bearer <token>
Content-Type: application/json
```

Request body:
```json
{
  "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
  "nama_kelas": "X TKJ 1"
}
```

Response sukses (200):
```json
{
  "success": true,
  "message": "Update kelas successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551111",
    "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X TKJ 1",
    "created_at": "2026-05-19 10:00:00",
    "updated_at": "2026-05-19 10:30:00"
  }
}
```

---

### 5. Hapus Kelas (Soft Delete)

**DELETE** `/api/kelas/:id`

Header yang diperlukan:
```
Authorization: Bearer <token>
```

Response sukses (200):
```json
{
  "success": true,
  "message": "Delete kelas successfully",
  "data": null
}
```

> **Penting untuk frontend:** Data yang dihapus tidak akan muncul di endpoint `GET /api/kelas` maupun `GET /api/kelas/:id`. Namun data masih ada di database dan bisa dipulihkan menggunakan endpoint restore.

---

### 6. Pulihkan Kelas

**PATCH** `/api/kelas/:id/restore`

Header yang diperlukan:
```
Authorization: Bearer <token>
```

Response sukses (200):
```json
{
  "success": true,
  "message": "Restore kelas successfully",
  "data": null
}
```

---

## Checklist Implementasi

Centang setiap item setelah selesai dikerjakan:

- [ ] `internal/modules/kelas/model/kelas_model.go` — dibuat
- [ ] `internal/database/migrate.go` — ditambahkan `&kelasmodel.Kelas{}`
- [ ] `internal/database/seeders/kelas_seeder.go` — dibuat
- [ ] `internal/database/seed.go` — ditambahkan `SeedKelas` di bagian paling akhir
- [ ] `internal/modules/kelas/dto/kelas_dto.go` — dibuat
- [ ] `internal/modules/kelas/repository/kelas_repository.go` — dibuat
- [ ] `internal/modules/kelas/service/kelas_service.go` — dibuat
- [ ] `internal/modules/kelas/controller/kelas_controller.go` — dibuat
- [ ] `internal/modules/kelas/routes/kelas_routes.go` — dibuat
- [ ] `cmd/server/main.go` — ditambahkan `SetupKelasRoutes`
- [ ] `go build ./...` — tidak ada error kompilasi
- [ ] Migration berhasil (tabel `kelas` muncul di database)
- [ ] Seeder berhasil (tabel `kelas` berisi data awal)
- [ ] Semua endpoint bisa diakses dan mengembalikan response yang benar
