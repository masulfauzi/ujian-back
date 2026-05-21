# Jadwal Kelas API Documentation

## Base URL

```
http://localhost:3000/api/jadwal-kelas
```

---

## Overview

Jadwal Kelas API mengelola penugasan kelas ke jadwal ujian. Tabel `jadwal_kelas` adalah tabel pivot yang menghubungkan jadwal ujian dengan kelas — menentukan kelas mana saja yang akan mengikuti jadwal ujian tertentu.

**Perbedaan penting dari modul lain:**
- **Hard delete** — data dihapus permanen, tidak ada restore
- **Tidak ada soft delete** — tidak ada kolom `deleted_at`
- **Unique constraint** — satu kelas tidak boleh didaftarkan dua kali ke jadwal yang sama

---

## Authentication

Beberapa endpoint memerlukan JWT token:

- **POST /api/jadwal-kelas** — Memerlukan JWT
- **PUT /api/jadwal-kelas/:id** — Memerlukan JWT
- **DELETE /api/jadwal-kelas/:id** — Memerlukan JWT

Kirim token di header:
```
Authorization: Bearer <token>
```

---

## Endpoints

### 1. GET /api/jadwal-kelas

Mengambil daftar semua penugasan kelas ke jadwal dengan pagination dan filter.

**Method:** `GET`

**Auth:** Tidak diperlukan

**Query Parameters:**

| Parameter | Type | Default | Deskripsi |
|-----------|------|---------|-----------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |
| `id_jadwal` | string (UUID) | _(opsional)_ | Filter berdasarkan jadwal tertentu |
| `id_kelas` | string (UUID) | _(opsional)_ | Filter berdasarkan kelas tertentu |

**Example Request:**

```bash
curl -X GET "http://localhost:3000/api/jadwal-kelas?page=1&page_size=10"
```

**Example Request dengan filter:**

```bash
curl -X GET "http://localhost:3000/api/jadwal-kelas?id_jadwal=b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7&page=1"
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Get all jadwal kelas successfully",
  "data": {
    "data": [
      {
        "id": "d5e6f7a8-b9c0-49d1-e2f3-a4b5c6d7e8f9",
        "id_jadwal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
        "id_kelas": "c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8",
        "nama_kelas": "XI-A",
        "nama_bank_soal": "Soal Ujian Matematika 2025",
        "wkt_mulai": "2025-08-01 08:00:00",
        "wkt_selesai": "2025-08-01 10:00:00",
        "created_at": "2025-05-21 14:30:25",
        "updated_at": "2025-05-21 14:30:25"
      },
      {
        "id": "e6f7a8b9-c0d1-50e2-f3a4-b5c6d7e8f9a0",
        "id_jadwal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
        "id_kelas": "d4e5f6a7-b8c9-50d1-e2f3-a4b5c6d7e8f9",
        "nama_kelas": "XI-B",
        "nama_bank_soal": "Soal Ujian Matematika 2025",
        "wkt_mulai": "2025-08-01 08:00:00",
        "wkt_selesai": "2025-08-01 10:00:00",
        "created_at": "2025-05-21 14:35:10",
        "updated_at": "2025-05-21 14:35:10"
      }
    ],
    "total": 2,
    "page": 1,
    "page_size": 10,
    "total_page": 1
  }
}
```

**Error Response (500):**

```json
{
  "success": false,
  "message": "Internal server error",
  "data": null
}
```

---

### 2. GET /api/jadwal-kelas/:id

Mengambil detail penugasan kelas ke jadwal berdasarkan ID.

**Method:** `GET`

**Auth:** Tidak diperlukan

**Path Parameters:**

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `id` | string (UUID) | ID penugasan jadwal-kelas |

**Example Request:**

```bash
curl -X GET "http://localhost:3000/api/jadwal-kelas/d5e6f7a8-b9c0-49d1-e2f3-a4b5c6d7e8f9"
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Get jadwal kelas successfully",
  "data": {
    "id": "d5e6f7a8-b9c0-49d1-e2f3-a4b5c6d7e8f9",
    "id_jadwal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
    "id_kelas": "c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8",
    "nama_kelas": "XI-A",
    "nama_bank_soal": "Soal Ujian Matematika 2025",
    "wkt_mulai": "2025-08-01 08:00:00",
    "wkt_selesai": "2025-08-01 10:00:00",
    "created_at": "2025-05-21 14:30:25",
    "updated_at": "2025-05-21 14:30:25"
  }
}
```

