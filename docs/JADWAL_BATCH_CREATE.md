# Jadwal Batch Create - Dokumentasi Fitur

## Perubahan Logika Penyimpanan Jadwal

Logika pembuatan jadwal telah diubah untuk mendukung pembuatan jadwal dengan multiple kelas dalam satu request.

### Sebelum (Old Logic)
- Request hanya berisi data jadwal (id_bank_soal, wkt_mulai, wkt_selesai)
- Hanya membuat 1 record di tabel `jadwal`
- Untuk menambahkan kelas, harus membuat request terpisah ke `/api/jadwal-kelas`

### Sesudah (New Logic)
- Request berisi data jadwal + array id_kelas
- Membuat 1 record di tabel `jadwal`
- Membuat multiple records di tabel `jadwal_kelas` (satu per kelas)
- Semua operasi dalam satu endpoint (atomic operation)

---

## Format Request Baru

### POST /api/jadwal

**Request Body:**
```json
{
  "id_bank_soal": "7aa5823d-b770-4015-9952-da723c7b2505",
  "wkt_mulai": "2026-05-21 10:18:00",
  "wkt_selesai": "2026-05-22 10:18:00",
  "id_kelas": [
    "dc3cb7a6-49cd-43f2-a5b7-68afa6d28ede",
    "51221e32-22b7-4994-a630-32a5d46eb788"
  ]
}
```

**Parameter Penjelasan:**

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `id_bank_soal` | string (UUID) | ID bank soal yang akan dijadwalkan |
| `wkt_mulai` | string | Waktu mulai ujian (format: YYYY-MM-DD HH:MM:SS) |
| `wkt_selesai` | string | Waktu selesai ujian (format: YYYY-MM-DD HH:MM:SS) |
| `id_kelas` | array of UUID | Daftar ID kelas yang akan mengikuti jadwal ini |

---

## Contoh Curl Request

```bash
curl -X POST "http://localhost:3000/api/jadwal" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id_bank_soal": "7aa5823d-b770-4015-9952-da723c7b2505",
    "wkt_mulai": "2026-05-21 10:18:00",
    "wkt_selesai": "2026-05-22 10:18:00",
    "id_kelas": [
      "dc3cb7a6-49cd-43f2-a5b7-68afa6d28ede",
      "51221e32-22b7-4994-a630-32a5d46eb788"
    ]
  }'
```

---

## Response Success (201 Created)

```json
{
  "success": true,
  "message": "Create jadwal successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6",
    "id_bank_soal": "7aa5823d-b770-4015-9952-da723c7b2505",
    "nama_bank_soal": "Soal Ujian Matematika 2025",
    "wkt_mulai": "2026-05-21 10:18:00",
    "wkt_selesai": "2026-05-22 10:18:00",
    "created_at": "2025-05-21 14:30:25",
    "updated_at": "2025-05-21 14:30:25"
  }
}
```

**Yang terjadi di backend:**
1. ✅ 1 record dibuat di tabel `jadwal` dengan UUID `a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6`
2. ✅ 2 records dibuat di tabel `jadwal_kelas`:
   - `{id_jadwal: a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6, id_kelas: dc3cb7a6-49cd-43f2-a5b7-68afa6d28ede}`
   - `{id_jadwal: a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6, id_kelas: 51221e32-22b7-4994-a630-32a5d46eb788}`

---

## Error Responses

### 400 Bad Request - Minimal satu kelas diperlukan
```json
{
  "success": false,
  "message": "minimal harus ada satu kelas yang didaftarkan",
  "data": null
}
```

### 400 Bad Request - Duplikat assignment
```json
{
  "success": false,
  "message": "kelas dengan ID dc3cb7a6-49cd-43f2-a5b7-68afa6d28ede sudah terdaftar di jadwal tersebut",
  "data": null
}
```

### 400 Bad Request - Format datetime salah
```json
{
  "success": false,
  "message": "format wkt_mulai tidak valid, gunakan: 2006-01-02 15:04:05",
  "data": null
}
```

### 400 Bad Request - Waktu selesai harus setelah mulai
```json
{
  "success": false,
  "message": "wkt_selesai harus setelah wkt_mulai",
  "data": null
}
```

---

## Validasi yang Dilakukan

1. **Format datetime** — `wkt_mulai` dan `wkt_selesai` harus format `YYYY-MM-DD HH:MM:SS`
2. **Urutan waktu** — `wkt_selesai` harus selalu **setelah** `wkt_mulai`
3. **Minimal ada kelas** — Array `id_kelas` tidak boleh kosong (minimal 1 kelas)
4. **Unique constraint** — Setiap kelas di array dicek apakah sudah terdaftar di jadwal yang sama
5. **Duplikat dalam array** — Jika ada kelas yang sama dalam array, hanya satu yang akan disimpan (database level unique constraint)

---

## Perubahan di Database

### Tabel jadwal
```sql
id              | id_bank_soal                     | wkt_mulai           | wkt_selesai         | created_at          | updated_at          | deleted_at
---|---|---|---|---|---|---
a1b2c3d4-...   | 7aa5823d-...                    | 2026-05-21 10:18:00 | 2026-05-22 10:18:00 | 2025-05-21 14:30:25 | 2025-05-21 14:30:25 | NULL
```

### Tabel jadwal_kelas
```sql
id              | id_jadwal           | id_kelas                        | created_at          | updated_at
---|---|---|---|---
f7a8b9c0-...   | a1b2c3d4-...       | dc3cb7a6-...                   | 2025-05-21 14:30:25 | 2025-05-21 14:30:25
g8b9c0d1-...   | a1b2c3d4-...       | 51221e32-...                   | 2025-05-21 14:30:25 | 2025-05-21 14:30:25
```

