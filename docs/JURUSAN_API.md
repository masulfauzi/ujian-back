# Dokumentasi API Jurusan (Department Management)

## 📌 Base URL
```
http://localhost:3000/api/jurusan
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

### 1️⃣ GET - Daftar Semua Jurusan (Dengan Pagination)

**Endpoint:** `GET /api/jurusan`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/jurusan?page=1&page_size=10"
```

**Success Response (200 OK):**
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
      },
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "nama_jurusan": "Rekayasa Perangkat Lunak",
        "created_at": "2026-05-19 10:00:00",
        "updated_at": "2026-05-19 10:00:00"
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

### 2️⃣ GET - Detail Jurusan by ID

**Endpoint:** `GET /api/jurusan/:id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID jurusan |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/jurusan/550e8400-e29b-41d4-a716-446655440000"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get jurusan successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "nama_jurusan": "Teknik Komputer dan Jaringan",
    "created_at": "2026-05-19 10:00:00",
    "updated_at": "2026-05-19 10:00:00"
  },
  "errors": null
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "data not found",
  "data": null,
  "errors": null
}
```

---

### 3️⃣ POST - Buat Jurusan Baru

**Endpoint:** `POST /api/jurusan`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
  "nama_jurusan": "Teknik Komputer dan Jaringan"
}
```

**Field Requirements:**
| Field | Type | Required | Max Length | Description |
|-------|------|----------|-----------|-------------|
| `nama_jurusan` | string | ✅ Yes | 255 | Nama jurusan (UNIQUE) |

**Request Example:**
```bash
curl -X POST "http://localhost:3000/api/jurusan" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "nama_jurusan": "Multimedia"
  }'
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Create jurusan successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "nama_jurusan": "Multimedia",
    "created_at": "2026-05-19 10:15:30",
    "updated_at": "2026-05-19 10:15:30"
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

### 4️⃣ PUT - Update Jurusan

**Endpoint:** `PUT /api/jurusan/:id`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID jurusan yang akan diupdate |

**Request Body:**
```json
{
  "nama_jurusan": "Multimedia Design"
}
```

**Request Example:**
```bash
curl -X PUT "http://localhost:3000/api/jurusan/550e8400-e29b-41d4-a716-446655440002" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "nama_jurusan": "Multimedia Design"
  }'
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Update jurusan successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "nama_jurusan": "Multimedia Design",
    "created_at": "2026-05-19 10:15:30",
    "updated_at": "2026-05-19 10:20:00"
  },
  "errors": null
}
```

---

### 5️⃣ DELETE - Soft Delete Jurusan

**Endpoint:** `DELETE /api/jurusan/:id`

**Authentication:** ✅ Required (JWT Token)

**Note:** Ini adalah **Soft Delete** - data tidak benar-benar dihapus, hanya ditandai dengan `deleted_at` timestamp.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID jurusan yang akan dihapus |

**Request Example:**
```bash
curl -X DELETE "http://localhost:3000/api/jurusan/550e8400-e29b-41d4-a716-446655440002" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Delete jurusan successfully",
  "data": null,
  "errors": null
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "data not found",
  "data": null,
  "errors": null
}
```

---

### 6️⃣ PATCH - Restore Deleted Jurusan

**Endpoint:** `PATCH /api/jurusan/:id/restore`

**Authentication:** ✅ Required (JWT Token)

**Description:** Mengembalikan jurusan yang sudah di-soft delete.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID jurusan yang akan di-restore |

**Request Example:**
```bash
curl -X PATCH "http://localhost:3000/api/jurusan/550e8400-e29b-41d4-a716-446655440002/restore" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Restore jurusan successfully",
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

### List Jurusan (Public)
```bash
curl "http://localhost:3000/api/jurusan?page=1&page_size=5"
```

### Create Jurusan (Require Auth)
```bash
TOKEN="your_token_here"
curl -X POST "http://localhost:3000/api/jurusan" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "nama_jurusan": "Akuntansi"
  }'
