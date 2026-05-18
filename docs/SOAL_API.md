# Dokumentasi API Soal (Questions)

## 📌 Base URL
```
http://localhost:3000/api/soal
```

---

## 🔐 Authentication
Endpoint POST, PUT, DELETE, PATCH memerlukan JWT Token di header:
```
Authorization: Bearer <JWT_TOKEN>
```

Endpoint GET bersifat public (tidak memerlukan authentication).

---

## 📚 Endpoints

### 1️⃣ GET - Daftar Semua Soal (Dengan Pagination)

**Endpoint:** `GET /api/soal`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/soal?page=1&page_size=10"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get all soal successfully",
  "data": {
    "data": [
      {
        "id": "abc123def456",
        "id_bank_soal": "5112e444-25d8-4ca6-859f-3d24099f45ce",
        "no_soal": 1,
        "soal": "Berapa hasil dari 2 + 2?",
        "gambar_soal": "https://example.com/soal.jpg",
        "opsi_a": "3",
        "opsi_b": "4",
        "opsi_c": "5",
        "opsi_d": "6",
        "opsi_e": "7",
        "gambar_a": "https://example.com/a.jpg",
        "gambar_b": "https://example.com/b.jpg",
        "gambar_c": "https://example.com/c.jpg",
        "gambar_d": "https://example.com/d.jpg",
        "gambar_e": "https://example.com/e.jpg",
        "kunci": "B",
        "created_at": "2026-05-18 14:00:00",
        "updated_at": "2026-05-18 14:00:00"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10,
    "total_page": 5
  },
  "errors": null
}
```

---

### 2️⃣ GET - Soal by Bank Soal (Filter by Bank Soal ID)

**Endpoint:** `GET /api/soal/bank/:bank_soal_id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `bank_soal_id` | string (UUID) | ID bank_soal |

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/soal/bank/5112e444-25d8-4ca6-859f-3d24099f45ce?page=1&page_size=10"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get soal by bank successfully",
  "data": {
    "data": [
      {
        "id": "abc123def456",
        "id_bank_soal": "5112e444-25d8-4ca6-859f-3d24099f45ce",
        "no_soal": 1,
        "soal": "Berapa hasil dari 2 + 2?",
        "gambar_soal": "https://example.com/soal.jpg",
        "opsi_a": "3",
        "opsi_b": "4",
        "opsi_c": "5",
        "opsi_d": "6",
        "opsi_e": "7",
        "gambar_a": "https://example.com/a.jpg",
        "gambar_b": "https://example.com/b.jpg",
        "gambar_c": "https://example.com/c.jpg",
        "gambar_d": "https://example.com/d.jpg",
        "gambar_e": "https://example.com/e.jpg",
        "kunci": "B",
        "created_at": "2026-05-18 14:00:00",
        "updated_at": "2026-05-18 14:00:00"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 10,
    "total_page": 2
  },
  "errors": null
}
```

---

### 3️⃣ GET - Detail Soal by ID

**Endpoint:** `GET /api/soal/:id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID soal |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/soal/abc123def456"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get soal successfully",
  "data": {
    "id": "abc123def456",
    "id_bank_soal": "5112e444-25d8-4ca6-859f-3d24099f45ce",
    "no_soal": 1,
    "soal": "Berapa hasil dari 2 + 2?",
    "gambar_soal": "https://example.com/soal.jpg",
    "opsi_a": "3",
    "opsi_b": "4",
    "opsi_c": "5",
    "opsi_d": "6",
    "opsi_e": "7",
    "gambar_a": "https://example.com/a.jpg",
    "gambar_b": "https://example.com/b.jpg",
    "gambar_c": "https://example.com/c.jpg",
    "gambar_d": "https://example.com/d.jpg",
    "gambar_e": "https://example.com/e.jpg",
    "kunci": "B",
    "created_at": "2026-05-18 14:00:00",
    "updated_at": "2026-05-18 14:00:00"
  },
  "errors": null
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Resource not found",
  "data": null,
  "errors": null
}
```

---

### 4️⃣ POST - Buat Soal Baru (dengan Upload Gambar)

**Endpoint:** `POST /api/soal`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: multipart/form-data
Authorization: Bearer <JWT_TOKEN>
```

**Form Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id_bank_soal` | string (UUID) | ✅ Yes | ID bank_soal |
| `no_soal` | integer | ✅ Yes | Nomor soal (min 1) |
| `soal` | string | ✅ Yes | Pertanyaan soal |
| `gambar_soal` | file (jpg, jpeg, png, gif, webp) | ❌ No | Gambar pertanyaan (max 5MB) |
| `opsi_a` | string | ✅ Yes | Opsi A |
| `opsi_b` | string | ✅ Yes | Opsi B |
| `opsi_c` | string | ✅ Yes | Opsi C |
| `opsi_d` | string | ❌ No | Opsi D |
| `opsi_e` | string | ❌ No | Opsi E |
| `gambar_a` | file | ❌ No | Gambar opsi A (max 5MB) |
| `gambar_b` | file | ❌ No | Gambar opsi B (max 5MB) |
| `gambar_c` | file | ❌ No | Gambar opsi C (max 5MB) |
| `gambar_d` | file | ❌ No | Gambar opsi D (max 5MB) |
| `gambar_e` | file | ❌ No | Gambar opsi E (max 5MB) |
| `kunci` | string | ✅ Yes | Jawaban benar (A/B/C/D/E) |

**Validasi:**
- Format gambar: jpg, jpeg, png, gif, webp
- Ukuran maksimal: 5MB per file
- Kunci harus A, B, C, D, atau E
- Jika kunci adalah D atau E, opsi yang dirujuk tidak boleh kosong

**Request Example:**
```bash
curl -X POST "http://localhost:3000/api/soal" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -F "id_bank_soal=5112e444-25d8-4ca6-859f-3d24099f45ce" \
  -F "no_soal=1" \
  -F "soal=Berapa hasil dari 2 + 2?" \
  -F "gambar_soal=@soal.jpg" \
  -F "opsi_a=3" \
  -F "opsi_b=4" \
  -F "opsi_c=5" \
  -F "opsi_d=6" \
  -F "opsi_e=7" \
  -F "gambar_a=@a.jpg" \
  -F "gambar_b=@b.jpg" \
  -F "gambar_c=@c.jpg" \
  -F "gambar_d=@d.jpg" \
  -F "gambar_e=@e.jpg" \
  -F "kunci=B"
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Create soal successfully",
  "data": {
    "id": "abc123def456",
    "id_bank_soal": "5112e444-25d8-4ca6-859f-3d24099f45ce",
    "no_soal": 1,
    "soal": "Berapa hasil dari 2 + 2?",
    "gambar_soal": "http://localhost:3000/uploads/soal/1684426848_abc12345.jpg",
    "opsi_a": "3",
    "opsi_b": "4",
    "opsi_c": "5",
    "opsi_d": "6",
    "opsi_e": "7",
    "gambar_a": "http://localhost:3000/uploads/soal/1684426849_def67890.jpg",
    "gambar_b": "http://localhost:3000/uploads/soal/1684426850_ghi11111.jpg",
    "gambar_c": "http://localhost:3000/uploads/soal/1684426851_jkl22222.jpg",
    "gambar_d": "http://localhost:3000/uploads/soal/1684426852_mno33333.jpg",
    "gambar_e": "http://localhost:3000/uploads/soal/1684426853_pqr44444.jpg",
    "kunci": "B",
    "created_at": "2026-05-18 14:00:00",
    "updated_at": "2026-05-18 14:00:00"
  },
  "errors": null
}
```

---

### 5️⃣ PUT - Update Soal (dengan Upload Gambar)

**Endpoint:** `PUT /api/soal/:id`

**Authentication:** ✅ Required (JWT Token)

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID soal yang akan diupdate |

**Request Headers:**
```
Content-Type: multipart/form-data
Authorization: Bearer <JWT_TOKEN>
```

**Form Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `no_soal` | integer | ✅ Yes | Nomor soal (min 1) |
| `soal` | string | ✅ Yes | Pertanyaan soal |
| `gambar_soal` | file | ❌ No | Gambar pertanyaan baru (max 5MB) |
| `opsi_a` | string | ✅ Yes | Opsi A |
| `opsi_b` | string | ✅ Yes | Opsi B |
| `opsi_c` | string | ✅ Yes | Opsi C |
| `opsi_d` | string | ❌ No | Opsi D |
| `opsi_e` | string | ❌ No | Opsi E |
| `gambar_a` | file | ❌ No | Gambar opsi A baru (max 5MB) |
| `gambar_b` | file | ❌ No | Gambar opsi B baru (max 5MB) |
| `gambar_c` | file | ❌ No | Gambar opsi C baru (max 5MB) |
| `gambar_d` | file | ❌ No | Gambar opsi D baru (max 5MB) |
| `gambar_e` | file | ❌ No | Gambar opsi E baru (max 5MB) |
| `kunci` | string | ✅ Yes | Jawaban benar (A/B/C/D/E) |

**Request Example (Partial Update):**
```bash
curl -X PUT "http://localhost:3000/api/soal/abc123def456" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -F "no_soal=2" \
  -F "soal=Berapa hasil dari 3 + 3?" \
  -F "opsi_a=5" \
  -F "opsi_b=6" \
  -F "opsi_c=7" \
  -F "gambar_soal=@soal-updated.jpg" \
  -F "kunci=B"
