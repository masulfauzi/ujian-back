# ISSUE: Implementasi Image Upload pada Modul Soal

## 📋 Deskripsi Fitur

Menambahkan kemampuan untuk menerima upload gambar pada modul soal, menyimpan file gambar ke folder lokal, dan menyimpan nama/path gambar ke database.

### Kolom yang Perlu Diupdate:
- `gambar_soal` - Gambar pertanyaan soal
- `gambar_a` - Gambar opsi A
- `gambar_b` - Gambar opsi B
- `gambar_c` - Gambar opsi C
- `gambar_d` - Gambar opsi D
- `gambar_e` - Gambar opsi E

---

## 🎯 Tujuan Implementasi

1. Menerima file gambar melalui HTTP multipart/form-data
2. Validasi file (format, ukuran)
3. Menyimpan file ke folder server
4. Menyimpan nama/path file ke database
5. Mengembalikan URL/path gambar dalam response API
6. Menangani penghapusan gambar lama saat update
7. Menangani soft delete gambar saat delete soal

---

## 🔄 Fase Implementasi (8 Tahap)

### **FASE 1: Setup Konfigurasi Upload**

**Tujuan:** Mempersiapkan konfigurasi dan konstanta untuk upload gambar

**Langkah-langkah:**

1. **Buat file config upload** di `internal/config/upload.go`:
   ```go
   package config
   
   import "path/filepath"
   
   type UploadConfig struct {
       // Folder tempat menyimpan gambar
       ImageUploadPath string
       // Ukuran maksimal file (bytes)
       MaxFileSize int64
       // Format/extension yang diizinkan
       AllowedFormats []string
       // Base URL untuk mengakses gambar
       ImageBaseURL string
   }
   
   func GetUploadConfig() UploadConfig {
       return UploadConfig{
           ImageUploadPath: "./uploads/soal",
           MaxFileSize: 5 * 1024 * 1024, // 5MB
           AllowedFormats: []string{"jpg", "jpeg", "png", "gif", "webp"},
           ImageBaseURL: "http://localhost:3000/uploads/soal",
       }
   }
   ```

2. **Buat folder uploads** jika belum ada
   - Struktur: `/uploads/soal/` - untuk menyimpan semua gambar soal
   - Pastikan folder memiliki read/write permission

3. **Update `.env`** jika diperlukan untuk konfigurasi path

**File yang Dimodifikasi:**
- `internal/config/upload.go` (baru)
- `.env` (opsional)

---

### **FASE 2: Buat Utility/Helper untuk Upload**

**Tujuan:** Membuat reusable function untuk handle file upload dan validasi

**Langkah-langkah:**

1. **Buat file utility** di `internal/utils/file_upload.go`:
   ```go
   package utils
   
   import (
       "fmt"
       "mime/multipart"
       "os"
       "path/filepath"
       "strings"
       "time"
       
       "backend/internal/config"
   )
   
   // SaveImage menyimpan file gambar dan mengembalikan filename
   func SaveImage(file *multipart.FileHeader, folder string) (string, error) {
       cfg := config.GetUploadConfig()
       
       // Validasi ukuran file
       if file.Size > cfg.MaxFileSize {
           return "", fmt.Errorf("file terlalu besar, maksimal %.2f MB", float64(cfg.MaxFileSize)/1024/1024)
       }
       
       // Validasi extension
       ext := strings.ToLower(filepath.Ext(file.Filename))
       ext = strings.TrimPrefix(ext, ".")
       
       if !isAllowedFormat(ext, cfg.AllowedFormats) {
           return "", fmt.Errorf("format file tidak diizinkan: %s", ext)
       }
       
       // Generate unique filename dengan timestamp
       timestamp := time.Now().Unix()
       randomStr := generateRandomString(8)
       newFilename := fmt.Sprintf("%d_%s.%s", timestamp, randomStr, ext)
       
       // Path lengkap untuk menyimpan file
       uploadPath := filepath.Join(cfg.ImageUploadPath, folder)
       
       // Buat folder jika tidak ada
       if err := os.MkdirAll(uploadPath, 0755); err != nil {
           return "", fmt.Errorf("gagal membuat folder: %v", err)
       }
       
       // Buka file dari request
       src, err := file.Open()
       if err != nil {
           return "", fmt.Errorf("gagal membuka file: %v", err)
       }
       defer src.Close()
       
       // Buat file destination
       filepath := filepath.Join(uploadPath, newFilename)
       dst, err := os.Create(filepath)
       if err != nil {
           return "", fmt.Errorf("gagal membuat file: %v", err)
       }
       defer dst.Close()
       
       // Copy file content
       _, err = dst.ReadFrom(src)
       if err != nil {
           // Hapus file jika gagal menyalin
           os.Remove(filepath)
           return "", fmt.Errorf("gagal menyimpan file: %v", err)
       }
       
       return newFilename, nil
   }
   
   // DeleteImage menghapus file gambar
   func DeleteImage(filename string, folder string) error {
       if filename == "" {
           return nil // Skip jika filename kosong
       }
       
       cfg := config.GetUploadConfig()
       filepath := filepath.Join(cfg.ImageUploadPath, folder, filename)
       
       // Cek apakah file ada
       if _, err := os.Stat(filepath); os.IsNotExist(err) {
           return nil // File tidak ada, skip
       }
       
       return os.Remove(filepath)
   }
   
   // Helper functions
   func isAllowedFormat(ext string, allowed []string) bool {
       for _, a := range allowed {
           if ext == a {
               return true
           }
       }
       return false
   }
   
   func generateRandomString(length int) string {
       // Implementation untuk generate random string
       // Atau gunakan library crypto/rand
       return generateRandomStringHelper(length)
   }
   
   // TODO: Implementasikan function generateRandomStringHelper
   ```

