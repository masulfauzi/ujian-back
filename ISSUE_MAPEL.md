# Issue: Implementasi Fitur Mapel (Subject Management)

## 📋 Deskripsi
Implementasi fitur manajemen mapel (mata pelajaran) termasuk:
1. Database migration untuk tabel mapel
2. Seeder untuk data awal mapel
3. CRUD API endpoints untuk mapel
4. Model dan Repository pattern
5. Unit tests

---

## 🎯 Objectives

✅ Membuat tabel `mapel` di database dengan soft delete support
✅ Membuat seeder untuk populate data mapel awal
✅ Setup CRUD endpoints untuk mapel
✅ Implementasi Repository pattern
✅ Dokumentasi API

---

## 📊 Database Schema

### Tabel: `mapel`

```sql
CREATE TABLE mapel (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_mapel VARCHAR(255) NOT NULL UNIQUE,
    kode_mapel VARCHAR(20) UNIQUE,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by UUID,
    updated_by UUID
);

CREATE INDEX idx_mapel_nama ON mapel(nama_mapel);
CREATE INDEX idx_mapel_deleted ON mapel(deleted_at);
```

### Field Details:
- **id** (UUID): Primary key
- **nama_mapel** (VARCHAR 255): Nama mata pelajaran, UNIQUE
- **kode_mapel** (VARCHAR 20): Kode singkat mapel (optional)
- **deskripsi** (TEXT): Deskripsi mapel (optional)
- **created_at** (TIMESTAMP): Waktu pembuatan record
- **updated_at** (TIMESTAMP): Waktu update terakhir
- **deleted_at** (TIMESTAMP NULL): Soft delete timestamp
- **created_by** (UUID): ID user yang membuat
- **updated_by** (UUID): ID user yang update terakhir

---

## 🛠️ Tahapan Implementasi

### **FASE 1: Database Setup**

#### Step 1.1: Buat Migration File
**File**: `migrations/[timestamp]_create_mapel_table.sql`

```sql
-- Up Migration
CREATE TABLE IF NOT EXISTS mapel (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_mapel VARCHAR(255) NOT NULL UNIQUE,
    kode_mapel VARCHAR(20) UNIQUE,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by UUID,
    updated_by UUID
);

CREATE INDEX IF NOT EXISTS idx_mapel_nama ON mapel(nama_mapel);
CREATE INDEX IF NOT EXISTS idx_mapel_deleted ON mapel(deleted_at);

-- Down Migration (comment out untuk safety)
-- DROP TABLE IF EXISTS mapel;
```

**Lokasi file**: `/migrations/[timestamp]_create_mapel_table.sql`

**Naming convention**: `[YYYYMMDDHHMMSS]_create_mapel_table.sql`
Contoh: `20260518100000_create_mapel_table.sql`

#### Step 1.2: Jalankan Migration

```bash
# Cek migrasi yang belum dijalankan
go run ./cmd/migrate migrate status

# Jalankan migration
go run ./cmd/migrate migrate up
```

**Expected Output**:
```
✓ Migration applied: 20260518100000_create_mapel_table.sql
```

---

### **FASE 2: Model & DTO**

#### Step 2.1: Buat Model File
**File**: `internal/modules/mapel/model/mapel_model.go`

```go
package model

import (
	"database/sql"
	"time"
)

type Mapel struct {
	ID        string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NamaMapel string       `gorm:"type:varchar(255);uniqueIndex" json:"nama_mapel"`
	KodeMapel string       `gorm:"type:varchar(20);uniqueIndex" json:"kode_mapel"`
	Deskripsi string       `gorm:"type:text" json:"deskripsi"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`
	CreatedBy string       `gorm:"type:uuid" json:"created_by"`
	UpdatedBy string       `gorm:"type:uuid" json:"updated_by"`
}

func (Mapel) TableName() string {
	return "mapel"
}
```

#### Step 2.2: Buat DTO File
**File**: `internal/modules/mapel/dto/mapel_dto.go`

