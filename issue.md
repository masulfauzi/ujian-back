# Issue: Auto-Generate Jawaban Kosong saat Mulai Ujian

## 📋 Deskripsi

Saat ini endpoint `POST /api/nilai/mulai-ujian/:id_jadwal` hanya **membuat record di tabel `nilai`**. Belum ada logika untuk membuat record jawaban kosong di tabel `jawaban`.

**Tujuan issue ini:** Setelah insert nilai berhasil (kasus pertama kali mulai ujian), endpoint harus **otomatis generate record jawaban kosong** untuk setiap soal yang ada di bank soal jadwal tersebut, dengan **urutan acak** (random).

Manfaatnya:
- Frontend langsung punya daftar soal urut + tempat menyimpan jawaban → tinggal PUT/PATCH jawaban saat peserta jawab.
- Setiap peserta dapat **urutan soal berbeda** (random) → mencegah contek-mencontek.
- Audit lengkap: jika peserta tidak menjawab soal X, record `jawaban` tetap ada dengan field `jawaban = NULL`.

---

## 🎯 Behavior yang Diharapkan

### Kondisi: Peserta pertama kali mulai ujian (`wkt_mulai` baru di-set)

1. Insert ke tabel `nilai` (sudah ada — pertahankan).
2. Query jadwal untuk dapatkan `id_bank_soal`.
3. Query semua soal di bank soal tersebut: `SELECT * FROM soal WHERE id_bank_soal = ? AND deleted_at IS NULL ORDER BY RANDOM()`.
4. Untuk setiap soal hasil query (urut sesuai random order), insert ke tabel `jawaban` dengan mapping:
   | Kolom | Diisi dari |
   |-------|-----------|
   | `id_nilai` | ID nilai yang baru di-insert |
   | `id_soal` | `soal.id` |
   | `id_peserta` | `idPeserta` (dari JWT) |
   | `no_urut` | Index urut hasil random (1, 2, 3, ...) |
   | `jawaban` | `NULL` |
   | `is_benar` | `NULL` |
5. Return detail nilai (sama seperti sekarang) **+** info tambahan: jumlah jawaban yang di-generate.

### Kondisi: Resume (record nilai sudah ada, `wkt_selesai` masih NULL)

- **JANGAN** generate ulang jawaban. Record jawaban sudah ada dari sebelumnya.
- Cukup return detail nilai existing (sama seperti sekarang).

### Kondisi: Sudah selesai (`wkt_selesai` != NULL)

- Return error 400 (sama seperti sekarang).

### Atomicity (PENTING!)

Semua operasi (insert nilai + query soal + bulk insert jawaban) **HARUS dalam satu transaction**. Kalau salah satu langkah gagal → rollback semua. Jangan sampai ada record nilai tanpa jawaban.

---

## 🗂️ Konteks: Relasi Antar Tabel

```
peserta ──┐
          │
jadwal ───┼──> nilai (id_peserta, id_jadwal, ...)
          │       │
          │       └──> jawaban (id_nilai, id_soal, id_peserta, no_urut, jawaban, is_benar)
          │                        │
          │                        └─> soal (id_bank_soal, ...)
          │
          └──> bank_soal (jadwal.id_bank_soal)
                  │
                  └──> soal (soal.id_bank_soal)
```

**Catatan penting:**
- Tabel `soal` **TIDAK PUNYA** kolom `id_jadwal`. Yang ada: `id_bank_soal`.
- Untuk dapatkan soal milik jadwal: ambil `id_bank_soal` dari `jadwal`, lalu query `soal WHERE id_bank_soal = ?`.

---

## ⚠️ Perubahan Skema yang Diperlukan

Tabel `jawaban` saat ini punya constraint:
- `jawaban varchar(1) NOT NULL`
- `is_benar boolean NOT NULL DEFAULT false`

Karena bulk-insert akan mengisi kedua kolom dengan `NULL`, kita harus **menghilangkan `NOT NULL`** dari kedua kolom tersebut.

**Solusi:**
- Ubah tipe field di model Go dari `string`/`bool` jadi `*string`/`*bool` (pointer → nullable).
- Tambah raw SQL di `migrate.go` untuk drop constraint NOT NULL (AutoMigrate tidak akan otomatis drop constraint existing).

---

## 📝 Referensi File yang Relevan