2. **Update DTO** untuk menerima file upload di `internal/modules/soal/dto/soal_dto.go`:
   ```go
   import "mime/multipart"
   
   type CreateSoalRequest struct {
       IdBankSoal  string                 `json:"id_bank_soal" form:"id_bank_soal" validate:"required"`
       NoSoal      int                    `json:"no_soal" form:"no_soal" validate:"required,min=1"`
       Soal        string                 `json:"soal" form:"soal" validate:"required"`
       GambarSoal  *multipart.FileHeader  `form:"gambar_soal"` // File upload
       OpsiA       string                 `json:"opsi_a" form:"opsi_a" validate:"required"`
       OpsiB       string                 `json:"opsi_b" form:"opsi_b" validate:"required"`
       OpsiC       string                 `json:"opsi_c" form:"opsi_c" validate:"required"`
       OpsiD       string                 `json:"opsi_d" form:"opsi_d"`
       OpsiE       string                 `json:"opsi_e" form:"opsi_e"`
       GambarA     *multipart.FileHeader  `form:"gambar_a"` // File upload
       GambarB     *multipart.FileHeader  `form:"gambar_b"` // File upload
       GambarC     *multipart.FileHeader  `form:"gambar_c"` // File upload
       GambarD     *multipart.FileHeader  `form:"gambar_d"` // File upload
       GambarE     *multipart.FileHeader  `form:"gambar_e"` // File upload
       Kunci       string                 `json:"kunci" form:"kunci" validate:"required,len=1"`
   }
   
   type UpdateSoalRequest struct {
       NoSoal      int                    `json:"no_soal" form:"no_soal" validate:"required,min=1"`
       Soal        string                 `json:"soal" form:"soal" validate:"required"`
       GambarSoal  *multipart.FileHeader  `form:"gambar_soal"` // File upload (optional)
       OpsiA       string                 `json:"opsi_a" form:"opsi_a" validate:"required"`
       OpsiB       string                 `json:"opsi_b" form:"opsi_b" validate:"required"`
       OpsiC       string                 `json:"opsi_c" form:"opsi_c" validate:"required"`
       OpsiD       string                 `json:"opsi_d" form:"opsi_d"`
       OpsiE       string                 `json:"opsi_e" form:"opsi_e"`
       GambarA     *multipart.FileHeader  `form:"gambar_a"` // File upload (optional)
       GambarB     *multipart.FileHeader  `form:"gambar_b"` // File upload (optional)
       GambarC     *multipart.FileHeader  `form:"gambar_c"` // File upload (optional)
       GambarD     *multipart.FileHeader  `form:"gambar_d"` // File upload (optional)
       GambarE     *multipart.FileHeader  `form:"gambar_e"` // File upload (optional)
       Kunci       string                 `json:"kunci" form:"kunci" validate:"required,len=1"`
   }
   ```

