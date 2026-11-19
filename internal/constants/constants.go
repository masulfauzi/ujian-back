package constants

const (
	// User Roles
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RoleGuest  = "guest"

	// Status
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusDeleted  = "deleted"
)

// Error Messages
const (
	ErrUnauthorized     = "Unauthorized"
	ErrForbidden        = "Forbidden"
	ErrNotFound         = "Resource not found"
	ErrValidation       = "Validation error"
	ErrInternalServer   = "Internal server error"
	ErrInvalidEmail     = "Invalid email format"
	ErrPasswordTooShort = "Password must be at least 6 characters"
	ErrDuplicateEmail   = "Email already registered"
)

// Success Messages
const (
	SuccessRegister = "Register successfully"
	SuccessLogin    = "Login successfully"
	SuccessGetUser  = "Get user successfully"
	SuccessCreate   = "Create successfully"
	SuccessUpdate   = "Update successfully"
	SuccessDelete   = "Delete successfully"
	SuccessGetAll   = "Get all data successfully"
)

// Jenis Soal
const (
	JenisSoalPG  = "pg"  // pilihan ganda biasa, kunci 1 huruf A-E
	JenisSoalPGK = "pgk" // pilihan ganda kompleks, kunci >= 2 huruf A-E (mis. "AC")
	JenisSoalBS  = "bs"  // benar/salah, 3 pernyataan di opsi_a-c, kunci 3 karakter B/S (mis. "BSB")
)
