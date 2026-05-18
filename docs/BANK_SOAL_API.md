# Dokumentasi API Bank Soal (Question Bank Management)

## 📌 Base URL
```
http://localhost:3000/api/bank-soal
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

### 1️⃣ GET - Daftar Semua Bank Soal (Dengan Pagination)

**Endpoint:** `GET /api/bank-soal`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/bank-soal?page=1&page_size=10"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get all bank soal successfully",
  "data": {
    "data": [
      {
        "id": "5112e444-25d8-4ca6-859f-3d24099f45ce",
        "nama_bank_soal": "Bank Soal Matematika Dasar",
        "id_mapel": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
        "jml_soal": 50,
        "deskripsi": "Kumpulan soal matematika level dasar",
        "created_at": "2026-05-18 10:39:07",
        "updated_at": "2026-05-18 10:39:07"
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

### 2️⃣ GET - Bank Soal by Mapel (Filter by Mata Pelajaran)

**Endpoint:** `GET /api/bank-soal/mapel/:mapel_id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `mapel_id` | string (UUID) | ID mata pelajaran |

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/bank-soal/mapel/f71b3390-beeb-4e06-8fd4-deb1778c8210?page=1&page_size=10"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get bank soal by mapel successfully",
  "data": {
    "data": [
      {
        "id": "5112e444-25d8-4ca6-859f-3d24099f45ce",
        "nama_bank_soal": "Bank Soal Matematika Dasar",
        "id_mapel": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
        "jml_soal": 50,
        "deskripsi": "Kumpulan soal matematika level dasar",
        "created_at": "2026-05-18 10:39:07",
        "updated_at": "2026-05-18 10:39:07"
      }
    ],
    "total": 2,
    "page": 1,
    "page_size": 10,
    "total_page": 1
  },
  "errors": null
}
```

---

### 3️⃣ GET - Detail Bank Soal by ID

**Endpoint:** `GET /api/bank-soal/:id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID bank soal |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/bank-soal/5112e444-25d8-4ca6-859f-3d24099f45ce"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get bank soal successfully",
  "data": {
    "id": "5112e444-25d8-4ca6-859f-3d24099f45ce",
    "nama_bank_soal": "Bank Soal Matematika Dasar",
    "id_mapel": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
    "jml_soal": 50,
    "deskripsi": "Kumpulan soal matematika level dasar",
    "created_at": "2026-05-18 10:39:07",
    "updated_at": "2026-05-18 10:39:07"
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

### 4️⃣ POST - Buat Bank Soal Baru

**Endpoint:** `POST /api/bank-soal`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
  "nama_bank_soal": "Bank Soal Seni Budaya",
  "id_mapel": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
  "jml_soal": 30,
  "deskripsi": "Kumpulan soal seni budaya"
}
```

**Field Requirements:**
| Field | Type | Required | Max Length | Description |
|-------|------|----------|-----------|-------------|
| `nama_bank_soal` | string | ✅ Yes | 255 | Nama bank soal (UNIQUE) |
| `id_mapel` | string | ✅ Yes | - | UUID mata pelajaran |
| `jml_soal` | integer | ✅ Yes | - | Jumlah soal (min: 0) |
| `deskripsi` | string | ❌ No | - | Deskripsi bank soal |

**Request Example:**
```bash
curl -X POST "http://localhost:3000/api/bank-soal" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "nama_bank_soal": "Bank Soal Seni Budaya",
    "id_mapel": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
    "jml_soal": 30,
    "deskripsi": "Kumpulan soal seni budaya"
  }'
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Create bank soal successfully",
  "data": {
    "id": "5112e444-25d8-4ca6-859f-3d24099f45ce",
    "nama_bank_soal": "Bank Soal Seni Budaya",
    "id_mapel": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
    "jml_soal": 30,
    "deskripsi": "Kumpulan soal seni budaya",
    "created_at": "2026-05-18 10:39:07",
    "updated_at": "2026-05-18 10:39:07"
  },
  "errors": null
}
```

---

### 5️⃣ PUT - Update Bank Soal

**Endpoint:** `PUT /api/bank-soal/:id`

**Authentication:** ✅ Required (JWT Token)

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID bank soal yang akan diupdate |

**Request Body:**
```json
{
  "nama_bank_soal": "Bank Soal Seni Budaya Updated",
  "id_mapel": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
  "jml_soal": 35,
  "deskripsi": "Updated deskripsi"
}
```

**Request Example:**
```bash
curl -X PUT "http://localhost:3000/api/bank-soal/5112e444-25d8-4ca6-859f-3d24099f45ce" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "nama_bank_soal": "Bank Soal Seni Budaya Updated",
    "id_mapel": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
    "jml_soal": 35,
    "deskripsi": "Updated deskripsi"
  }'
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Update bank soal successfully",
  "data": {
    "id": "5112e444-25d8-4ca6-859f-3d24099f45ce",
    "nama_bank_soal": "Bank Soal Seni Budaya Updated",
    "id_mapel": "f71b3390-beeb-4e06-8fd4-deb1778c8210",
    "jml_soal": 35,
    "deskripsi": "Updated deskripsi",
    "created_at": "2026-05-18 10:39:07",
    "updated_at": "2026-05-18 10:40:00"
  },
  "errors": null
}
```

---

### 6️⃣ DELETE - Soft Delete Bank Soal

**Endpoint:** `DELETE /api/bank-soal/:id`

**Authentication:** ✅ Required (JWT Token)

**Note:** Ini adalah **Soft Delete** - data tidak benar-benar dihapus, hanya ditandai dengan `deleted_at` timestamp.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID bank soal yang akan dihapus |

**Request Example:**
```bash
curl -X DELETE "http://localhost:3000/api/bank-soal/5112e444-25d8-4ca6-859f-3d24099f45ce" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Delete bank soal successfully",
  "data": null,
  "errors": null
}
```

---

### 7️⃣ PATCH - Restore Deleted Bank Soal

**Endpoint:** `PATCH /api/bank-soal/:id/restore`

**Authentication:** ✅ Required (JWT Token)

**Description:** Mengembalikan bank soal yang sudah di-soft delete.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID bank soal yang akan di-restore |

**Request Example:**
```bash
curl -X PATCH "http://localhost:3000/api/bank-soal/5112e444-25d8-4ca6-859f-3d24099f45ce/restore" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Restore bank soal successfully",
  "data": null,
  "errors": null
}
```

---

## 🧪 Testing dengan cURL

### Get All Bank Soal
```bash
curl "http://localhost:3000/api/bank-soal?page=1&page_size=5"
```

### Get by Mapel
```bash
curl "http://localhost:3000/api/bank-soal/mapel/{mapel_id}?page=1&page_size=10"
```

### Create Bank Soal (with JWT)
```bash
TOKEN="your_jwt_token_here"
curl -X POST "http://localhost:3000/api/bank-soal" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "nama_bank_soal": "Bank Soal Test",
    "id_mapel": "mapel-uuid",
    "jml_soal": 25,
    "deskripsi": "Test bank soal"
  }'
