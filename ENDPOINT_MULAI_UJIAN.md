# Dokumentasi API: Endpoint Mulai Ujian

## 📋 Ringkasan

Endpoint ini digunakan untuk **memulai atau melanjutkan ujian**. Saat peserta pertama kali mulai ujian:
1. Endpoint otomatis membuat record **nilai** peserta
2. Backend otomatis membuat record **jawaban kosong** untuk setiap soal di bank soal jadwal, dengan urutan **acak (random)** — mencegah contek antar peserta
3. Semua operasi berjalan atomik dalam satu transaction

Peserta dapat melanjutkan ujian yang belum diselesaikan tanpa perlu re-generate jawaban.

---

## 📌 Informasi Endpoint

### URL
```
POST /api/nilai/mulai-ujian/:id_jadwal
```

### Method
`POST`

### Authentication
**Wajib** — JWT Bearer Token (dari login peserta)

### Content-Type
`application/json` (tidak ada body, tapi header harus ada)

---

## 🔐 Autentikasi

### Header yang Diperlukan
```http
Authorization: Bearer <JWT_TOKEN>
```

**Catatan:**
- Token diperoleh saat login peserta
- Token harus valid (tidak expired)
- Jika tidak ada token atau token invalid → HTTP 401

---

## 📤 Request

### URL Parameter

| Parameter | Tipe | Wajib | Keterangan |
|-----------|------|-------|-----------|
| `id_jadwal` | string (UUID) | ✅ | UUID jadwal ujian yang akan dimulai |

### Contoh Request

**cURL:**
```bash
curl -X POST http://localhost:3000/api/nilai/mulai-ujian/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

**JavaScript (Fetch):**
```javascript
const jadwalId = "550e8400-e29b-41d4-a716-446655440000";
const token = localStorage.getItem("token"); // ambil token dari localStorage