```

### Get Detail Jurusan (Public)
```bash
curl "http://localhost:3000/api/jurusan/550e8400-e29b-41d4-a716-446655440000"
```

### Update Jurusan (Require Auth)
```bash
TOKEN="your_token_here"
curl -X PUT "http://localhost:3000/api/jurusan/550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "nama_jurusan": "Akuntansi dan Keuangan"
  }'
```

### Delete Jurusan (Require Auth)
```bash
TOKEN="your_token_here"
curl -X DELETE "http://localhost:3000/api/jurusan/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer $TOKEN"
```

### Restore Jurusan (Require Auth)
```bash
TOKEN="your_token_here"
curl -X PATCH "http://localhost:3000/api/jurusan/550e8400-e29b-41d4-a716-446655440000/restore" \
  -H "Authorization: Bearer $TOKEN"
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
| 201 | Create jurusan successfully | Jurusan berhasil dibuat |
| 400 | Invalid request format | Format request tidak sesuai |
| 400 | data not found | Jurusan tidak ditemukan |
| 401 | Unauthorized | JWT token tidak valid/expired |
| 500 | Internal server error | Error di server |

---

## 📋 Field Details

### ID (UUID)
- Auto-generated oleh database
- Format: `550e8400-e29b-41d4-a716-446655440000`
- Tidak bisa dimodifikasi setelah dibuat

### Nama Jurusan
- Max 255 karakter
- Harus UNIQUE (tidak boleh duplikat)
- Wajib diisi (required)
- Contoh: "Teknik Komputer dan Jaringan", "Rekayasa Perangkat Lunak", "Multimedia"

### Timestamps
- Format: `YYYY-MM-DD HH:MM:SS`
- Contoh: `2026-05-19 10:00:00`
- `created_at`: Auto-set saat record dibuat, tidak bisa diubah
- `updated_at`: Auto-update saat record diubah
- `deleted_at`: NULL saat aktif, diisi saat soft delete

---

## 🔄 Data Flow Example

```
1. Register user → Get JWT token
   POST /api/auth/register
   Response: { token: "..." }

2. Create jurusan → Get jurusan ID
   POST /api/jurusan (with token)
   Response: { id: "...", nama_jurusan: "..." }

3. List jurusan → Lihat semua jurusan
   GET /api/jurusan
   Response: { data: [...], total: 5 }

4. Get detail → Lihat detail jurusan
   GET /api/jurusan/{id}
   Response: { id: "...", nama_jurusan: "..." }

5. Update jurusan → Edit jurusan
   PUT /api/jurusan/{id} (with token)
   Response: { id: "...", nama_jurusan: "..." (updated) }

6. Delete jurusan → Soft delete jurusan
   DELETE /api/jurusan/{id} (with token)
   Response: { success: true }

7. Restore jurusan → Undo soft delete
   PATCH /api/jurusan/{id}/restore (with token)
   Response: { success: true }
```

---

## 💡 Tips & Best Practices

1. **Pagination**: Selalu gunakan pagination pada GET /api/jurusan untuk performa optimal
   - Gunakan `page_size=100` jika ingin ambil semua data untuk dropdown
   
2. **Error Handling**: Selalu check `success` field dalam response sebelum mengakses `data`
   
3. **Authentication**: Simpan JWT token dengan aman, jangan expose di client-side
   
4. **Soft Delete**: Record yang sudah didelete masih bisa di-restore kapan saja
   - Data yang didelete tidak akan muncul di GET requests
   - Gunakan restore endpoint untuk membawa kembali data
   
5. **Validation**: 
   - Field `nama_jurusan` tidak boleh duplikat
   - `nama_jurusan` tidak boleh kosong
   
6. **Frontend Integration**:
   - Untuk dropdown select: `GET /api/jurusan?page=1&page_size=100`
   - For single select by ID: `GET /api/jurusan/{id}`
   - Always handle error responses gracefully

---

## 🔗 Related Endpoints

- **Auth**: `/api/auth/register`, `/api/auth/login`, `/api/auth/logout`
- **User**: `/api/user/{id}`
- **Mapel**: `/api/mapel`

---

## 📞 Support

Untuk pertanyaan atau issue, buat issue di repository atau hubungi tim development.

