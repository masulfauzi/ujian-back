package repository

import (
	"backend/internal/modules/bank_soal/model"

	"gorm.io/gorm"
)

type BankSoalRepository interface {
	Create(bankSoal *model.BankSoal) error
	GetByID(id string) (*model.BankSoal, error)
	GetAll(page, pageSize int) ([]model.BankSoal, int64, error)
	GetByMapelID(mapelID string, page, pageSize int) ([]model.BankSoal, int64, error)
	Update(bankSoal *model.BankSoal) error
	Delete(id string) error
	Restore(id string) error
	HardDelete(id string) error
}

type bankSoalRepository struct {
	db *gorm.DB
}

func NewBankSoalRepository(db *gorm.DB) BankSoalRepository {
	return &bankSoalRepository{db: db}
}

func (r *bankSoalRepository) Create(bankSoal *model.BankSoal) error {
	return r.db.Create(bankSoal).Error
}

func (r *bankSoalRepository) GetByID(id string) (*model.BankSoal, error) {
	var bankSoal model.BankSoal
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&bankSoal).Error
	if err != nil {
		return nil, err
	}
	return &bankSoal, nil
}

func (r *bankSoalRepository) GetAll(page, pageSize int) ([]model.BankSoal, int64, error) {
	var bankSoals []model.BankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.BankSoal{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Find(&bankSoals).Error

	return bankSoals, total, err
}

func (r *bankSoalRepository) GetByMapelID(mapelID string, page, pageSize int) ([]model.BankSoal, int64, error) {
	var bankSoals []model.BankSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	err := r.db.
		Model(&model.BankSoal{}).
		Where("id_mapel = ? AND deleted_at IS NULL", mapelID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Where("id_mapel = ? AND deleted_at IS NULL", mapelID).
		Offset(offset).
		Limit(pageSize).
		Find(&bankSoals).Error

	return bankSoals, total, err
}

func (r *bankSoalRepository) Update(bankSoal *model.BankSoal) error {
	return r.db.Save(bankSoal).Error
}

func (r *bankSoalRepository) Delete(id string) error {
	return r.db.Delete(&model.BankSoal{}, "id = ?", id).Error
}

func (r *bankSoalRepository) Restore(id string) error {
	return r.db.Table("bank_soal").Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *bankSoalRepository) HardDelete(id string) error {
	return r.db.Unscoped().Delete(&model.BankSoal{}, "id = ?", id).Error
}
