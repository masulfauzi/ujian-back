package database

import (
	banksoalmodel "backend/internal/modules/bank_soal/model"
	jurusanmodel "backend/internal/modules/jurusan/model"
	mapelmodel "backend/internal/modules/mapel/model"
	soalmodel "backend/internal/modules/soal/model"
	usermodel "backend/internal/modules/user/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Drop old non-partial unique index on jurusan.nama_jurusan before migrating
	// so AutoMigrate can create the correct partial index (where:deleted_at IS NULL)
	db.Exec("DROP INDEX IF EXISTS idx_jurusans_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS uni_jurusan_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS idx_jurusan_nama_jurusan")

	return db.AutoMigrate(
		&usermodel.User{},
		&mapelmodel.Mapel{},
		&banksoalmodel.BankSoal{},
		&soalmodel.Soal{},
		&jurusanmodel.Jurusan{},
	)
}