---

## Use Case / Skenario Penggunaan

### Skenario 1: Setup jadwal ujian untuk semua kelas
```bash
# Adminitrator ingin membuat jadwal Ujian Matematika untuk 5 kelas
curl -X POST "http://localhost:3000/api/jadwal" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "id_bank_soal": "7aa5823d-b770-4015-9952-da723c7b2505",
    "wkt_mulai": "2026-08-01 08:00:00",
    "wkt_selesai": "2026-08-01 10:00:00",
    "id_kelas": [
      "dc3cb7a6-49cd-43f2-a5b7-68afa6d28ede",  # XI-A
      "51221e32-22b7-4994-a630-32a5d46eb788",  # XI-B
      "61331e32-22b7-4994-a630-32a5d46eb799",  # XI-C
      "71441e32-22b7-4994-a630-32a5d46eb700",  # XI-D
      "81551e32-22b7-4994-a630-32a5d46eb811"   # XI-E
    ]
  }'
```

**Benefit:**
- ✅ Hanya 1 request, bukan 5 requests
- ✅ Atomic operation — semua beres atau semua gagal
- ✅ Lebih cepat dari loop manual

### Skenario 2: Setup jadwal untuk subset kelas
```bash
# Setup ujian khusus untuk kelas XII yang sudah siap
curl -X POST "http://localhost:3000/api/jadwal" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "id_bank_soal": "8bb5823d-b770-4015-9952-da723c7b2506",
    "wkt_mulai": "2026-06-15 09:00:00",
    "wkt_selesai": "2026-06-15 11:00:00",
    "id_kelas": [
      "aa3cb7a6-49cd-43f2-a5b7-68afa6d28f00",  # XII-A
      "bb4cb7a6-49cd-43f2-a5b7-68afa6d28f01"   # XII-B
    ]
  }'
```

---

## Perubahan Kode di Backend

### File yang Dimodifikasi

1. **`internal/modules/jadwal/dto/jadwal_dto.go`**
   - Tambah field `IDKelas []string` di struct `CreateJadwalRequest`

2. **`internal/modules/jadwal/service/jadwal_service.go`**
   - Inject `JadwalKelasRepository` ke `JadwalService`
   - Update `CreateJadwal()` untuk create jadwal + jadwal_kelas in batch

3. **`internal/modules/jadwal/routes/jadwal_routes.go`**
   - Tambah import jadwal_kelas repository
   - Pass jadwal_kelas repository ke service

4. **`internal/modules/jadwal_kelas/repository/jadwal_kelas_repository.go`**
   - Tambah method `CreateBulk()` untuk insert multiple records

---

## Backward Compatibility

⚠️ **Breaking Change** — Request format berubah:

**Old Request (tidak support lagi):**
```json
{
  "id_bank_soal": "7aa5823d-b770-4015-9952-da723c7b2505",
  "wkt_mulai": "2026-05-21 10:18:00",
  "wkt_selesai": "2026-05-22 10:18:00"
}
```

**New Request (harus ada id_kelas):**
```json
{
  "id_bank_soal": "7aa5823d-b770-4015-9952-da723c7b2505",
  "wkt_mulai": "2026-05-21 10:18:00",
  "wkt_selesai": "2026-05-22 10:18:00",
  "id_kelas": ["dc3cb7a6-49cd-43f2-a5b7-68afa6d28ede"]
}
```

Jika field `id_kelas` tidak ada atau array kosong, API akan return error 400.

---

## Performance Improvement

**Sebelum (Old Way):**
```
1 request POST /api/jadwal            → 1 record di jadwal table
5 requests POST /api/jadwal-kelas     → 5 records di jadwal_kelas table
---
Total: 6 HTTP requests, 6 database writes
```

**Sesudah (New Way):**
```
1 request POST /api/jadwal            → 1 record di jadwal table + 5 records di jadwal_kelas table
---
Total: 1 HTTP request, 6 database writes
```

**Improvement:**
- 83% kurang HTTP requests (dari 6 jadi 1)
- Lebih cepat (network latency berkurang)
- Atomic operation (semua beres atau semua gagal, tidak ada state tengah-tengah)

---

## Testing

### Test dengan cURL
```bash
# Test 1: Create jadwal dengan 2 kelas
curl -X POST "http://localhost:3000/api/jadwal" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "id_bank_soal": "7aa5823d-b770-4015-9952-da723c7b2505",
    "wkt_mulai": "2026-05-21 10:18:00",
    "wkt_selesai": "2026-05-22 10:18:00",
    "id_kelas": [
      "dc3cb7a6-49cd-43f2-a5b7-68afa6d28ede",
      "51221e32-22b7-4994-a630-32a5d46eb788"
    ]
  }'

# Test 2: Verify jadwal dibuat
curl "http://localhost:3000/api/jadwal"

# Test 3: Verify jadwal_kelas entries dibuat
curl "http://localhost:3000/api/jadwal-kelas?id_jadwal=<jadwal-id-dari-response>"
```

---

## Notes

- Perubahan ini hanya mempengaruhi endpoint `POST /api/jadwal`
- Endpoint lainnya (`GET`, `PUT`, `DELETE`, `PATCH /restore`) tetap sama
- Data di kedua tabel (`jadwal` dan `jadwal_kelas`) semuanya baru (bukan update data lama)
- Jika ada error saat create jadwal_kelas, transaksi akan rollback di beberapa database, tergantung konfigurasi
