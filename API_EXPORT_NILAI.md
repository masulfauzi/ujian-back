# API Documentation - Export Nilai Ujian

## 📌 Overview

Endpoint ini digunakan untuk mengekspor nilai ujian berdasarkan jadwal ujian (`id_jadwal`).
Sistem akan menghasilkan **satu file Excel per kelas** yang terdaftar pada jadwal tersebut,
kemudian semua file dikemas dalam **satu file ZIP** yang siap diunduh.

---

## 🔌 Endpoint

```
GET /api/nilai/export/:id_jadwal
```

---

## 🔐 Authentication

**Wajib menggunakan JWT Token**

Header yang diperlukan:
```
Authorization: Bearer <JWT_TOKEN>
```

Cara mendapatkan token: login terlebih dahulu via endpoint `/api/auth/login`

---

## 📥 Request

### Path Parameter

| Parameter | Tipe | Wajib | Keterangan |
|-----------|------|-------|------------|
| `id_jadwal` | string (UUID) | ✅ Ya | ID jadwal ujian yang akan diekspor (ambil dari tabel `jadwal`) |

### Request Example

```bash
curl -X GET \
  "http://localhost:3000/api/nilai/export/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 📤 Response

### Success Response (200 OK)

**Headers:**
```
Content-Type: application/zip
Content-Disposition: attachment; filename="export_nilai_Ujian_Tengah_Semester.zip"
Content-Length: <file_size>
```

**Body:** Binary ZIP file

**Struktur Isi ZIP:**
```
export_nilai_Ujian_Tengah_Semester.zip
├── X_RPL_1.xlsx          (Sheet: "Nilai Ujian")
├── X_RPL_2.xlsx          (Sheet: "Nilai Ujian")
├── X_TKJ_1.xlsx          (Sheet: "Nilai Ujian")
└── XI_RPL_1.xlsx         (Sheet: "Nilai Ujian")
```

### Excel File Format

Setiap file Excel memiliki struktur berikut:

| No | Nama Peserta | Username | Nilai | Waktu Mulai | Waktu Selesai |
|----|-------------|----------|-------|-------------|---------------|
| 1  | Andi Prasetyo | andi.p | 85.0 | 2024-01-15 08:00:00 | 2024-01-15 09:30:00 |
| 2  | Budi Santoso | budi.s | 72.5 | 2024-01-15 08:05:00 | 2024-01-15 09:28:00 |
| 3  | Citra Dewi | citra.d | 90.0 | 2024-01-15 08:10:00 | 2024-01-15 09:40:00 |

**Penjelasan Kolom:**
- **No**: Nomor urut peserta dalam kelas
- **Nama Peserta**: Nama lengkap dari tabel `peserta`
- **Username**: Username peserta login
- **Nilai**: Nilai akhir ujian (0-100)
- **Waktu Mulai**: Waktu peserta mulai mengerjakan (format: `YYYY-MM-DD HH:MM:SS`, zona: WIB/UTC+7)
- **Waktu Selesai**: Waktu peserta selesai mengerjakan (sama format), atau `-` jika NULL

**Data hanya mencakup peserta yang sudah ada di tabel `nilai`** (yang sudah mengerjakan ujian).

---

## ❌ Error Response

### 400 Bad Request - Jadwal Tidak Ditemukan

**Kondisi:** `id_jadwal` tidak ada atau sudah dihapus

```json
{
  "status": false,
  "message": "jadwal tidak ditemukan",
  "data": null
}
```

### 400 Bad Request - Tidak Ada Kelas Terdaftar

**Kondisi:** Jadwal tidak memiliki kelas yang terdaftar di `jadwal_kelas`

```json
{
  "status": false,
  "message": "tidak ada kelas yang terdaftar pada jadwal ini",
  "data": null
}
```

### 400 Bad Request - ID Jadwal Kosong

**Kondisi:** Path parameter `id_jadwal` tidak dikirim

```json
{
  "status": false,
  "message": "id_jadwal tidak boleh kosong",
  "data": null
}
```

### 401 Unauthorized

**Kondisi:** Token tidak dikirim atau tidak valid

```json
{
  "status": false,
  "message": "Unauthorized",
  "data": null
}
```

### 500 Internal Server Error

**Kondisi:** Error pada server (DB, file system, etc)

```json
{
  "status": false,
  "message": "gagal ambil data kelas X_RPL_1: <error_detail>",
  "data": null
}
```

---

## 💻 Frontend Implementation

### 1️⃣ Vanilla JavaScript (Fetch API)

```javascript
/**
 * Download nilai ujian sebagai ZIP
 * @param {string} idJadwal - UUID jadwal ujian
 * @param {string} token - JWT token dari login
 */
