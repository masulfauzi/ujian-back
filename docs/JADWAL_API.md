# Jadwal API Documentation

## Base URL

```
http://localhost:3000/api/jadwal
```

---

## Overview

Jadwal API mengelola jadwal ujian. Setiap jadwal berelasi ke satu bank soal dan mencatat waktu mulai serta waktu selesai ujian.

---

## Authentication

Beberapa endpoint memerlukan JWT token:

- **POST /api/jadwal** — Memerlukan JWT
- **PUT /api/jadwal/:id** — Memerlukan JWT
- **DELETE /api/jadwal/:id** — Memerlukan JWT
- **PATCH /api/jadwal/:id/restore** — Memerlukan JWT

Kirim token di header:
```
Authorization: Bearer <token>
```

---

## Endpoints

### 1. GET /api/jadwal

Mengambil daftar semua jadwal dengan pagination.

**Method:** `GET`

**Auth:** Tidak diperlukan

**Query Parameters:**

| Parameter | Type | Default | Deskripsi |
|-----------|------|---------|-----------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |

**Example Request:**

```bash
curl -X GET "http://localhost:3000/api/jadwal?page=1&page_size=10"
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Get all jadwal successfully",
  "data": {
    "data": [
      {
        "id": "a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6",
        "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
        "nama_bank_soal": "Soal Ujian Matematika 2025",
        "wkt_mulai": "2025-08-01 08:00:00",
        "wkt_selesai": "2025-08-01 10:00:00",
        "created_at": "2025-05-21 14:30:25",
        "updated_at": "2025-05-21 14:30:25"
      },
      {
        "id": "c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8",
        "id_bank_soal": "d4e5f6a7-b8c9-50d1-e2f3-a4b5c6d7e8f9",
        "nama_bank_soal": "Soal Ujian Fisika 2025",
        "wkt_mulai": "2025-08-02 08:00:00",
        "wkt_selesai": "2025-08-02 10:00:00",
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

### 2. GET /api/jadwal/:id

Mengambil detail jadwal berdasarkan ID.

**Method:** `GET`

**Auth:** Tidak diperlukan

**Path Parameters:**

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `id` | string (UUID) | ID jadwal |

**Example Request:**

```bash
curl -X GET "http://localhost:3000/api/jadwal/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6"
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Get jadwal successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6",
    "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
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

### 3. GET /api/jadwal/bank-soal/:bank_soal_id

Mengambil daftar jadwal berdasarkan bank soal tertentu.

**Method:** `GET`

**Auth:** Tidak diperlukan

**Path Parameters:**

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `bank_soal_id` | string (UUID) | ID bank soal |

**Query Parameters:**

| Parameter | Type | Default | Deskripsi |
|-----------|------|---------|-----------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |

**Example Request:**

```bash
curl -X GET "http://localhost:3000/api/jadwal/bank-soal/b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7?page=1&page_size=10"
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Get jadwal by bank soal successfully",
  "data": {
    "data": [
      {
        "id": "a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6",
        "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
        "nama_bank_soal": "Soal Ujian Matematika 2025",
        "wkt_mulai": "2025-08-01 08:00:00",
        "wkt_selesai": "2025-08-01 10:00:00",
        "created_at": "2025-05-21 14:30:25",
        "updated_at": "2025-05-21 14:30:25"
      }
    ],
    "total": 1,
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

### 4. POST /api/jadwal

Membuat jadwal ujian baru.

**Method:** `POST`

**Auth:** JWT diperlukan

**Request Body:**

```json
{
  "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
  "wkt_mulai": "2025-08-01 08:00:00",
  "wkt_selesai": "2025-08-01 10:00:00"
}
```

**Request Body Parameters:**

| Parameter | Type | Required | Deskripsi |
|-----------|------|----------|-----------|
| `id_bank_soal` | string (UUID) | Ya | ID bank soal yang valid |
| `wkt_mulai` | string | Ya | Waktu mulai (format: `YYYY-MM-DD HH:MM:SS`) |
| `wkt_selesai` | string | Ya | Waktu selesai (format: `YYYY-MM-DD HH:MM:SS`) |

**Example Request:**

```bash
curl -X POST "http://localhost:3000/api/jadwal" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
    "wkt_mulai": "2025-08-01 08:00:00",
    "wkt_selesai": "2025-08-01 10:00:00"
  }'
```

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "Create jadwal successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6",
    "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
    "nama_bank_soal": "Soal Ujian Matematika 2025",
    "wkt_mulai": "2025-08-01 08:00:00",
    "wkt_selesai": "2025-08-01 10:00:00",
    "created_at": "2025-05-21 14:30:25",
    "updated_at": "2025-05-21 14:30:25"
  }
}
```

**Error Response (400 Bad Request - Invalid Format):**

```json
{
  "success": false,
  "message": "format wkt_mulai tidak valid, gunakan: 2006-01-02 15:04:05",
  "data": null
}
```

**Error Response (400 Bad Request - Invalid Time Order):**