**File yang Dimodifikasi/Dibuat:**
- `internal/utils/file_upload.go` (baru)
- `internal/modules/soal/dto/soal_dto.go` (update)

---

### **FASE 3: Update Service Layer untuk Handle Upload**

**Tujuan:** Menambahkan logika upload gambar di service layer

**Langkah-langkah:**

1. **Update `soalService.CreateSoal`** untuk handle file upload:
   ```go
   func (s *soalService) CreateSoal(req *dto.CreateSoalRequest) (*dto.SoalResponse, error) {
       if err := s.validateKunci(...); err != nil {
           return nil, err
       }
   
       // Process upload gambar soal
       var gambarSoalName string
       if req.GambarSoal != nil {
           filename, err := utils.SaveImage(req.GambarSoal, "soal")
           if err != nil {
               return nil, errors.New("gagal upload gambar_soal: " + err.Error())
           }
           gambarSoalName = filename
       }
       
       // Process upload gambar opsi A, B, C, D, E dengan cara yang sama
       // ...
   
       soal := &model.Soal{
           IdBankSoal: req.IdBankSoal,
           NoSoal:     req.NoSoal,
           Soal:       req.Soal,
           GambarSoal: gambarSoalName, // Simpan nama file, bukan URL
           OpsiA:      req.OpsiA,
           OpsiB:      req.OpsiB,
           OpsiC:      req.OpsiC,
           OpsiD:      req.OpsiD,
           OpsiE:      req.OpsiE,
           GambarA:    gambarAName,
           GambarB:    gambarBName,
           GambarC:    gambarCName,
           GambarD:    gambarDName,
           GambarE:    gambarEName,
           Kunci:      req.Kunci,
       }
       
       if err := s.repo.Create(soal); err != nil {
           // Jika gagal create, hapus semua file yang sudah diupload
           utils.DeleteImage(gambarSoalName, "soal")
           utils.DeleteImage(gambarAName, "opsi")
           // ... delete yang lain
           return nil, err
       }
       
       return s.modelToResponse(soal), nil
   }
   ```

2. **Update `soalService.UpdateSoal`** untuk handle file upload:
   ```go
   func (s *soalService) UpdateSoal(id string, req *dto.UpdateSoalRequest) (*dto.SoalResponse, error) {
       soal, err := s.repo.GetByID(id)
       if err != nil {
           return nil, errors.New(constants.ErrNotFound)
       }
       
       // Handle update gambar_soal
       if req.GambarSoal != nil {
           newFilename, err := utils.SaveImage(req.GambarSoal, "soal")
           if err != nil {
               return nil, errors.New("gagal upload gambar_soal: " + err.Error())
           }
           // Hapus gambar lama
           if soal.GambarSoal != "" {
               utils.DeleteImage(soal.GambarSoal, "soal")
           }
           soal.GambarSoal = newFilename
       }
       
       // Handle update gambar opsi dengan cara yang sama
       // ...
       
       soal.NoSoal = req.NoSoal
       soal.Soal = req.Soal
       soal.OpsiA = req.OpsiA
       soal.OpsiB = req.OpsiB
       soal.OpsiC = req.OpsiC
       soal.OpsiD = req.OpsiD
       soal.OpsiE = req.OpsiE
       soal.Kunci = req.Kunci
       
       if err := s.repo.Update(soal); err != nil {
           return nil, err
       }
       
       return s.modelToResponse(soal), nil
   }
   ```

3. **Update `soalService.DeleteSoal`** untuk hapus file gambar:
   ```go
   func (s *soalService) DeleteSoal(id string) error {
       soal, err := s.repo.GetByID(id)
       if err != nil {
           return errors.New(constants.ErrNotFound)
       }
       
       // Hapus semua file gambar
       utils.DeleteImage(soal.GambarSoal, "soal")
       utils.DeleteImage(soal.GambarA, "opsi")
       utils.DeleteImage(soal.GambarB, "opsi")
       utils.DeleteImage(soal.GambarC, "opsi")
       utils.DeleteImage(soal.GambarD, "opsi")
       utils.DeleteImage(soal.GambarE, "opsi")
       
       return s.repo.Delete(soal.ID)
   }
   ```

