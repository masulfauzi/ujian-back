# Issue: Tambah Kolom `tingkat` pada Modul Kelas

## Deskripsi

Tambahkan kolom `tingkat` ke tabel `kelas` untuk menyimpan tingkat/grade secara eksplisit (X, XI, XII). Saat ini informasi tingkat hanya tersimpan secara implisit di dalam string `nama_kelas` (contoh: "X - TKJ"). Dengan kolom ini, filtering kelas per tingkat menjadi lebih mudah dan API lebih ekspresif untuk kebutuhan frontend.

**Ringkasan perubahan:**
1. Tambah kolom `tingkat` pada model dan migration
2. Update seeder untuk mengisi kolom `tingkat`
3. Drop tabel `kelas`, jalankan ulang migration dan seeder
4. Sesuaikan semua layer API (DTO, Repository, Service, Controller)
5. Update dokumentasi API untuk frontend

---

## Database Schema (Setelah Perubahan)

### Tabel: `kelas`

```sql
CREATE TABLE kelas (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_jurusan UUID NOT NULL REFERENCES jurusan(id),
    nama_kelas VARCHAR(255) NOT NULL,
    tingkat    VARCHAR(10) NOT NULL,  -- ← BARU: nilai: "X", "XI", "XII"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_kelas_deleted    ON kelas(deleted_at);
CREATE INDEX idx_kelas_id_jurusan ON kelas(id_jurusan);
CREATE INDEX idx_kelas_tingkat    ON kelas(tingkat);   -- ← BARU
```

### Keterangan Field Baru:
| Field | Tipe | Keterangan |
|---|---|---|
| `tingkat` | VARCHAR(10) | Tingkat/grade kelas. Nilai yang valid: `"X"`, `"XI"`, `"XII"` |

---

## Tahapan Implementasi

Ikuti urutan tahapan ini secara berurutan. Jangan lewati satu pun.

---

### TAHAP 1 — Update Model

**File:** `internal/modules/kelas/model/kelas_model.go`

Tambahkan field `Tingkat` ke struct `Kelas`:

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
	Tingkat   string         `gorm:"type:varchar(10);not null;index" json:"tingkat"`  // ← BARU
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Kelas) TableName() string {
	return "kelas"
}
```

> **Penjelasan:** Field `Tingkat` ditambahkan dengan tag `gorm:"type:varchar(10);not null;index"` agar GORM membuat kolom yang tepat dan membuat index untuk mempercepat query filter berdasarkan tingkat.

---

### TAHAP 2 — Update DTO

**File:** `internal/modules/kelas/dto/kelas_dto.go`

Tambahkan field `Tingkat` ke semua struct DTO yang relevan:

```go
package dto

type CreateKelasRequest struct {
	IDJurusan string `json:"id_jurusan" validate:"required"`
	NamaKelas string `json:"nama_kelas" validate:"required"`
	Tingkat   string `json:"tingkat" validate:"required"`  // ← BARU
}

type UpdateKelasRequest struct {
	IDJurusan string `json:"id_jurusan" validate:"required"`
	NamaKelas string `json:"nama_kelas" validate:"required"`
	Tingkat   string `json:"tingkat" validate:"required"`  // ← BARU
}

