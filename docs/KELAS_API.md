# Dokumentasi API Kelas (Class Management)

## 📌 Base URL
```
http://localhost:3000/api/kelas
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

### 1️⃣ GET - Daftar Semua Kelas (Dengan Pagination & Filter)

**Endpoint:** `GET /api/kelas`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |
| `id_jurusan` | string (UUID) | _(optional)_ | Filter kelas berdasarkan ID jurusan |
| `tingkat` | string | _(optional)_ | Filter kelas berdasarkan tingkat (X, XI, XII) |

**Request Examples:**
```bash
# Ambil semua kelas (default pagination)
curl -X GET "http://localhost:3000/api/kelas"

# Filter berdasarkan tingkat
curl -X GET "http://localhost:3000/api/kelas?tingkat=X"

# Filter berdasarkan jurusan
curl -X GET "http://localhost:3000/api/kelas?id_jurusan=550e8400-e29b-41d4-a716-446655440000"

# Filter berdasarkan tingkat dan jurusan sekaligus
curl -X GET "http://localhost:3000/api/kelas?id_jurusan=550e8400-e29b-41d4-a716-446655440000&tingkat=X"

# Dengan custom pagination
curl -X GET "http://localhost:3000/api/kelas?page=2&page_size=20&tingkat=XI"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get all kelas successfully",
  "data": {
    "data": [
      {
        "id": "661e9511-f30c-52e5-b827-557766551111",
        "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
        "nama_kelas": "X - Teknik Komputer dan Jaringan",
        "tingkat": "X",
        "nama_jurusan": "Teknik Komputer dan Jaringan",
        "created_at": "2026-05-19 10:00:00",
        "updated_at": "2026-05-19 10:00:00"
      },
      {
        "id": "661e9511-f30c-52e5-b827-557766551112",
        "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
        "nama_kelas": "XI - Teknik Komputer dan Jaringan",
        "tingkat": "XI",
        "nama_jurusan": "Teknik Komputer dan Jaringan",
        "created_at": "2026-05-19 10:00:00",
        "updated_at": "2026-05-19 10:00:00"
      }
    ],
    "total": 15,
    "page": 1,
    "page_size": 10,
    "total_page": 2
  },
  "errors": null
}
```

**Error Response (500 Internal Server Error):**
```json
{
  "success": false,
  "message": "Internal server error message",
  "data": null,
  "errors": null
}
```

---

### 2️⃣ GET - Detail Kelas by ID

**Endpoint:** `GET /api/kelas/:id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID kelas |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/kelas/661e9511-f30c-52e5-b827-557766551111"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get kelas successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551111",
    "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X - Teknik Komputer dan Jaringan",
    "tingkat": "X",
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

### 3️⃣ POST - Buat Kelas Baru

**Endpoint:** `POST /api/kelas`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
  "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
  "nama_kelas": "X TKJ 1",
  "tingkat": "X"
}
```

**Field Requirements:**
| Field | Type | Required | Length | Description |
|-------|------|----------|--------|-------------|
| `id_jurusan` | string (UUID) | ✅ Yes | - | UUID dari tabel jurusan |
| `nama_kelas` | string | ✅ Yes | max 255 | Nama kelas (contoh: "X TKJ 1") |
| `tingkat` | string | ✅ Yes | max 10 | Tingkat kelas: "X", "XI", atau "XII" |

**Request Example:**
```bash
curl -X POST "http://localhost:3000/api/kelas" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X TKJ 1",
    "tingkat": "X"
  }'
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Create kelas successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551111",
    "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X TKJ 1",
    "tingkat": "X",
    "nama_jurusan": "Teknik Komputer dan Jaringan",
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

### 4️⃣ PUT - Update Kelas

**Endpoint:** `PUT /api/kelas/:id`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID kelas yang akan diupdate |

**Request Body:**
```json
{
  "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
  "nama_kelas": "X TKJ 2",
  "tingkat": "X"
}
```

