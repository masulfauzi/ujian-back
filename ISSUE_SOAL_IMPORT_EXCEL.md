# ISSUE: Implementasi Import Soal dari Excel

## 📋 Deskripsi Fitur

Membuat endpoint API baru untuk import/bulk insert soal dari file Excel ke dalam tabel `soal`. Endpoint ini akan membaca file Excel dan melakukan mapping kolom Excel ke field database sesuai dengan struktur yang telah ditentukan.

### Tujuan Implementasi:
1. Menerima file Excel (.xls, .xlsx) melalui HTTP multipart/form-data
2. Membaca dan parse kolom Excel sesuai mapping yang ditentukan
3. Validasi data sebelum insert ke database
4. Bulk insert data soal ke tabel `soal` dengan `id_bank_soal` yang sesuai
5. Menangani error dengan response yang informatif
6. Mengembalikan laporan insert (success count, failed count, error details)

---

## 📊 Mapping Kolom Excel → Database

| Kolom Excel | Column | Field Database | Tipe Data | Keterangan |
|-------------|--------|----------------|-----------|------------|
| A | 1 | `no_soal` | Integer | Nomor urut soal |
| B | 2 | `soal` | Text | Pertanyaan soal |
| C | 3 | `opsi_a` | Text | Opsi jawaban A |
| D | 4 | `opsi_b` | Text | Opsi jawaban B |
| E | 5 | `opsi_c` | Text | Opsi jawaban C |
| F | 6 | `opsi_d` | Text | Opsi jawaban D |
| G | 7 | `opsi_e` | Text | Opsi jawaban E |
| H | 8 | `kunci` | VARCHAR(1) | Kunci jawaban (A/B/C/D/E) |
| I | 9 | - | - | **TIDAK DIGUNAKAN** |
| J | 10 | `gambar_soal` | VARCHAR(500) | Nama file gambar soal |
| K | 11 | - | - | **TIDAK DIGUNAKAN** |
| L | 12 | `gambar_a` | VARCHAR(500) | Nama file gambar opsi A |
| M | 13 | `gambar_b` | VARCHAR(500) | Nama file gambar opsi B |
| N | 14 | `gambar_c` | VARCHAR(500) | Nama file gambar opsi C |
| O | 15 | `gambar_d` | VARCHAR(500) | Nama file gambar opsi D |
| P | 16 | `gambar_e` | VARCHAR(500) | Nama file gambar opsi E |

---

## 🔌 API Endpoint Specification

### **POST /api/bank-soal/import**

#### Request Format:

**Content-Type:** `multipart/form-data`

```
POST /api/bank-soal/import
Content-Type: multipart/form-data

id_bank_soal: "550e8400-e29b-41d4-a716-446655440000"
file: [Excel file binary]
```

**Request Parameters:**
- `id_bank_soal` (string, required): UUID dari bank soal tujuan
- `file` (file, required): File Excel dengan format .xls atau .xlsx

---