fetch(`http://localhost:3000/api/nilai/mulai-ujian/${jadwalId}`, {
  method: "POST",
  headers: {
    "Authorization": `Bearer ${token}`,
    "Content-Type": "application/json"
  }
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error("Error:", error));
```

**Axios (TypeScript/JavaScript):**
```typescript
import axios from "axios";

const startExam = async (jadwalId: string, token: string) => {
  try {
    const response = await axios.post(
      `http://localhost:3000/api/nilai/mulai-ujian/${jadwalId}`,
      {}, // body kosong
      {
        headers: {
          "Authorization": `Bearer ${token}`,
          "Content-Type": "application/json"
        }
      }
    );
    return response.data;
  } catch (error) {
    console.error("Error starting exam:", error);
    throw error;
  }
};
```

---

## 📥 Response

### Skenario 1: Mulai Ujian Pertama Kali (Belum Pernah Mulai)

**HTTP Status:** `201 Created`

```json
{
  "success": true,
  "message": "Mulai ujian successfully",
  "data": {
    "id": "12345678-1234-1234-1234-123456789012",
    "id_peserta": "87654321-4321-4321-4321-210987654321",
    "nama_peserta": "Budi Santoso",
    "id_jadwal": "550e8400-e29b-41d4-a716-446655440000",
    "nama_ujian": "UTS Matematika Kelas X",
    "nilai": 0,
    "wkt_mulai": "2026-05-25T08:15:30Z",
    "aktivitas_terakhir": "2026-05-25T08:15:30Z",
    "wkt_selesai": null,
    "created_at": "2026-05-25T08:15:30Z",
    "updated_at": "2026-05-25T08:15:30Z"
  }
}
```

**Penjelasan Field:**
| Field | Keterangan |
|-------|-----------|
| `id` | ID unik record nilai yang baru dibuat |
| `id_peserta` | UUID peserta (dari JWT) |
| `nama_peserta` | Nama lengkap peserta |
| `id_jadwal` | UUID jadwal ujian |
| `nama_ujian` | Nama ujian (misal: "UTS Matematika") |
| `nilai` | Nilai saat ini (selalu 0 di awal) |
| `wkt_mulai` | Waktu peserta mulai ujian (sekarang) |
| `aktivitas_terakhir` | Waktu aktivitas terakhir (sama dengan wkt_mulai saat mulai) |
| `wkt_selesai` | Waktu selesai ujian (null = belum selesai) |
| `created_at` | Timestamp record dibuat |
| `updated_at` | Timestamp record terakhir update |

---

### Skenario 2: Lanjutkan Ujian (Sudah Mulai, Belum Selesai)

Jika peserta menjalankan endpoint yang sama **kedua kalinya** sebelum menyelesaikan ujian:

**HTTP Status:** `200 OK`

```json
{
  "success": true,
  "message": "Lanjutkan ujian successfully",
  "data": {
    "id": "12345678-1234-1234-1234-123456789012",
    "id_peserta": "87654321-4321-4321-4321-210987654321",
    "nama_peserta": "Budi Santoso",
    "id_jadwal": "550e8400-e29b-41d4-a716-446655440000",
    "nama_ujian": "UTS Matematika Kelas X",
    "nilai": 0,
    "wkt_mulai": "2026-05-25T08:15:30Z",
    "aktivitas_terakhir": "2026-05-25T08:15:30Z",
    "wkt_selesai": null,
    "created_at": "2026-05-25T08:15:30Z",
    "updated_at": "2026-05-25T08:15:30Z"
  }
}
```

**Catatan:**
- Message berbeda: `"Lanjutkan ujian successfully"` (vs `"Mulai ujian successfully"`)
- HTTP status 200 (vs 201)
- Data sama seperti sebelumnya (TIDAK ada update)
- Frontend dapat membedakan scenario dengan mengecek `message` atau HTTP status code

---

### Skenario 3: Ujian Sudah Selesai (Ditolak)

Jika peserta sudah menyelesaikan ujian sebelumnya (field `wkt_selesai` sudah terisi):

**HTTP Status:** `400 Bad Request`

```json
{
  "success": false,
  "message": "Ujian sudah pernah dilakukan",
  "data": null
}
```

**Penjelasan:**
- Peserta tidak boleh memulai ulang ujian yang sudah pernah diselesaikan
- Frontend harus menampilkan pesan error kepada peserta
- Jika ingin reset ujian, perlu admin yang melakukan via endpoint lain

---

### Skenario 4: Error - Tidak Ada Authorization

**HTTP Status:** `401 Unauthorized`

```json
{
  "success": false,
  "message": "Unauthorized",
  "data": null
}
```

**Penyebab:**
- Header `Authorization` tidak ada
- Token sudah expired
- Token invalid/corrupt

---

### Skenario 5: Error - ID Jadwal Invalid

**HTTP Status:** `400 Bad Request`

```json
{
  "success": false,
  "message": "invalid input syntax for type uuid: \"bukan-uuid\"",
  "data": null
}
```

**Penyebab:**
- `id_jadwal` parameter bukan format UUID yang valid
- Contoh UUID yang valid: `550e8400-e29b-41d4-a716-446655440000`

---

## 📝 Side Effect: Jawaban Otomatis Di-generate

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

**⚠️ Catatan penting untuk frontend:**
- `jawaban` & `is_benar` sekarang **bisa `null`**. Handle null check di UI:
  ```jsx
  <span>{jawaban.jawaban ?? "Belum dijawab"}</span>
  <Icon color={jawaban.is_benar === null ? "gray" : (jawaban.is_benar ? "green" : "red")} />
  ```
- **Setiap peserta dapat urutan soal BERBEDA** (random per peserta) → anti-contek
- Endpoint `GET /api/jawaban/nilai/:id_nilai` otomatis `ORDER BY soal.no_soal ASC`. Jika frontend ingin display urutan random sesuai design, perlu re-sort by `no_urut` di client-side.

---

## 🔄 Flow Diagram

```
┌─────────────────────────────────┐
│  Peserta Klik "Mulai Ujian"    │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────┐
│ POST /api/nilai/mulai-ujian/:id_jadwal      │
│ Header: Authorization: Bearer <TOKEN>       │
└────────────┬────────────────────────────────┘
             │
             ▼
      ┌──────────────┐
      │ Server Check │
      └──────┬───────┘
             │
    ┌────────┴──────────┐
    │                   │
    ▼                   ▼
[Belum Mulai]     [Sudah Ada]
    │                   │
    ▼                   ▼
[INSERT baru]    [Cek wkt_selesai]
    │                   │
    │             ┌─────┴──────┐
    │             │            │
    ▼             ▼            ▼
  201 OK      [NULL]      [NOT NULL]
  Mulai       200 OK         400 Error
  Ujian       Lanjutkan      Sudah Selesai
    │             │            │
    └─────┬───────┘            │
          │                    │
          ▼                    ▼
    ┌──────────────┐    ┌─────────────┐
    │  Show Timer  │    │ Show Alert  │
    │  Show Soal   │    │ Disable UI  │
    └──────────────┘    └─────────────┘
```

---

## 💻 Contoh Implementasi Frontend

### React dengan TypeScript

```typescript
import { useState, useEffect } from "react";
import axios from "axios";

interface NilaiData {
  id: string;
  id_peserta: string;
  nama_peserta: string;
  id_jadwal: string;
  nama_ujian: string;
  nilai: number;
  wkt_mulai: string;
  aktivitas_terakhir: string;
  wkt_selesai: string | null;
  created_at: string;
  updated_at: string;
}

interface ExamStartResponse {
  success: boolean;
  message: string;
  data: NilaiData | null;
}

export const ExamStartPage: React.FC<{ jadwalId: string }> = ({ jadwalId }) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [nilaiData, setNilaiData] = useState<NilaiData | null>(null);
  const [isNew, setIsNew] = useState(false);

  const handleStartExam = async () => {
    setLoading(true);
    setError(null);

    try {
      const token = localStorage.getItem("token");
      if (!token) {
        throw new Error("Token not found. Please login first.");
      }

      const response = await axios.post<ExamStartResponse>(
        `http://localhost:3000/api/nilai/mulai-ujian/${jadwalId}`,
        {},
        {
          headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json"
          }
        }
      );

      if (response.status === 201) {
        // Mulai ujian baru
        setIsNew(true);
        setNilaiData(response.data.data);
      } else if (response.status === 200) {
        // Lanjutkan ujian
        setIsNew(false);
        setNilaiData(response.data.data);
      }
    } catch (err: any) {
      if (err.response?.status === 400) {
        setError(err.response.data.message || "Gagal memulai ujian");
      } else if (err.response?.status === 401) {
        setError("Sesi expired. Silakan login ulang.");
        // Redirect ke login page
        window.location.href = "/login";
      } else {
        setError("Terjadi kesalahan. Coba lagi nanti.");
      }
    } finally {
      setLoading(false);
    }
  };

  if (error) {
    return (
      <div className="error-container">
        <p className="error-message">{error}</p>
        <button onClick={() => window.history.back()}>Kembali</button>
      </div>
    );
  }

  if (nilaiData) {
    return (
      <div className="exam-container">
        <div className="exam-header">
          <h1>{nilaiData.nama_ujian}</h1>
          <p>Peserta: {nilaiData.nama_peserta}</p>
          <p className="status">
            {isNew ? "✨ Mulai Ujian Baru" : "📖 Lanjutkan Ujian"}
          </p>
        </div>

        <div className="exam-info">
          <div>
            <label>Mulai:</label>
            <span>{new Date(nilaiData.wkt_mulai).toLocaleString("id-ID")}</span>
          </div>
          <div>
            <label>Status:</label>
            <span>
              {nilaiData.wkt_selesai ? "Selesai" : "Sedang Berlangsung"}
            </span>
          </div>
        </div>

        <div className="exam-content">
          {/* Render soal-soal di sini */}
          <ExamQuestions nilaiId={nilaiData.id} />
        </div>
      </div>
    );
  }

  return (
    <div className="exam-start-page">
      <button
        onClick={handleStartExam}
        disabled={loading}
        className="start-button"
      >
        {loading ? "Loading..." : "Mulai Ujian"}
      </button>
    </div>
  );
};
```

### Vue.js dengan TypeScript

```vue
<template>
  <div class="exam-start-container">
    <!-- Error Alert -->
    <div v-if="error" class="alert alert-danger">
      <p>{{ error }}</p>
      <button @click="error = null">Tutup</button>
    </div>

    <!-- Exam Info (setelah mulai) -->
    <div v-if="nilaiData" class="exam-info-card">
      <h2>{{ nilaiData.nama_ujian }}</h2>
      <p><strong>Peserta:</strong> {{ nilaiData.nama_peserta }}</p>
      <p v-if="isNew" class="badge-new">✨ Ujian Baru Dimulai</p>
      <p v-else class="badge-resume">📖 Lanjutkan Ujian</p>

      <div class="exam-details">
        <div>
          <strong>Mulai:</strong>
          {{ formatDateTime(nilaiData.wkt_mulai) }}
        </div>
        <div>
          <strong>Waktu Aktivitas Terakhir:</strong>
          {{ formatDateTime(nilaiData.aktivitas_terakhir) }}
        </div>
        <div v-if="nilaiData.wkt_selesai">
          <strong>Selesai:</strong>
          {{ formatDateTime(nilaiData.wkt_selesai) }}
        </div>
        <div v-else class="status-ongoing">
          <strong>Status:</strong> Sedang Berlangsung
        </div>
      </div>

      <!-- Komponen ujian -->
      <ExamQuestions :nilaiId="nilaiData.id" />
    </div>

    <!-- Button Mulai Ujian -->
    <div v-else class="exam-start-section">
      <button
        @click="handleStartExam"
        :disabled="loading"
        class="btn btn-primary btn-lg"
      >
        <span v-if="loading" class="spinner"></span>
        {{ loading ? "Sedang Loading..." : "Mulai Ujian Sekarang" }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import axios from "axios";

interface NilaiData {
  id: string;
  id_peserta: string;
  nama_peserta: string;
  id_jadwal: string;
  nama_ujian: string;
  nilai: number;
  wkt_mulai: string;
  aktivitas_terakhir: string;
  wkt_selesai: string | null;
  created_at: string;
  updated_at: string;
}

const props = defineProps<{
  jadwalId: string;
}>();

const loading = ref(false);
const error = ref<string | null>(null);
const nilaiData = ref<NilaiData | null>(null);
const isNew = ref(false);

const handleStartExam = async () => {
  loading.value = true;
  error.value = null;

  try {
    const token = localStorage.getItem("token");
    if (!token) {
      throw new Error("Token tidak ditemukan. Silakan login terlebih dahulu.");
    }

    const response = await axios.post(
      `http://localhost:3000/api/nilai/mulai-ujian/${props.jadwalId}`,
      {},
      {
        headers: {
          "Authorization": `Bearer ${token}`,
          "Content-Type": "application/json"
        }
      }
    );

    // Tentukan apakah ini ujian baru atau resume
    isNew.value = response.status === 201;
    nilaiData.value = response.data.data;

    // Simpan nilaiId ke sessionStorage untuk komponen selanjutnya
    if (nilaiData.value) {
      sessionStorage.setItem("nilaiId", nilaiData.value.id);
    }
  } catch (err: any) {
    if (err.response?.status === 400) {
      error.value =
        err.response.data.message || "Gagal memulai ujian. Coba lagi.";
    } else if (err.response?.status === 401) {
      error.value = "Sesi Anda telah berakhir. Silakan login kembali.";
      setTimeout(() => {
        window.location.href = "/login";
      }, 2000);
    } else {
      error.value =
        "Terjadi kesalahan jaringan. Periksa koneksi dan coba lagi.";
    }
  } finally {
    loading.value = false;
  }
};

const formatDateTime = (isoString: string): string => {
  return new Date(isoString).toLocaleString("id-ID", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit"
  });
};
</script>