**Error Response (404 Not Found):**

```json
{
  "success": false,
  "message": "Resource not found",
  "data": null
}
```

---

### 3. POST /api/jadwal-kelas

Mendaftarkan sebuah kelas ke jadwal ujian tertentu (membuat penugasan baru).

**Method:** `POST`

**Auth:** JWT diperlukan

**Request Body:**

```json
{
  "id_jadwal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
  "id_kelas": "c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8"
}
```

**Request Body Parameters:**

| Parameter | Type | Required | Deskripsi |
|-----------|------|----------|-----------|
| `id_jadwal` | string (UUID) | Ya | ID jadwal yang valid (harus ada di tabel jadwal) |
| `id_kelas` | string (UUID) | Ya | ID kelas yang valid (harus ada di tabel kelas) |

**Example Request:**

```bash
curl -X POST "http://localhost:3000/api/jadwal-kelas" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id_jadwal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
    "id_kelas": "c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8"
  }'
```

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "Create jadwal kelas successfully",
  "data": {
    "id": "d5e6f7a8-b9c0-49d1-e2f3-a4b5c6d7e8f9",
    "id_jadwal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
    "id_kelas": "c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8",
    "nama_kelas": "XI-A",
    "nama_bank_soal": "Soal Ujian Matematika 2025",
    "wkt_mulai": "2025-08-01 08:00:00",
    "wkt_selesai": "2025-08-01 10:00:00",
    "created_at": "2025-05-21 14:30:25",
    "updated_at": "2025-05-21 14:30:25"
  }
}
```

**Error Response (400 Bad Request - Duplicate Assignment):**

```json
{
  "success": false,
  "message": "kelas ini sudah terdaftar di jadwal tersebut",
  "data": null
}
```

**Error Response (400 Bad Request - Invalid Request):**

```json
{
  "success": false,
  "message": "Invalid request format",
  "data": null
}
```

---

### 4. PUT /api/jadwal-kelas/:id

Mengupdate penugasan kelas ke jadwal (mengubah jadwal atau kelas yang ditugaskan).

**Method:** `PUT`

**Auth:** JWT diperlukan

**Path Parameters:**

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `id` | string (UUID) | ID penugasan jadwal-kelas |

**Request Body:**

```json
{
  "id_jadwal": "d4e5f6a7-b8c9-50d1-e2f3-a4b5c6d7e8f9",
  "id_kelas": "e5f6a7b8-c9d0-51e2-f3a4-b5c6d7e8f9a0"
}
```

**Request Body Parameters:**

| Parameter | Type | Required | Deskripsi |
|-----------|------|----------|-----------|
| `id_jadwal` | string (UUID) | Ya | ID jadwal yang valid |
| `id_kelas` | string (UUID) | Ya | ID kelas yang valid |

**Example Request:**

```bash
curl -X PUT "http://localhost:3000/api/jadwal-kelas/d5e6f7a8-b9c0-49d1-e2f3-a4b5c6d7e8f9" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id_jadwal": "d4e5f6a7-b8c9-50d1-e2f3-a4b5c6d7e8f9",
    "id_kelas": "e5f6a7b8-c9d0-51e2-f3a4-b5c6d7e8f9a0"
  }'
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Update jadwal kelas successfully",
  "data": {
    "id": "d5e6f7a8-b9c0-49d1-e2f3-a4b5c6d7e8f9",
    "id_jadwal": "d4e5f6a7-b8c9-50d1-e2f3-a4b5c6d7e8f9",
    "id_kelas": "e5f6a7b8-c9d0-51e2-f3a4-b5c6d7e8f9a0",
    "nama_kelas": "XI-B",
    "nama_bank_soal": "Soal Ujian Fisika 2025",
    "wkt_mulai": "2025-08-02 08:00:00",
    "wkt_selesai": "2025-08-02 10:00:00",
    "created_at": "2025-05-21 14:30:25",
    "updated_at": "2025-05-21 15:45:30"
  }
}
```

**Error Response (404 Not Found):**

```json
{
  "success": false,
  "message": "Resource not found",
  "data": null
}
```

**Error Response (400 Bad Request - Duplicate After Update):**

```json
{
  "success": false,
  "message": "kelas ini sudah terdaftar di jadwal tersebut",
  "data": null
}
```

---

### 5. DELETE /api/jadwal-kelas/:id

Menghapus penugasan kelas ke jadwal secara permanen (hard delete). **Data yang dihapus tidak bisa dikembalikan.**

**Method:** `DELETE`

**Auth:** JWT diperlukan

**Path Parameters:**

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `id` | string (UUID) | ID penugasan jadwal-kelas |

**Example Request:**

```bash
curl -X DELETE "http://localhost:3000/api/jadwal-kelas/d5e6f7a8-b9c0-49d1-e2f3-a4b5c6d7e8f9" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Delete jadwal kelas successfully",
  "data": null
}
```

**Error Response (404 Not Found):**

```json
{
  "success": false,
  "message": "Resource not found",
  "data": null
}
```

**Error Response (400 Bad Request):**

```json
{
  "success": false,
  "message": "Some error occurred",
  "data": null
}
```

---

## HTTP Status Codes

| Code | Deskripsi |
|------|-----------|
| 200 | OK — Request berhasil |
| 201 | Created — Resource berhasil dibuat |
| 400 | Bad Request — Request tidak valid (format, validasi bisnis) |
| 404 | Not Found — Resource tidak ditemukan |
| 500 | Internal Server Error — Kesalahan server |

---

## Business Rules

### 1. Unique Constraint pada (id_jadwal, id_kelas)

Kombinasi `(id_jadwal, id_kelas)` harus **unik** di database. Satu kelas tidak boleh didaftarkan lebih dari satu kali ke jadwal yang sama.

**Valid:**
- Kelas XI-A dapat didaftarkan ke Jadwal Matematika
- Kelas XI-A dapat didaftarkan ke Jadwal Fisika (berbeda jadwal) ✅

**Invalid:**
- Kelas XI-A didaftarkan ke Jadwal Matematika dua kali ❌

Jika mencoba membuat atau mengupdate dengan pasangan `(id_jadwal, id_kelas)` yang sudah ada, API mengembalikan error 400 dengan pesan: `"kelas ini sudah terdaftar di jadwal tersebut"`

### 2. Hard Delete — Data Terhapus Permanen

`DELETE /api/jadwal-kelas/:id` **menghapus baris secara permanen** dari database. Tidak ada cara untuk mengembalikan data yang sudah dihapus.

- Tidak ada endpoint restore
- Tidak ada soft delete
- Data tidak bisa dipulihkan

### 3. Validasi Duplikat di Service

Pengecekan duplikat dilakukan di level service (bukan hanya database constraint) untuk memberikan error message yang lebih informatif kepada client.

### 4. Response Mencakup Data Relasi

Response API otomatis menyertakan data dari tabel relasi melalui JOIN:
- `nama_kelas` — dari tabel `kelas`
- `nama_bank_soal`, `wkt_mulai`, `wkt_selesai` — dari tabel `jadwal` → `bank_soal`

Ini memudahkan client mengakses informasi lengkap tanpa perlu call API terpisah.

### 5. Update Duplikat Hanya Jika Nilai Berubah

Saat update, pengecekan duplikat hanya dilakukan **jika ada perubahan** pada `id_jadwal` atau `id_kelas`. Jika nilai tetap sama, tidak perlu cek duplikat ulang.

---

## Response Structure

Semua response mengikuti format konsisten:

**Success:**
```json
{
  "success": true,
  "message": "Deskripsi operasi",
  "data": { /* resource atau list */ }
}
```

**Error:**
```json
{
  "success": false,
  "message": "Deskripsi error",
  "data": null
}
```

---

## Pagination

Untuk endpoint list (`GET /api/jadwal-kelas`), gunakan query parameter:

- `page` — Default: 1, Minimal: 1
- `page_size` — Default: 10, Rekomendasi: ≤ 100

Response akan mencakup:
- `total` — Total jumlah record di database
- `page` — Halaman saat ini
- `page_size` — Jumlah item per halaman
- `total_page` — Total halaman yang tersedia
- `data` — Array penugasan jadwal-kelas

**Contoh:**
```bash
# Halaman pertama, 10 item per halaman
curl "http://localhost:3000/api/jadwal-kelas?page=1&page_size=10"