#### Response Format:

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Import soal berhasil",
  "data": {
    "total_processed": 50,
    "total_success": 48,
    "total_failed": 2,
    "import_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2026-05-19T10:30:00Z",
    "summary": {
      "inserted": 48,
      "skipped": 0,
      "errors": 2
    },
    "errors": [
      {
        "row": 5,
        "error": "kunci harus berupa huruf A-E"
      },
      {
        "row": 12,
        "error": "soal tidak boleh kosong"
      }
    ]
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Import soal gagal",
  "data": {
    "error": "File tidak valid atau id_bank_soal tidak ditemukan",
    "details": "File harus berupa Excel (.xls/.xlsx)"
  }
}
```

**Error Response (500 Internal Server Error):**
```json
{
  "success": false,
  "message": "Terjadi kesalahan pada server",
  "data": {
    "error": "Database error atau file system error"
  }
}
```

---

## 🎯 Validasi Data

### Validasi di Level Row/Data:

1. **no_soal** (required)
   - Harus berupa integer
   - Harus > 0
   - Tidak boleh duplikat dalam satu import

2. **soal** (required)
   - Tidak boleh kosong/blank
   - Max length: 5000 karakter

3. **opsi_a, opsi_b, opsi_c, opsi_d, opsi_e** (required)
   - Tidak boleh kosong/blank
   - Max length: 5000 karakter setiap opsi

4. **kunci** (required)
   - Harus berupa single character: A, B, C, D, atau E (case insensitive)
   - Akan di-uppercase otomatis

5. **gambar_* fields** (optional)
   - Menerima string (nama file atau path)
   - Tidak perlu validasi existence di folder
   - Bisa kosong jika tidak ada gambar

### Validasi di Level File:

1. File harus berupa Excel (.xls atau .xlsx)
2. File size max 10 MB
3. Minimal memiliki header row
4. Minimal memiliki 1 data row
5. Bank soal dengan `id_bank_soal` harus exist di database

---

## 🔄 Fase Implementasi (7 Tahap)

### **FASE 1: Setup Dependencies dan Utility**

**Tujuan:** Menambahkan library untuk membaca Excel dan membuat utility functions

**Langkah-langkah:**

1. **Tambah dependency Excel library ke `go.mod`**
   - Gunakan library `github.com/xuri/excelize/v2` untuk membaca Excel file
   - Alternatif: `github.com/360EntSecGroup-Skylar/excelize` (versi lama)
   - Run: `go get github.com/xuri/excelize/v2@latest`

2. **Buat file validation utility** di `internal/utils/excel_validator.go`:
   ```go
   package utils
   
   import (
       "errors"
       "strings"
       "unicode"
   )
   
   // ValidateSoalRow memvalidasi satu row dari excel
   func ValidateSoalRow(row ExcelSoalRow, rowIndex int) []error {
       var errors []error
       
       // Validate no_soal
       if row.NoSoal <= 0 {
           errors = append(errors, "no_soal harus lebih dari 0")
       }
       
       // Validate soal
       if strings.TrimSpace(row.Soal) == "" {
           errors = append(errors, "soal tidak boleh kosong")
       }
       
       // Validate opsi
       if strings.TrimSpace(row.OpsiA) == "" {
           errors = append(errors, "opsi_a tidak boleh kosong")
       }
       // ... validate opsi_b, c, d, e similarly
       
       // Validate kunci
       kunci := strings.ToUpper(strings.TrimSpace(row.Kunci))
       if !isValidKunci(kunci) {
           errors = append(errors, "kunci harus berupa A, B, C, D, atau E")
       }
       
       return errors
   }
   
   func isValidKunci(k string) bool {
       return k == "A" || k == "B" || k == "C" || k == "D" || k == "E"
   }
   ```

3. **Buat file utility untuk excel parsing** di `internal/utils/excel_parser.go`:
   ```go
   package utils
   
   import (
       "fmt"
       "strconv"
       "strings"
   )
   
   type ExcelSoalRow struct {
       RowIndex   int
       NoSoal     int
       Soal       string
       OpsiA      string
       OpsiB      string
       OpsiC      string
       OpsiD      string
       OpsiE      string
       Kunci      string
       GambarSoal string
       GambarA    string
       GambarB    string
       GambarC    string
       GambarD    string
       GambarE    string
   }
   
   // ParseExcelRow mengextract data dari row excel dan convert ke struct
   func ParseExcelRow(values []interface{}, rowIndex int) (*ExcelSoalRow, error) {
       if len(values) < 16 {
           return nil, fmt.Errorf("row harus memiliki minimal 16 kolom")
       }
       
       // Parse no_soal (Column A, index 0)
       noSoal, err := strconv.Atoi(fmt.Sprintf("%v", values[0]))
       if err != nil {
           return nil, fmt.Errorf("no_soal harus berupa number")
       }
       
       return &ExcelSoalRow{
           RowIndex:   rowIndex,
           NoSoal:     noSoal,
           Soal:       toString(values[1]),
           OpsiA:      toString(values[2]),
           OpsiB:      toString(values[3]),
           OpsiC:      toString(values[4]),
           OpsiD:      toString(values[5]),
           OpsiE:      toString(values[6]),
           Kunci:      toUpperString(values[7]),
           // Column I (index 8) skipped
           GambarSoal: toString(values[9]),
           // Column K (index 10) skipped
           GambarA:    toString(values[11]),
           GambarB:    toString(values[12]),
           GambarC:    toString(values[13]),
           GambarD:    toString(values[14]),
           GambarE:    toString(values[15]),
       }, nil
   }
   
   func toString(v interface{}) string {
       return strings.TrimSpace(fmt.Sprintf("%v", v))
   }
   
   func toUpperString(v interface{}) string {
       return strings.ToUpper(strings.TrimSpace(fmt.Sprintf("%v", v)))
   }
   ```

**File yang Dibuat:**
- `internal/utils/excel_validator.go` (baru)
- `internal/utils/excel_parser.go` (baru)

**File yang Dimodifikasi:**
- `go.mod` (tambah dependency)

---

### **FASE 2: Update DTO (Data Transfer Object)**

**Tujuan:** Membuat DTO untuk import request dan response

**Langkah-langkah:**

1. **Update atau buat file** `internal/modules/soal/dto/soal_dto.go`:

Tambahkan struct-struct baru:
```go
// ImportSoalRequest adalah request untuk import soal dari excel
type ImportSoalRequest struct {
    IdBankSoal string `form:"id_bank_soal" validate:"required"`
    File       *multipart.FileHeader `form:"file" validate:"required"`
}