async function downloadNilaiExport(idJadwal, token) {
  try {
    // 1. Validasi input
    if (!idJadwal || !token) {
      throw new Error('idJadwal dan token harus diisi');
    }

    // 2. Fetch dari API
    const response = await fetch(`/api/nilai/export/${idJadwal}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    // 3. Check response status
    if (!response.ok) {
      // Parse error response
      const errorData = await response.json();
      throw new Error(errorData.message || 'Gagal mengunduh file');
    }

    // 4. Ambil filename dari Content-Disposition header
    const disposition = response.headers.get('Content-Disposition');
    let filename = 'export_nilai.zip';
    
    if (disposition) {
      const match = disposition.match(/filename="(.+)"/);
      if (match) {
        filename = match[1];
      }
    }

    // 5. Convert response ke Blob dan trigger download
    const blob = await response.blob();
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);

    console.log('File berhasil diunduh:', filename);
  } catch (error) {
    console.error('Error mengunduh file:', error);
    alert('Gagal mengunduh file: ' + error.message);
  }
}

// Penggunaan:
// const idJadwal = '550e8400-e29b-41d4-a716-446655440000';
// const token = localStorage.getItem('token');
// downloadNilaiExport(idJadwal, token);
```

### 2️⃣ Axios

```javascript
import axios from 'axios';

/**
 * Download nilai ujian menggunakan Axios
 * @param {string} idJadwal - UUID jadwal ujian
 * @param {string} token - JWT token
 */
async function downloadNilaiExportAxios(idJadwal, token) {
  try {
    if (!idJadwal || !token) {
      throw new Error('idJadwal dan token harus diisi');
    }

    // PENTING: responseType harus 'blob' untuk file download
    const response = await axios.get(
      `/api/nilai/export/${idJadwal}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        responseType: 'blob',  // ⚠️ CRITICAL
      }
    );

    // Ambil filename dari header
    const disposition = response.headers['content-disposition'];
    let filename = 'export_nilai.zip';
    
    if (disposition) {
      const match = disposition.match(/filename="(.+)"/);
      if (match) {
        filename = match[1];
      }
    }

    // Trigger download
    const blob = new Blob([response.data], { type: 'application/zip' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    link.click();
    URL.revokeObjectURL(url);

    console.log('File berhasil diunduh:', filename);
  } catch (error) {
    console.error('Error:', error);
    alert('Gagal mengunduh: ' + error.message);
  }
}

// Penggunaan:
// downloadNilaiExportAxios(idJadwal, token);
```

### 3️⃣ React Hook

```javascript
import { useState } from 'react';

function useDownloadNilai() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const download = async (idJadwal, token) => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`/api/nilai/export/${idJadwal}`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || 'Gagal download');
      }

      const blob = await response.blob();
      const disposition = response.headers.get('Content-Disposition');
      const filename = disposition?.match(/filename="(.+)"/)?.[1] || 'export_nilai.zip';

      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      link.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return { download, loading, error };
}

// Component Usage:
export function NilaiExportButton({ idJadwal, token }) {
  const { download, loading, error } = useDownloadNilai();

  return (
    <div>
      <button 
        onClick={() => download(idJadwal, token)}
        disabled={loading}
      >
        {loading ? 'Mengunduh...' : 'Download Nilai (ZIP)'}
      </button>
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </div>
  );
}
```

### 4️⃣ Vue 3 Composition API

```javascript
import { ref } from 'vue';

export function useDownloadNilai() {
  const loading = ref(false);
  const error = ref(null);

  const download = async (idJadwal, token) => {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetch(`/api/nilai/export/${idJadwal}`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || 'Gagal download');
      }

      const blob = await response.blob();
      const disposition = response.headers.get('Content-Disposition');
      const match = disposition?.match(/filename="(.+)"/);
      const filename = match?.[1] || 'export_nilai.zip';

      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      link.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      error.value = err.message;
    } finally {
      loading.value = false;
    }
  };

  return { download, loading, error };
}

// Component Usage (Template):
// <button @click="download(idJadwal, token)" :disabled="loading">
//   {{ loading ? 'Mengunduh...' : 'Download Nilai' }}
// </button>
// <p v-if="error" style="color: red;">{{ error }}</p>
```

### 5️⃣ Angular Service

```typescript
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class NilaiExportService {
  private apiUrl = '/api/nilai/export';

  constructor(private http: HttpClient) {}

  downloadNilai(idJadwal: string, token: string): Observable<Blob> {
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
    });

    return this.http.get<Blob>(
      `${this.apiUrl}/${idJadwal}`,
      {
        headers,
        responseType: 'blob' as 'json',
      }
    ).pipe(
      catchError(error => {
        console.error('Error:', error);
        return throwError(() => new Error('Gagal mengunduh file'));
      })
    );
  }

  triggerDownload(blob: Blob, filename: string) {
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename || 'export_nilai.zip';
    link.click();
    window.URL.revokeObjectURL(url);
  }
}