<style scoped>
.exam-start-container {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.exam-info-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 2rem;
  background-color: #f9f9f9;
}

.badge-new {
  display: inline-block;
  background-color: #4caf50;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  font-size: 0.9rem;
  margin-top: 1rem;
}

.badge-resume {
  display: inline-block;
  background-color: #2196f3;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  font-size: 0.9rem;
  margin-top: 1rem;
}

.exam-details {
  margin-top: 1.5rem;
  padding: 1rem;
  background-color: white;
  border-radius: 4px;
}

.exam-details > div {
  margin-bottom: 0.5rem;
}

.status-ongoing {
  color: #ff9800;
  font-weight: bold;
}

.btn-primary {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 1rem 2rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
}

.btn-primary:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.alert-danger {
  background-color: #f8d7da;
  color: #721c24;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
}
</style>
```

---

## ⏱️ Catatan Waktu

### Format Waktu
Semua field waktu menggunakan format **ISO 8601** dengan timezone UTC (Z):
```
2026-05-25T08:15:30Z
```

### Parsing di Frontend

**JavaScript:**
```javascript
const wktMulai = new Date("2026-05-25T08:15:30Z");
console.log(wktMulai.toLocaleString("id-ID")); // Format lokal Indonesia
```

**TypeScript:**
```typescript
const isoTime: string = nilaiData.wkt_mulai;
const dateObj = new Date(isoTime);
const formattedTime = dateObj.toLocaleString("id-ID");
```

---

## 📊 Daftar Periksa Implementasi Frontend

**Endpoint handling:**
- [ ] Handle HTTP 201 (ujian baru — jawaban auto-generated)
- [ ] Handle HTTP 200 (resume ujian — jawaban sudah ada, tidak di-generate ulang)
- [ ] Handle HTTP 400 (ujian sudah selesai)
- [ ] Handle HTTP 401 (token expired/invalid)

**UI/UX:**
- [ ] Tampilkan pesan error yang user-friendly
- [ ] Tampilkan badge/status berbeda untuk ujian baru vs resume
- [ ] Format waktu dengan timezone lokal peserta
- [ ] Disable tombol saat loading
- [ ] Redirect ke login jika 401

**Data handling:**
- [ ] Simpan `nilaiId` dari response untuk endpoint soal selanjutnya
- [ ] Handle `jawaban = null` & `is_benar = null` (belum dijawab) dengan proper null-check
- [ ] Display `no_urut` (urutan random) saat render soal, bukan `no_soal` asli
- [ ] Update jawaban via PUT ke `/api/jawaban/:id` (jangan INSERT baru)

**Fitur keamanan & anti-contek:**
- [ ] Pastikan setiap peserta hanya bisa akses soal mereka sendiri (`id_nilai` mereka)
- [ ] Validasi `id_jadwal` format sebelum submit (optional, tapi lebih baik)
- [ ] Handle urutan soal random per peserta (feature, bukan bug)

---

## 🔗 Endpoint Terkait

Setelah mulai ujian, frontend bisa langsung lanjut ke workflow berikut:

1. **Ambil Daftar Soal (urut random + jawaban kosong):**
   - Endpoint: `GET /api/jawaban/nilai/:id_nilai`
   - Gunakan `id_nilai` dari response endpoint mulai-ujian
   - **Response sudah berisi soal urut random dengan jawaban kosong** (auto-generated saat mulai ujian)
   - Tidak perlu manual generate jawaban

2. **Submit Jawaban Peserta (UPDATE jawaban kosong dengan jawaban sebenarnya):**
   - Endpoint: `PUT /api/jawaban/:id` (atau `PATCH`)
   - Gunakan `id` jawaban dari Step 1
   - Kirim: `id_nilai`, `id_soal`, `id_peserta`, `no_urut`, `jawaban` (A-E)

3. **Selesaikan Ujian:**
   - Endpoint: `PUT /api/nilai/:id` (update `wkt_selesai`)
   - Gunakan `id` dari response mulai-ujian
   - Opsional: hitung nilai akhir berdasarkan jawaban

---

## ❓ FAQ

**Q: Apa bedanya HTTP 201 vs 200?**
A: HTTP 201 = ujian baru dibuat (pertama kali mulai). HTTP 200 = ujian sudah ada (melanjutkan). Frontend bisa menampilkan message berbeda.

**Q: Bagaimana kalau token expired saat ujian sedang berlangsung?**
A: Request akan return HTTP 401. Frontend harus redirect ke login page. Peserta perlu login ulang (nilai/jawaban sudah tersimpan di database).

**Q: Apakah bisa reset ujian yang sudah selesai?**
A: Tidak melalui endpoint ini. Admin perlu menghubungi backend untuk manual reset (update `wkt_selesai` jadi NULL di database).

**Q: Format UUID apa yang benar?**
A: Format standar UUID v4, contoh: `550e8400-e29b-41d4-a716-446655440000`. Pastikan lowercase & ada 5 segmen dipisah dash.

**Q: Apa yang disimpan saat mulai ujian?**
A: ID Peserta, ID Jadwal, Nilai = 0, Waktu Mulai = sekarang, PLUS record jawaban kosong untuk setiap soal (auto-generated). Jadi frontend langsung bisa fetch soal tanpa manual generate.

**Q: Bagaimana jika bank soal kosong (tidak punya soal)?**
A: Record nilai tetap dibuat. Tidak ada record jawaban yang di-generate (0 jawaban). Peserta akan melihat ujian dengan 0 soal.

**Q: Apakah setiap peserta dapat urutan soal yang SAMA?**
A: Tidak! Setiap peserta dapat urutan soal BERBEDA (random per peserta). Ini fitur untuk anti-contek. Peserta A melihat soal urut: 5-2-10-1-3. Peserta B melihat urut: 7-1-4-9-2 (berbeda).

**Q: Dapatkah peserta melihat jawaban soal lain?**
A: Tidak. Masing-masing peserta hanya melihat soal & jawaban mereka sendiri (via `GET /api/jawaban/nilai/:id_nilai`). Jawaban peserta lain tersimpan di record dengan `id_nilai` mereka sendiri.

**Q: Apakah jawaban kosong bisa dihapus?**
A: Tidak bisa dihapus, hanya bisa di-UPDATE dengan jawaban sebenarnya. Jika peserta tidak menjawab soal tertentu, record jawaban tetap ada dengan `jawaban = null`. Ini untuk audit trail lengkap.

---

## 📞 Support

Jika ada error atau pertanyaan:
1. Cek format `id_jadwal` (harus UUID)
2. Cek token validity (tidak expired)
3. Lihat response `message` field untuk detail error
4. Hubungi backend developer jika diperlukan