// ImportSoalErrorDetail adalah detail error per row
type ImportSoalErrorDetail struct {
    Row    int    `json:"row"`
    Error  string `json:"error"`
}

// ImportSoalResponse adalah response dari import
type ImportSoalResponse struct {
    TotalProcessed int                         `json:"total_processed"`
    TotalSuccess   int                         `json:"total_success"`
    TotalFailed    int                         `json:"total_failed"`
    ImportID       string                      `json:"import_id"`
    Timestamp      time.Time                   `json:"timestamp"`
    Summary        map[string]int              `json:"summary"`
    Errors         []ImportSoalErrorDetail     `json:"errors"`
}

// BulkCreateSoalRequest untuk batch insert
type BulkCreateSoalRequest struct {
    IdBankSoal string
    NoSoal     int
    Soal       string
    OpsiA      string
    OpsiB      string
    OpsiC      string
    OpsiD      string
    OpsiE      string
    Kunci      string
    GambarSoal string
    GambarA    string
    GambarB    string
    GambarC    string
    GambarD    string
    GambarE    string
}
```

**File yang Dimodifikasi:**
- `internal/modules/soal/dto/soal_dto.go`

---

### **FASE 3: Update Repository Layer**

**Tujuan:** Menambahkan method untuk bulk insert soal

**Langkah-langkah:**

1. **Update file** `internal/modules/soal/repository/soal_repository.go`:

Tambahkan method baru:
```go
// BulkCreateSoal melakukan bulk insert multiple soal
func (r *SoalRepository) BulkCreateSoal(ctx context.Context, soals []model.Soal) error {
    return r.db.WithContext(ctx).CreateInBatches(soals, 100).Error
    // Note: CreateInBatches secara otomatis batch insert dengan batch size 100
}

