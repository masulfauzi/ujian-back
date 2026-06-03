# 🚀 Quick Reference - Export Nilai API

## Copy-Paste Solutions

### Vanilla JS (Paling Simple)

```javascript
async function downloadNilai(idJadwal) {
  const token = localStorage.getItem('token');
  
  const res = await fetch(`/api/nilai/export/${idJadwal}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  
  if (!res.ok) {
    const err = await res.json();
    alert('Error: ' + err.message);
    return;
  }
  
  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = res.headers.get('content-disposition')?.match(/filename="(.+)"/)?.[1] || 'export.zip';
  a.click();
  URL.revokeObjectURL(url);
}

// Gunakan:
// <button onclick="downloadNilai('uuid-here')">Download</button>
```

### React Functional Component

```jsx
import { useState } from 'react';

export function ExportButton({ idJadwal }) {
  const [loading, setLoading] = useState(false);
  
  const handleExport = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem('token');
      const res = await fetch(`/api/nilai/export/${idJadwal}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      
      if (!res.ok) throw new Error((await res.json()).message);
      
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      Object.assign(document.createElement('a'), {
        href: url,
        download: 'export.zip'
      }).click();
      URL.revokeObjectURL(url);
    } catch (err) {
      alert('Error: ' + err.message);
    } finally {
      setLoading(false);
    }
  };
  
  return (
    <button onClick={handleExport} disabled={loading}>
      {loading ? '⏳ Downloading...' : '📥 Export'}
    </button>
  );
}
```

### React + Axios

```jsx
import axios from 'axios';

const handleExport = async (idJadwal) => {
  try {
    const { data, headers } = await axios.get(
      `/api/nilai/export/${idJadwal}`,
      {
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        responseType: 'blob'  // ⚠️ IMPORTANT
      }
    );
    
    const url = URL.createObjectURL(data);
    Object.assign(document.createElement('a'), {
      href: url,
      download: headers['content-disposition']?.match(/filename="(.+)"/)?.[1] || 'export.zip'
    }).click();
    URL.revokeObjectURL(url);
  } catch (err) {
    alert('Error: ' + err.response?.data?.message || err.message);
  }
};
```

### Vue 3

```vue
<template>
  <button @click="exportNilai" :disabled="loading">
    {{ loading ? 'Downloading...' : 'Download Nilai' }}
  </button>
</template>

<script setup>
import { ref } from 'vue';

const loading = ref(false);
const idJadwal = 'your-id-here';

const exportNilai = async () => {
  loading.value = true;
  try {
    const res = await fetch(`/api/nilai/export/${idJadwal}`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    });
    
    if (!res.ok) throw new Error((await res.json()).message);
    
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'export.zip';
    a.click();
    URL.revokeObjectURL(url);
  } catch (err) {
    alert('Error: ' + err.message);
  } finally {
    loading.value = false;
  }
};
</script>
```

---

## ⚡ Key Points

| Item | Detail |
|------|--------|
| **Endpoint** | `GET /api/nilai/export/:id_jadwal` |
| **Auth** | Bearer token di header `Authorization` |
| **Response Type** | `Blob` (binary ZIP file) |
| **Axios Setting** | Harus `responseType: 'blob'` |
| **Success Status** | `200 OK` |
| **Error Format** | JSON: `{status: false, message: "..."}` |
| **File Type** | ZIP (berisi N file Excel .xlsx) |

---

## 🆘 Instant Fixes

| Problem | Fix |
|---------|-----|
| "Blob is not defined" | Gunakan browser environment, bukan Node.js |
| File tidak download | Pastikan `responseType: 'blob'` (Axios) atau `await res.blob()` (Fetch) |
| 401 Unauthorized | Login dulu, pastikan token valid: `localStorage.getItem('token')` |
| 400 "jadwal tidak ditemukan" | Periksa ID jadwal valid dan tidak dihapus di database |
| Filename salah | Jangan hardcode, ambil dari header: `content-disposition` |
| CORS Error | Backend sudah setup CORS, pastikan origin cocok |

---

## 📋 HTML Button Example

```html
<!-- Minimal -->
<button onclick="downloadNilai('550e8400-e29b-41d4-a716-446655440000')">
  Download
</button>

<!-- Bootstrap -->
<button class="btn btn-primary" onclick="downloadNilai(idJadwal)">
  <i class="bi bi-download"></i> Export Nilai
</button>

<!-- Tailwind -->
<button class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        onclick="downloadNilai(idJadwal)">
  📥 Export Nilai
</button>

<!-- Material Design -->
<button class="mdc-button mdc-button--raised" @click="exportNilai">
  <span class="mdc-button__label">Export</span>
</button>
```

---

## 🔄 Flow Diagram

```
User Click Button
       ↓
Get Token from Storage
       ↓
POST /api/nilai/export/:id_jadwal
       ↓
Server Create Excel Files + ZIP
       ↓
Response: ZIP Binary (200 OK)
       ↓
Create Blob & Download
       ↓
✅ File Saved to Downloads
```

---

## 📦 Response Structure

**Success (200):**
```
Headers:
  Content-Type: application/zip
  Content-Disposition: attachment; filename="export_nilai_Ujian.zip"

Body: Binary ZIP file
```

**Error (4xx/5xx):**
```json
{
  "status": false,
  "message": "Error description here",
  "data": null
}
```

---

## 🎓 Learning Path

1. **Beginners** → Gunakan "Vanilla JS (Paling Simple)"
2. **React Dev** → Gunakan "React Functional Component"
3. **Advanced** → Customize dengan loading state, error handling, retry logic

---

## 💾 Token Management

```javascript
// Get Token
const token = localStorage.getItem('token');

// Cek Token Expired
if (!token) {
  // Redirect to login
  window.location.href = '/login';
}

// Refresh Token (jika ada endpoint)
const newToken = await refreshToken(); // Implement this
localStorage.setItem('token', newToken);
```

---

## 📱 Mobile Browser Support

```javascript
// Android Chrome, iOS Safari, Firefox - semua supported ✅
// File akan diunduh ke folder Downloads
```

---

## 🚀 Production Checklist

- [ ] Token diambil dari secure storage
- [ ] Loading indicator ditampilkan
- [ ] Error message user-friendly
- [ ] Filename dinamis dari server
- [ ] Browser compatibility tested
- [ ] Slow network tested
- [ ] Cancel download supported (optional)

---

## 🔗 Other Useful Endpoints

```
// Get all jadwal
GET /api/jadwal

// Get jadwal by ID
GET /api/jadwal/:id

// Get nilai by jadwal
GET /api/nilai/jadwal/:id_jadwal

// Get nilai by peserta
GET /api/nilai/peserta/:id_peserta
```

---

## 📞 When to Call API

✅ **User clicks "Export Nilai" button**  
✅ **User confirms export in modal**  
✅ **Auto-export at scheduled time**  

❌ **Page load (ambil data biasa, jangan file download)**  
❌ **Background service (user tidak perlu, gunakan webhook)**

---

**Last Updated:** 2024-01-15  
**Copy-paste ready:** ✅ Yes  
**Need more help?** → Check API_EXPORT_NILAI.md
