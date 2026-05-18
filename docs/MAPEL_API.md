# Dokumentasi API Mapel (Subject Management)

## 📌 Base URL
```
http://localhost:3000/api/mapel
```

---

## 🔐 Authentication
Semua endpoint POST, PUT, DELETE, PATCH memerlukan JWT Token di header:
```
Authorization: Bearer <JWT_TOKEN>
```

Endpoint GET bersifat public (tidak memerlukan authentication).

---

## 📚 Endpoints

### 1️⃣ GET - Daftar Semua Mapel (Dengan Pagination)

**Endpoint:** `GET /api/mapel`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/mapel?page=1&page_size=10"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get all mapel successfully",
  "data": {
    "data": [
      {
        "id": "916f68b8-e24d-4595-84c8-2775076c54e0",
        "nama_mapel": "Matematika",
        "kode_mapel": "MAT",
        "deskripsi": "Pelajaran Matematika",
        "created_at": "2026-05-18 09:58:20",
        "updated_at": "2026-05-18 09:58:20"
      },
      {
        "id": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
        "nama_mapel": "Bahasa Indonesia",
        "kode_mapel": "IND",
        "deskripsi": "Pelajaran Bahasa Indonesia",
        "created_at": "2026-05-18 09:58:20",
        "updated_at": "2026-05-18 09:58:20"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10,
    "total_page": 1
  },
  "errors": null
}
```

---

### 2️⃣ GET - Detail Mapel by ID

**Endpoint:** `GET /api/mapel/:id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID mapel |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/mapel/916f68b8-e24d-4595-84c8-2775076c54e0"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get mapel successfully",
  "data": {
    "id": "916f68b8-e24d-4595-84c8-2775076c54e0",
    "nama_mapel": "Matematika",
    "kode_mapel": "MAT",
    "deskripsi": "Pelajaran Matematika",
    "created_at": "2026-05-18 09:58:20",
    "updated_at": "2026-05-18 09:58:20"
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

### 3️⃣ POST - Buat Mapel Baru

**Endpoint:** `POST /api/mapel`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
  "nama_mapel": "Seni Budaya",
  "kode_mapel": "SBD",
  "deskripsi": "Pelajaran Seni Budaya"
}
```

**Field Requirements:**
| Field | Type | Required | Max Length | Description |
|-------|------|----------|-----------|-------------|
| `nama_mapel` | string | ✅ Yes | 255 | Nama mata pelajaran (UNIQUE) |
| `kode_mapel` | string | ✅ Yes | 20 | Kode singkat mapel (UNIQUE) |
| `deskripsi` | string | ❌ No | - | Deskripsi mapel |

**Request Example:**
```bash
curl -X POST "http://localhost:3000/api/mapel" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "nama_mapel": "Seni Budaya",
    "kode_mapel": "SBD",
    "deskripsi": "Pelajaran Seni Budaya"
  }'
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Create mapel successfully",
  "data": {
    "id": "b19530a1-783a-4c9d-a232-353f6f5aa766",
    "nama_mapel": "Seni Budaya",
    "kode_mapel": "SBD",
    "deskripsi": "Pelajaran Seni Budaya",
    "created_at": "2026-05-18 09:59:37",
    "updated_at": "2026-05-18 09:59:37"
  },
  "errors": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Invalid request format",
  "data": null,
  "errors": null
}
```

---

### 4️⃣ PUT - Update Mapel

**Endpoint:** `PUT /api/mapel/:id`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID mapel yang akan diupdate |

**Request Body:**
```json
{
  "nama_mapel": "Seni Budaya Updated",
  "kode_mapel": "SBD",
  "deskripsi": "Deskripsi yang sudah diupdate"
}
```

**Request Example:**
```bash
curl -X PUT "http://localhost:3000/api/mapel/b19530a1-783a-4c9d-a232-353f6f5aa766" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "nama_mapel": "Seni Budaya Updated",
    "kode_mapel": "SBD",
    "deskripsi": "Updated deskripsi"
  }'
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Update mapel successfully",
  "data": {
    "id": "b19530a1-783a-4c9d-a232-353f6f5aa766",
    "nama_mapel": "Seni Budaya Updated",
    "kode_mapel": "SBD",
    "deskripsi": "Updated deskripsi",
    "created_at": "2026-05-18 09:59:37",
    "updated_at": "2026-05-18 10:05:00"
  },
  "errors": null
}
```

---

### 5️⃣ DELETE - Soft Delete Mapel

**Endpoint:** `DELETE /api/mapel/:id`

**Authentication:** ✅ Required (JWT Token)

**Note:** Ini adalah **Soft Delete** - data tidak benar-benar dihapus, hanya ditandai dengan `deleted_at` timestamp.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID mapel yang akan dihapus |

**Request Example:**
```bash
curl -X DELETE "http://localhost:3000/api/mapel/b19530a1-783a-4c9d-a232-353f6f5aa766" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Delete mapel successfully",
  "data": null,
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