// GetBankSoalByID untuk validasi bank_soal existence
func (r *SoalRepository) GetBankSoalExists(ctx context.Context, bankSoalID string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&model.BankSoal{}).
        Where("id = ?", bankSoalID).
        Count(&count).Error
    
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
```

**File yang Dimodifikasi:**
- `internal/modules/soal/repository/soal_repository.go`

---

### **FASE 4: Update Service Layer**

**Tujuan:** Membuat business logic untuk import soal dari excel

**Langkah-langkah:**

1. **Update file** `internal/modules/soal/service/soal_service.go`:

Tambahkan method baru:
```go
// ImportSoalFromExcel melakukan import soal dari file excel
func (s *SoalService) ImportSoalFromExcel(ctx context.Context, req *dto.ImportSoalRequest) (*dto.ImportSoalResponse, error) {
    // 1. Validasi bank_soal exists
    exists, err := s.repo.GetBankSoalExists(ctx, req.IdBankSoal)
    if err != nil || !exists {
        return nil, errors.New("bank_soal tidak ditemukan")
    }
    
    // 2. Buka file dari request
    file, err := req.File.Open()
    if err != nil {
        return nil, errors.New("gagal membuka file")
    }
    defer file.Close()
    
    // 3. Parse excel file
    xlsx, err := excelize.OpenReader(file)
    if err != nil {
        return nil, errors.New("file bukan format excel yang valid")
    }
    defer xlsx.Close()
    
    // 4. Get sheet pertama
    sheetName := xlsx.GetSheetName(0)
    rows, err := xlsx.GetRows(sheetName)
    if err != nil {
        return nil, errors.New("gagal membaca sheet excel")
    }
    
    // 5. Parse dan validasi data
    var soals []model.Soal
    var errorDetails []dto.ImportSoalErrorDetail
    var processedCount, successCount, failedCount int
    
    // Skip header row (start dari index 1)
    for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
        row := rows[rowIndex]
        processedCount++
        
        // Parse row
        excelRow, err := utils.ParseExcelRow(row, rowIndex+1)
        if err != nil {
            failedCount++
            errorDetails = append(errorDetails, dto.ImportSoalErrorDetail{
                Row:   rowIndex + 1,
                Error: err.Error(),
            })
            continue
        }
        
        // Validate row
        validationErrors := utils.ValidateSoalRow(*excelRow, rowIndex+1)
        if len(validationErrors) > 0 {
            failedCount++
            errorDetails = append(errorDetails, dto.ImportSoalErrorDetail{
                Row:   rowIndex + 1,
                Error: strings.Join(validationErrors, "; "),
            })
            continue
        }
        
        // Create soal model
        soal := model.Soal{
            IdBankSoal: req.IdBankSoal,
            NoSoal:     excelRow.NoSoal,
            Soal:       excelRow.Soal,
            OpsiA:      excelRow.OpsiA,
            OpsiB:      excelRow.OpsiB,
            OpsiC:      excelRow.OpsiC,
            OpsiD:      excelRow.OpsiD,
            OpsiE:      excelRow.OpsiE,
            Kunci:      excelRow.Kunci,
            GambarSoal: excelRow.GambarSoal,
            GambarA:    excelRow.GambarA,
            GambarB:    excelRow.GambarB,
            GambarC:    excelRow.GambarC,
            GambarD:    excelRow.GambarD,
            GambarE:    excelRow.GambarE,
        }
        
        soals = append(soals, soal)
        successCount++
    }
    
    // 6. Bulk insert ke database
    if len(soals) > 0 {
        err = s.repo.BulkCreateSoal(ctx, soals)
        if err != nil {
            return nil, errors.New("gagal menyimpan data ke database: " + err.Error())
        }
    }
    
    // 7. Return response
    return &dto.ImportSoalResponse{
        TotalProcessed: processedCount,
        TotalSuccess:   successCount,
        TotalFailed:    failedCount,
        ImportID:       req.IdBankSoal,
        Timestamp:      time.Now(),
        Summary: map[string]int{
            "inserted": successCount,
            "skipped":  0,
            "errors":   failedCount,
        },
        Errors: errorDetails,
    }, nil
}
```

**Catatan Implementasi:**
- Gunakan context untuk handle cancellation
- Error handling harus informatif dan user-friendly
- Limit error details yang ditampilkan (max 100 errors misalnya)
- Gunakan transaction jika ingin atomic operation (all or nothing)

**File yang Dimodifikasi:**
- `internal/modules/soal/service/soal_service.go`

---

### **FASE 5: Update Controller Layer**

**Tujuan:** Membuat HTTP handler untuk endpoint import

**Langkah-langkah:**

1. **Update file** `internal/modules/soal/controller/soal_controller.go`:

Tambahkan method baru di struct `SoalController`:
```go
// ImportSoalFromExcel handle import soal dari excel
func (c *SoalController) ImportSoalFromExcel(ctx *fiber.Ctx) error {
    // 1. Parse multipart form
    file, err := ctx.FormFile("file")
    if err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "File tidak ditemukan", map[string]string{
            "error": "Silakan upload file excel",
        })
    }
    
    // 2. Validate file size (max 10MB)
    const maxFileSize = 10 * 1024 * 1024
    if file.Size > maxFileSize {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "File terlalu besar", map[string]string{
            "error": "Max file size adalah 10MB",
        })
    }
    
    // 3. Validate file extension
    ext := filepath.Ext(file.Filename)
    if ext != ".xls" && ext != ".xlsx" {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Format file tidak valid", map[string]string{
            "error": "File harus berupa .xls atau .xlsx",
        })
    }
    
    // 4. Build request
    req := &dto.ImportSoalRequest{
        IdBankSoal: ctx.FormValue("id_bank_soal"),
        File:       file,
    }
    
    // 5. Validate required fields
    if req.IdBankSoal == "" {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_bank_soal tidak ditemukan", nil)
    }
    
    // 6. Call service
    resp, err := c.service.ImportSoalFromExcel(ctx.Context(), req)
    if err != nil {
        return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Import soal gagal", map[string]string{
            "error": err.Error(),
        })
    }
    
    // 7. Return success response
    return helpers.SuccessResponse(ctx, fiber.StatusOK, "Import soal berhasil", resp)
}
```

**File yang Dimodifikasi:**
- `internal/modules/soal/controller/soal_controller.go`

---

### **FASE 6: Update Routes**

**Tujuan:** Menambahkan route untuk endpoint import

**Langkah-langkah:**

1. **Update file** `internal/modules/soal/routes/soal_routes.go`:

Tambahkan route baru (pastikan menggunakan middleware yang sesuai):
```go
// Di dalam function SetupSoalRoutes, tambahkan:
soal.Post("/import", middleware.JWTAuth(), ctrl.ImportSoalFromExcel)
```

Struktur routes yang seharusnya:
```go
func SetupSoalRoutes(app *fiber.App, db *gorm.DB) {
    repo := repository.NewSoalRepository(db)
    svc := service.NewSoalService(repo)
    ctrl := controller.NewSoalController(svc)
    
    api := app.Group("/api")
    soal := api.Group("/soal")
    
    // POST routes
    soal.Post("/", middleware.JWTAuth(), ctrl.CreateSoal)
    soal.Post("/import", middleware.JWTAuth(), ctrl.ImportSoalFromExcel) // ← TAMBAHKAN INI
    
    // GET routes
    soal.Get("/", ctrl.GetAllSoal)
    soal.Get("/bank/:bank_soal_id", ctrl.GetSoalByBankSoal)
    soal.Get("/:id", ctrl.GetSoalByID)
    
    // PUT routes
    soal.Put("/:id", middleware.JWTAuth(), ctrl.UpdateSoal)
    
    // DELETE routes
    soal.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteSoal)
    soal.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreSoal)
}
```

**File yang Dimodifikasi:**
- `internal/modules/soal/routes/soal_routes.go`

---

### **FASE 7: Testing & Documentation**

**Tujuan:** Test endpoint dan dokumentasi API

**Langkah-langkah:**

1. **Test menggunakan curl atau Postman**:
   ```bash
   curl -X POST http://localhost:8000/api/soal/import \
     -F "id_bank_soal=550e8400-e29b-41d4-a716-446655440000" \
     -F "file=@/Applications/XAMPP/xamppfiles/htdocs/UJIAN-NEW/ujian-back/tmp/upload_excel.xls" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

