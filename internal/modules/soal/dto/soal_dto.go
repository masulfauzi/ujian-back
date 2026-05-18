package dto

type CreateSoalRequest struct {
	IdBankSoal string `json:"id_bank_soal" validate:"required"`
	Soal       string `json:"soal" validate:"required"`
	GambarSoal string `json:"gambar_soal"`
	OpsiA      string `json:"opsi_a" validate:"required"`
	OpsiB      string `json:"opsi_b" validate:"required"`
	OpsiC      string `json:"opsi_c" validate:"required"`
	OpsiD      string `json:"opsi_d"`
	OpsiE      string `json:"opsi_e"`
	GambarA    string `json:"gambar_a"`
	GambarB    string `json:"gambar_b"`
	GambarC    string `json:"gambar_c"`
	GambarD    string `json:"gambar_d"`
	GambarE    string `json:"gambar_e"`
	Kunci      string `json:"kunci" validate:"required,len=1"`
}

type UpdateSoalRequest struct {
	Soal       string `json:"soal" validate:"required"`
	GambarSoal string `json:"gambar_soal"`
	OpsiA      string `json:"opsi_a" validate:"required"`
	OpsiB      string `json:"opsi_b" validate:"required"`
	OpsiC      string `json:"opsi_c" validate:"required"`
	OpsiD      string `json:"opsi_d"`
	OpsiE      string `json:"opsi_e"`
	GambarA    string `json:"gambar_a"`
	GambarB    string `json:"gambar_b"`
	GambarC    string `json:"gambar_c"`
	GambarD    string `json:"gambar_d"`
	GambarE    string `json:"gambar_e"`
	Kunci      string `json:"kunci" validate:"required,len=1"`
}

type SoalResponse struct {
	ID         string `json:"id"`
	IdBankSoal string `json:"id_bank_soal"`
	Soal       string `json:"soal"`
	GambarSoal string `json:"gambar_soal"`
	OpsiA      string `json:"opsi_a"`
	OpsiB      string `json:"opsi_b"`
	OpsiC      string `json:"opsi_c"`
	OpsiD      string `json:"opsi_d"`
	OpsiE      string `json:"opsi_e"`
	GambarA    string `json:"gambar_a"`
	GambarB    string `json:"gambar_b"`
	GambarC    string `json:"gambar_c"`
	GambarD    string `json:"gambar_d"`
	GambarE    string `json:"gambar_e"`
	Kunci      string `json:"kunci"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type SoalListResponse struct {
	Data      []SoalResponse `json:"data"`
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"page_size"`
	TotalPage int            `json:"total_page"`
}