Sebelum mulai, baca file-file ini:

- [internal/modules/jawaban/model/jawaban_model.go](internal/modules/jawaban/model/jawaban_model.go) — model yang akan diubah jadi nullable
- [internal/modules/jawaban/dto/jawaban_dto.go](internal/modules/jawaban/dto/jawaban_dto.go) — response DTO yang ikut berubah
- [internal/modules/jawaban/repository/jawaban_repository.go](internal/modules/jawaban/repository/jawaban_repository.go) — `JawabanWithDetail` struct yang ikut berubah, + tambah `BulkCreate`
- [internal/modules/jawaban/service/jawaban_service.go](internal/modules/jawaban/service/jawaban_service.go) — sesuaikan pemanggilan setelah model nullable
- [internal/modules/soal/model/soal_model.go](internal/modules/soal/model/soal_model.go) — sumber field `id_bank_soal`
- [internal/modules/jadwal/model/jadwal_model.go](internal/modules/jadwal/model/jadwal_model.go) — sumber field `id_bank_soal`
- [internal/modules/nilai/service/nilai_service.go](internal/modules/nilai/service/nilai_service.go) — method `MulaiUjian` yang akan ditambah logika baru
- [internal/modules/nilai/routes/nilai_routes.go](internal/modules/nilai/routes/nilai_routes.go) — wiring constructor service yang berubah
- [internal/database/migrate.go](internal/database/migrate.go) — tempat tambah ALTER TABLE

---

## 🪜 Tahapan Implementasi

> **Kerjakan urut dari Tahap 1 → Tahap 10.** Setiap tahap selesai, jalankan `go build ./...` untuk pastikan tidak ada error compile.

### Tahap 1 — Ubah Model `Jawaban` Jadi Nullable

**File:** [internal/modules/jawaban/model/jawaban_model.go](internal/modules/jawaban/model/jawaban_model.go)

Ubah dua field:
- `Jawaban string` → `Jawaban *string` (hapus `not null`)
- `IsBenar bool` → `IsBenar *bool` (hapus `not null;default:false`)

**Sebelum:**
```go
Jawaban   string     `gorm:"type:varchar(1);not null" json:"jawaban"`
IsBenar   bool       `gorm:"type:boolean;not null;default:false" json:"is_benar"`
```

**Sesudah:**
```go
Jawaban   *string    `gorm:"type:varchar(1)" json:"jawaban"`
IsBenar   *bool      `gorm:"type:boolean" json:"is_benar"`
```

---

### Tahap 2 — Update Migration: Drop NOT NULL Constraint

**File:** [internal/database/migrate.go](internal/database/migrate.go)

Tambahkan raw SQL **SEBELUM** `AutoMigrate` untuk drop constraint NOT NULL & default value dari kolom existing. AutoMigrate **tidak akan** drop constraint NOT NULL secara otomatis.

Tambahkan di dalam `RunMigrations`, di antara `DROP INDEX...` dan `AutoMigrate(...)`:

```go
// Drop NOT NULL & DEFAULT dari kolom jawaban supaya bisa di-set NULL saat bulk-insert
db.Exec("ALTER TABLE jawaban ALTER COLUMN jawaban DROP NOT NULL")
db.Exec("ALTER TABLE jawaban ALTER COLUMN is_benar DROP NOT NULL")
db.Exec("ALTER TABLE jawaban ALTER COLUMN is_benar DROP DEFAULT")
```

> **Catatan:** `db.Exec` di sini sengaja tidak di-handle errornya — kalau tabel belum ada (fresh install), error wajar dan diabaikan. Saat tabel sudah ada, perintah ini berhasil.

---

### Tahap 3 — Update DTO Response Jawaban Jadi Nullable

**File:** [internal/modules/jawaban/dto/jawaban_dto.go](internal/modules/jawaban/dto/jawaban_dto.go)

Hanya `JawabanResponse` yang berubah. **Request DTO TIDAK BERUBAH** (user CRUD tetap kirim `jawaban` sebagai string A-E).

**Sebelum:**
```go
type JawabanResponse struct {
    ...
    Jawaban     string `json:"jawaban"`
    IsBenar     bool   `json:"is_benar"`
    ...
}
```

**Sesudah:**
```go
type JawabanResponse struct {
    ...
    Jawaban     *string `json:"jawaban"`
    IsBenar     *bool   `json:"is_benar"`
    ...
}
```