```

### Update Bank Soal (with JWT)
```bash
TOKEN="your_jwt_token_here"
curl -X PUT "http://localhost:3000/api/bank-soal/{id}" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "nama_bank_soal": "Updated Name",
    "id_mapel": "mapel-uuid",
    "jml_soal": 30,
    "deskripsi": "Updated deskripsi"
  }'
```

### Delete Bank Soal (with JWT)
```bash
TOKEN="your_jwt_token_here"
curl -X DELETE "http://localhost:3000/api/bank-soal/{id}" \
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
| 201 | Create bank soal successfully | Bank soal berhasil dibuat |
| 400 | Invalid request format | Format request tidak sesuai |
| 400 | Resource not found | Bank soal tidak ditemukan |
| 401 | Unauthorized | JWT token tidak valid/expired |
| 500 | Internal server error | Error di server |

---

## 📋 Field Details

### ID (UUID)
- Auto-generated oleh database
- Format: `550e8400-e29b-41d4-a716-446655440000`

### Nama Bank Soal
- Max 255 karakter
- Harus UNIQUE (tidak boleh duplikat)

### ID Mapel
- Foreign key reference ke tabel mapel
- Harus valid mapel ID yang ada di database

### Jumlah Soal
- Integer, default: 0
- Min: 0

### Deskripsi
- Opsional (boleh kosong/tidak diisi)
- Tidak ada batasan panjang

### Timestamps
- Format: `YYYY-MM-DD HH:MM:SS`
- `created_at`: Auto-set saat record dibuat
- `updated_at`: Auto-update saat record diubah
- `deleted_at`: NULL saat aktif, diisi saat soft delete

---

## 🔄 Data Flow Example

```
1. Get all mapel → Lihat list mapel
   GET /api/mapel

2. Create bank soal → Buat bank soal untuk mapel tertentu
   POST /api/bank-soal (with mapel_id)

3. Get by mapel → Filter bank soal per mapel
   GET /api/bank-soal/mapel/{mapel_id}

4. Get detail → Lihat detail bank soal
   GET /api/bank-soal/{id}

5. Update → Edit bank soal
   PUT /api/bank-soal/{id}

6. Delete → Soft delete bank soal
   DELETE /api/bank-soal/{id}

7. Restore → Undo soft delete
   PATCH /api/bank-soal/{id}/restore
```

---

## 💡 Important Notes

1. **Relasi dengan Mapel**: Setiap bank soal harus memiliki id_mapel yang valid
2. **Soft Delete**: Data yang didelete tidak benar-benar dihapus dari database
3. **Pagination**: Gunakan pagination pada GET /api/bank-soal untuk performa optimal
4. **Filter by Mapel**: Endpoint GET /api/bank-soal/mapel/:mapel_id sangat penting untuk filter soal per mata pelajaran
5. **Authorization**: Semua write operations (POST, PUT, DELETE, PATCH) memerlukan JWT token

---

Generated with Claude Code 🤖