// Component Usage:
// this.nilaiService.downloadNilai(idJadwal, token).subscribe(
//   (blob) => {
//     this.nilaiService.triggerDownload(blob, 'export_nilai.zip');
//   },
//   (error) => {
//     console.error('Error:', error);
//   }
// );
```

---

## 🎯 Common Use Cases

### Use Case 1: Export dari Table/List Jadwal

```html
<table>
  <thead>
    <tr>
      <th>Nama Ujian</th>
      <th>Tingkat</th>
      <th>Tanggal</th>
      <th>Aksi</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Ujian Tengah Semester</td>
      <td>X</td>
      <td>15-01-2024</td>
      <td>
        <button onclick="downloadNilai('550e8400-e29b-41d4-a716-446655440000')">
          📥 Export Excel
        </button>
      </td>
    </tr>
  </tbody>
</table>
```

### Use Case 2: Modal Dialog dengan Konfirmasi

```javascript
function showExportConfirmation(idJadwal, namaUjian) {
  const confirmed = confirm(
    `Export nilai untuk "${namaUjian}"?\n\nFile akan diunduh sebagai ZIP.`
  );
  
  if (confirmed) {
    const token = localStorage.getItem('token');
    downloadNilaiExport(idJadwal, token);
  }
}
```

### Use Case 3: Toast Notification

```javascript
async function downloadWithNotification(idJadwal, token) {
  try {
    showToast('Mempersiapkan file...', 'info');
    await downloadNilaiExport(idJadwal, token);
    showToast('File berhasil diunduh!', 'success');
  } catch (error) {
    showToast('Gagal: ' + error.message, 'error');
  }
}

function showToast(message, type) {
  // Implementasi sesuai UI library (e.g., Toastify, React Toastify, etc)
  console.log(`[${type.toUpperCase()}] ${message}`);
}
```

---

## 🔍 Troubleshooting

### ❌ Problem: File tidak ter-download

**Kemungkinan Penyebab:**
1. `responseType` tidak diset ke `'blob'` (Axios)
2. Response error tapi tetap mencoba download
3. Browser memblokir popup/download

**Solusi:**
```javascript
// ✅ Selalu check response.ok sebelum download
if (!response.ok) {
  const error = await response.json();
  throw new Error(error.message);
}

// ✅ Set responseType dengan benar
responseType: 'blob'  // Axios
// atau
response.blob()       // Fetch API
```

### ❌ Problem: Error 401 Unauthorized

**Penyebab:** Token tidak valid atau expired

**Solusi:**
```javascript
// Cek token sebelum request
if (!token) {
  // Redirect ke login
  window.location.href = '/login';
  return;
}

// Atau refresh token jika expired
const newToken = await refreshToken();
downloadNilaiExport(idJadwal, newToken);
```

### ❌ Problem: Error 400 "jadwal tidak ditemukan"

**Penyebab:** ID jadwal tidak valid atau sudah dihapus

**Solusi:**
```javascript
// Validasi id_jadwal format
const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

if (!uuidRegex.test(idJadwal)) {
  alert('Format ID jadwal tidak valid');
  return;
}
```

### ❌ Problem: Filename tidak benar

**Penyebab:** Header Content-Disposition tidak terbaca

**Solusi:**
```javascript
// Fallback ke default filename
const disposition = response.headers.get('Content-Disposition') || '';
const filename = disposition.match(/filename="(.+)"/)?.[1] || 'export_nilai.zip';
```

---

## 📊 Integration Checklist

Sebelum deploy ke production, pastikan:

- [ ] Token JWT tersimpan dengan aman (localStorage/sessionStorage/cookie)
- [ ] Error handling sudah implementasi untuk semua status code
- [ ] Loading indicator ditampilkan saat download berlangsung
- [ ] Filename dari server digunakan (jangan hardcode)
- [ ] CORS sudah dikonfigurasi dengan benar di backend
- [ ] Browser support untuk Blob dan URL.createObjectURL
- [ ] Mobile browser support sudah ditest

---

## 📱 Mobile Considerations

Untuk aplikasi mobile (React Native, Flutter Web, etc):

```javascript
// React Native - Download ke local storage
import { documentDirectory, downloadAsync } from 'expo-file-system';

async function downloadNilaiMobile(idJadwal, token) {
  try {
    const fileUri = documentDirectory + 'export_nilai.zip';
    
    await downloadAsync(
      `http://your-api/api/nilai/export/${idJadwal}`,
      fileUri,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    
    console.log('File saved to:', fileUri);
    // Share file atau buka file manager
  } catch (error) {
    console.error('Download failed:', error);
  }
}
```

---

## 🔗 Related Endpoints

| Endpoint | Method | Deskripsi |
|----------|--------|-----------|
| `/api/nilai` | GET | Daftar semua nilai (dengan pagination) |
| `/api/nilai/jadwal/:id_jadwal` | GET | Daftar nilai per jadwal |
| `/api/jadwal` | GET | Daftar semua jadwal ujian |
| `/api/auth/login` | POST | Login & dapatkan JWT token |

---

## 📞 Support & Questions

Jika ada pertanyaan atau masalah, hubungi:
- Backend Developer: <contact>
- Reference: ISSUE_EXPORT_NILAI.md

---

**Last Updated:** 2024-01-15  
**Version:** 1.0.0  
**Status:** ✅ Production Ready