4. **Update `modelToResponse`** untuk include image URLs:
   ```go
   func (s *soalService) modelToResponse(soal *model.Soal) *dto.SoalResponse {
       cfg := config.GetUploadConfig()
       
       return &dto.SoalResponse{
           ID:         soal.ID,
           IdBankSoal: soal.IdBankSoal,
           NoSoal:     soal.NoSoal,
           Soal:       soal.Soal,
           GambarSoal: buildImageURL(soal.GambarSoal, "soal", cfg.ImageBaseURL),
           OpsiA:      soal.OpsiA,
           OpsiB:      soal.OpsiB,
           OpsiC:      soal.OpsiC,
           OpsiD:      soal.OpsiD,
           OpsiE:      soal.OpsiE,
           GambarA:    buildImageURL(soal.GambarA, "opsi", cfg.ImageBaseURL),
           GambarB:    buildImageURL(soal.GambarB, "opsi", cfg.ImageBaseURL),
           GambarC:    buildImageURL(soal.GambarC, "opsi", cfg.ImageBaseURL),
           GambarD:    buildImageURL(soal.GambarD, "opsi", cfg.ImageBaseURL),
           GambarE:    buildImageURL(soal.GambarE, "opsi", cfg.ImageBaseURL),
           Kunci:      soal.Kunci,
           CreatedAt:  soal.CreatedAt.Format("2006-01-02 15:04:05"),
           UpdatedAt:  soal.UpdatedAt.Format("2006-01-02 15:04:05"),
       }
   }
   
   func buildImageURL(filename, folder, baseURL string) string {
       if filename == "" {
           return ""
       }
       return fmt.Sprintf("%s/%s/%s", baseURL, folder, filename)
   }
   ```

**File yang Dimodifikasi:**
- `internal/modules/soal/service/soal_service.go` (update)

---

### **FASE 4: Update Controller untuk Menerima Multipart Form**

**Tujuan:** Mengubah controller untuk menerima multipart/form-data bukan hanya JSON

**Langkah-langkah:**

1. **Update `soalController.CreateSoal`**:
   ```go
   func (c *soalController) CreateSoal(ctx *fiber.Ctx) error {
       req := new(dto.CreateSoalRequest)
       
       // Parse multipart form (untuk file upload)
       form, err := ctx.MultipartForm()
       if err != nil {
           return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
               "success": false,
               "message": "Invalid request format",
           })
       }
       
       // Extract form fields
       req.IdBankSoal = ctx.FormValue("id_bank_soal")
       req.Soal = ctx.FormValue("soal")
       req.OpsiA = ctx.FormValue("opsi_a")
       req.OpsiB = ctx.FormValue("opsi_b")
       req.OpsiC = ctx.FormValue("opsi_c")
       req.OpsiD = ctx.FormValue("opsi_d")
       req.OpsiE = ctx.FormValue("opsi_e")
       req.Kunci = ctx.FormValue("kunci")
       
       // Parse NoSoal
       noSoalStr := ctx.FormValue("no_soal")
       if noSoal, err := strconv.Atoi(noSoalStr); err == nil {
           req.NoSoal = noSoal
       }
       
       // Extract file upload
       if file, err := ctx.FormFile("gambar_soal"); err == nil {
           req.GambarSoal = file
       }
       // Extract file gambar opsi A, B, C, D, E dengan cara yang sama
       // ...
       
       // Validate request
       if err := c.validate.Struct(req); err != nil {
           return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
               "success": false,
               "message": "Validation error",
               "errors": err.Error(),
           })
       }
       
       resp, err := c.service.CreateSoal(req)
       if err != nil {
           return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
               "success": false,
               "message": err.Error(),
           })
       }
       
       return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
           "success": true,
           "message": "Create soal successfully",
           "data": resp,
           "errors": nil,
       })
   }
   ```

2. **Update `soalController.UpdateSoal`** dengan cara yang sama untuk multipart form

**File yang Dimodifikasi:**
- `internal/modules/soal/controller/soal_controller.go` (update)

---

### **FASE 5: Setup Static File Serving**

**Tujuan:** Membuat endpoint untuk serve file gambar yang sudah diupload

**Langkah-langkah:**

