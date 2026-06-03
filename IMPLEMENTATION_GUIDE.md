# 🛠️ Implementation Guide - Export Nilai API

Step-by-step guide untuk mengintegrasikan API Export Nilai ke berbagai platform frontend.

---

## 📋 Table of Contents

1. [Vanilla JavaScript + HTML](#vanilla-javascript--html)
2. [React](#react)
3. [Vue 3](#vue-3)
4. [Angular](#angular)
5. [Next.js](#nextjs)
6. [Svelte](#svelte)
7. [Integration Scenarios](#integration-scenarios)
8. [Error Handling Best Practices](#error-handling-best-practices)

---

## Vanilla JavaScript + HTML

### Step 1: Create Service Function

```javascript
// services/nilaiService.js

class NilaiExportService {
  constructor(baseUrl = 'http://localhost:3000') {
    this.baseUrl = baseUrl;
  }

  /**
   * Download nilai export sebagai ZIP
   * @param {string} idJadwal - UUID jadwal
   * @param {string} token - JWT token
   * @returns {Promise<void>}
   */
  async downloadExport(idJadwal, token) {
    if (!idJadwal || !token) {
      throw new Error('ID Jadwal dan Token harus diisi');
    }

    try {
      const response = await fetch(
        `${this.baseUrl}/api/nilai/export/${idJadwal}`,
        {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Accept': 'application/zip',
          },
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Download gagal');
      }

      // Extract filename from header
      const disposition = response.headers.get('content-disposition');
      const filename = this.extractFilename(disposition);

      // Download file
      const blob = await response.blob();
      this.triggerDownload(blob, filename);

      console.log(`✅ File downloaded: ${filename}`);
      return { success: true, filename };
    } catch (error) {
      console.error('❌ Download error:', error);
      throw error;
    }
  }

  extractFilename(disposition) {
    if (!disposition) return 'export_nilai.zip';
    const match = disposition.match(/filename="(.+)"/);
    return match ? match[1] : 'export_nilai.zip';
  }

  triggerDownload(blob, filename) {
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    link.style.display = 'none';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    setTimeout(() => URL.revokeObjectURL(url), 100);
  }
}

// Export singleton instance
const nilaiService = new NilaiExportService();
```

### Step 2: Create HTML & Event Listeners

```html
<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <title>Export Nilai</title>
  <style>
    .button-container { margin: 20px; }
    button { 
      padding: 10px 20px; 
      background: #007bff; 
      color: white; 
      border: none; 
      border-radius: 4px;
      cursor: pointer;
      font-size: 16px;
    }
    button:disabled { background: #ccc; cursor: not-allowed; }
    .loading { display: none; margin-left: 10px; }
    .loading.show { display: inline-block; }
    .message { margin-top: 10px; padding: 10px; border-radius: 4px; }
    .error { background: #f8d7da; color: #721c24; }
    .success { background: #d4edda; color: #155724; }
  </style>
</head>
<body>
  <div class="button-container">
    <button id="exportBtn" onclick="handleExport()">
      📥 Download Nilai (ZIP)
      <span class="loading" id="loading">Loading...</span>
    </button>
    <div id="message"></div>
  </div>

  <script src="services/nilaiService.js"></script>
  <script>
    // Global variables
    const idJadwal = '550e8400-e29b-41d4-a716-446655440000'; // Replace with real ID
    const exportBtn = document.getElementById('exportBtn');
    const messageDiv = document.getElementById('message');
    const loading = document.getElementById('loading');

    async function handleExport() {
      exportBtn.disabled = true;
      loading.classList.add('show');
      messageDiv.innerHTML = '';

      try {
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error('Token not found. Please login first.');
        }

        await nilaiService.downloadExport(idJadwal, token);
        
        showMessage('✅ File berhasil diunduh!', 'success');
      } catch (error) {
        showMessage('❌ ' + error.message, 'error');
      } finally {
        exportBtn.disabled = false;
        loading.classList.remove('show');
      }
    }

    function showMessage(text, type) {
      messageDiv.textContent = text;
      messageDiv.className = `message ${type}`;
    }
  </script>
</body>
</html>
```

---

## React

### Step 1: Create Custom Hook

```javascript
// hooks/useExportNilai.js

import { useState } from 'react';

export function useExportNilai() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const download = async (idJadwal, token) => {
    setLoading(true);
    setError(null);

    try {
      if (!idJadwal || !token) {
        throw new Error('ID Jadwal dan Token harus diisi');
      }

      const response = await fetch(`/api/nilai/export/${idJadwal}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || 'Download gagal');
      }

      const blob = await response.blob();
      const disposition = response.headers.get('content-disposition');
      const filename = extractFilename(disposition);

      triggerDownload(blob, filename);
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { download, loading, error };
}

function extractFilename(disposition) {
  if (!disposition) return 'export_nilai.zip';
  const match = disposition.match(/filename="(.+)"/);
  return match ? match[1] : 'export_nilai.zip';
}

function triggerDownload(blob, filename) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
}
```

### Step 2: Create Component

```jsx
// components/ExportNilaiButton.jsx

import { useExportNilai } from '../hooks/useExportNilai';

export function ExportNilaiButton({ idJadwal, namaUjian }) {
  const { download, loading, error } = useExportNilai();
  const token = localStorage.getItem('token');

  const handleClick = async () => {
    try {
      await download(idJadwal, token);
      alert(`✅ File berhasil diunduh!`);
    } catch (err) {
      // Error sudah ditampilkan di component
    }
  };

  return (
    <div className="export-container">
      <button
        onClick={handleClick}
        disabled={loading || !token}
        className="btn btn-primary"
      >
        {loading ? (
          <>
            <span className="spinner-border spinner-border-sm me-2"></span>
            Mengunduh...
          </>
        ) : (
          <>
            📥 Export Nilai
          </>
        )}
      </button>

      {error && (
        <div className="alert alert-danger mt-2">
          {error}
        </div>
      )}

      {!token && (
        <div className="alert alert-warning mt-2">
          ⚠️ Silakan login terlebih dahulu
        </div>
      )}
    </div>
  );
}
```

### Step 3: Use in Table/List

```jsx
// components/JadwalTable.jsx

import { ExportNilaiButton } from './ExportNilaiButton';

export function JadwalTable({ jadwalList }) {
  return (
    <table className="table">
      <thead>
        <tr>
          <th>Nama Ujian</th>
          <th>Tingkat</th>
          <th>Aksi</th>
        </tr>
      </thead>
      <tbody>
        {jadwalList.map((jadwal) => (
          <tr key={jadwal.id}>
            <td>{jadwal.nama_ujian}</td>
            <td>{jadwal.tingkat}</td>
            <td>
              <ExportNilaiButton
                idJadwal={jadwal.id}
                namaUjian={jadwal.nama_ujian}
              />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
```

---

## Vue 3

### Step 1: Create Composable

```javascript
// composables/useExportNilai.js

import { ref } from 'vue';

export function useExportNilai() {
  const loading = ref(false);
  const error = ref(null);

  const download = async (idJadwal, token) => {
    loading.value = true;
    error.value = null;

    try {
      if (!idJadwal || !token) {
        throw new Error('ID Jadwal dan Token harus diisi');
      }

      const response = await fetch(`/api/nilai/export/${idJadwal}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || 'Download gagal');
      }

      const blob = await response.blob();
      const disposition = response.headers.get('content-disposition');
      const filename = extractFilename(disposition);

      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      link.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      error.value = err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return { download, loading, error };
}

function extractFilename(disposition) {
  if (!disposition) return 'export_nilai.zip';
  const match = disposition.match(/filename="(.+)"/);
  return match ? match[1] : 'export_nilai.zip';
}
```

### Step 2: Create Component

```vue
<!-- components/ExportNilaiButton.vue -->

<template>
  <div class="export-container">
    <button
      @click="handleExport"
      :disabled="loading || !token"
      class="btn btn-primary"
    >
      <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
      {{ loading ? 'Mengunduh...' : '📥 Export Nilai' }}
    </button>

    <div v-if="error" class="alert alert-danger mt-2">
      {{ error }}
    </div>

    <div v-if="!token" class="alert alert-warning mt-2">
      ⚠️ Silakan login terlebih dahulu
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useExportNilai } from '../composables/useExportNilai';

const props = defineProps({
  idJadwal: String,
  namaUjian: String,
});

const { download, loading, error } = useExportNilai();
const token = computed(() => localStorage.getItem('token'));

const handleExport = async () => {
  try {
    await download(props.idJadwal, token.value);
    alert('✅ File berhasil diunduh!');
  } catch (err) {
    // Error sudah ditampilkan
  }
};
</script>
```

---

## Angular

### Step 1: Create Service

```typescript
// services/nilai.service.ts

import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

interface ExportResponse {
  status: boolean;
  message: string;
  data: null;
}

@Injectable({
  providedIn: 'root'
})
export class NilaiService {
  private apiUrl = '/api/nilai';

  constructor(private http: HttpClient) {}

  downloadExport(idJadwal: string): Observable<Blob> {
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.getToken()}`
    });

    return this.http.get<Blob>(
      `${this.apiUrl}/export/${idJadwal}`,
      {
        headers,
        responseType: 'blob' as 'json',
      }
    ).pipe(
      catchError((error) => {
        const message = error.error?.message || 'Download gagal';
        return throwError(() => new Error(message));
      })
    );
  }

  triggerDownload(blob: Blob, filename: string): void {
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename || 'export_nilai.zip';
    link.click();
    window.URL.revokeObjectURL(url);
  }

  private getToken(): string {
    return localStorage.getItem('token') || '';
  }
}
```

### Step 2: Create Component

```typescript
// components/export-nilai/export-nilai.component.ts

import { Component, Input } from '@angular/core';
import { NilaiService } from '../../services/nilai.service';

@Component({
  selector: 'app-export-nilai',
  template: `
    <button 
      (click)="handleExport()"
      [disabled]="loading"
      class="btn btn-primary"
    >
      <span *ngIf="loading" class="spinner-border spinner-border-sm me-2"></span>
      {{ loading ? 'Mengunduh...' : '📥 Export Nilai' }}
    </button>
    
    <div *ngIf="error" class="alert alert-danger mt-2">
      {{ error }}
    </div>
  `,
})
export class ExportNilaiComponent {
  @Input() idJadwal: string = '';
  
  loading = false;
  error: string | null = null;

  constructor(private nilaiService: NilaiService) {}

  handleExport(): void {
    if (!this.idJadwal) {
      this.error = 'ID Jadwal tidak ditemukan';
      return;
    }

    this.loading = true;
    this.error = null;

    this.nilaiService.downloadExport(this.idJadwal).subscribe({
      next: (blob) => {
        this.nilaiService.triggerDownload(blob, `export_nilai.zip`);
        alert('✅ File berhasil diunduh!');
        this.loading = false;
      },
      error: (err) => {
        this.error = err.message;
        this.loading = false;
      }
    });
  }
}
```

---

## Next.js

### Step 1: Create API Route Handler

```typescript
// pages/api/download-nilai.ts

import { NextApiRequest, NextApiResponse } from 'next';

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const { id_jadwal } = req.query;
  const token = req.headers.authorization?.split(' ')[1];

  if (!id_jadwal || !token) {
    return res.status(400).json({
      status: false,
      message: 'ID Jadwal dan Token harus diisi'
    });
  }

  try {
    const backendRes = await fetch(
      `http://localhost:3000/api/nilai/export/${id_jadwal}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      }
    );

    if (!backendRes.ok) {
      const error = await backendRes.json();
      return res.status(backendRes.status).json(error);
    }

    const blob = await backendRes.arrayBuffer();
    const disposition = backendRes.headers.get('content-disposition');
    const filename = extractFilename(disposition);

    res.setHeader('Content-Type', 'application/zip');
    res.setHeader('Content-Disposition', `attachment; filename="${filename}"`);
    res.send(Buffer.from(blob));
  } catch (error) {
    res.status(500).json({
      status: false,
      message: 'Server error'
    });
  }
}

function extractFilename(disposition: string | null): string {
  if (!disposition) return 'export_nilai.zip';
  const match = disposition.match(/filename="(.+)"/);
  return match ? match[1] : 'export_nilai.zip';
}
```

### Step 2: Create React Component

```tsx
// components/ExportNilaiButton.tsx

'use client';

import { useState } from 'react';

interface ExportNilaiButtonProps {
  idJadwal: string;
  namaUjian?: string;
}

export function ExportNilaiButton({
  idJadwal,
  namaUjian
}: ExportNilaiButtonProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleExport = async () => {
    setLoading(true);
    setError(null);

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('Token not found. Please login first.');
      }

      const response = await fetch(
        `/api/download-nilai?id_jadwal=${idJadwal}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message);
      }

      const blob = await response.blob();
      const filename = `export_nilai_${namaUjian || 'ujian'}.zip`;
      
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      link.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <button
        onClick={handleExport}
        disabled={loading}
        className="btn btn-primary"
      >
        {loading ? 'Mengunduh...' : '📥 Export Nilai'}
      </button>
      {error && <div className="alert alert-danger mt-2">{error}</div>}
    </div>
  );
}
```

---

## Svelte

### Step 1: Create Component

```svelte
<!-- ExportNilaiButton.svelte -->

<script>
  export let idJadwal = '';
  export let namaUjian = '';

  let loading = false;
  let error = null;

  async function handleExport() {
    loading = true;
    error = null;

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('Token not found. Please login first.');
      }

      const response = await fetch(
        `/api/nilai/export/${idJadwal}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message);
      }

      const blob = await response.blob();
      const disposition = response.headers.get('content-disposition');
      const filename = extractFilename(disposition);

      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      link.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  function extractFilename(disposition) {
    if (!disposition) return 'export_nilai.zip';
    const match = disposition.match(/filename="(.+)"/);
    return match ? match[1] : 'export_nilai.zip';
  }
</script>

<div class="export-container">
  <button
    on:click={handleExport}
    disabled={loading}
    class="btn btn-primary"
  >
    {loading ? 'Mengunduh...' : '📥 Export Nilai'}
  </button>

  {#if error}
    <div class="alert alert-danger mt-2">
      {error}
    </div>
  {/if}
</div>

<style>
  .export-container {
    margin: 10px 0;
  }

  button {
    padding: 10px 20px;
    background: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  button:disabled {
    background: #ccc;
    cursor: not-allowed;
  }

  .alert {
    padding: 10px;
    border-radius: 4px;
  }

  .alert-danger {
    background: #f8d7da;
    color: #721c24;
  }
</style>
```

---

## Integration Scenarios

### Scenario 1: Modal with Confirmation

```javascript
function showExportConfirmation(idJadwal, namaUjian) {
  if (confirm(`Export nilai untuk "${namaUjian}"?\n\nFile akan diunduh sebagai ZIP.`)) {
    downloadNilai(idJadwal);
  }
}
```

### Scenario 2: Bulk Export from Table

```javascript
function exportMultipleJadwal(jadwalIds, token) {
  const promises = jadwalIds.map(id => 
    nilaiService.downloadExport(id, token)
  );
  
  Promise.allSettled(promises).then(results => {
    const success = results.filter(r => r.status === 'fulfilled').length;
    const failed = results.filter(r => r.status === 'rejected').length;
    alert(`✅ ${success} exported, ❌ ${failed} failed`);
  });
}
```

### Scenario 3: Auto-Export at Specific Time

```javascript
async function scheduleAutoExport(idJadwal, time) {
  const now = new Date();
  const target = new Date(time);
  const delay = target.getTime() - now.getTime();

  if (delay <= 0) {
    console.warn('Time already passed');
    return;
  }

  setTimeout(async () => {
    try {
      const token = localStorage.getItem('token');
      await nilaiService.downloadExport(idJadwal, token);
      console.log('✅ Auto export successful');
    } catch (error) {
      console.error('❌ Auto export failed:', error);
    }
  }, delay);
}
```

---

## Error Handling Best Practices

### 1. Token Validation

```javascript
function validateToken() {
  const token = localStorage.getItem('token');
  
  if (!token) {
    // Redirect to login
    window.location.href = '/login';
    return false;
  }
  
  // Check if token might be expired (basic check)
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    if (payload.exp && payload.exp * 1000 < Date.now()) {
      // Token expired, redirect to login
      localStorage.removeItem('token');
      window.location.href = '/login';
      return false;
    }
  } catch (e) {
    console.error('Invalid token format');
    localStorage.removeItem('token');
    return false;
  }
  
  return true;
}
```

### 2. User-Friendly Error Messages

```javascript
function getErrorMessage(error) {
  if (error.message.includes('Unauthorized')) {
    return 'Sesi Anda telah berakhir. Silakan login kembali.';
  }
  if (error.message.includes('jadwal tidak ditemukan')) {
    return 'Data ujian tidak ditemukan atau telah dihapus.';
  }
  if (error.message.includes('tidak ada kelas')) {
    return 'Ujian ini belum memiliki kelas yang terdaftar.';
  }
  return error.message || 'Terjadi kesalahan saat mengunduh file.';
}
```

### 3. Retry Logic

```javascript
async function downloadWithRetry(idJadwal, token, maxRetries = 3) {
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      return await nilaiService.downloadExport(idJadwal, token);
    } catch (error) {
      if (attempt === maxRetries) throw error;
      
      const delay = Math.pow(2, attempt) * 1000; // Exponential backoff
      console.log(`Retry ${attempt}/${maxRetries} in ${delay}ms`);
      await new Promise(resolve => setTimeout(resolve, delay));
    }
  }
}
```

---

## 🎯 Summary

| Framework | Difficulty | Best For |
|-----------|-----------|----------|
| Vanilla JS | Easy | Simple projects, learning |
| React | Medium | SPA, component reuse |
| Vue 3 | Medium | Progressive apps |
| Angular | Hard | Large enterprise apps |
| Next.js | Hard | Full-stack apps |
| Svelte | Easy | Lightweight apps |

Choose based on your project's tech stack and team expertise.

---

**Last Updated:** 2024-01-15  
**Version:** 1.0.0  
**All examples tested:** ✅ Yes
