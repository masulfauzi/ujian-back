# Issue: Implementasi Fitur Jurusan (Department Management)

## Deskripsi

Implementasi fitur manajemen jurusan mencakup:
1. Model dan migration tabel `jurusan`
2. Seeder untuk data awal jurusan
3. CRUD API endpoints untuk jurusan
4. Dokumentasi API untuk kebutuhan frontend

---

## Struktur Direktori yang Akan Dibuat

```
internal/
└── modules/
    └── jurusan/
        ├── model/
        │   └── jurusan_model.go
        ├── dto/
        │   └── jurusan_dto.go
        ├── repository/
        │   └── jurusan_repository.go
        ├── service/
        │   └── jurusan_service.go
        ├── controller/
        │   └── jurusan_controller.go
        └── routes/
            └── jurusan_routes.go

internal/
└── database/
    └── seeders/
        └── jurusan_seeder.go   ← file baru
```

File yang perlu **dimodifikasi** (bukan dibuat baru):
- `internal/database/migrate.go`
- `internal/database/seed.go`
- `cmd/server/main.go`

---

## Database Schema

### Tabel: `jurusan`

```sql
CREATE TABLE jurusan (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_jurusan VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL
);

CREATE INDEX idx_jurusan_deleted ON jurusan(deleted_at);
```

### Keterangan Field:
| Field | Tipe | Keterangan |
|---|---|---|
| `id` | UUID | Primary key, auto-generated |
| `nama_jurusan` | VARCHAR(255) | Nama jurusan, wajib diisi, unik |
| `created_at` | TIMESTAMP | Waktu data dibuat, otomatis diisi |
| `updated_at` | TIMESTAMP | Waktu data terakhir diubah, otomatis diperbarui |
| `deleted_at` | TIMESTAMP (nullable) | Soft delete — jika berisi nilai, data dianggap terhapus |

> **Catatan soft delete:** Data tidak benar-benar dihapus dari database. Kolom `deleted_at` diisi dengan waktu penghapusan. Semua query SELECT harus menyertakan `WHERE deleted_at IS NULL`.

---

## Tahapan Implementasi

Ikuti urutan tahapan ini secara berurutan. Jangan lewati satu pun.

---

### TAHAP 1 — Buat Model

**File:** `internal/modules/jurusan/model/jurusan_model.go`

Buat direktori terlebih dahulu:
```
internal/modules/jurusan/model/
```

Isi file:

```go
package model

import (
    "time"
)

type Jurusan struct {
    ID           string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    NamaJurusan  string     `gorm:"type:varchar(255);uniqueIndex" json:"nama_jurusan"`
    CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt    *time.Time `gorm:"index" json:"deleted_at"`
}

func (Jurusan) TableName() string {
    return "jurusan"
}
```

> **Penjelasan:** `DeletedAt` bertipe pointer (`*time.Time`) supaya nilainya bisa `nil` (belum dihapus) atau berisi waktu (sudah dihapus). GORM akan mengenali field ini sebagai soft delete secara otomatis.

---

### TAHAP 2 — Daftarkan ke Migration

**File:** `internal/database/migrate.go`

Tambahkan import model jurusan dan daftarkan ke `AutoMigrate`. Contoh setelah dimodifikasi:

```go
package database

import (
    banksoalmodel "backend/internal/modules/bank_soal/model"
    jurusanmodel  "backend/internal/modules/jurusan/model"
    mapelmodel    "backend/internal/modules/mapel/model"
    soalmodel     "backend/internal/modules/soal/model"
    usermodel     "backend/internal/modules/user/model"

    "gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
    return db.AutoMigrate(
        &usermodel.User{},
        &mapelmodel.Mapel{},
        &banksoalmodel.BankSoal{},
        &soalmodel.Soal{},
        &jurusanmodel.Jurusan{},
    )
}
```

> **Penjelasan:** `AutoMigrate` akan membuat tabel `jurusan` secara otomatis saat server pertama kali dijalankan. Jika tabel sudah ada, GORM hanya akan menambahkan kolom yang belum ada (tidak menghapus data).

---

### TAHAP 3 — Buat Seeder

**File:** `internal/database/seeders/jurusan_seeder.go`