type KelasResponse struct {
	ID        string `json:"id"`
	IDJurusan string `json:"id_jurusan"`
	NamaKelas string `json:"nama_kelas"`
	Tingkat   string `json:"tingkat"`   // ← BARU
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

### TAHAP 3 — Update Repository

**File:** `internal/modules/kelas/repository/kelas_repository.go`

Gabungkan `GetAll` dan `GetByJurusan` menjadi satu method yang lebih fleksibel, dan tambahkan filter `tingkat`:

```go
package repository

import (
	"backend/internal/modules/kelas/model"

	"gorm.io/gorm"
)

type KelasRepository interface {
	Create(kelas *model.Kelas) error
	GetByID(id string) (*model.Kelas, error)
	GetAll(page, pageSize int, idJurusan string, tingkat string) ([]model.Kelas, int64, error)  // ← DIUBAH
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

// GetAll menggantikan GetAll dan GetByJurusan sebelumnya.
// idJurusan dan tingkat bersifat opsional — kosongkan string untuk menonaktifkan filter.
func (r *kelasRepository) GetAll(page, pageSize int, idJurusan string, tingkat string) ([]model.Kelas, int64, error) {
	var kelasList []model.Kelas
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := r.db.Model(&model.Kelas{}).Where("deleted_at IS NULL")

	if idJurusan != "" {
		query = query.Where("id_jurusan = ?", idJurusan)
	}
	if tingkat != "" {
		query = query.Where("tingkat = ?", tingkat)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(pageSize).Find(&kelasList).Error

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

> **Penjelasan perubahan:** `GetByJurusan` dihapus. `GetAll` kini menerima dua parameter filter opsional (`idJurusan` dan `tingkat`). Jika parameter dikosongkan (`""`), filter tidak diterapkan. Ini menyederhanakan interface tanpa kehilangan fungsionalitas.

---

### TAHAP 4 — Update Service

**File:** `internal/modules/kelas/service/kelas_service.go`

Sesuaikan service untuk menghapus `GetByJurusan` dan menambahkan `tingkat` ke semua operasi:

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
	GetAllKelas(page, pageSize int, idJurusan string, tingkat string) (*dto.KelasListResponse, error)  // ← DIUBAH
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
		Tingkat:   req.Tingkat,   // ← BARU
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

func (s *kelasService) GetAllKelas(page, pageSize int, idJurusan string, tingkat string) (*dto.KelasListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	kelasList, total, err := s.repo.GetAll(page, pageSize, idJurusan, tingkat)  // ← DIUBAH
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
	kelas.Tingkat   = req.Tingkat   // ← BARU

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
		Tingkat:   kelas.Tingkat,   // ← BARU
		CreatedAt: kelas.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: kelas.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
```

---

### TAHAP 5 — Update Controller

**File:** `internal/modules/kelas/controller/kelas_controller.go`

Tambahkan parsing query param `tingkat` di `GetAllKelas`:

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
	page     := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idJurusan := ctx.Query("id_jurusan", "")
	tingkat   := ctx.Query("tingkat", "")   // ← BARU

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllKelas(pageNum, pageSizeNum, idJurusan, tingkat)  // ← DIUBAH
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

> **Catatan:** File `routes/kelas_routes.go` tidak perlu diubah karena endpoint URL tidak berubah.

---

### TAHAP 6 — Update Seeder

**File:** `internal/database/seeders/kelas_seeder.go`

Ganti seluruh isi file. Seeder kini mengisi kolom `Tingkat` secara eksplisit:

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
				Tingkat:   tingkat,   // ← BARU: isi kolom tingkat secara eksplisit
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	return db.CreateInBatches(kelasList, 100).Error
}
```

---

### TAHAP 7 — Drop Tabel Kelas di Database

> **PENTING:** Kolom `NOT NULL` tidak bisa ditambahkan ke tabel yang sudah berisi data via `AutoMigrate` tanpa nilai default. Cara paling aman adalah drop tabel lama, lalu biarkan `AutoMigrate` membuat ulang tabel dengan schema baru.

Jalankan perintah SQL berikut di PostgreSQL (via psql, pgAdmin, atau tool database lain yang digunakan):

```sql
DROP TABLE IF EXISTS kelas;
```

> **Peringatan:** Perintah ini akan menghapus seluruh data di tabel `kelas`. Data akan diisi ulang oleh seeder pada tahap berikutnya. Pastikan tidak ada data production penting di tabel ini sebelum menjalankan perintah ini.

---

### TAHAP 8 — Jalankan Migration dan Seeder

Setelah TAHAP 1–7 selesai, jalankan server untuk menjalankan migration dan seeder secara otomatis:

```bash
make run
# atau
go run cmd/server/main.go
```

Server akan:
1. Menjalankan `RunMigrations` → GORM membuat ulang tabel `kelas` dengan schema baru (termasuk kolom `tingkat`)
2. Menjalankan `RunSeeders` → `SeedKelas` mengisi tabel `kelas` dengan data awal

---

### TAHAP 9 — Verifikasi

**Cek kompilasi dulu:**

```bash
go build ./...
```

Tidak boleh ada output error. Jika ada error, baca pesan error dan perbaiki sebelum lanjut.

**Cek tabel dan data di database:**

```sql
-- Pastikan kolom tingkat ada
SELECT column_name, data_type
FROM information_schema.columns
WHERE table_name = 'kelas'
ORDER BY ordinal_position;