```json
{
  "success": false,
  "message": "wkt_selesai harus setelah wkt_mulai",
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

### 5. PUT /api/jadwal/:id

Mengupdate data jadwal.

**Method:** `PUT`

**Auth:** JWT diperlukan

**Path Parameters:**

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `id` | string (UUID) | ID jadwal |

**Request Body:**

```json
{
  "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
  "wkt_mulai": "2025-08-01 09:00:00",
  "wkt_selesai": "2025-08-01 11:00:00"
}
```

**Request Body Parameters:**

| Parameter | Type | Required | Deskripsi |
|-----------|------|----------|-----------|
| `id_bank_soal` | string (UUID) | Ya | ID bank soal yang valid |
| `wkt_mulai` | string | Ya | Waktu mulai (format: `YYYY-MM-DD HH:MM:SS`) |
| `wkt_selesai` | string | Ya | Waktu selesai (format: `YYYY-MM-DD HH:MM:SS`) |

**Example Request:**

```bash
curl -X PUT "http://localhost:3000/api/jadwal/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
    "wkt_mulai": "2025-08-01 09:00:00",
    "wkt_selesai": "2025-08-01 11:00:00"
  }'
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Update jadwal successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6",
    "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
    "nama_bank_soal": "Soal Ujian Matematika 2025",
    "wkt_mulai": "2025-08-01 09:00:00",
    "wkt_selesai": "2025-08-01 11:00:00",
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

**Error Response (400 Bad Request - Invalid Time):**

```json
{
  "success": false,
  "message": "wkt_selesai harus setelah wkt_mulai",
  "data": null
}
```

---

### 6. DELETE /api/jadwal/:id

Soft delete jadwal (tidak benar-benar menghapus, hanya menandai sebagai dihapus).

**Method:** `DELETE`

**Auth:** JWT diperlukan

**Path Parameters:**

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `id` | string (UUID) | ID jadwal |

**Example Request:**

```bash
curl -X DELETE "http://localhost:3000/api/jadwal/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Delete jadwal successfully",
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

---

### 7. PATCH /api/jadwal/:id/restore

Mengembalikan jadwal yang sudah dihapus (undo soft delete).

**Method:** `PATCH`

**Auth:** JWT diperlukan

**Path Parameters:**

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `id` | string (UUID) | ID jadwal |

**Example Request:**

```bash
curl -X PATCH "http://localhost:3000/api/jadwal/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6/restore" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Restore jadwal successfully",
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

### 1. Format Datetime

`wkt_mulai` dan `wkt_selesai` harus dikirim dalam format:
```
YYYY-MM-DD HH:MM:SS
```

Contoh: `2025-08-01 08:00:00`

### 2. Validasi Urutan Waktu

`wkt_selesai` **harus selalu setelah** `wkt_mulai`. Jika tidak, API mengembalikan error 400.

**Valid:**
- `wkt_mulai: 2025-08-01 08:00:00`
- `wkt_selesai: 2025-08-01 10:00:00` ✅

**Invalid:**
- `wkt_mulai: 2025-08-01 10:00:00`
- `wkt_selesai: 2025-08-01 08:00:00` ❌

### 3. Soft Delete

Data jadwal tidak benar-benar dihapus dari database. Ketika didelete:
- Field `deleted_at` diisi dengan timestamp waktu penghapusan
- Data tidak akan muncul di endpoint GET
- Data dapat dikembalikan dengan endpoint PATCH restore

### 4. Relasi ke Bank Soal

Setiap jadwal wajib memiliki `id_bank_soal` yang valid dan mengacu ke tabel `bank_soal`.

Response API otomatis menyertakan `nama_bank_soal` melalui JOIN query.

### 5. Response Format

- Response sukses selalu mengembalikan:
  - `success: true`
  - `message`: Deskripsi operasi
  - `data`: Data jadwal atau list jadwal

- Response error selalu mengembalikan:
  - `success: false`
  - `message`: Deskripsi error
  - `data: null`

---

## Common Use Cases

### Use Case 1: Membuat jadwal ujian untuk suatu bank soal

```bash
# 1. Dapatkan ID bank soal terlebih dahulu
curl http://localhost:3000/api/bank-soal

# 2. Buat jadwal dengan bank soal yang dipilih
curl -X POST "http://localhost:3000/api/jadwal" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
    "wkt_mulai": "2025-08-01 08:00:00",
    "wkt_selesai": "2025-08-01 10:00:00"
  }'
```

### Use Case 2: Melihat semua jadwal ujian

```bash
curl http://localhost:3000/api/jadwal?page=1&page_size=10
```

### Use Case 3: Melihat jadwal ujian untuk bank soal tertentu

```bash
curl http://localhost:3000/api/jadwal/bank-soal/b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7
```

### Use Case 4: Mengubah jadwal ujian

```bash
curl -X PUT "http://localhost:3000/api/jadwal/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id_bank_soal": "b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7",
    "wkt_mulai": "2025-08-01 09:00:00",
    "wkt_selesai": "2025-08-01 11:00:00"
  }'
```

### Use Case 5: Menghapus jadwal ujian

```bash
curl -X DELETE "http://localhost:3000/api/jadwal/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Use Case 6: Mengembalikan jadwal yang sudah dihapus

```bash
curl -X PATCH "http://localhost:3000/api/jadwal/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6/restore" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Rate Limiting

Saat ini tidak ada rate limiting untuk Jadwal API. Namun, gunakan secara bijak untuk menghindari beban server.

---

## Pagination

Untuk endpoint yang mengembalikan list (GET /api/jadwal, GET /api/jadwal/bank-soal/:id), gunakan query parameter:

- `page` — Default: 1
- `page_size` — Default: 10, Max: 100 (tidak divalidasi, tapi sebaiknya ≤ 100)

Response akan mencakup:
- `total` — Total jumlah data di database
- `page` — Halaman saat ini
- `page_size` — Jumlah item per halaman
- `total_page` — Total jumlah halaman
- `data` — Array jadwal

---

## Notes

- Setiap ID menggunakan UUID format
- Timestamp otomatis diisi oleh sistem (created_at, updated_at)
- Token JWT harus valid dan tidak expired
- Semua field pada request harus diisi (tidak ada field optional)

---

## Support

Untuk pertanyaan atau laporan bug, hubungi tim development.