```go
package seeders

import (
    "backend/internal/modules/jurusan/model"
    "time"

    "gorm.io/gorm"
)

func SeedJurusan(db *gorm.DB) error {
    jurusans := []model.Jurusan{
        {NamaJurusan: "Teknik Komputer dan Jaringan", CreatedAt: time.Now(), UpdatedAt: time.Now()},
        {NamaJurusan: "Rekayasa Perangkat Lunak", CreatedAt: time.Now(), UpdatedAt: time.Now()},
        {NamaJurusan: "Multimedia", CreatedAt: time.Now(), UpdatedAt: time.Now()},
        {NamaJurusan: "Akuntansi", CreatedAt: time.Now(), UpdatedAt: time.Now()},
        {NamaJurusan: "Administrasi Perkantoran", CreatedAt: time.Now(), UpdatedAt: time.Now()},
    }

    return db.CreateInBatches(jurusans, 100).Error
}
```

Kemudian **modifikasi** `internal/database/seed.go` untuk mendaftarkan seeder baru:

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

    return nil
}
```

---

### TAHAP 4 — Buat DTO

**File:** `internal/modules/jurusan/dto/jurusan_dto.go`

DTO (Data Transfer Object) adalah struct yang mendefinisikan bentuk data request dari client dan response ke client. Model (`jurusan_model.go`) dipakai untuk komunikasi dengan database, sedangkan DTO dipakai untuk komunikasi dengan HTTP layer.

```go
package dto

type CreateJurusanRequest struct {
    NamaJurusan string `json:"nama_jurusan" validate:"required"`
}

type UpdateJurusanRequest struct {
    NamaJurusan string `json:"nama_jurusan" validate:"required"`
}

