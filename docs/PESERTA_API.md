# Dokumentasi API Peserta (Student Management)

## 📌 Base URL
```
http://localhost:3000/api/peserta
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

### 1️⃣ GET - Daftar Semua Peserta (Dengan Pagination & Filter)

**Endpoint:** `GET /api/peserta`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Nomor halaman |
| `page_size` | integer | 10 | Jumlah data per halaman |
| `id_kelas` | string (UUID) | _(optional)_ | Filter peserta berdasarkan ID kelas |

**Request Examples:**
```bash
# Ambil semua peserta (default pagination)
curl -X GET "http://localhost:3000/api/peserta"

# Filter berdasarkan kelas
curl -X GET "http://localhost:3000/api/peserta?id_kelas=550e8400-e29b-41d4-a716-446655440000"

# Dengan custom pagination
curl -X GET "http://localhost:3000/api/peserta?page=2&page_size=20"

# Kombinasi filter dan pagination
curl -X GET "http://localhost:3000/api/peserta?id_kelas=550e8400-e29b-41d4-a716-446655440000&page=1&page_size=50"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get all peserta successfully",
  "data": {
    "data": [
      {
        "id": "661e9511-f30c-52e5-b827-557766551111",
        "nama": "Peserta 1 - X - Teknik Komputer dan Jaringan",
        "id_kelas": "550e8400-e29b-41d4-a716-446655440000",
        "nama_kelas": "X - Teknik Komputer dan Jaringan",
        "username": "peserta_550e8400_1",
        "created_at": "2026-05-20 10:00:00",
        "updated_at": "2026-05-20 10:00:00"
      },
      {
        "id": "661e9511-f30c-52e5-b827-557766551112",
        "nama": "Peserta 2 - X - Teknik Komputer dan Jaringan",
        "id_kelas": "550e8400-e29b-41d4-a716-446655440000",
        "nama_kelas": "X - Teknik Komputer dan Jaringan",
        "username": "peserta_550e8400_2",
        "created_at": "2026-05-20 10:00:00",
        "updated_at": "2026-05-20 10:00:00"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 10,
    "total_page": 3
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

### 2️⃣ GET - Detail Peserta by ID

**Endpoint:** `GET /api/peserta/:id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID peserta |

**Request Example:**
```bash
curl -X GET "http://localhost:3000/api/peserta/661e9511-f30c-52e5-b827-557766551111"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Get peserta successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551111",
    "nama": "Peserta 1 - X - Teknik Komputer dan Jaringan",
    "id_kelas": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X - Teknik Komputer dan Jaringan",
    "username": "peserta_550e8400_1",
    "created_at": "2026-05-20 10:00:00",
    "updated_at": "2026-05-20 10:00:00"
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

### 3️⃣ POST - Buat Peserta Baru

**Endpoint:** `POST /api/peserta`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
  "nama": "Andi Pratama",
  "id_kelas": "550e8400-e29b-41d4-a716-446655440000",
  "username": "andi_pratama",
  "password": "password123"
}
```

**Field Requirements:**
| Field | Type | Required | Length | Description |
|-------|------|----------|--------|-------------|
| `nama` | string | ✅ Yes | max 255 | Nama lengkap peserta |
| `id_kelas` | string (UUID) | ✅ Yes | - | UUID dari tabel kelas |
| `username` | string | ✅ Yes | max 100 | Username unik untuk login |
| `password` | string | ✅ Yes | min 6 | Password (akan di-hash bcrypt) |

**Request Example:**
```bash
curl -X POST "http://localhost:3000/api/peserta" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "nama": "Andi Pratama",
    "id_kelas": "550e8400-e29b-41d4-a716-446655440000",
    "username": "andi_pratama",
    "password": "password123"
  }'
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Create peserta successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551113",
    "nama": "Andi Pratama",
    "id_kelas": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X - Teknik Komputer dan Jaringan",
    "username": "andi_pratama",
    "created_at": "2026-05-20 14:30:00",
    "updated_at": "2026-05-20 14:30:00"
  },
  "errors": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "username sudah digunakan",
  "data": null,
  "errors": null
}
```

---

### 4️⃣ PUT - Update Peserta

**Endpoint:** `PUT /api/peserta/:id`

**Authentication:** ✅ Required (JWT Token)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID peserta yang akan diupdate |

**Request Body:**
```json
{
  "nama": "Andi Pratama Updated",
  "id_kelas": "550e8400-e29b-41d4-a716-446655440000",
  "username": "andi_pratama_new",
  "password": ""
}
```

**Field Notes:**
- Semua field wajib ada di request body
- Field `password` bersifat opsional — jika kosong atau tidak dikirim, password lama dipertahankan
- Jika ingin mengubah password, kirim password baru dalam plaintext (akan di-hash otomatis)

**Request Example:**
```bash
curl -X PUT "http://localhost:3000/api/peserta/661e9511-f30c-52e5-b827-557766551113" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "nama": "Andi Pratama Updated",
    "id_kelas": "550e8400-e29b-41d4-a716-446655440000",
    "username": "andi_pratama_new",
    "password": ""
  }'
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Update peserta successfully",
  "data": {
    "id": "661e9511-f30c-52e5-b827-557766551113",
    "nama": "Andi Pratama Updated",
    "id_kelas": "550e8400-e29b-41d4-a716-446655440000",
    "nama_kelas": "X - Teknik Komputer dan Jaringan",
    "username": "andi_pratama_new",
    "created_at": "2026-05-20 14:30:00",
    "updated_at": "2026-05-20 14:35:00"
  },
  "errors": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "username sudah digunakan",
  "data": null,
  "errors": null
}
```

---

### 5️⃣ DELETE - Soft Delete Peserta

**Endpoint:** `DELETE /api/peserta/:id`

**Authentication:** ✅ Required (JWT Token)

**Note:** Ini adalah **Soft Delete** - data tidak benar-benar dihapus, hanya ditandai dengan `deleted_at` timestamp. Data yang sudah dihapus tidak akan muncul di endpoint GET.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID peserta yang akan dihapus |

**Request Example:**
```bash
curl -X DELETE "http://localhost:3000/api/peserta/661e9511-f30c-52e5-b827-557766551113" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Delete peserta successfully",
  "data": null,
  "errors": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Resource not found",
  "data": null,
  "errors": null
}
```

---

### 6️⃣ PATCH - Restore Deleted Peserta

**Endpoint:** `PATCH /api/peserta/:id/restore`

**Authentication:** ✅ Required (JWT Token)

**Description:** Mengembalikan peserta yang sudah di-soft delete (menghapus timestamp `deleted_at`).

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | ID peserta yang akan di-restore |

**Request Example:**
```bash
curl -X PATCH "http://localhost:3000/api/peserta/661e9511-f30c-52e5-b827-557766551113/restore" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Restore peserta successfully",
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

