package repository

import (
	"backend/internal/modules/jadwal/model"
	"time"

	"gorm.io/gorm"
)

type JadwalWithBankSoal struct {
	ID           string  `gorm:"column:id"`
	IDBankSoal   string  `gorm:"column:id_bank_soal"`
	NamaBankSoal string  `gorm:"column:nama_bank_soal"`
	Tingkat      string  `gorm:"column:tingkat"`
	WktMulai     string  `gorm:"column:wkt_mulai"`
	WktSelesai   string  `gorm:"column:wkt_selesai"`
	CreatedAt    string  `gorm:"column:created_at"`
	UpdatedAt    string  `gorm:"column:updated_at"`
}

type KelasDetail struct {
	ID        string `gorm:"column:id"`
	IDKelas   string `gorm:"column:id_kelas"`
	NamaKelas string `gorm:"column:nama_kelas"`
	IDJurusan string `gorm:"column:id_jurusan"`
}

type JurusanDetail struct {
	ID          string `gorm:"column:id"`
	IDJurusan   string `gorm:"column:id_jurusan"`
	NamaJurusan string `gorm:"column:nama_jurusan"`
}

type JadwalWithKelas struct {
	ID           string           `json:"id"`
	IDBankSoal   string           `json:"id_bank_soal"`
	NamaBankSoal string           `json:"nama_bank_soal"`
	Tingkat      string           `json:"tingkat"`
	WktMulai     string           `json:"wkt_mulai"`
	WktSelesai   string           `json:"wkt_selesai"`
	IDKelas      []KelasDetail    `json:"id_kelas"`
	IDJurusan    []JurusanDetail  `json:"id_jurusan"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
}

type JadwalRepository interface {
	Create(jadwal *model.Jadwal) error
	GetByID(id string) (*model.Jadwal, error)
	GetByIDWithBankSoal(id string) (*JadwalWithBankSoal, error)
	GetByIDWithKelas(id string) (*JadwalWithKelas, error)
	GetAllWithBankSoal(page, pageSize int) ([]JadwalWithBankSoal, int64, error)
	GetByBankSoalID(bankSoalID string, page, pageSize int) ([]JadwalWithBankSoal, int64, error)
	Update(jadwal *model.Jadwal) error
	Delete(id string) error
	Restore(id string) error
}

type jadwalRepository struct {
	db *gorm.DB
}

func NewJadwalRepository(db *gorm.DB) JadwalRepository {
	return &jadwalRepository{db: db}
}

func (r *jadwalRepository) Create(jadwal *model.Jadwal) error {
	return r.db.Create(jadwal).Error
}

func (r *jadwalRepository) GetByID(id string) (*model.Jadwal, error) {
	var jadwal model.Jadwal
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&jadwal).Error
	if err != nil {
		return nil, err
	}
	return &jadwal, nil
}

func (r *jadwalRepository) GetByIDWithBankSoal(id string) (*JadwalWithBankSoal, error) {
	var jadwal JadwalWithBankSoal
	err := r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.tingkat, jadwal.wkt_mulai, jadwal.wkt_selesai, jadwal.created_at, jadwal.updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id = ? AND jadwal.deleted_at IS NULL", id).
		First(&jadwal).Error
	if err != nil {
		return nil, err
	}
	return &jadwal, nil
}

func (r *jadwalRepository) GetByIDWithKelas(id string) (*JadwalWithKelas, error) {
	var jadwal JadwalWithBankSoal
	err := r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.tingkat, jadwal.wkt_mulai, jadwal.wkt_selesai, jadwal.created_at, jadwal.updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id = ? AND jadwal.deleted_at IS NULL", id).
		First(&jadwal).Error
	if err != nil {
		return nil, err
	}

	var kelasList []KelasDetail
	err = r.db.
		Table("jadwal_kelas").
		Select("jadwal_kelas.id, jadwal_kelas.id_kelas, kelas.nama_kelas, kelas.id_jurusan").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Where("jadwal_kelas.id_jadwal = ?", id).
		Scan(&kelasList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if kelasList == nil {
		kelasList = []KelasDetail{}
	}

	var jurusanList []JurusanDetail
	err = r.db.
		Table("jurusan").
		Select("DISTINCT jurusan.id, kelas.id_jurusan, jurusan.nama_jurusan").
		Joins("INNER JOIN kelas ON jurusan.id = kelas.id_jurusan").
		Joins("INNER JOIN jadwal_kelas ON kelas.id = jadwal_kelas.id_kelas").
		Where("jadwal_kelas.id_jadwal = ? AND jurusan.deleted_at IS NULL", id).
		Scan(&jurusanList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if jurusanList == nil {
		jurusanList = []JurusanDetail{}
	}

	return &JadwalWithKelas{
		ID:           jadwal.ID,
		IDBankSoal:   jadwal.IDBankSoal,
		NamaBankSoal: jadwal.NamaBankSoal,
		Tingkat:      jadwal.Tingkat,
		WktMulai:     jadwal.WktMulai,
		WktSelesai:   jadwal.WktSelesai,
		IDKelas:      kelasList,
		IDJurusan:    jurusanList,
		CreatedAt:    jadwal.CreatedAt,
		UpdatedAt:    jadwal.UpdatedAt,
	}, nil
}

func (r *jadwalRepository) GetAllWithBankSoal(page, pageSize int) ([]JadwalWithBankSoal, int64, error) {
	var jadwalList []JadwalWithBankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Table("jadwal").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.tingkat, jadwal.wkt_mulai, jadwal.wkt_selesai, jadwal.created_at, jadwal.updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Scan(&jadwalList).Error

	return jadwalList, total, err
}

func (r *jadwalRepository) GetByBankSoalID(bankSoalID string, page, pageSize int) ([]JadwalWithBankSoal, int64, error) {
	var jadwalList []JadwalWithBankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Table("jadwal").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id_bank_soal = ? AND jadwal.deleted_at IS NULL", bankSoalID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Table("jadwal").
		Select("jadwal.id, jadwal.id_bank_soal, bank_soal.nama_bank_soal, jadwal.tingkat, jadwal.wkt_mulai, jadwal.wkt_selesai, jadwal.created_at, jadwal.updated_at").
		Joins("INNER JOIN bank_soal ON jadwal.id_bank_soal = bank_soal.id").
		Where("jadwal.id_bank_soal = ? AND jadwal.deleted_at IS NULL", bankSoalID).
		Offset(offset).
		Limit(pageSize).
		Scan(&jadwalList).Error

	return jadwalList, total, err
}

func (r *jadwalRepository) Update(jadwal *model.Jadwal) error {
	return r.db.Save(jadwal).Error
}

func (r *jadwalRepository) Delete(id string) error {
	now := time.Now()
	return r.db.Model(&model.Jadwal{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *jadwalRepository) Restore(id string) error {
	return r.db.Model(&model.Jadwal{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
}