2. **Test cases yang harus dicek**:
   - ✅ Import file excel valid dengan 50 soal (verify semua berhasil insert)
   - ✅ Import file dengan beberapa error di row tertentu (verify error details)
   - ✅ Import file kosong/hanya header (verify handled gracefully)
   - ✅ Import dengan id_bank_soal tidak exist (verify error message)
   - ✅ Import file non-excel (verify rejected)
   - ✅ Import file > 10MB (verify size validation)
   - ✅ Verify database records created dengan fields yang benar

3. **Dokumentasi**:
   - Update `docs/api.md` atau API documentation dengan endpoint baru
   - Tambahkan contoh request/response
   - Dokumentasikan validation rules

**File yang Perlu Dicek/Update:**
- `docs/` folder
- Test file (optional, jika ada)

---

## 📝 Contoh File Excel (Reference)

File yang ada di `tmp/upload_excel.xls` memiliki struktur:

| A | B | C | D | E | F | G | H | I | J | K | L | M | N | O | P |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| 1 | Apa ibu kota Indonesia? | Jakarta | Bandung | Medan | Surabaya | Yogyakarta | A | - | - | - | - | - | - | - | - |
| 2 | Berapa 2+2? | 4 | 5 | 6 | 7 | 8 | A | - | - | - | - | - | - | - | - |
| ... | ... | ... | ... | ... | ... | ... | ... | ... | ... | ... | ... | ... | ... | ... | ... | ... |

