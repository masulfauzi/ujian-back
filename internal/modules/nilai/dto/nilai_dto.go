package dto

type CreateNilaiRequest struct {
	IDPeserta         string  `json:"id_peserta" validate:"required"`
	IDJadwal          string  `json:"id_jadwal" validate:"required"`
	Nilai             float64 `json:"nilai" validate:"required,min=0,max=100"`
	WktMulai          *string `json:"wkt_mulai"`
	AktivitasTerakhir *string `json:"aktivitas_terakhir"`
	WktSelesai        *string `json:"wkt_selesai"`
}

type UpdateNilaiRequest struct {
	IDPeserta         *string  `json:"id_peserta"`
	IDJadwal          *string  `json:"id_jadwal"`
	Nilai             *float64 `json:"nilai"`
	WktMulai          *string  `json:"wkt_mulai"`
	AktivitasTerakhir *string  `json:"aktivitas_terakhir"`
	WktSelesai        *string  `json:"wkt_selesai"`
}

type NilaiResponse struct {
	ID                string  `json:"id"`
	IDPeserta         string  `json:"id_peserta"`
	NamaPeserta       string  `json:"nama_peserta"`
	IDJadwal          string  `json:"id_jadwal"`
	NamaUjian         string  `json:"nama_ujian"`
	Nilai             float64 `json:"nilai"`
	WktMulai          *string `json:"wkt_mulai"`
	AktivitasTerakhir *string `json:"aktivitas_terakhir"`
	WktSelesai        *string `json:"wkt_selesai"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type NilaiListResponse struct {
	Data      []NilaiResponse `json:"data"`
	Total     int64           `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalPage int             `json:"total_page"`
}