```go
package dto

type CreateMapelRequest struct {
	NamaMapel string `json:"nama_mapel" validate:"required"`
	KodeMapel string `json:"kode_mapel" validate:"required,max=20"`
	Deskripsi string `json:"deskripsi"`
}

type UpdateMapelRequest struct {
	NamaMapel string `json:"nama_mapel" validate:"required"`
	KodeMapel string `json:"kode_mapel" validate:"required,max=20"`
	Deskripsi string `json:"deskripsi"`
}

type MapelResponse struct {
	ID        string `json:"id"`
	NamaMapel string `json:"nama_mapel"`
	KodeMapel string `json:"kode_mapel"`
	Deskripsi string `json:"deskripsi"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MapelListResponse struct {
	Data      []MapelResponse `json:"data"`
	Total     int64           `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalPage int             `json:"total_page"`
}
```

---

### **FASE 3: Repository & Service**

#### Step 3.1: Buat Repository Interface
**File**: `internal/modules/mapel/repository/mapel_repository.go`

```go
package repository

import (
	"backend/internal/modules/mapel/model"
	"context"
)

type MapelRepository interface {
	Create(ctx context.Context, mapel *model.Mapel) error
	GetByID(ctx context.Context, id string) (*model.Mapel, error)
	GetAll(ctx context.Context, page, pageSize int) ([]model.Mapel, int64, error)
	Update(ctx context.Context, mapel *model.Mapel) error
	Delete(ctx context.Context, id string) error // Soft delete
	Restore(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
}

type mapelRepository struct {
	db interface{}
}

// Implement semua methods dari interface
// Gunakan GORM untuk database operations
// Untuk soft delete, gunakan: db.Model(&mapel).Update("deleted_at", time.Now())
```

#### Step 3.2: Buat Service Layer
**File**: `internal/modules/mapel/service/mapel_service.go`

```go
package service

import (
	"backend/internal/modules/mapel/dto"
	"backend/internal/modules/mapel/model"
	"backend/internal/modules/mapel/repository"
	"context"
	"errors"
)

type MapelService interface {
	CreateMapel(ctx context.Context, req *dto.CreateMapelRequest) (*dto.MapelResponse, error)
	GetMapelByID(ctx context.Context, id string) (*dto.MapelResponse, error)
	GetAllMapel(ctx context.Context, page, pageSize int) (*dto.MapelListResponse, error)
	UpdateMapel(ctx context.Context, id string, req *dto.UpdateMapelRequest) (*dto.MapelResponse, error)
	DeleteMapel(ctx context.Context, id string) error
	RestoreMapel(ctx context.Context, id string) error
}

// Implement service methods
// Business logic: validasi, transformation, error handling
```

---

### **FASE 4: Controller & Routes**

#### Step 4.1: Buat Controller
**File**: `internal/modules/mapel/controller/mapel_controller.go`

```go
package controller

import (
	"backend/internal/helpers"
	"backend/internal/modules/mapel/dto"
	"backend/internal/modules/mapel/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type MapelController struct {
	service service.MapelService
}

func NewMapelController(service service.MapelService) *MapelController {
	return &MapelController{service: service}
}

// Endpoints:
// POST /api/mapel - Create
// GET /api/mapel - List
// GET /api/mapel/:id - Get by ID
// PUT /api/mapel/:id - Update
// DELETE /api/mapel/:id - Delete (soft)
// PATCH /api/mapel/:id/restore - Restore
```

#### Step 4.2: Setup Routes
**File**: `internal/modules/mapel/routes/mapel_routes.go`

```go
package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/mapel/controller"
	"backend/internal/modules/mapel/repository"
	"backend/internal/modules/mapel/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupMapelRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewMapelRepository(db)
	svc := service.NewMapelService(repo)
	ctrl := controller.NewMapelController(svc)

	api := app.Group("/api")
	mapel := api.Group("/mapel")

	// Routes
	mapel.Post("/", middleware.JWTAuth(), ctrl.CreateMapel)
	mapel.Get("/", ctrl.GetAllMapel)
	mapel.Get("/:id", ctrl.GetMapelByID)
	mapel.Put("/:id", middleware.JWTAuth(), ctrl.UpdateMapel)
	mapel.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteMapel)
	mapel.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreMapel)
}
```

#### Step 4.3: Register Routes di Main
**File**: `cmd/server/main.go`

Di function `setupRoutes()`, tambahkan:
```go
mapelroutes.SetupMapelRoutes(app, database.DB)
```

---

### **FASE 5: Seeder**

#### Step 5.1: Buat Seeder File
**File**: `internal/database/seeders/mapel_seeder.go`

```go
package seeders

import (
	"backend/internal/modules/mapel/model"
	"gorm.io/gorm"
	"time"
)

func SeedMapel(db *gorm.DB) error {
	mapels := []model.Mapel{
		{
			NamaMapel: "Matematika",
			KodeMapel: "MAT",
			Deskripsi: "Pelajaran Matematika",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			NamaMapel: "Bahasa Indonesia",
			KodeMapel: "IND",
			Deskripsi: "Pelajaran Bahasa Indonesia",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			NamaMapel: "Bahasa Inggris",
			KodeMapel: "ENG",
			Deskripsi: "Pelajaran Bahasa Inggris",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			NamaMapel: "IPA",
			KodeMapel: "IPA",
			Deskripsi: "Ilmu Pengetahuan Alam",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			NamaMapel: "IPS",
			KodeMapel: "IPS",
			Deskripsi: "Ilmu Pengetahuan Sosial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return db.CreateInBatches(mapels, 100).Error
}
```

#### Step 5.2: Register Seeder
**File**: `internal/database/seed.go`

```go
// Tambahkan ke function yang menjalankan semua seeders
if err := seeders.SeedMapel(db); err != nil {
	return fmt.Errorf("failed to seed mapel: %w", err)
}
```

#### Step 5.3: Jalankan Seeder

```bash
go run ./cmd/seed/main.go
# atau jika ada command untuk seed
go run ./cmd/migrate migrate seed
```

---

### **FASE 6: Testing**

#### Step 6.1: Test Migration
```bash
# Verify tabel dibuat
psql -U postgres -d ujian -c "\dt mapel"

# Expected output:
# public | mapel | table | postgres
```

#### Step 6.2: Test Seeder
```bash
psql -U postgres -d ujian -c "SELECT COUNT(*) FROM mapel;"

# Expected output: 5 records
```

#### Step 6.3: Test API Endpoints

**Create Mapel:**
```bash
curl -X POST http://localhost:3000/api/mapel \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "nama_mapel": "Seni Budaya",
    "kode_mapel": "SBD",
    "deskripsi": "Pelajaran Seni Budaya"
  }'
```

**Get All Mapel:**
```bash
curl http://localhost:3000/api/mapel?page=1&page_size=10
```

**Get By ID:**
```bash
curl http://localhost:3000/api/mapel/{id}
```

**Update Mapel:**
```bash
curl -X PUT http://localhost:3000/api/mapel/{id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "nama_mapel": "Seni Budaya Updated",
    "kode_mapel": "SBD",
    "deskripsi": "Updated deskripsi"
  }'
```

**Soft Delete:**
```bash
curl -X DELETE http://localhost:3000/api/mapel/{id} \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

**Restore (Undo Soft Delete):**
```bash
curl -X PATCH http://localhost:3000/api/mapel/{id}/restore \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

## 📁 Struktur File yang Dibuat

```
project-root/
├── migrations/
│   └── 20260518100000_create_mapel_table.sql
├── internal/
│   └── modules/
│       └── mapel/
│           ├── controller/
│           │   └── mapel_controller.go
│           ├── dto/
│           │   └── mapel_dto.go
│           ├── model/
│           │   └── mapel_model.go
│           ├── repository/
│           │   └── mapel_repository.go
│           ├── service/
│           │   └── mapel_service.go
│           └── routes/
│               └── mapel_routes.go
└── internal/
    └── database/
        └── seeders/
            └── mapel_seeder.go
```

---

## ✅ Checklist Implementasi

- [ ] Migration file dibuat dan dijalankan
- [ ] Model Mapel dibuat dengan soft delete support
- [ ] DTO untuk request/response dibuat
- [ ] Repository interface dan implementasi dibuat
- [ ] Service layer dibuat dengan business logic
- [ ] Controller dengan semua endpoints dibuat
- [ ] Routes di-register di main.go
- [ ] Seeder dibuat dan dijalankan (5 data awal)
- [ ] API tested dengan curl/Postman
- [ ] Soft delete berfungsi dengan baik
- [ ] Restore functionality berfungsi
- [ ] Database schema sudah dikonfigurasi di production
- [ ] Documentation (API docs/Swagger) ditambahkan

---

## 🔧 Catatan Teknis

### Soft Delete Implementation
- Gunakan `gorm.io` dengan `DeletedAt` field
- Query otomatis exclude record dengan `deleted_at` NOT NULL
- Untuk query termasuk deleted records, gunakan `.Unscoped()`

### UUID Generation
- PostgreSQL: `gen_random_uuid()` (require `pgcrypto` extension)
- GORM: `default:gen_random_uuid()`

### Timestamps
- `created_at`: Set otomatis saat record dibuat
- `updated_at`: Set otomatis saat record diupdate
- `deleted_at`: Set saat soft delete, NULL saat masih aktif

### Indexing
- `nama_mapel`: Indexed untuk fast search
- `deleted_at`: Indexed untuk efficient soft delete queries

---

## 🚀 Next Steps

Setelah semua tahapan selesai:

1. **Integration Testing** - Test full flow dari API endpoint
2. **API Documentation** - Update Swagger/OpenAPI docs
3. **Frontend Integration** - Connect dengan Vue.js frontend
4. **Validation & Error Handling** - Improve error messages
5. **Performance Optimization** - Query optimization jika diperlukan

---

## 📞 Questions & Support

Jika ada pertanyaan saat implementasi:
- Cek structure project yang sudah ada (auth, user modules)
- Follow naming convention & pattern yang sudah ada
- Refer ke `internal/modules/auth` sebagai reference implementation