### 1. Ambil Peserta dalam Satu Kelas
```bash
KELAS_ID="550e8400-e29b-41d4-a716-446655440000"
curl -X GET "http://localhost:3000/api/peserta?id_kelas=$KELAS_ID&page=1&page_size=100"
```

### 2. Pagination untuk Daftar Peserta Besar
```bash
# Halaman pertama dengan 20 peserta per halaman
curl -X GET "http://localhost:3000/api/peserta?page=1&page_size=20"

# Halaman kedua
curl -X GET "http://localhost:3000/api/peserta?page=2&page_size=20"
```

### 3. Filter Peserta dalam Kelas Spesifik dengan Pagination
```bash
KELAS_ID="550e8400-e29b-41d4-a716-446655440000"
curl -X GET "http://localhost:3000/api/peserta?id_kelas=$KELAS_ID&page=1&page_size=50"
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

Simpan token sebagai `TOKEN`.

### Step 2: Get Kelas ID
```bash
curl -X GET "http://localhost:3000/api/kelas?page=1&page_size=100"
```

Simpan salah satu ID kelas sebagai `KELAS_ID`.

### Step 3: Create Peserta
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
KELAS_ID="550e8400-e29b-41d4-a716-446655440000"

curl -X POST "http://localhost:3000/api/peserta" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"nama\": \"Andi Pratama\",
    \"id_kelas\": \"$KELAS_ID\",
    \"username\": \"andi_pratama\",
    \"password\": \"password123\"
  }"
```

Response akan berisi ID peserta baru. Simpan sebagai `PESERTA_ID`.

### Step 4: Get All Peserta
```bash
curl -X GET "http://localhost:3000/api/peserta"
```

### Step 5: Get Detail Peserta
```bash
PESERTA_ID="661e9511-f30c-52e5-b827-557766551113"
curl -X GET "http://localhost:3000/api/peserta/$PESERTA_ID"
```

### Step 6: Update Peserta
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
PESERTA_ID="661e9511-f30c-52e5-b827-557766551113"
KELAS_ID="550e8400-e29b-41d4-a716-446655440000"

curl -X PUT "http://localhost:3000/api/peserta/$PESERTA_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"nama\": \"Andi Pratama Updated\",
    \"id_kelas\": \"$KELAS_ID\",
    \"username\": \"andi_pratama_updated\",
    \"password\": \"\"
  }"
