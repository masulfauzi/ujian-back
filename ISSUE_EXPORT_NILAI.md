# Issue: Export Nilai Ujian ke Excel per Kelas (ZIP Download)

## Deskripsi

Tambahkan endpoint baru untuk mengekspor nilai ujian berdasarkan `id_jadwal`.
Sistem akan membuat **satu file Excel per kelas** yang terdaftar pada jadwal tersebut,
kemudian mengemas seluruh file Excel ke dalam **satu file ZIP** yang langsung bisa diunduh oleh browser.

---

## Konteks Proyek

Proyek ini menggunakan:
- **Bahasa:** Go (Golang), module name: `backend`
- **Framework HTTP:** [Fiber v2](https://gofiber.io/) — `github.com/gofiber/fiber/v2`
- **ORM:** GORM dengan driver PostgreSQL
- **Library Excel:** `github.com/xuri/excelize/v2` — **sudah ada di `go.mod`**, tidak perlu `go get`
- **Library ZIP:** `archive/zip` — paket standar Go, tidak perlu instalasi tambahan
- **Struktur Modul:** `internal/modules/<nama_modul>/` berisi `model/`, `dto/`, `repository/`, `service/`, `controller/`, `routes/`
- **Auth Middleware:** `backend/internal/middleware` — `middleware.JWTAuth()`
- **Helper Response:** `backend/internal/helpers` — `helpers.ErrorResponse()` / `helpers.SuccessResponse()`

Sebelum mulai, pelajari file-file berikut untuk memahami pola yang digunakan:
- `internal/modules/nilai/repository/nilai_repository.go` — contoh query JOIN multi-tabel
- `internal/modules/nilai/service/nilai_service.go` — contoh service dengan akses `*gorm.DB` langsung
- `internal/modules/nilai/controller/nilai_controller.go` — contoh handler Fiber
- `internal/modules/nilai/routes/nilai_routes.go` — cara registrasi route

---

## Relasi Tabel yang Digunakan

```
jadwal (id) ──< jadwal_kelas (id_jadwal, id_kelas) >── kelas (id, nama_kelas)
                                                              ^
nilai (id_jadwal) ──> [jadwal]             peserta (id_kelas) ┘
nilai (id_peserta) ──> peserta (id, nama, username, id_kelas)
```

### Tabel `jadwal_kelas`
| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID | Primary key |
| `id_jadwal` | UUID | FK ke `jadwal.id` |
| `id_kelas` | UUID | FK ke `kelas.id` |

### Tabel `kelas`
| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID | Primary key |
| `nama_kelas` | varchar(255) | Nama kelas, misal "X RPL 1" |

### Tabel `peserta`
| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID | Primary key |
| `nama` | varchar(255) | Nama lengkap |
| `username` | varchar(100) | Username peserta |
| `id_kelas` | UUID | FK ke `kelas.id` |

### Tabel `nilai`
| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID | Primary key |
| `id_peserta` | UUID | FK ke `peserta.id` |
| `id_jadwal` | UUID | FK ke `jadwal.id` |
| `nilai` | float | Nilai akhir (0–100) |
| `wkt_mulai` | timestamp | Waktu peserta mulai ujian |
| `wkt_selesai` | timestamp | Waktu peserta selesai ujian |
| `deleted_at` | timestamp | NULL = aktif |

---

## Alur Logika Endpoint

```
Client → GET /api/nilai/export/:id_jadwal
           │
           ▼
1. Query jadwal_kelas WHERE id_jadwal = ? → dapat list id_kelas
2. Untuk tiap id_kelas:
     a. Query kelas WHERE id = id_kelas → dapat nama_kelas
     b. Query nilai JOIN peserta
           WHERE nilai.id_jadwal = ?
             AND peserta.id_kelas = ?
             AND nilai.deleted_at IS NULL
           ORDER BY peserta.nama ASC
     c. Buat file Excel dari hasil query
        - Nama file: {nama_kelas}.xlsx (spasi → underscore)
        - Sheet: "Nilai Ujian"
        - Header: No | Nama Peserta | Username | Nilai | Waktu Mulai | Waktu Selesai
3. Kemas semua .xlsx ke dalam satu .zip
4. Return ZIP sebagai binary response (Content-Type: application/zip)
```

---

## Tahapan Implementasi

### Tahap 1 — Tambah Method Repository

**File yang dimodifikasi:** `internal/modules/nilai/repository/nilai_repository.go`

**Tambahkan struct baru** di bawah struct `NilaiWithDetail` yang sudah ada:

```go
type NilaiExportRow struct {
    NamaPeserta string  `gorm:"column:nama_peserta"`
    Username    string  `gorm:"column:username"`
    Nilai       float64 `gorm:"column:nilai"`
    WktMulai    *string `gorm:"column:wkt_mulai"`
    WktSelesai  *string `gorm:"column:wkt_selesai"`
}
```

**Tambahkan method ke interface `NilaiRepository`:**

```go
// di dalam interface NilaiRepository
GetByJadwalAndKelas(idJadwal, idKelas string) ([]NilaiExportRow, error)
```

**Tambahkan implementasinya** di struct `nilaiRepository`:

```go
func (r *nilaiRepository) GetByJadwalAndKelas(idJadwal, idKelas string) ([]NilaiExportRow, error) {
    var results []NilaiExportRow
    err := r.db.
        Table("nilai").
        Select(`
            peserta.nama AS nama_peserta,
            peserta.username,
            nilai.nilai,
            TO_CHAR(nilai.wkt_mulai AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD HH24:MI:SS') AS wkt_mulai,
            TO_CHAR(nilai.wkt_selesai AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD HH24:MI:SS') AS wkt_selesai
        `).
        Joins("INNER JOIN peserta ON nilai.id_peserta = peserta.id").
        Where("nilai.id_jadwal = ? AND peserta.id_kelas = ? AND nilai.deleted_at IS NULL", idJadwal, idKelas).
        Order("peserta.nama ASC").
        Scan(&results).Error
    return results, err
}
```

---

### Tahap 2 — Tambah Method Service

**File yang dimodifikasi:** `internal/modules/nilai/service/nilai_service.go`

**Tambahkan import** yang diperlukan di bagian atas file (tambahkan ke blok `import` yang sudah ada):

```go
import (
    // ... import yang sudah ada ...
    "archive/zip"
    "bytes"
    "fmt"
    "strings"

    "github.com/xuri/excelize/v2"
)
```

**Tambahkan struct hasil export** — letakkan di bawah struct atau konstanta yang sudah ada:

```go
type ExportResult struct {
    ZipBytes  []byte
    NamaUjian string
}
```

**Tambahkan method ke interface `NilaiService`:**

```go
// di dalam interface NilaiService
ExportNilaiByJadwal(idJadwal string) (*ExportResult, error)
```

**Tambahkan implementasinya** di struct `nilaiService`:

```go
func (s *nilaiService) ExportNilaiByJadwal(idJadwal string) (*ExportResult, error) {
    // 1. Ambil nama ujian dari tabel jadwal
    var namaUjian string
    if err := s.db.Table("jadwal").
        Select("nama_ujian").
        Where("id = ? AND deleted_at IS NULL", idJadwal).
        Scan(&namaUjian).Error; err != nil || namaUjian == "" {
        return nil, errors.New("jadwal tidak ditemukan")
    }

    // 2. Ambil semua kelas yang terdaftar pada jadwal ini
    type KelasRow struct {
        IDKelas   string `gorm:"column:id_kelas"`
        NamaKelas string `gorm:"column:nama_kelas"`
    }
    var kelasList []KelasRow
    if err := s.db.Table("jadwal_kelas").
        Select("jadwal_kelas.id_kelas, kelas.nama_kelas").
        Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
        Where("jadwal_kelas.id_jadwal = ?", idJadwal).
        Scan(&kelasList).Error; err != nil {
        return nil, err
    }
    if len(kelasList) == 0 {
        return nil, errors.New("tidak ada kelas yang terdaftar pada jadwal ini")
    }

    // 3. Buat ZIP di memory
    var zipBuf bytes.Buffer
    zipWriter := zip.NewWriter(&zipBuf)

    for _, kelas := range kelasList {
        // 3a. Ambil data nilai peserta untuk kelas ini
        rows, err := s.repo.GetByJadwalAndKelas(idJadwal, kelas.IDKelas)
        if err != nil {
            return nil, fmt.Errorf("gagal ambil data kelas %s: %w", kelas.NamaKelas, err)
        }

        // 3b. Buat file Excel
        xlsx := excelize.NewFile()
        sheet := "Nilai Ujian"
        xlsx.SetSheetName("Sheet1", sheet)

        // Header
        headers := []string{"No", "Nama Peserta", "Username", "Nilai", "Waktu Mulai", "Waktu Selesai"}
        for col, h := range headers {
            cell, _ := excelize.CoordinatesToCellName(col+1, 1)
            xlsx.SetCellValue(sheet, cell, h)
        }

        // Data rows
        for i, row := range rows {
            rowNum := i + 2
            xlsx.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), i+1)
            xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), row.NamaPeserta)
            xlsx.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), row.Username)
            xlsx.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), row.Nilai)

            wktMulai := "-"
            if row.WktMulai != nil {
                wktMulai = *row.WktMulai
            }
            xlsx.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), wktMulai)

            wktSelesai := "-"
            if row.WktSelesai != nil {
                wktSelesai = *row.WktSelesai
            }
            xlsx.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), wktSelesai)
        }

        // 3c. Tulis Excel ke buffer lalu masukkan ke ZIP
        var xlsBuf bytes.Buffer
        if err := xlsx.Write(&xlsBuf); err != nil {
            return nil, fmt.Errorf("gagal tulis excel kelas %s: %w", kelas.NamaKelas, err)
        }

        safeNama := strings.ReplaceAll(kelas.NamaKelas, " ", "_")
        zipEntry, err := zipWriter.Create(fmt.Sprintf("%s.xlsx", safeNama))
        if err != nil {
            return nil, err
        }
        if _, err := zipEntry.Write(xlsBuf.Bytes()); err != nil {
            return nil, err
        }
    }

    if err := zipWriter.Close(); err != nil {
        return nil, err
    }

    return &ExportResult{
        ZipBytes:  zipBuf.Bytes(),
        NamaUjian: namaUjian,
    }, nil
}
```

---

### Tahap 3 — Tambah Handler Controller

**File yang dimodifikasi:** `internal/modules/nilai/controller/nilai_controller.go`

**Tambahkan import** berikut ke blok import (jika belum ada):

```go
import (
    // ... import yang sudah ada ...
    "fmt"
    "strings"
)
```

**Tambahkan method `ExportNilai`** di akhir struct `NilaiController`:

```go
func (c *NilaiController) ExportNilai(ctx *fiber.Ctx) error {
    idJadwal := ctx.Params("id_jadwal")
    if idJadwal == "" {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_jadwal tidak boleh kosong", nil)
    }

    result, err := c.service.ExportNilaiByJadwal(idJadwal)
    if err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
    }

    safeNama := strings.ReplaceAll(result.NamaUjian, " ", "_")
    filename := fmt.Sprintf("export_nilai_%s.zip", safeNama)

    ctx.Set("Content-Type", "application/zip")
    ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
    return ctx.Send(result.ZipBytes)
}
```

---

### Tahap 4 — Daftarkan Route

**File yang dimodifikasi:** `internal/modules/nilai/routes/nilai_routes.go`

Tambahkan satu baris route baru di dalam fungsi `SetupNilaiRoutes`, **sebelum** route `nilai.Get("/:id", ...)`:

```go
// Export nilai per jadwal ke ZIP berisi file Excel per kelas
nilai.Get("/export/:id_jadwal", middleware.JWTAuth(), ctrl.ExportNilai)
```

Posisi penambahan di dalam fungsi (urutannya penting agar tidak bentrok dengan `/:id`):

```go
func SetupNilaiRoutes(app *fiber.App, db *gorm.DB) {
    // ... inisialisasi repo/svc/ctrl yang sudah ada ...

    api   := app.Group("/api")
    nilai := api.Group("/nilai")

    // Dua-segmen path: aman dari konflik dengan /:id
    nilai.Get("/peserta/:id_peserta", ctrl.GetNilaiByPeserta)
    nilai.Get("/jadwal/:id_jadwal", ctrl.GetNilaiByJadwal)
    nilai.Get("/export/:id_jadwal", middleware.JWTAuth(), ctrl.ExportNilai)  // ← TAMBAHKAN DI SINI
    nilai.Post("/mulai-ujian/:id_jadwal", middleware.JWTAuth(), ctrl.MulaiUjian)

    nilai.Post("/", middleware.JWTAuth(), ctrl.CreateNilai)
    nilai.Get("/", ctrl.GetAllNilai)
    nilai.Get("/:id", ctrl.GetNilaiByID)
    nilai.Put("/:id", middleware.JWTAuth(), ctrl.UpdateNilai)
    nilai.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteNilai)
    nilai.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreNilai)
}
```

---

### Tahap 5 — Verifikasi Build

Setelah semua perubahan selesai, jalankan perintah berikut dari root proyek untuk memastikan tidak ada error kompilasi:

```bash
go build ./...
```

Jika ada error `undefined`, pastikan semua import sudah ditambahkan dengan benar.

---

### Tahap 6 — Test Manual

Gunakan `curl` atau Postman untuk menguji endpoint:

```bash
# Ganti TOKEN dengan JWT yang valid (login dulu)
# Ganti ID_JADWAL dengan UUID jadwal yang ada di database

curl -X GET \
  "http://localhost:3000/api/nilai/export/<ID_JADWAL>" \
  -H "Authorization: Bearer <TOKEN>" \
  --output export_nilai.zip

# Verifikasi isi ZIP
unzip -l export_nilai.zip
```

**Kondisi yang harus ditest:**
1. `id_jadwal` valid dengan beberapa kelas → harus menghasilkan ZIP berisi N file .xlsx
2. `id_jadwal` valid tapi tidak ada kelas terdaftar → harus return 400 dengan pesan error
3. `id_jadwal` tidak ditemukan di tabel `jadwal` → harus return 400 dengan pesan error
4. Request tanpa token Authorization → harus return 401 Unauthorized
5. Kelas yang ada tapi tidak ada peserta yang mengerjakan ujian → file Excel tetap dibuat dengan hanya header (baris data kosong)

---

## Dokumentasi API untuk Frontend

### Endpoint

```
GET /api/nilai/export/:id_jadwal
```

### Authentication

Wajib. Kirim JWT token di header:
```
Authorization: Bearer <token>
```

### Path Parameter

| Parameter | Tipe | Wajib | Keterangan |
|-----------|------|-------|------------|
| `id_jadwal` | string (UUID) | Ya | ID jadwal ujian yang akan diekspor |

### Response Sukses

- **Status:** `200 OK`
- **Content-Type:** `application/zip`
- **Content-Disposition:** `attachment; filename="export_nilai_{nama_ujian}.zip"`
- **Body:** Binary data file ZIP

Isi file ZIP:
```
export_nilai_Ujian_Tengah_Semester.zip
├── X_RPL_1.xlsx
├── X_RPL_2.xlsx
└── XI_TKJ.xlsx
```

Setiap file Excel memiliki format:

| No | Nama Peserta | Username | Nilai | Waktu Mulai | Waktu Selesai |
|----|-------------|----------|-------|-------------|---------------|
| 1  | Andi Prasetyo | andi.p | 85.0 | 2024-01-15 08:00:00 | 2024-01-15 09:30:00 |
| 2  | Budi Santoso | budi.s | 72.5 | 2024-01-15 08:05:00 | 2024-01-15 09:28:00 |

> Peserta yang belum mengerjakan ujian **tidak muncul** di file Excel (hanya peserta yang sudah ada recordnya di tabel `nilai`).
> Kolom Waktu Mulai / Waktu Selesai berisi tanda `-` jika nilainya NULL.

### Response Error

| Status | Kondisi | Contoh Body |
|--------|---------|-------------|
| 400 | `id_jadwal` tidak ditemukan | `{"status": false, "message": "jadwal tidak ditemukan", "data": null}` |
| 400 | Tidak ada kelas terdaftar | `{"status": false, "message": "tidak ada kelas yang terdaftar pada jadwal ini", "data": null}` |
| 401 | Token tidak ada / tidak valid | `{"status": false, "message": "Unauthorized", "data": null}` |
| 500 | Error server | `{"status": false, "message": "...", "data": null}` |

### Contoh Implementasi Frontend (JavaScript / Fetch API)

```javascript
async function exportNilai(idJadwal, token) {
    const response = await fetch(`/api/nilai/export/${idJadwal}`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Gagal mengunduh file');
    }

    // Ambil nama file dari header Content-Disposition
    const disposition = response.headers.get('Content-Disposition');
    const filename = disposition
        ? disposition.match(/filename="(.+)"/)?.[1] ?? 'export_nilai.zip'
        : 'export_nilai.zip';

    // Trigger download
    const blob = await response.blob();
    const url  = URL.createObjectURL(blob);
    const a    = document.createElement('a');
    a.href     = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
}

// Penggunaan:
// exportNilai('uuid-jadwal-disini', 'jwt-token-disini');
```

### Contoh Implementasi Frontend (Axios)

```javascript
async function exportNilai(idJadwal, token) {
    const response = await axios.get(`/api/nilai/export/${idJadwal}`, {
        headers: { Authorization: `Bearer ${token}` },
        responseType: 'blob',  // PENTING: harus 'blob' agar bisa diunduh
    });

    const disposition = response.headers['content-disposition'];
    const filename = disposition
        ? disposition.match(/filename="(.+)"/)?.[1] ?? 'export_nilai.zip'
        : 'export_nilai.zip';

    const url = URL.createObjectURL(new Blob([response.data]));
    const a   = document.createElement('a');
    a.href    = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
}
```

---

## Ringkasan File yang Diubah

| File | Jenis Perubahan |
|------|----------------|
| `internal/modules/nilai/repository/nilai_repository.go` | Tambah struct `NilaiExportRow`, tambah method `GetByJadwalAndKelas` ke interface dan implementasinya |
| `internal/modules/nilai/service/nilai_service.go` | Tambah struct `ExportResult`, tambah method `ExportNilaiByJadwal` ke interface dan implementasinya; tambah import `archive/zip`, `bytes`, `fmt`, `strings`, `excelize` |
| `internal/modules/nilai/controller/nilai_controller.go` | Tambah handler `ExportNilai`; tambah import `fmt`, `strings` |
| `internal/modules/nilai/routes/nilai_routes.go` | Daftarkan route `GET /export/:id_jadwal` |

**Tidak ada file baru yang perlu dibuat.** Semua perubahan dilakukan di file yang sudah ada dalam modul `nilai`.

---

## Catatan Penting

1. **Library excelize sudah ada** — tidak perlu `go get`, cukup tambahkan import di file service.
2. **Library archive/zip adalah stdlib** — tidak perlu instalasi.
3. **Urutan route di Fiber penting** — `/export/:id_jadwal` harus didaftarkan **sebelum** `/:id` agar tidak tertangkap oleh route `:id`.
4. **Waktu ditampilkan dalam WIB (UTC+7)** — query SQL sudah mengonversi via `AT TIME ZONE 'Asia/Jakarta'`.
5. **Peserta yang tidak mengerjakan ujian tidak muncul** — hanya peserta yang ada di tabel `nilai` yang ditampilkan.
6. **Nama file Excel menggunakan underscore** — spasi pada `nama_kelas` diubah ke `_` agar kompatibel dengan semua OS.