**Request Example:**
```bash
curl -X PUT "http://localhost:3000/api/kelas/661e9511-f30c-52e5-b827-557766551111" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X TKJ 2",
    "tingkat": "X"
  }'
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Update kelas successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551111",
    "id_jurusan": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X TKJ 2",
    "tingkat": "X",
    "nama_jurusan": "Teknik Komputer dan Jaringan",
    "created_at": "2026-05-19 10:15:30",
    "updated_at": "2026-05-19 10:20:00"
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

### 5️⃣ DELETE - Soft Delete Kelas

**Endpoint:** `DELETE /api/kelas/:id`

**Authentication:** ✅ Required (JWT Token)

**Note:** Ini adalah **Soft Delete** - data tidak benar-benar dihapus, hanya ditandai dengan `deleted_at` timestamp. Data yang sudah dihapus tidak akan muncul di endpoint GET.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID kelas yang akan dihapus |

**Request Example:**
```bash
curl -X DELETE "http://localhost:3000/api/kelas/661e9511-f30c-52e5-b827-557766551111" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Delete kelas successfully",
  "data": null,
  "errors": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "data not found",
  "data": null,
  "errors": null
}
```

---

### 6️⃣ PATCH - Restore Deleted Kelas

**Endpoint:** `PATCH /api/kelas/:id/restore`

**Authentication:** ✅ Required (JWT Token)

**Description:** Mengembalikan kelas yang sudah di-soft delete (menghapus timestamp `deleted_at`).

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID kelas yang akan di-restore |

**Request Example:**
```bash
curl -X PATCH "http://localhost:3000/api/kelas/661e9511-f30c-52e5-b827-557766551111/restore" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Restore kelas successfully",
  "data": null,
  "errors": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Error message",
  "data": null,
  "errors": null
}
```

---

## 🎯 Use Cases & Tips

### 1. Dropdown Kelas Berdasarkan Tingkat
Untuk membuat dropdown kelas tingkat X di frontend:
```bash
curl -X GET "http://localhost:3000/api/kelas?tingkat=X&page=1&page_size=100"
```

### 2. Dropdown Kelas Berdasarkan Jurusan
Untuk membuat dropdown kelas berdasarkan jurusan yang dipilih:
```bash
curl -X GET "http://localhost:3000/api/kelas?id_jurusan=550e8400-e29b-41d4-a716-446655440000&page=1&page_size=100"
```

### 3. Dropdown Kelas Berdasarkan Jurusan & Tingkat
Kombinasi kedua filter untuk hasil yang spesifik:
```bash
curl -X GET "http://localhost:3000/api/kelas?id_jurusan=550e8400-e29b-41d4-a716-446655440000&tingkat=XI&page=1&page_size=100"
```

### 4. List Semua Kelas dengan Pagination
```bash
curl -X GET "http://localhost:3000/api/kelas?page=1&page_size=20"
```

---

## 🧪 Testing dengan cURL (Complete Workflow)

### Step 1: Register & Get Token
```bash
curl -X POST "http://localhost:3000/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@example.com",
    "password": "password123"
  }'
```

Response akan berisi JWT token:
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

Simpan token dari response di atas sebagai `TOKEN`.

### Step 2: Get Jurusan ID
Kita perlu ID jurusan untuk membuat kelas:
```bash
curl -X GET "http://localhost:3000/api/jurusan?page=1&page_size=100"
```

Simpan salah satu ID jurusan sebagai `JURUSAN_ID`.

### Step 3: Create Kelas
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
JURUSAN_ID="550e8400-e29b-41d4-a716-446655440000"

curl -X POST "http://localhost:3000/api/kelas" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"id_jurusan\": \"$JURUSAN_ID\",
    \"nama_kelas\": \"X TKJ 1\",
    \"tingkat\": \"X\"
  }"
```

Response akan berisi ID kelas baru. Simpan sebagai `KELAS_ID`.

### Step 4: Get All Kelas
```bash
curl -X GET "http://localhost:3000/api/kelas"
```

### Step 5: Get Detail Kelas
```bash
KELAS_ID="661e9511-f30c-52e5-b827-557766551111"
curl -X GET "http://localhost:3000/api/kelas/$KELAS_ID"
```

### Step 6: Update Kelas
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
KELAS_ID="661e9511-f30c-52e5-b827-557766551111"
JURUSAN_ID="550e8400-e29b-41d4-a716-446655440000"