> **Kenapa request tidak berubah?** Karena saat user CRUD jawaban (POST/PUT), dia WAJIB kirim jawaban yang valid (A-E). Yang NULL hanya record yang di-generate otomatis (belum dijawab).

---

### Tahap 4 — Update `JawabanWithDetail` Struct di Repository

**File:** [internal/modules/jawaban/repository/jawaban_repository.go](internal/modules/jawaban/repository/jawaban_repository.go)

Ubah dua field di struct `JawabanWithDetail`:

**Sebelum:**
```go
Jawaban     string `gorm:"column:jawaban"`
IsBenar     bool   `gorm:"column:is_benar"`
```

**Sesudah:**
```go
Jawaban     *string `gorm:"column:jawaban"`
IsBenar     *bool   `gorm:"column:is_benar"`
```

> SELECT query tidak berubah (kolom yang sama tetap di-select). GORM otomatis handle nullable saat tipe pointer.

---

### Tahap 5 — Tambah Method `BulkCreate` di Repository Jawaban

**File:** [internal/modules/jawaban/repository/jawaban_repository.go](internal/modules/jawaban/repository/jawaban_repository.go)

Ini method baru yang menerima `*gorm.DB` (supaya bisa dipanggil dengan tx dari nilai service). Tambah di interface dan implementasi.

**Tambah di interface `JawabanRepository`:**
```go
BulkCreateWithTx(tx *gorm.DB, jawabans []model.Jawaban) error
```

**Tambah implementasinya (letakkan setelah method `Create`):**
```go
func (r *jawabanRepository) BulkCreateWithTx(tx *gorm.DB, jawabans []model.Jawaban) error {
    if len(jawabans) == 0 {
        return nil
    }
    return tx.Create(&jawabans).Error
}
```

---

### Tahap 6 — Sesuaikan Service Jawaban dengan Tipe Nullable

**File:** [internal/modules/jawaban/service/jawaban_service.go](internal/modules/jawaban/service/jawaban_service.go)

Karena `model.Jawaban.Jawaban` & `IsBenar` sekarang pointer, beberapa baris perlu disesuaikan.

**Di `CreateJawaban`:**

**Sebelum:**
```go
row := &model.Jawaban{
    IDNilai:   req.IDNilai,
    IDSoal:    req.IDSoal,
    IDPeserta: req.IDPeserta,
    NoUrut:    req.NoUrut,
    Jawaban:   jawaban,
    IsBenar:   isBenar,
}
```

**Sesudah:**
```go
row := &model.Jawaban{
    IDNilai:   req.IDNilai,
    IDSoal:    req.IDSoal,
    IDPeserta: req.IDPeserta,
    NoUrut:    req.NoUrut,
    Jawaban:   &jawaban,    // pointer
    IsBenar:   &isBenar,    // pointer
}
```

**Di `UpdateJawaban`:**

**Sebelum:**
```go
existing.Jawaban   = jawaban
existing.IsBenar   = isBenar
```

**Sesudah:**
```go
existing.Jawaban   = &jawaban
existing.IsBenar   = &isBenar
```

**Di `detailToResponse`:** sudah pointer (karena `JawabanWithDetail` & `JawabanResponse` sudah disesuaikan di Tahap 3 & 4) — tidak perlu ubah apa-apa.

---

### Tahap 7 — Modifikasi `MulaiUjian` di Service Nilai

**File:** [internal/modules/nilai/service/nilai_service.go](internal/modules/nilai/service/nilai_service.go)

Ini bagian terpenting. Method `MulaiUjian` perlu refactor besar:

1. Tambah dependency baru di service (db langsung untuk transaction, plus jawaban repo).
2. Bungkus operasi insert nilai + bulk insert jawaban dalam `db.Transaction(...)`.
3. Resume scenario tetap di luar transaction.

**Langkah:**

**(a) Update import di paling atas file:**
```go
import (
    "errors"
    "math"
    "time"

    "backend/internal/constants"
    jawabanmodel "backend/internal/modules/jawaban/model"
    jawabanrepo  "backend/internal/modules/jawaban/repository"
    "backend/internal/modules/nilai/dto"
    "backend/internal/modules/nilai/model"
    "backend/internal/modules/nilai/repository"
    soalmodel "backend/internal/modules/soal/model"
    jadwalmodel "backend/internal/modules/jadwal/model"

    "gorm.io/gorm"
)
```