1. **Update `cmd/server/main.go`** untuk serve uploads folder:
   ```go
   func setupRoutes(app *fiber.App) {
       app.Get("/health", func(ctx *fiber.Ctx) error {
           return ctx.JSON(fiber.Map{
               "status": "ok",
               "service": "Fiber Backend API",
           })
       })
       
       // Serve static files dari folder uploads
       app.Static("/uploads", "./uploads")
       
       // Setup routes modules
       authroutes.SetupAuthRoutes(app, database.DB)
       // ... routes lainnya
   }
   ```

2. **Pastikan folder uploads accessible** dan memiliki proper permissions

**File yang Dimodifikasi:**
- `cmd/server/main.go` (update)

---

### **FASE 6: Update API Documentation**

**Tujuan:** Dokumentasi endpoint untuk upload gambar

**Langkah-langkah:**

1. **Update `docs/SOAL_API.md`** dengan contoh multipart form:

```markdown
### POST - Buat Soal Baru (dengan Upload Gambar)

**Endpoint:** `POST /api/soal`

**Content-Type:** `multipart/form-data` (bukan application/json)

**Form Fields:**
- `id_bank_soal` (string, required) - UUID bank soal
- `no_soal` (integer, required) - Nomor soal, min 1
- `soal` (string, required) - Pertanyaan soal
- `gambar_soal` (file, optional) - Gambar pertanyaan (jpg, jpeg, png, gif, webp, max 5MB)
- `opsi_a` (string, required) - Opsi A
- `opsi_b` (string, required) - Opsi B
- `opsi_c` (string, required) - Opsi C
- `opsi_d` (string, optional) - Opsi D
- `opsi_e` (string, optional) - Opsi E
- `gambar_a` (file, optional) - Gambar opsi A (max 5MB)
- `gambar_b` (file, optional) - Gambar opsi B (max 5MB)
- `gambar_c` (file, optional) - Gambar opsi C (max 5MB)
- `gambar_d` (file, optional) - Gambar opsi D (max 5MB)
- `gambar_e` (file, optional) - Gambar opsi E (max 5MB)
- `kunci` (string, required) - Jawaban benar (A/B/C/D/E)

**Request Example (using curl):**
\`\`\`bash
curl -X POST "http://localhost:3000/api/soal" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -F "id_bank_soal=5112e444-25d8-4ca6-859f-3d24099f45ce" \
  -F "no_soal=1" \
  -F "soal=Berapa hasil dari 2 + 2?" \
  -F "gambar_soal=@/path/to/soal.jpg" \
  -F "opsi_a=3" \
  -F "opsi_b=4" \
  -F "opsi_c=5" \
  -F "gambar_a=@/path/to/a.jpg" \
  -F "gambar_b=@/path/to/b.jpg" \
  -F "gambar_c=@/path/to/c.jpg" \
  -F "kunci=B"
\`\`\`

**Success Response (201 Created):**
\`\`\`json
{
  "success": true,
  "message": "Create soal successfully",
  "data": {
    "id": "abc123def456",
    "id_bank_soal": "5112e444-25d8-4ca6-859f-3d24099f45ce",
    "no_soal": 1,
    "soal": "Berapa hasil dari 2 + 2?",
    "gambar_soal": "http://localhost:3000/uploads/soal/1684426848_abc12345.jpg",
    "opsi_a": "3",
    "opsi_b": "4",
    "opsi_c": "5",
    "opsi_d": "",
    "opsi_e": "",
    "gambar_a": "http://localhost:3000/uploads/soal/1684426848_def67890.jpg",
    "gambar_b": "http://localhost:3000/uploads/soal/1684426849_ghi11111.jpg",
    "gambar_c": "http://localhost:3000/uploads/soal/1684426850_jkl22222.jpg",
    "gambar_d": "",
    "gambar_e": "",
    "kunci": "B",
    "created_at": "2026-05-18 14:00:00",
    "updated_at": "2026-05-18 14:00:00"
  },
  "errors": null
}
\`\`\`
```

**File yang Dimodifikasi:**
- `docs/SOAL_API.md` (update)

---

### **FASE 7: Testing Image Upload Functionality**

**Tujuan:** Test semua fitur upload gambar

**Langkah-langkah:**

1. **Test POST dengan gambar soal saja**
   - Create soal dengan hanya upload gambar_soal
   - Verify file tersimpan di folder
   - Verify response berisi URL gambar