curl -X PUT "http://localhost:3000/api/kelas/$KELAS_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"id_jurusan\": \"$JURUSAN_ID\",
    \"nama_kelas\": \"X TKJ 2\",
    \"tingkat\": \"X\"
  }"
```

### Step 7: Filter by Tingkat
```bash
curl -X GET "http://localhost:3000/api/kelas?tingkat=X"
```

### Step 8: Filter by Jurusan & Tingkat
```bash
JURUSAN_ID="550e8400-e29b-41d4-a716-446655440000"
curl -X GET "http://localhost:3000/api/kelas?id_jurusan=$JURUSAN_ID&tingkat=XI"
```

### Step 9: Delete Kelas
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
KELAS_ID="661e9511-f30c-52e5-b827-557766551111"

curl -X DELETE "http://localhost:3000/api/kelas/$KELAS_ID" \
  -H "Authorization: Bearer $TOKEN"
```

### Step 10: Restore Kelas
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
KELAS_ID="661e9511-f30c-52e5-b827-557766551111"

curl -X PATCH "http://localhost:3000/api/kelas/$KELAS_ID/restore" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📋 Response Status Codes

| Status Code | Meaning | Description |
|---|---|---|
| 200 | OK | Request berhasil, data diterima |
| 201 | Created | Resource baru berhasil dibuat |
| 400 | Bad Request | Format request tidak valid |
| 404 | Not Found | Resource tidak ditemukan |
| 500 | Internal Server Error | Error di server |

---

## ⚠️ Common Errors

### Invalid UUID Format
```json
{
  "success": false,
  "message": "Invalid request format",
  "errors": null
}
```
**Solution:** Pastikan `id_jurusan` dan `id` adalah UUID yang valid.

### Missing Required Fields
```json
{
  "success": false,
  "message": "Invalid request format",
  "errors": null
}
```
**Solution:** Pastikan semua required fields (`id_jurusan`, `nama_kelas`, `tingkat`) ada di request body.

### Invalid Tingkat Value
Request akan diterima, namun tingkat harus salah satu dari: `"X"`, `"XI"`, atau `"XII"`.

### Unauthorized (Missing/Invalid Token)
```json
{
  "success": false,
  "message": "Unauthorized",
  "errors": null
}
```
**Solution:** Pastikan JWT token valid dan dikirim di header `Authorization: Bearer <TOKEN>`.

---

## 🔗 Related Endpoints

- **Jurusan API:** `/api/jurusan` - Lihat [JURUSAN_API.md](./JURUSAN_API.md)
- **Mapel API:** `/api/mapel` - Lihat [MAPEL_API.md](./MAPEL_API.md)

---

## 📝 Catatan Penting

1. **Nama Jurusan (JOIN):** Field `nama_jurusan` dalam response diambil dari tabel `jurusan` melalui LEFT JOIN dengan `id_jurusan`. Field ini hanya untuk display di frontend dan tidak bisa diupdate melalui API kelas.

2. **Soft Delete:** Ketika kelas dihapus, data tetap ada di database dengan kolom `deleted_at` terisi. Data tidak muncul di hasil GET.

3. **Filter Opsional:** Filter `id_jurusan` dan `tingkat` bersifat opsional. Jika tidak dikirim, endpoint mengembalikan semua kelas.

4. **Pagination:** Default page adalah 1, default page_size adalah 10. Maksimal page_size adalah 100 untuk performa optimal.

5. **Authentication Required:** Hanya endpoint GET yang public. Semua operasi modifikasi (POST, PUT, DELETE, PATCH) memerlukan JWT token.

6. **Tingkat Values:** Nilai valid untuk `tingkat` adalah:
   - `"X"` - Kelas 10 (Grade 10)
   - `"XI"` - Kelas 11 (Grade 11)
   - `"XII"` - Kelas 12 (Grade 12)

7. **Database Query:** Semua endpoint GET menggunakan LEFT JOIN dengan tabel `jurusan`. Jika `id_jurusan` tidak valid atau data jurusan dihapus, field `nama_jurusan` akan null.
