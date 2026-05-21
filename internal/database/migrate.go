package database

import (
	banksoalmodel    "backend/internal/modules/bank_soal/model"
	jadwalmodel      "backend/internal/modules/jadwal/model"
	jadwalkelasmodel "backend/internal/modules/jadwal_kelas/model"
	jurusanmodel     "backend/internal/modules/jurusan/model"
	kelasmodel       "backend/internal/modules/kelas/model"
	mapelmodel       "backend/internal/modules/mapel/model"
	soalmodel        "backend/internal/modules/soal/model"
	usermodel        "backend/internal/modules/user/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Drop old non-partial unique index on jurusan.nama_jurusan before migrating
	// so AutoMigrate can create the correct partial index (where:deleted_at IS NULL)
	db.Exec("DROP INDEX IF EXISTS idx_jurusans_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS uni_jurusan_nama_jurusan")
	db.Exec("DROP INDEX IF EXISTS idx_jurusan_nama_jurusan")

	if err := db.AutoMigrate(
		&usermodel.User{},
		&mapelmodel.Mapel{},
		&banksoalmodel.BankSoal{},
		&soalmodel.Soal{},
		&jurusanmodel.Jurusan{},
		&kelasmodel.Kelas{},
		&jadwalmodel.Jadwal{},
		&jadwalkelasmodel.JadwalKelas{},
	); err != nil {
		return err
	}

	// Buat unique constraint untuk mencegah duplikasi assignment kelas ke jadwal
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_jadwal_kelas_unique ON jadwal_kelas(id_jadwal, id_kelas)")

	return nil
}
