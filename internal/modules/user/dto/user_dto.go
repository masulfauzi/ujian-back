package dto

import "time"

// CreateUserRequest mendaftarkan akun admin/guru.
// Password kosong = login memakai username & password sianida (disarankan untuk guru).
// Password diisi = akun lokal (mis. akun cadangan saat sianida tidak bisa diakses).
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"omitempty,email"`
	Role     string `json:"role" validate:"omitempty,oneof=admin"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

// UpdateUserRequest: field kosong/tidak dikirim berarti tidak diubah.
// Password: tidak dikirim (null) = tidak diubah, "" = pindah ke login sianida, diisi = set password lokal.
type UpdateUserRequest struct {
	Name     string  `json:"name" validate:"omitempty,max=255"`
	Email    string  `json:"email" validate:"omitempty,email"`
	Role     string  `json:"role" validate:"omitempty,oneof=admin"`
	Status   string  `json:"status" validate:"omitempty,oneof=active inactive"`
	Password *string `json:"password" validate:"omitempty"`
}

type UserResponse struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	LoginSianida bool      `json:"login_sianida"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
