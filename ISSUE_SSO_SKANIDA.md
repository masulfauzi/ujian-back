# Issue: Endpoint login otomatis (SSO) dari aplikasi Skanida Student

## Latar belakang
Aplikasi Android **Skanida Student** membuka situs ujian di WebView (menu Ujian Online). Siswa sudah login di aplikasi
(akun sianida), jadi siswa harus langsung masuk ke dashboard tanpa mengetik password.

Alurnya:
```
Aplikasi  -> sianida POST /api/ujian/sso-code           -> code (acak, 60 detik, sekali pakai)
WebView   -> https://ujian.smkn2semarang.sch.id/sso?code=XXX       (ujian-front, lihat issue-sso.md)
ujian-front -> ujian-back POST /api/auth/sso { code }              <- ENDPOINT INI
ujian-back  -> sianida POST /api/ujian/sso-verify { code } + header X-SSO-Secret -> { username, name }
ujian-back  -> cari peserta by username -> JWT peserta (sama seperti /auth/login)
```

`username` dari sianida = `peserta.username` (login peserta saat ini juga divalidasi ke sianida dengan username yang sama,
lihat `validateWithExternalAPI` di `internal/modules/auth/service/auth_service.go`).

## Yang dikerjakan

### 1. Konfigurasi — `configs/sso.go` (baru)
```go
package configs

type SSOConfig struct {
	VerifyURL string
	Secret    string
}

func GetSSOConfig() *SSOConfig {
	return &SSOConfig{
		VerifyURL: getEnv("SIANIDA_SSO_VERIFY_URL", "https://apps.smkn2semarang.sch.id/api/ujian/sso-verify"),
		Secret:    getEnv("SIANIDA_SSO_SECRET", ""),
	}
}
```
`.env` (jangan di-commit), nilainya **sama** dengan `UJIAN_SSO_SECRET` di sianida:
```env
SIANIDA_SSO_SECRET=...
```
Tambahkan juga ke `.env.example` (tanpa nilai) bila ada.

### 2. DTO — `internal/modules/auth/dto/auth_dto.go`
```go
type SSOLoginRequest struct {
	Code string `json:"code" validate:"required,max=100"`
}
```
Validator: tambah `ValidateSSOLogin` di `internal/modules/auth/validator/auth_validator.go` (sama seperti `ValidateLogin`).

### 3. Service — `internal/modules/auth/service/auth_service.go`
- Tambah `LoginSSO(code string) (*dto.AuthResponse, error)` ke interface `AuthService`.
- Refactor: pindahkan pembuatan `tempUser` + `GenerateToken` + `AuthResponse` untuk peserta di `Login` ke helper
  `issuePesertaToken(peserta *pesertamodel.Peserta) (*dto.AuthResponse, error)` dan pakai di `Login` & `LoginSSO`
  (agar isi JWT identik).

```go
var (
	ErrSSOInvalid       = errors.New("Sesi login otomatis tidak valid atau kedaluwarsa")
	ErrSSONotRegistered = errors.New("Anda tidak terdaftar sebagai peserta ujian")
)

func (s *authService) LoginSSO(code string) (*dto.AuthResponse, error) {
	cfg := configs.GetSSOConfig()
	if cfg.Secret == "" {
		return nil, errors.New("SSO belum dikonfigurasi")
	}

	body, _ := json.Marshal(map[string]string{"code": code})
	req, err := http.NewRequest(http.MethodPost, cfg.VerifyURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-SSO-Secret", cfg.Secret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrSSOInvalid // 404 = kode salah/kedaluwarsa/sudah dipakai, 401 = secret salah
	}

	var result struct {
		Username string `json:"username"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.Username == "" {
		return nil, ErrSSOInvalid
	}

	peserta, err := s.repo.GetPesertaByUsername(result.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSSONotRegistered
		}
		return nil, err
	}

	return s.issuePesertaToken(peserta)
}
```
Catatan: secret salah di sisi server juga berujung 401 ke siswa; log status code sianida (tanpa kode/secret) agar mudah dilacak.

### 4. Controller — `internal/modules/auth/controller/auth_controller.go`
```go
// SSOLogin godoc
// @Summary Login otomatis peserta dari aplikasi Skanida Student
// @Description Menukar kode SSO sekali pakai (diterbitkan sianida) menjadi JWT peserta.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SSOLoginRequest true "Kode SSO"
// @Success 200 {object} helpers.Response{data=dto.AuthResponse}
// @Failure 401 {object} helpers.Response
// @Router /auth/sso [post]
func (c *AuthController) SSOLogin(ctx *fiber.Ctx) error {
	var req dto.SSOLoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}
	if err := validator.ValidateSSOLogin(req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Validation error", nil)
	}

	resp, err := c.authService.LoginSSO(req.Code)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, err.Error(), nil)
	}
	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Login successfully", resp)
}
```

### 5. Route — `internal/modules/auth/routes/auth_routes.go`
```go
auth.Post("/sso", ctrl.SSOLogin)
```
Bila ada middleware rate limit (fiber `limiter`), pasang di route ini (mis. 20 request/menit per IP).
Regenerate swagger (`swag init` sesuai Makefile).

## Pengujian
1. Ambil kode dari sianida (lihat `sianida/issue-ujian-sso.md`).
2. ```bash
   curl -s -X POST https://ujian.smkn2semarang.sch.id/api/auth/sso \
     -H "Content-Type: application/json" -d '{"code":"<CODE>"}'
   ```
   → 200, `data.token` JWT dengan `role: "peserta"`; `GET /api/auth/me` dengan token itu mengembalikan peserta yang benar.
3. Kode yang sama dipakai lagi → 401 "Sesi login otomatis tidak valid...".
4. Siswa yang tidak ada di tabel `peserta` → 401 "Anda tidak terdaftar sebagai peserta ujian".
5. `/api/auth/login` biasa tetap berfungsi (regresi refactor `issuePesertaToken`).