**(b) Update struct `nilaiService` dan constructor:**

**Sebelum:**
```go
type nilaiService struct {
    repo repository.NilaiRepository
}

func NewNilaiService(repo repository.NilaiRepository) NilaiService {
    return &nilaiService{repo: repo}
}
```

**Sesudah:**
```go
type nilaiService struct {
    repo        repository.NilaiRepository
    jawabanRepo jawabanrepo.JawabanRepository
    db          *gorm.DB
}

func NewNilaiService(repo repository.NilaiRepository, jawabanRepo jawabanrepo.JawabanRepository, db *gorm.DB) NilaiService {
    return &nilaiService{
        repo:        repo,
        jawabanRepo: jawabanRepo,
        db:          db,
    }
}
```

**(c) Refactor method `MulaiUjian`:**

Hapus implementasi lama, ganti dengan ini:

```go
func (s *nilaiService) MulaiUjian(idPeserta, idJadwal string) (*dto.NilaiResponse, bool, error) {
    // 1. Cek apakah record sudah ada (di luar transaction — read only)
    existing, err := s.repo.GetByPesertaAndJadwal(idPeserta, idJadwal)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, false, err
    }

    // 2. Jika sudah ada → cek wkt_selesai
    if existing != nil {
        if existing.WktSelesai != nil {
            return nil, false, errors.New("Ujian sudah pernah dilakukan")
        }
        // Resume — return detail, JANGAN re-generate jawaban
        detail, err := s.repo.GetByIDWithDetail(existing.ID)
        if err != nil {
            return nil, false, err
        }
        return detailToResponse(detail), false, nil
    }

    // 3. Belum ada → transaction: insert nilai + bulk insert jawaban
    var newNilaiID string
    err = s.db.Transaction(func(tx *gorm.DB) error {
        // 3a. Get jadwal untuk dapatkan id_bank_soal
        var jadwal jadwalmodel.Jadwal
        if err := tx.Where("id = ? AND deleted_at IS NULL", idJadwal).First(&jadwal).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return errors.New("jadwal tidak ditemukan")
            }
            return err
        }

        // 3b. Insert nilai baru
        now := time.Now()
        nilai := &model.Nilai{
            IDPeserta:         idPeserta,
            IDJadwal:          idJadwal,
            Nilai:             0,
            WktMulai:          &now,
            AktivitasTerakhir: &now,
            WktSelesai:        nil,
        }
        if err := tx.Create(nilai).Error; err != nil {
            return err
        }
        newNilaiID = nilai.ID

        // 3c. Query soal random by bank_soal
        var soals []soalmodel.Soal
        if err := tx.
            Where("id_bank_soal = ? AND deleted_at IS NULL", jadwal.IDBankSoal).
            Order("RANDOM()").
            Find(&soals).Error; err != nil {
            return err
        }

        // 3d. Build & bulk insert jawaban kosong
        if len(soals) > 0 {
            jawabans := make([]jawabanmodel.Jawaban, len(soals))
            for i, soal := range soals {
                jawabans[i] = jawabanmodel.Jawaban{
                    IDNilai:   nilai.ID,
                    IDSoal:    soal.ID,
                    IDPeserta: idPeserta,
                    NoUrut:    i + 1, // 1-based index sesuai urutan random
                    Jawaban:   nil,
                    IsBenar:   nil,
                }
            }
            if err := s.jawabanRepo.BulkCreateWithTx(tx, jawabans); err != nil {
                return err
            }
        }

        return nil
    })

    if err != nil {
        return nil, false, err
    }

    // 4. Ambil detail nilai yang baru di-insert (di luar transaction, read)
    created, err := s.repo.GetByIDWithDetail(newNilaiID)
    if err != nil {
        return nil, false, err
    }
    return detailToResponse(created), true, nil
}
```

> **Catatan penting:**
> - Resume scenario **TIDAK MASUK transaction** — tidak ada perubahan data, cukup read.
> - Pakai `tx` di dalam transaction, **bukan** `s.db` (kalau pakai `s.db` di dalam Transaction, query terjadi di luar tx → kehilangan atomicity).
> - `Order("RANDOM()")` adalah syntax PostgreSQL. Kalau pindah ke MySQL, ganti jadi `Order("RAND()")`.
> - `no_urut` mulai dari 1 (bukan 0) — lebih natural untuk display.