```

### Step 7: Update Password Peserta
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
PESERTA_ID="661e9511-f30c-52e5-b827-557766551113"
KELAS_ID="550e8400-e29b-41d4-a716-446655440000"

curl -X PUT "http://localhost:3000/api/peserta/$PESERTA_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"nama\": \"Andi Pratama\",
    \"id_kelas\": \"$KELAS_ID\",
    \"username\": \"andi_pratama\",
    \"password\": \"newpassword123\"
  }"
```

### Step 8: Filter Peserta by Kelas
```bash
KELAS_ID="550e8400-e29b-41d4-a716-446655440000"
curl -X GET "http://localhost:3000/api/peserta?id_kelas=$KELAS_ID"
```

### Step 9: Delete Peserta
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
PESERTA_ID="661e9511-f30c-52e5-b827-557766551113"

curl -X DELETE "http://localhost:3000/api/peserta/$PESERTA_ID" \
  -H "Authorization: Bearer $TOKEN"
```

### Step 10: Restore Peserta
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
PESERTA_ID="661e9511-f30c-52e5-b827-557766551113"

curl -X PATCH "http://localhost:3000/api/peserta/$PESERTA_ID/restore" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📋 Response Status Codes

| Status Code | Meaning | Description |
|---|---|---|
| 200 | OK | Request berhasil, data diterima |
| 201 | Created | Resource baru berhasil dibuat |
| 400 | Bad Request | Format request tidak valid atau error bisnis |
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
**Solution:** Pastikan `id_kelas` dan `id` adalah UUID yang valid.

### Username Sudah Digunakan
```json
{
  "success": false,
  "message": "username sudah digunakan",
  "errors": null
}
```
**Solution:** Gunakan username yang belum pernah digunakan oleh peserta aktif lainnya.

### Missing Required Fields
```json
{
  "success": false,
  "message": "Invalid request format",
  "errors": null
}
```
**Solution:** Pastikan semua required fields (`nama`, `id_kelas`, `username`, `password`) ada di request body.

### Password Terlalu Pendek
```json
{
  "success": false,
  "message": "Invalid request format",
  "errors": null
}
```
**Solution:** Password minimal 6 karakter.

### Unauthorized (Missing/Invalid Token)
```json
{
  "success": false,
  "message": "Unauthorized",
  "errors": null
}
```
**Solution:** Pastikan JWT token valid dan dikirim di header `Authorization: Bearer <TOKEN>`.

### Peserta Not Found
```json
{
  "success": false,
  "message": "Resource not found",
  "errors": null
}
```
**Solution:** Pastikan ID peserta benar dan peserta belum di-soft delete.

---

## 🔗 Related Endpoints

- **Kelas API:** `/api/kelas` - Lihat [KELAS_API.md](./KELAS_API.md)
- **Jurusan API:** `/api/jurusan` - Lihat [JURUSAN_API.md](./JURUSAN_API.md)
- **Auth API:** `/api/auth` - Untuk mendapatkan JWT token

---

## 📝 Catatan Penting

1. **Password tidak dikembalikan:** Field `password` tidak akan pernah muncul di response API. Password hanya disimpan dalam bentuk hash bcrypt di database.

2. **Soft Delete:** Ketika peserta dihapus, data tetap ada di database dengan kolom `deleted_at` terisi. Data tidak muncul di hasil GET kecuali sudah di-restore.

3. **Filter Opsional:** Filter `id_kelas` bersifat opsional. Jika tidak dikirim, endpoint mengembalikan semua peserta.

4. **Pagination:** Default page adalah 1, default page_size adalah 10.

5. **Authentication Required:** Hanya endpoint GET yang public. Semua operasi modifikasi (POST, PUT, DELETE, PATCH) memerlukan JWT token.

6. **Username Uniqueness:** Username harus unik di antara peserta yang aktif (tidak di-soft delete). Peserta yang sudah dihapus bisa memiliki username yang sama dengan peserta baru karena menggunakan partial unique index.

7. **Password Update:** Saat update, jika field `password` kosong atau tidak dikirim, password lama dipertahankan. Ini memungkinkan update data peserta lain tanpa perlu mengubah password.

8. **Relasi Kelas:** Setiap peserta wajib terikat ke satu kelas. Response API selalu menyertakan `nama_kelas` hasil JOIN dari tabel `kelas` untuk kemudahan display di frontend.

9. **Database Query:** Semua endpoint GET menggunakan LEFT JOIN dengan tabel `kelas`. Jika `id_kelas` tidak valid atau data kelas dihapus, field `nama_kelas` akan null.

10. **Seeder Default:** Saat aplikasi pertama kali dijalankan, seeder otomatis membuat 5 peserta per kelas dengan password default `password123`. Username dibuat dengan format `peserta_<first_8_chars_of_kelas_id>_<nomor>`.
