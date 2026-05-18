# CORS Error Fix Guide

## ❌ Error
```
strict-origin-when-cross-origin
```

## 🔍 Penyebab
- Frontend berjalan di: `http://localhost:5173`
- Backend berjalan di: `http://localhost:3000`
- Browser memblokir cross-origin request karena CORS belum dikonfigurasi

---

## ✅ Solusi: Setup CORS di Backend

### **Jika Backend Node.js/Express**

**File:** `server.js` atau `app.js`

```javascript
const express = require('express');
const cors = require('cors');
const app = express();

// Setup CORS
app.use(cors({
    origin: 'http://localhost:5173',
    credentials: true,
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS', 'PATCH'],
    allowedHeaders: ['Content-Type', 'Authorization'],
}));

// Middleware lainnya
app.use(express.json());

// Routes
app.post('/api/auth/register', (req, res) => {
    // ... handler code
});

app.post('/api/auth/login', (req, res) => {
    // ... handler code
});

// Start server
app.listen(3000, () => {
    console.log('Server running on http://localhost:3000');
});
```

**Install cors package jika belum:**
```bash
npm install cors
```

---

### **Jika Backend PHP**

**File:** Awal file routes atau middleware

```php
<?php
// Set CORS Headers
header('Access-Control-Allow-Origin: http://localhost:5173');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH');
header('Access-Control-Allow-Headers: Content-Type, Authorization');
header('Access-Control-Allow-Credentials: true');
header('Content-Type: application/json');

// Handle preflight requests
if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// Routes
if ($_SERVER['REQUEST_METHOD'] === 'POST' && $_SERVER['REQUEST_URI'] === '/api/auth/register') {
    // ... handler code
}
?>
```

---

### **Jika Backend Python/Flask**

```python
from flask import Flask
from flask_cors import CORS

app = Flask(__name__)

# Setup CORS
CORS(app, resources={
    r"/api/*": {
        "origins": ["http://localhost:5173"],
        "methods": ["GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"],
        "allow_headers": ["Content-Type", "Authorization"],
        "supports_credentials": True
    }
})

@app.route('/api/auth/register', methods=['POST'])
def register():
    # ... handler code
    pass

@app.route('/api/auth/login', methods=['POST'])
def login():
    # ... handler code
    pass

if __name__ == '__main__':
    app.run(host='localhost', port=3000, debug=True)
```

**Install flask-cors jika belum:**
```bash
pip install flask-cors
```

---

### **Jika Backend Laravel**

**File:** `config/cors.php`

```php
<?php

return [
    'paths' => ['api/*'],
    'allowed_methods' => ['*'],
    'allowed_origins' => ['http://localhost:5173'],
    'allowed_origins_patterns' => [],
    'allowed_headers' => ['*'],
    'exposed_headers' => [],
    'max_age' => 0,
    'supports_credentials' => true,
];
```

---

## 🧪 Testing CORS Fix

### Step 1: Restart Backend Server
```bash
# Node.js
npm start
# atau
node server.js

# Python
python app.py

# PHP (jika pakai built-in server)
php -S localhost:3000
```

### Step 2: Restart Frontend Dev Server
```bash
npm run dev
```

### Step 3: Test Register
1. Buka http://localhost:5173
2. Click "Create Account"
3. Isi form dengan data test
4. Click "Create Account"

### Step 4: Check Network Tab
- DevTools → Network tab
- Cari request ke `/api/auth/register`
- Lihat Response Headers:
  ```
  access-control-allow-origin: http://localhost:5173
  access-control-allow-credentials: true
  ```

---

## 🔍 Debugging

### Jika masih error:

1. **Pastikan backend running:**
   ```bash
   curl http://localhost:3000
   ```

2. **Check browser console untuk detail error:**
   - DevTools → Console
   - Lihat error message lengkapnya

3. **Test dengan curl:**
   ```bash
   curl -X POST http://localhost:3000/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"name":"Test","email":"test@test.com","password":"123456"}'
   ```

4. **Check backend logs:**
   - Lihat apakah request diterima backend
   - Cek error messages

---

## 📋 Checklist

- [ ] Backend CORS configuration ditambah
- [ ] Backend server di-restart
- [ ] Frontend dev server di-restart
- [ ] Test register berhasil tanpa CORS error
- [ ] Response status 200 atau 400 (bukan 0 atau blocked)
- [ ] Token diterima di response

---

## ⚡ Production Deployment

Jangan gunakan `*` (wildcard) untuk origin di production!

**Development:**
```javascript
origin: 'http://localhost:5173'
```

**Production:**
```javascript
origin: 'https://yourdomain.com'
```

---

## 📞 Support

Jika masih error:
1. Tentukan framework backend yang digunakan
2. Share backend code snippet
3. Share error message dari browser console
4. Share request/response dari Network tab