2. **Test POST dengan semua gambar**
   - Create soal dengan upload semua gambar
   - Verify semua file tersimpan
   - Verify response berisi semua URL gambar

3. **Test POST validation**
   - Test upload file terlalu besar (>5MB)
   - Test upload format tidak diizinkan (misalnya .pdf)
   - Verify error message yang sesuai

4. **Test PUT dengan update gambar**
   - Update soal dengan gambar baru
   - Verify gambar lama dihapus
   - Verify gambar baru tersimpan

5. **Test PUT partial update**
   - Update soal dengan hanya update gambar_soal (gambar lain tidak diubah)
   - Verify hanya gambar_soal yang berubah

6. **Test DELETE**
   - Delete soal dengan gambar
   - Verify semua gambar terhapus dari folder

7. **Test GET endpoints**
   - Verify GET returns image URLs
   - Test image URL dapat di-access

**Test Script Example:**
```bash
# Login untuk mendapat JWT token
TOKEN=$(curl -s -X POST "http://localhost:3000/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' | jq -r '.data.token')

# Test POST dengan file
curl -X POST "http://localhost:3000/api/soal" \
  -H "Authorization: Bearer $TOKEN" \
  -F "id_bank_soal=<bank_soal_id>" \
  -F "no_soal=1" \
  -F "soal=Test Question?" \
  -F "opsi_a=A" \
  -F "opsi_b=B" \
  -F "opsi_c=C" \
  -F "gambar_soal=@test-image.jpg" \
  -F "kunci=B" | jq .
```

---

### **FASE 8: Cleanup dan Error Handling**

**Tujuan:** Handle edge cases dan cleanup file

**Langkah-langkah:**

1. **Implementasi hard delete untuk gambar**
   - Saat repository hard delete soal, semua gambar harus dihapus permanent
   - Update `SoalRepository.HardDelete` method

2. **Handle file cleanup saat server error**
   - Jika database insert fail setelah upload, semua file harus dihapus
   - Implementasi sudah di service layer (lihat Fase 3)

3. **Implement file cleanup background job** (optional untuk production)
   - Delete orphaned files yang tidak ada di database
   - Schedule job untuk clear upload folder secara periodik

4. **Add logging**
   - Log semua upload activity (success/failure)
   - Log file deletion
   - Helpful untuk debugging

5. **Add metrics/monitoring** (optional)
   - Track upload success rate
   - Track file storage usage
   - Monitor disk space

6. **Update migration documentation**
   - Document folder structure requirements
   - Document file permission requirements
   - Document disk space requirements

---

## 📊 Summary Perubahan File

| File | Status | Keterangan |
|------|--------|-----------|
| `internal/config/upload.go` | NEW | Konfigurasi upload |
| `internal/utils/file_upload.go` | NEW | Utility untuk handle file |
| `internal/modules/soal/dto/soal_dto.go` | UPDATE | Add multipart.FileHeader fields |
| `internal/modules/soal/service/soal_service.go` | UPDATE | Add upload logic |
| `internal/modules/soal/controller/soal_controller.go` | UPDATE | Handle multipart form |
| `cmd/server/main.go` | UPDATE | Serve static files |
| `docs/SOAL_API.md` | UPDATE | Document upload endpoints |
| `uploads/soal/` | NEW | Folder untuk menyimpan gambar |

---

## ✅ Checklist Implementasi

- [ ] Fase 1: Setup konfigurasi upload
- [ ] Fase 2: Buat utility untuk upload
- [ ] Fase 3: Update service layer
- [ ] Fase 4: Update controller
- [ ] Fase 5: Setup static file serving
- [ ] Fase 6: Update dokumentasi API
- [ ] Fase 7: Testing fitur upload
- [ ] Fase 8: Cleanup dan error handling

---

## 🔍 Key Points untuk Diingat

1. **Simpan nama file, bukan full path** di database
2. **Validate extension dan size** di utility function
3. **Generate unique filename** untuk menghindari conflict
4. **Hapus file lama** saat update atau delete
5. **Rollback file** jika database transaction gagal
6. **Serve static files** melalui endpoint `/uploads`
7. **Return full URL** di response, bukan hanya filename
8. **Handle error gracefully** dengan pesan yang jelas

---

Generated with Claude Code 🤖