---

### Tahap 8 — Update Wiring di Routes Nilai

**File:** [internal/modules/nilai/routes/nilai_routes.go](internal/modules/nilai/routes/nilai_routes.go)

Karena constructor `NewNilaiService` berubah signature, perlu pass dependency baru.

**Sebelum:**
```go
import (
    "backend/internal/middleware"
    "backend/internal/modules/nilai/controller"
    "backend/internal/modules/nilai/repository"
    "backend/internal/modules/nilai/service"

    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

func SetupNilaiRoutes(app *fiber.App, db *gorm.DB) {
    repo := repository.NewNilaiRepository(db)
    svc  := service.NewNilaiService(repo)
    ctrl := controller.NewNilaiController(svc)
    ...
}
```

**Sesudah:**
```go
import (
    "backend/internal/middleware"
    jawabanrepo "backend/internal/modules/jawaban/repository"
    "backend/internal/modules/nilai/controller"
    "backend/internal/modules/nilai/repository"
    "backend/internal/modules/nilai/service"

    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

func SetupNilaiRoutes(app *fiber.App, db *gorm.DB) {
    repo         := repository.NewNilaiRepository(db)
    jawabanRepository := jawabanrepo.NewJawabanRepository(db)
    svc          := service.NewNilaiService(repo, jawabanRepository, db)
    ctrl         := controller.NewNilaiController(svc)
    ...
}
```

---

### Tahap 9 — Build & Smoke Test

```bash
# Pastikan compile bersih
go build ./...

# Jalankan server (migration + ALTER TABLE jalan otomatis)
go run cmd/server/main.go
```

Cek log startup. Kalau ada error `ERROR: column "jawaban" of relation "jawaban" contains null values` → berarti masih ada record lama. Hapus dulu: `DELETE FROM jawaban` (HATI-HATI, hanya kalau dev/test env).

---

### Tahap 10 — Manual Testing

**Persiapan:**
- Pastikan ada minimal 1 jadwal valid dengan `id_bank_soal` yang ada di tabel `bank_soal`.
- Pastikan bank soal tersebut punya minimal 3 soal aktif.
- Login sebagai peserta valid, simpan token sebagai `$TOKEN`.
- Simpan `id_jadwal` sebagai `$JADWAL_ID`.

**Test 1 — Mulai ujian pertama kali (record nilai + jawaban di-generate):**
```bash
curl -X POST http://localhost:3000/api/nilai/mulai-ujian/$JADWAL_ID \
  -H "Authorization: Bearer $TOKEN"
```
Ekspektasi: HTTP 201, message `"Mulai ujian successfully"`.

Lalu cek tabel jawaban via SQL atau endpoint:
```bash
# Ganti $NILAI_ID dengan id dari response sebelumnya
curl "http://localhost:3000/api/jawaban/nilai/$NILAI_ID"
```
Ekspektasi:
- `data` array berisi N record jawaban (N = jumlah soal di bank_soal).
- Semua record punya `jawaban: null`, `is_benar: null`.
- `no_urut` mulai dari 1, urut, tapi `no_soal` (dari soal asli) berurutan acak — bukti random order berjalan.

**Test 2 — Resume (panggil lagi):**
```bash
curl -X POST http://localhost:3000/api/nilai/mulai-ujian/$JADWAL_ID \
  -H "Authorization: Bearer $TOKEN"
```
Ekspektasi:
- HTTP 200, message `"Lanjutkan ujian successfully"`.
- Cek lagi `GET /api/jawaban/nilai/$NILAI_ID` → **jumlah dan urutan record TIDAK BERUBAH** (tidak ada bulk-insert lagi).

**Test 3 — Update salah satu jawaban:**
```bash
# Ambil salah satu jawaban_id dari Test 1
curl -X PUT http://localhost:3000/api/jawaban/$JAWABAN_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id_nilai":   "'$NILAI_ID'",
    "id_soal":    "'$SOAL_ID'",
    "id_peserta": "'$PESERTA_ID'",
    "no_urut":    1,
    "jawaban":    "B"
  }'
```
Ekspektasi: HTTP 200, response `jawaban: "B"`, `is_benar: true/false` (sesuai kunci).

