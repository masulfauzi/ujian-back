package database

import (
	banksoalmodel "backend/internal/modules/bank_soal/model"
	mapelmodel "backend/internal/modules/mapel/model"
	soalmodel "backend/internal/modules/soal/model"
	usermodel "backend/internal/modules/user/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&usermodel.User{},
		&mapelmodel.Mapel{},
		&banksoalmodel.BankSoal{},
		&soalmodel.Soal{},
	)
}