-- Cek data seeder masuk dengan kolom tingkat
SELECT k.nama_kelas, k.tingkat, j.nama_jurusan
FROM kelas k
JOIN jurusan j ON j.id = k.id_jurusan
WHERE k.deleted_at IS NULL
ORDER BY j.nama_jurusan, k.tingkat;
```

**Cek API endpoint dengan curl:**

```bash
# Ambil semua kelas
curl http://localhost:3000/api/kelas

# Filter berdasarkan tingkat X
curl "http://localhost:3000/api/kelas?tingkat=X"

# Filter berdasarkan tingkat dan jurusan sekaligus
curl "http://localhost:3000/api/kelas?tingkat=XI&id_jurusan=<uuid-jurusan>"
```

---

## API Endpoints (Setelah Perubahan)

Base URL: `http://localhost:3000`

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| `POST` | `/api/kelas` | JWT Required | Tambah kelas baru |
| `GET` | `/api/kelas` | Public | Ambil semua kelas (pagination + filter opsional) |
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
  "nama_kelas": "X TKJ 1",
  "tingkat": "X"
}
```

> **Nilai `tingkat` yang valid:** `"X"`, `"XI"`, `"XII"`

Response sukses (201):
```json
{
  "success": true,
  "message": "Create kelas successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551111",
    "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X TKJ 1",
    "tingkat": "X",
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

Query parameters (semua opsional):
| Parameter | Default | Keterangan |
|---|---|---|
| `page` | `1` | Nomor halaman |
| `page_size` | `10` | Jumlah data per halaman |
| `id_jurusan` | _(kosong)_ | Filter kelas berdasarkan UUID jurusan tertentu |
| `tingkat` | _(kosong)_ | Filter kelas berdasarkan tingkat (`X`, `XI`, `XII`) |

Contoh request:
```
GET /api/kelas
GET /api/kelas?tingkat=X
GET /api/kelas?tingkat=XI&page=1&page_size=10
GET /api/kelas?id_jurusan=550e8400-e29b-41d4-a716-446655440000
GET /api/kelas?id_jurusan=550e8400-e29b-41d4-a716-446655440000&tingkat=X
GET /api/kelas?id_jurusan=550e8400-e29b-41d4-a716-446655440000&tingkat=X&page=1&page_size=100
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
        "tingkat": "X",
        "created_at": "2026-05-19 10:00:00",
        "updated_at": "2026-05-19 10:00:00"
      },
      {
        "id": "661e9511-f30c-52e5-b827-557766551112",
        "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
        "nama_kelas": "XI - Teknik Komputer dan Jaringan",
        "tingkat": "XI",
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

> **Tips untuk frontend:**
> - Dropdown kelas berdasarkan jurusan: `GET /api/kelas?id_jurusan=<uuid>&page=1&page_size=100`
> - Dropdown kelas berdasarkan tingkat: `GET /api/kelas?tingkat=X&page=1&page_size=100`
> - Kombinasi filter: `GET /api/kelas?id_jurusan=<uuid>&tingkat=X&page=1&page_size=100`

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
    "tingkat": "X",
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
  "nama_kelas": "X TKJ 1",
  "tingkat": "X"
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
    "tingkat": "X",
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

- [ ] `internal/modules/kelas/model/kelas_model.go` — tambahkan field `Tingkat`
- [ ] `internal/modules/kelas/dto/kelas_dto.go` — tambahkan field `Tingkat` di semua struct
- [ ] `internal/modules/kelas/repository/kelas_repository.go` — gabung `GetAll`+`GetByJurusan` jadi satu method dengan 2 filter
- [ ] `internal/modules/kelas/service/kelas_service.go` — update signature `GetAllKelas`, tambahkan `Tingkat` di Create/Update
- [ ] `internal/modules/kelas/controller/kelas_controller.go` — parse query param `tingkat` di `GetAllKelas`
- [ ] `internal/database/seeders/kelas_seeder.go` — tambahkan `Tingkat` di setiap entri kelas
- [ ] `DROP TABLE IF EXISTS kelas;` — jalankan di database **sebelum** menjalankan server
- [ ] `go build ./...` — tidak ada error kompilasi
- [ ] `make run` — migration dan seeder berhasil (tabel `kelas` muncul dengan kolom `tingkat`)
- [ ] Verifikasi SQL: tabel berisi data dengan kolom `tingkat` terisi
- [ ] Verifikasi API: `GET /api/kelas?tingkat=X` mengembalikan data yang difilter dengan benar