```

**Request Example (Update dengan Semua Gambar):**
```bash
curl -X PUT "http://localhost:3000/api/soal/abc123def456" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -F "no_soal=2" \
  -F "soal=Berapa hasil dari 3 + 3?" \
  -F "opsi_a=5" \
  -F "opsi_b=6" \
  -F "opsi_c=7" \
  -F "opsi_d=8" \
  -F "opsi_e=9" \
  -F "gambar_soal=@soal-updated.jpg" \
  -F "gambar_a=@a-updated.jpg" \
  -F "gambar_b=@b-updated.jpg" \
  -F "gambar_c=@c-updated.jpg" \
  -F "gambar_d=@d-updated.jpg" \
  -F "gambar_e=@e-updated.jpg" \
  -F "kunci=B"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Update soal successfully",
  "data": {
    "id": "abc123def456",
    "id_bank_soal": "5112e444-25d8-4ca6-859f-3d24099f45ce",
    "no_soal": 2,
    "soal": "Berapa hasil dari 3 + 3?",
    "gambar_soal": "http://localhost:3000/uploads/soal/1684426854_xyz55555.jpg",
    "opsi_a": "5",
    "opsi_b": "6",
    "opsi_c": "7",
    "opsi_d": "8",
    "opsi_e": "9",
    "gambar_a": "http://localhost:3000/uploads/soal/1684426855_abc66666.jpg",
    "gambar_b": "http://localhost:3000/uploads/soal/1684426856_def77777.jpg",
    "gambar_c": "http://localhost:3000/uploads/soal/1684426857_ghi88888.jpg",
    "gambar_d": "http://localhost:3000/uploads/soal/1684426858_jkl99999.jpg",
    "gambar_e": "http://localhost:3000/uploads/soal/1684426859_mno00000.jpg",
    "kunci": "B",
    "created_at": "2026-05-18 14:00:00",
    "updated_at": "2026-05-18 14:05:00"
  },
  "errors": null
}
```

**Notes:**
- Gambar lama akan otomatis dihapus saat update dengan gambar baru
- Jika update tanpa file gambar baru, gambar lama tetap dipertahankan
- Setiap upload file akan menghasilkan unique filename dengan timestamp

---

### 6️⃣ DELETE - Soft Delete Soal

**Endpoint:** `DELETE /api/soal/:id`

**Authentication:** ✅ Required (JWT Token)

**Note:** Ini adalah **Soft Delete** - data tidak benar-benar dihapus, hanya ditandai dengan `deleted_at` timestamp.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID soal yang akan dihapus |

**Request Example:**
```bash
curl -X DELETE "http://localhost:3000/api/soal/abc123def456" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Delete soal successfully",
  "data": null,
  "errors": null
}
```

---

### 7️⃣ PATCH - Restore Deleted Soal

**Endpoint:** `PATCH /api/soal/:id/restore`

**Authentication:** ✅ Required (JWT Token)

**Description:** Mengembalikan soal yang sudah di-soft delete.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID soal yang akan di-restore |

**Request Example:**
```bash
curl -X PATCH "http://localhost:3000/api/soal/abc123def456/restore" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Restore soal successfully",
  "data": null,
  "errors": null
}
```

---

## 📊 Response Format

Semua response mengikuti format standard:

**Success Response:**
```json
{
  "success": true,
  "message": "Success message",
  "data": { /* response data */ },
  "errors": null
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error message",
  "data": null,
  "errors": null
}
```

---

## 🚨 Error Messages

| Status Code | Message | Description |
|------------|---------|-------------|
| 200 | Success message | Request berhasil |
| 201 | Create soal successfully | Soal berhasil dibuat |
| 400 | Invalid request format | Format request tidak sesuai |
| 400 | Resource not found | Soal tidak ditemukan |
| 400 | kunci harus A, B, C, D, atau E | Validasi kunci gagal |
| 400 | opsi D/E tidak boleh kosong jika kunci D/E | Validasi opsi gagal |
| 401 | Unauthorized | JWT token tidak valid/expired |
| 404 | Resource not found | Soal tidak ditemukan |
| 500 | Internal server error | Error di server |

---

## 💡 Notes

1. **Soft Delete**: Data yang didelete tidak benar-benar dihapus dari database, hanya ditandai dengan deleted_at
2. **Pagination**: Gunakan pagination pada GET /api/soal untuk performa optimal
3. **Filter by Bank**: Endpoint GET /api/soal/bank/:id sangat penting untuk filter soal per bank_soal
4. **Kunci Validation**: Ensure kunci merujuk ke opsi yang valid (tidak kosong)
5. **Authorization**: Semua write operations (POST, PUT, DELETE, PATCH) memerlukan JWT token
6. **Images**: Field gambar_* bersifat opsional dan dapat berupa URL

---

Generated with Claude Code 🤖