# Halaman kedua, 20 item per halaman
curl "http://localhost:3000/api/jadwal-kelas?page=2&page_size=20"
```

---

## Filtering

Gunakan query parameter untuk filter:

### Filter berdasarkan jadwal tertentu
```bash
curl "http://localhost:3000/api/jadwal-kelas?id_jadwal=b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7"
```

Menampilkan semua kelas yang terdaftar ke jadwal tersebut.

### Filter berdasarkan kelas tertentu
```bash
curl "http://localhost:3000/api/jadwal-kelas?id_kelas=c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8"
```

Menampilkan semua jadwal yang terdaftar untuk kelas tersebut.

### Kombinasi filter dan pagination
```bash
curl "http://localhost:3000/api/jadwal-kelas?id_jadwal=b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7&page=1&page_size=20"
```

---

## Common Use Cases

### Use Case 1: Daftarkan semua kelas ke jadwal ujian

```bash
# 1. Ambil semua kelas
curl http://localhost:3000/api/kelas

# 2. Ambil detail jadwal
curl http://localhost:3000/api/jadwal/b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7

# 3. Daftarkan setiap kelas satu per satu
for kelas_id in $(curl -s http://localhost:3000/api/kelas | jq -r '.data[].id'); do
  curl -X POST "http://localhost:3000/api/jadwal-kelas" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer TOKEN" \
    -d "{\"id_jadwal\": \"b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7\", \"id_kelas\": \"$kelas_id\"}"
done
```

### Use Case 2: Lihat semua kelas untuk jadwal tertentu

```bash
curl "http://localhost:3000/api/jadwal-kelas?id_jadwal=b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7"
```

Menampilkan daftar lengkap kelas yang akan mengikuti jadwal tersebut.

### Use Case 3: Lihat semua jadwal untuk kelas tertentu

```bash
curl "http://localhost:3000/api/jadwal-kelas?id_kelas=c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8"
```

Menampilkan daftar lengkap jadwal yang akan diikuti kelas tersebut.

### Use Case 4: Pindahkan kelas ke jadwal berbeda

```bash
curl -X PUT "http://localhost:3000/api/jadwal-kelas/d5e6f7a8-b9c0-49d1-e2f3-a4b5c6d7e8f9" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "id_jadwal": "d4e5f6a7-b8c9-50d1-e2f3-a4b5c6d7e8f9",
    "id_kelas": "c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8"
  }'
```

### Use Case 5: Batalkan penugasan kelas ke jadwal

```bash
curl -X DELETE "http://localhost:3000/api/jadwal-kelas/d5e6f7a8-b9c0-49d1-e2f3-a4b5c6d7e8f9" \
  -H "Authorization: Bearer TOKEN"
```

**Catatan:** Data dihapus permanen, tidak bisa dipulihkan.

---

## Data Model Relationships

```
jadwal_kelas
├── id_jadwal ─────→ jadwal
│                      ├── id_bank_soal ──→ bank_soal
│                      │                      └── nama_bank_soal
│                      ├── wkt_mulai
│                      └── wkt_selesai
│
└── id_kelas ──────→ kelas
                       └── nama_kelas
```

**Response mencakup:**
- Dari `jadwal_kelas`: `id`, `id_jadwal`, `id_kelas`, `created_at`, `updated_at`
- Dari `kelas`: `nama_kelas`
- Dari `jadwal` → `bank_soal`: `nama_bank_soal`, `wkt_mulai`, `wkt_selesai`

---

## Error Handling

Semua error response mengikuti format:

```json
{
  "success": false,
  "message": "error message",
  "data": null
}
```

**Common error messages:**

| Error | Penyebab | Solusi |
|-------|----------|--------|
| `"kelas ini sudah terdaftar di jadwal tersebut"` | Duplikat assignment | Gunakan jadwal atau kelas yang berbeda |
| `"Resource not found"` | ID tidak ditemukan | Verifikasi ID yang dikirim |
| `"Invalid request format"` | Format request salah | Periksa JSON structure dan field names |
| `"Internal server error"` | Kesalahan server | Hubungi admin |

---

## Rate Limiting

Saat ini tidak ada rate limiting untuk Jadwal Kelas API. Namun, gunakan secara bijak untuk menghindari beban server.

---

## Notes

- Setiap ID menggunakan UUID format
- Timestamp otomatis diisi oleh sistem (`created_at`, `updated_at`)
- Token JWT harus valid dan tidak expired
- Semua field pada request harus diisi (tidak ada field optional)
- Tidak ada endpoint restore — hapus bersifat permanent
- Pengecekan duplikat dilakukan **sebelum** create dan update untuk validasi lebih baik

---

## Support

Untuk pertanyaan atau laporan bug, hubungi tim development.