**Test 4 — Jadwal dengan bank_soal kosong:**
Buat jadwal yang `id_bank_soal`-nya kosong / tidak punya soal. Panggil mulai-ujian.
Ekspektasi: HTTP 201, record nilai tetap dibuat, tidak ada record jawaban yang dibuat (0 jawaban).

**Test 5 — Jadwal tidak valid:**
```bash
curl -X POST http://localhost:3000/api/nilai/mulai-ujian/00000000-0000-0000-0000-000000000000 \
  -H "Authorization: Bearer $TOKEN"
```
Ekspektasi: HTTP 400, message `"jadwal tidak ditemukan"`. **Tidak ada** record nilai/jawaban yang dibuat (transaction rolled back).

---

## 📊 Ringkasan File yang Disentuh

| File | Perubahan |
|------|-----------|
| [internal/modules/jawaban/model/jawaban_model.go](internal/modules/jawaban/model/jawaban_model.go) | `Jawaban` & `IsBenar` jadi pointer (nullable) |
| [internal/modules/jawaban/dto/jawaban_dto.go](internal/modules/jawaban/dto/jawaban_dto.go) | `JawabanResponse.Jawaban` & `IsBenar` jadi pointer |
| [internal/modules/jawaban/repository/jawaban_repository.go](internal/modules/jawaban/repository/jawaban_repository.go) | `JawabanWithDetail.Jawaban` & `IsBenar` jadi pointer; tambah `BulkCreateWithTx` |
| [internal/modules/jawaban/service/jawaban_service.go](internal/modules/jawaban/service/jawaban_service.go) | Sesuaikan create/update jadi `&jawaban`, `&isBenar` |
| [internal/modules/nilai/service/nilai_service.go](internal/modules/nilai/service/nilai_service.go) | Tambah dependency `jawabanRepo` & `db`; refactor `MulaiUjian` pakai transaction |
| [internal/modules/nilai/routes/nilai_routes.go](internal/modules/nilai/routes/nilai_routes.go) | Pass `jawabanRepo` & `db` ke constructor service |
| [internal/database/migrate.go](internal/database/migrate.go) | Tambah ALTER TABLE drop NOT NULL untuk kolom jawaban & is_benar |

**TIDAK PERLU DIUBAH:**
- Controller jawaban — DTO request tetap sama
- DTO request jawaban — user CRUD tetap wajib kirim A-E
- Routes jawaban — tidak ada endpoint baru
- main.go — tidak ada modul baru

---

## ⚠️ Aturan & Catatan Penting

1. **Transaction WAJIB.** Insert nilai + bulk insert jawaban harus atomic. Jangan refactor jadi 2 query terpisah tanpa tx.

2. **Resume jangan re-generate.** Saat user panggil mulai-ujian kedua kali (record nilai sudah ada, belum selesai), JANGAN bulk-insert jawaban lagi.

3. **`no_urut` 1-based, bukan 0-based.** Lebih natural untuk display di frontend ("Soal 1 dari 50").

4. **Urutan random per peserta.** Setiap call `mulai-ujian` (untuk peserta berbeda) menghasilkan urutan acak yang berbeda — ini fitur, bukan bug.

5. **Bank soal kosong = jawaban kosong, bukan error.** Kalau jadwal punya bank_soal tapi bank_soal-nya nggak punya soal aktif → record nilai tetap dibuat, jumlah jawaban = 0.

6. **PostgreSQL specific.** `ORDER BY RANDOM()` adalah syntax PostgreSQL. Jangan diubah jadi `RAND()` kecuali pindah ke MySQL.

7. **Backward compat existing data.** Saat migrate jalan di DB yang sudah punya record jawaban lama (dengan `jawaban = 'A'` dll), record tersebut tetap valid — perubahan `NOT NULL → NULL` hanya membolehkan NULL, tidak memaksa.

8. **Unique constraint tetap berlaku.** `idx_jawaban_nilai_soal_unique` masih aktif. Karena bulk-insert melalui 1 nilai_id × N soal_id (semua unique), tidak akan conflict.

---

## 📖 Update Dokumentasi Frontend

Setelah implementasi selesai, **update file** [ENDPOINT_MULAI_UJIAN.md](ENDPOINT_MULAI_UJIAN.md) dengan informasi baru berikut. Tambahkan section ini **setelah** section "Response":