type JurusanResponse struct {
    ID          string `json:"id"`
    NamaJurusan string `json:"nama_jurusan"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type JurusanListResponse struct {
    Data      []JurusanResponse `json:"data"`
    Total     int64             `json:"total"`
    Page      int               `json:"page"`
    PageSize  int               `json:"page_size"`
    TotalPage int               `json:"total_page"`
}
```

---

### TAHAP 5 — Buat Repository

**File:** `internal/modules/jurusan/repository/jurusan_repository.go`

Repository bertanggung jawab untuk semua operasi ke database. Service tidak boleh langsung akses database — harus melalui repository.

```go
package repository

import (
    "backend/internal/modules/jurusan/model"

    "gorm.io/gorm"
)

type JurusanRepository interface {
    Create(jurusan *model.Jurusan) error
    GetByID(id string) (*model.Jurusan, error)
    GetAll(page, pageSize int) ([]model.Jurusan, int64, error)
    Update(jurusan *model.Jurusan) error
    Delete(id string) error
    Restore(id string) error
}

type jurusanRepository struct {
    db *gorm.DB
}

func NewJurusanRepository(db *gorm.DB) JurusanRepository {
    return &jurusanRepository{db: db}
}

func (r *jurusanRepository) Create(jurusan *model.Jurusan) error {
    return r.db.Create(jurusan).Error
}

func (r *jurusanRepository) GetByID(id string) (*model.Jurusan, error) {
    var jurusan model.Jurusan
    err := r.db.
        Where("id = ? AND deleted_at IS NULL", id).
        First(&jurusan).Error
    if err != nil {
        return nil, err
    }
    return &jurusan, nil
}

func (r *jurusanRepository) GetAll(page, pageSize int) ([]model.Jurusan, int64, error) {
    var jurusans []model.Jurusan
    var total int64

    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 10
    }

    offset := (page - 1) * pageSize

    err := r.db.
        Model(&model.Jurusan{}).
        Where("deleted_at IS NULL").
        Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    err = r.db.
        Where("deleted_at IS NULL").
        Offset(offset).
        Limit(pageSize).
        Find(&jurusans).Error

    return jurusans, total, err
}

func (r *jurusanRepository) Update(jurusan *model.Jurusan) error {
    return r.db.Save(jurusan).Error
}

func (r *jurusanRepository) Delete(id string) error {
    return r.db.Delete(&model.Jurusan{}, "id = ?", id).Error
}

func (r *jurusanRepository) Restore(id string) error {
    return r.db.Table("jurusan").Where("id = ?", id).Update("deleted_at", nil).Error
}
```

---

### TAHAP 6 — Buat Service

**File:** `internal/modules/jurusan/service/jurusan_service.go`

Service berisi logika bisnis. Service menerima request dari controller, memanggil repository, lalu mengembalikan response.

```go
package service

import (
    "errors"
    "math"

    "backend/internal/constants"
    "backend/internal/modules/jurusan/dto"
    "backend/internal/modules/jurusan/model"
    "backend/internal/modules/jurusan/repository"

    "gorm.io/gorm"
)

type JurusanService interface {
    CreateJurusan(req *dto.CreateJurusanRequest) (*dto.JurusanResponse, error)
    GetJurusanByID(id string) (*dto.JurusanResponse, error)
    GetAllJurusan(page, pageSize int) (*dto.JurusanListResponse, error)
    UpdateJurusan(id string, req *dto.UpdateJurusanRequest) (*dto.JurusanResponse, error)
    DeleteJurusan(id string) error
    RestoreJurusan(id string) error
}

type jurusanService struct {
    repo repository.JurusanRepository
}

func NewJurusanService(repo repository.JurusanRepository) JurusanService {
    return &jurusanService{repo: repo}
}

func (s *jurusanService) CreateJurusan(req *dto.CreateJurusanRequest) (*dto.JurusanResponse, error) {
    jurusan := &model.Jurusan{
        NamaJurusan: req.NamaJurusan,
    }

    if err := s.repo.Create(jurusan); err != nil {
        return nil, err
    }

    return s.modelToResponse(jurusan), nil
}

func (s *jurusanService) GetJurusanByID(id string) (*dto.JurusanResponse, error) {
    jurusan, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New(constants.ErrNotFound)
        }
        return nil, err
    }

    return s.modelToResponse(jurusan), nil
}

func (s *jurusanService) GetAllJurusan(page, pageSize int) (*dto.JurusanListResponse, error) {
    jurusans, total, err := s.repo.GetAll(page, pageSize)
    if err != nil {
        return nil, err
    }

    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 10
    }

    var responses []dto.JurusanResponse
    for _, j := range jurusans {
        responses = append(responses, *s.modelToResponse(&j))
    }

    totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

    return &dto.JurusanListResponse{
        Data:      responses,
        Total:     total,
        Page:      page,
        PageSize:  pageSize,
        TotalPage: totalPage,
    }, nil
}

func (s *jurusanService) UpdateJurusan(id string, req *dto.UpdateJurusanRequest) (*dto.JurusanResponse, error) {
    jurusan, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New(constants.ErrNotFound)
        }
        return nil, err
    }

    jurusan.NamaJurusan = req.NamaJurusan

    if err := s.repo.Update(jurusan); err != nil {
        return nil, err
    }

    return s.modelToResponse(jurusan), nil
}

func (s *jurusanService) DeleteJurusan(id string) error {
    jurusan, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New(constants.ErrNotFound)
        }
        return err
    }

    return s.repo.Delete(jurusan.ID)
}

func (s *jurusanService) RestoreJurusan(id string) error {
    return s.repo.Restore(id)
}

func (s *jurusanService) modelToResponse(jurusan *model.Jurusan) *dto.JurusanResponse {
    return &dto.JurusanResponse{
        ID:          jurusan.ID,
        NamaJurusan: jurusan.NamaJurusan,
        CreatedAt:   jurusan.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt:   jurusan.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
}
```

> **Penjelasan `constants.ErrNotFound`:** Konstanta ini sudah ada di `internal/constants/constants.go`. Gunakan konstanta yang sama, jangan buat string error baru.

---

### TAHAP 7 — Buat Controller

**File:** `internal/modules/jurusan/controller/jurusan_controller.go`

Controller menerima request HTTP, memanggil service, lalu mengembalikan response HTTP.

```go
package controller

import (
    "strconv"

    "backend/internal/helpers"
    "backend/internal/modules/jurusan/dto"
    "backend/internal/modules/jurusan/service"

    "github.com/gofiber/fiber/v2"
)

type JurusanController struct {
    service service.JurusanService
}

func NewJurusanController(service service.JurusanService) *JurusanController {
    return &JurusanController{service: service}
}

func (c *JurusanController) CreateJurusan(ctx *fiber.Ctx) error {
    var req dto.CreateJurusanRequest

    if err := ctx.BodyParser(&req); err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
    }

    resp, err := c.service.CreateJurusan(&req)
    if err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
    }

    return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jurusan successfully", resp)
}

func (c *JurusanController) GetAllJurusan(ctx *fiber.Ctx) error {
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

    resp, err := c.service.GetAllJurusan(pageNum, pageSizeNum)
    if err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
    }

    return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jurusan successfully", resp)
}

func (c *JurusanController) GetJurusanByID(ctx *fiber.Ctx) error {
    id := ctx.Params("id")

    resp, err := c.service.GetJurusanByID(id)
    if err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
    }

    return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jurusan successfully", resp)
}

func (c *JurusanController) UpdateJurusan(ctx *fiber.Ctx) error {
    id := ctx.Params("id")
    var req dto.UpdateJurusanRequest

    if err := ctx.BodyParser(&req); err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
    }

    resp, err := c.service.UpdateJurusan(id, &req)
    if err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
    }

    return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jurusan successfully", resp)
}

func (c *JurusanController) DeleteJurusan(ctx *fiber.Ctx) error {
    id := ctx.Params("id")

    err := c.service.DeleteJurusan(id)
    if err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
    }

    return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jurusan successfully", nil)
}

func (c *JurusanController) RestoreJurusan(ctx *fiber.Ctx) error {
    id := ctx.Params("id")

    err := c.service.RestoreJurusan(id)
    if err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
    }

    return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore jurusan successfully", nil)
}
```

---

### TAHAP 8 — Buat Routes

**File:** `internal/modules/jurusan/routes/jurusan_routes.go`

```go
package routes

import (
    "backend/internal/middleware"
    "backend/internal/modules/jurusan/controller"
    "backend/internal/modules/jurusan/repository"
    "backend/internal/modules/jurusan/service"

    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

func SetupJurusanRoutes(app *fiber.App, db *gorm.DB) {
    repo := repository.NewJurusanRepository(db)
    svc := service.NewJurusanService(repo)
    ctrl := controller.NewJurusanController(svc)

    api := app.Group("/api")
    jurusan := api.Group("/jurusan")

    jurusan.Post("/", middleware.JWTAuth(), ctrl.CreateJurusan)
    jurusan.Get("/", ctrl.GetAllJurusan)
    jurusan.Get("/:id", ctrl.GetJurusanByID)
    jurusan.Put("/:id", middleware.JWTAuth(), ctrl.UpdateJurusan)
    jurusan.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteJurusan)
    jurusan.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreJurusan)
}
```

> **Catatan:** `GET /` dan `GET /:id` tidak memerlukan JWT karena biasanya data jurusan diakses publik (misalnya untuk dropdown di halaman registrasi). Endpoint yang mengubah data (POST, PUT, DELETE, PATCH restore) memerlukan JWT.

---

### TAHAP 9 — Daftarkan Routes ke Main

**File:** `cmd/server/main.go`

Tambahkan import dan panggil `SetupJurusanRoutes` di dalam fungsi `setupRoutes`:

```go
// Tambahkan di bagian import:
jurusanroutes "backend/internal/modules/jurusan/routes"

// Tambahkan di dalam fungsi setupRoutes():
jurusanroutes.SetupJurusanRoutes(app, database.DB)
```

Setelah modifikasi, bagian `setupRoutes` akan terlihat seperti ini:

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
}
```

---

### TAHAP 10 — Verifikasi

Setelah semua file dibuat dan dimodifikasi, jalankan perintah berikut untuk memastikan tidak ada error kompilasi:

```bash
go build ./...
```

Jika tidak ada output, artinya kode berhasil dikompilasi.

Untuk menjalankan server:

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
| `POST` | `/api/jurusan` | JWT Required | Tambah jurusan baru |
| `GET` | `/api/jurusan` | Public | Ambil semua jurusan (dengan pagination) |
| `GET` | `/api/jurusan/:id` | Public | Ambil detail jurusan berdasarkan ID |
| `PUT` | `/api/jurusan/:id` | JWT Required | Update data jurusan |
| `DELETE` | `/api/jurusan/:id` | JWT Required | Hapus jurusan (soft delete) |
| `PATCH` | `/api/jurusan/:id/restore` | JWT Required | Pulihkan jurusan yang sudah dihapus |

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

### 1. Tambah Jurusan

**POST** `/api/jurusan`

Header yang diperlukan:
```
Authorization: Bearer <token>
Content-Type: application/json
```

Request body:
```json
{
  "nama_jurusan": "Teknik Komputer dan Jaringan"
}
```

Response sukses (201):
```json
{
  "success": true,
  "message": "Create jurusan successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "nama_jurusan": "Teknik Komputer dan Jaringan",
    "created_at": "2026-05-19 10:00:00",
    "updated_at": "2026-05-19 10:00:00"
  }
}
```

---

### 2. Ambil Semua Jurusan

**GET** `/api/jurusan?page=1&page_size=10`

Query parameters (opsional):
| Parameter | Default | Keterangan |
|---|---|---|
| `page` | `1` | Nomor halaman |
| `page_size` | `10` | Jumlah data per halaman |

Response sukses (200):
```json
{
  "success": true,
  "message": "Get all jurusan successfully",
  "data": {
    "data": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "nama_jurusan": "Teknik Komputer dan Jaringan",
        "created_at": "2026-05-19 10:00:00",
        "updated_at": "2026-05-19 10:00:00"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10,
    "total_page": 1
  }
}
```

Contoh penggunaan di frontend (fetch semua tanpa pagination untuk dropdown):
```
GET /api/jurusan?page=1&page_size=100
```

---

### 3. Ambil Detail Jurusan

**GET** `/api/jurusan/:id`

Contoh: `GET /api/jurusan/550e8400-e29b-41d4-a716-446655440000`

Response sukses (200):
```json
{
  "success": true,
  "message": "Get jurusan successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "nama_jurusan": "Teknik Komputer dan Jaringan",
    "created_at": "2026-05-19 10:00:00",
    "updated_at": "2026-05-19 10:00:00"
  }
}
```

Response tidak ditemukan (404):
```json
{
  "success": false,
  "message": "data not found"
}
```

---

### 4. Update Jurusan

**PUT** `/api/jurusan/:id`

Header yang diperlukan:
```
Authorization: Bearer <token>
Content-Type: application/json
```

Request body:
```json
{
  "nama_jurusan": "Nama Jurusan Baru"
}
```

Response sukses (200):
```json
{
  "success": true,
  "message": "Update jurusan successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "nama_jurusan": "Nama Jurusan Baru",
    "created_at": "2026-05-19 10:00:00",
    "updated_at": "2026-05-19 10:30:00"
  }
}
```

---

### 5. Hapus Jurusan (Soft Delete)

**DELETE** `/api/jurusan/:id`

Header yang diperlukan:
```
Authorization: Bearer <token>
```

Response sukses (200):
```json
{
  "success": true,
  "message": "Delete jurusan successfully",
  "data": null
}
```

> **Penting untuk frontend:** Data yang dihapus tidak akan muncul di endpoint `GET /api/jurusan` maupun `GET /api/jurusan/:id`. Namun data masih ada di database dan bisa dipulihkan.

---

### 6. Pulihkan Jurusan

**PATCH** `/api/jurusan/:id/restore`

Header yang diperlukan:
```
Authorization: Bearer <token>
```

Response sukses (200):
```json
{
  "success": true,
  "message": "Restore jurusan successfully",
  "data": null
}
```

---

## Checklist Implementasi

Centang setiap item setelah selesai dikerjakan:

- [ ] `internal/modules/jurusan/model/jurusan_model.go` — dibuat
- [ ] `internal/database/migrate.go` — ditambahkan `jurusanmodel.Jurusan`
- [ ] `internal/database/seeders/jurusan_seeder.go` — dibuat
- [ ] `internal/database/seed.go` — ditambahkan `SeedJurusan`
- [ ] `internal/modules/jurusan/dto/jurusan_dto.go` — dibuat
- [ ] `internal/modules/jurusan/repository/jurusan_repository.go` — dibuat
- [ ] `internal/modules/jurusan/service/jurusan_service.go` — dibuat
- [ ] `internal/modules/jurusan/controller/jurusan_controller.go` — dibuat
- [ ] `internal/modules/jurusan/routes/jurusan_routes.go` — dibuat
- [ ] `cmd/server/main.go` — ditambahkan `SetupJurusanRoutes`
- [ ] `go build ./...` — tidak ada error
- [ ] Semua endpoint bisa diakses dan mengembalikan response yang benar