### 6️⃣ PATCH - Restore Deleted Mapel

**Endpoint:** `PATCH /api/mapel/:id/restore`

**Authentication:** ✅ Required (JWT Token)

**Description:** Mengembalikan mapel yang sudah di-soft delete.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID mapel yang akan di-restore |

**Request Example:**
```bash
curl -X PATCH "http://localhost:3000/api/mapel/b19530a1-783a-4c9d-a232-353f6f5aa766/restore" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Restore mapel successfully",
  "data": null,
  "errors": null
}
```

---

## 🧪 Testing dengan cURL

### Register & Get Token
```bash
curl -X POST "http://localhost:3000/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@test.com",
    "password": "password123"
  }'
```

Simpan `token` dari response untuk digunakan di request berikutnya.

### List Mapel (Public)
```bash
curl "http://localhost:3000/api/mapel?page=1&page_size=5"
```

### Create Mapel (Require Auth)
```bash
TOKEN="your_token_here"
curl -X POST "http://localhost:3000/api/mapel" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "nama_mapel": "Olahraga",
    "kode_mapel": "OLH",
    "deskripsi": "Pelajaran Olahraga"
  }'
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
| 201 | Create mapel successfully | Mapel berhasil dibuat |
| 400 | Invalid request format | Format request tidak sesuai |
| 400 | Resource not found | Mapel tidak ditemukan |
| 401 | Unauthorized | JWT token tidak valid/expired |
| 500 | Internal server error | Error di server |

---

## 📋 Field Details

### ID (UUID)
- Auto-generated oleh database
- Format: `550e8400-e29b-41d4-a716-446655440000`
- Contoh: `916f68b8-e24d-4595-84c8-2775076c54e0`

### Nama Mapel
- Max 255 karakter
- Harus UNIQUE (tidak boleh duplikat)
- Contoh: "Matematika", "Bahasa Indonesia"

### Kode Mapel
- Max 20 karakter
- Harus UNIQUE (tidak boleh duplikat)
- Contoh: "MAT", "IND", "ENG"

### Deskripsi
- Opsional (boleh kosong/tidak diisi)
- Tidak ada batasan panjang
- Contoh: "Pelajaran Matematika dasar"

### Timestamps
- Format: `YYYY-MM-DD HH:MM:SS`
- Contoh: `2026-05-18 09:58:20`
- `created_at`: Auto-set saat record dibuat
- `updated_at`: Auto-update saat record diubah
- `deleted_at`: NULL saat aktif, diisi saat soft delete

---

## 🔄 Data Flow Example

```
1. Register user → Get JWT token
   POST /api/auth/register
   Response: { token: "..." }

2. Create mapel → Get mapel ID
   POST /api/mapel (with token)
   Response: { id: "...", nama_mapel: "..." }

3. List mapel → Lihat semua mapel
   GET /api/mapel
   Response: { data: [...], total: 5 }

4. Get detail → Lihat detail mapel
   GET /api/mapel/{id}
   Response: { id: "...", nama_mapel: "..." }

5. Update mapel → Edit mapel
   PUT /api/mapel/{id} (with token)
   Response: { id: "...", nama_mapel: "..." (updated) }

6. Delete mapel → Soft delete mapel
   DELETE /api/mapel/{id} (with token)
   Response: { success: true }

7. Restore mapel → Undo soft delete
   PATCH /api/mapel/{id}/restore (with token)
   Response: { success: true }
```

---

## 💡 Tips & Best Practices

1. **Pagination**: Selalu gunakan pagination pada GET /api/mapel untuk performa optimal
2. **Error Handling**: Selalu check `success` field dalam response
3. **Authentication**: Simpan JWT token dengan aman, jangan expose di client-side
4. **Soft Delete**: Record yang sudah didelete masih bisa di-restore kapan saja
5. **Validation**: Field `nama_mapel` dan `kode_mapel` tidak boleh duplikat

---

## 🔗 Related Endpoints

- **Auth**: `/api/auth/register`, `/api/auth/login`
- **User**: `/api/user/{id}`

---

## 📞 Support

Untuk pertanyaan atau issue, buat issue di repository atau hubungi tim development.