---

### 📝 Side Effect: Jawaban Otomatis Di-generate

Saat endpoint ini berhasil membuat record nilai baru (HTTP 201), backend **otomatis** membuat record jawaban kosong untuk **setiap soal** di bank soal jadwal tersebut, dengan urutan **acak (random)**.

**Mapping yang dibuat untuk setiap soal:**
| Field | Nilai |
|-------|-------|
| `id_nilai` | ID nilai yang baru dibuat |
| `id_soal` | ID soal |
| `id_peserta` | ID peserta (dari JWT) |
| `no_urut` | Urutan tampil di frontend (1, 2, 3, ...) — sesuai urutan random |
| `jawaban` | `null` (belum dijawab) |
| `is_benar` | `null` (belum bisa dinilai) |

**Implikasi untuk Frontend:**

Setelah call `POST /api/nilai/mulai-ujian/:id_jadwal`, frontend **langsung bisa** fetch daftar jawaban (yang berisi soal urut random):

```javascript
// Step 1: Mulai ujian
const startResp = await axios.post(`/api/nilai/mulai-ujian/${jadwalId}`, {}, {
  headers: { Authorization: `Bearer ${token}` }
});
const nilaiId = startResp.data.data.id;

// Step 2: Ambil daftar jawaban (yang sudah terisi soal urut random)
const jawabanResp = await axios.get(`/api/jawaban/nilai/${nilaiId}`);
const jawabanList = jawabanResp.data.data;
// jawabanList[0].no_urut = 1, jawabanList[0].soal = "<teks soal>", jawabanList[0].jawaban = null
// jawabanList[1].no_urut = 2, ...
```

**Saat peserta menjawab soal**, frontend tinggal PUT ke endpoint update jawaban:

```javascript
await axios.put(`/api/jawaban/${jawabanId}`, {
  id_nilai:   nilaiId,
  id_soal:    soalId,
  id_peserta: pesertaId,
  no_urut:    1,           // tetap dari record
  jawaban:    "B"          // jawaban peserta
}, {
  headers: { Authorization: `Bearer ${token}` }
});
```

**Behavior saat Resume (HTTP 200):**
- Tidak ada record jawaban baru yang dibuat.
- Daftar jawaban yang lama (beserta jawaban yang sudah diisi peserta) tetap utuh.
- Frontend cukup re-fetch `GET /api/jawaban/nilai/:id_nilai` untuk melanjutkan ujian.

**Format Response Jawaban setelah perubahan:**
```json
{
  "id": "uuid-jawaban",
  "id_nilai": "uuid-nilai",
  "id_soal": "uuid-soal",
  "no_urut": 1,
  "no_soal": 17,
  "soal": "Berapakah hasil 2 + 2?",
  "kunci": "B",
  "id_peserta": "uuid-peserta",
  "nama_peserta": "Budi",
  "jawaban": null,          // <-- bisa null (belum dijawab)
  "is_benar": null,         // <-- bisa null (belum bisa dinilai)
  "created_at": "...",
  "updated_at": "..."
}
```

**Catatan penting untuk frontend:**
- `jawaban` & `is_benar` sekarang **bisa `null`**. Handle null check di UI:
  ```jsx
  <span>{jawaban.jawaban ?? "Belum dijawab"}</span>
  <Icon color={jawaban.is_benar === null ? "gray" : (jawaban.is_benar ? "green" : "red")} />
  ```
- Urutan soal yang ditampilkan ke peserta = `ORDER BY no_urut ASC`. Backend sudah `ORDER BY soal.no_soal ASC` di endpoint `GET /api/jawaban/nilai/:id_nilai`, jadi frontend perlu **re-sort** by `no_urut` jika ingin urutan random sesuai design.

---

## 🔗 Referensi Pola di Kode

- Pola transaction GORM: docs https://gorm.io/docs/transactions.html
- Pola `Order("RANDOM()")` PostgreSQL: docs GORM `Order` chain
- Pola repository menerima `*gorm.DB` (untuk dipakai dengan tx): pattern industry-standard, kita pakai di method `BulkCreateWithTx`
- Pola constructor dengan multiple dependency: lihat [internal/modules/jadwal/routes/jadwal_routes.go](internal/modules/jadwal/routes/jadwal_routes.go) — `NewJadwalService` sudah pakai 3 dependency