---

## 🚀 Checklist Implementasi

Gunakan checklist ini untuk track progress implementasi:

### Persiapan:
- [ ] Pahami struktur project dan database schema
- [ ] Review file excel contoh di `tmp/upload_excel.xls`
- [ ] Siapkan test data/excel file untuk testing

### Coding:
- [ ] FASE 1: Setup dependencies dan utils
  - [ ] Add excel library ke go.mod
  - [ ] Create `internal/utils/excel_validator.go`
  - [ ] Create `internal/utils/excel_parser.go`
  - [ ] Test parsing dan validation logic

- [ ] FASE 2: Update DTO
  - [ ] Add structs ke `soal_dto.go`
  - [ ] Verify field names match dengan response spec

- [ ] FASE 3: Update Repository
  - [ ] Add `BulkCreateSoal` method
  - [ ] Add `GetBankSoalExists` method
  - [ ] Test repository methods

- [ ] FASE 4: Update Service
  - [ ] Add `ImportSoalFromExcel` method
  - [ ] Implement full logic (parse, validate, insert)
  - [ ] Test dengan berbagai scenario

- [ ] FASE 5: Update Controller
  - [ ] Add `ImportSoalFromExcel` handler
  - [ ] Implement file validation
  - [ ] Test HTTP endpoint

- [ ] FASE 6: Update Routes
  - [ ] Add POST /api/soal/import route
  - [ ] Verify route accessible via HTTP

- [ ] FASE 7: Testing
  - [ ] Test dengan excel valid
  - [ ] Test dengan excel yang ada errors
  - [ ] Test dengan edge cases
  - [ ] Verify database records
  - [ ] Test dengan postman/curl

### QA:
- [ ] Semua test cases passed
- [ ] Response format sesuai spec
- [ ] Error handling informatif
- [ ] No SQL injection atau security issues
- [ ] Code style sesuai project convention

---

## 💡 Tips Implementasi

1. **Mulai dari FASE 1**: Setup utilities dulu, test parsing & validation logic secara isolated
2. **Test incrementally**: Test setiap phase sebelum lanjut ke phase berikutnya
3. **Gunakan context**: Untuk cancellation dan timeout handling
4. **Error handling**: Selalu return informative error messages
5. **Database transaction**: Consider menggunakan transaction untuk atomic operation
6. **Batch size**: Gunakan batch insert (CreateInBatches) untuk performa, jangan insert satu-satu
7. **Logging**: Add logging untuk debugging, terutama di parsing dan insert step
8. **Performance**: Jika file sangat besar (10K+ rows), consider background job/queue

---

## 🔗 Reference Links

- Excel Library Docs: https://pkg.go.dev/github.com/xuri/excelize/v2
- Fiber Framework: https://docs.gofiber.io/
- GORM Batch Operations: https://gorm.io/docs/create.html#Batch-Insert

---

## 📞 Pertanyaan yang Sering Diajukan

**Q: Apakah gambar_* fields bisa kosong?**
A: Ya, fields gambar_* adalah optional. Jika kolom di excel kosong, bisa langsung insert tanpa nilai.

**Q: Apakah perlu validasi bahwa gambar file sudah ada di folder?**
A: Tidak perlu. Cukup simpan string/path yang ada di excel. Tidak perlu validasi existence file di folder uploads.

**Q: Bagaimana jika ada duplicate no_soal dalam satu import?**
A: Validasi bisa reject row dengan no_soal yang duplicate dalam import yang sama. Atau allow saja, tergending business logic.

**Q: Apakah perlu middleware authentication?**
A: Ya, gunakan JWTAuth middleware untuk security. Hanya authenticated user yang bisa import.

**Q: Berapa max error detail yang ditampilkan di response?**
A: Limit ke max 100 errors untuk response size. Jika lebih dari 100 errors, tampilkan "dan N more errors".

---

## 📋 Additional Notes

- Fields gambar_* hanya menyimpan string (nama file atau path) tanpa validasi existence
- Tidak perlu check folder uploads saat import
- Consider implementing retry logic atau background job untuk file processing besar
- Monitor database performance saat bulk insert besar (>10K rows)
- Log semua import activity untuk audit trail
