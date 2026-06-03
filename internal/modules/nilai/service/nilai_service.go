package service

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"backend/internal/constants"
	jawabanmodel "backend/internal/modules/jawaban/model"
	jawabanrepo "backend/internal/modules/jawaban/repository"
	jadwalmodel "backend/internal/modules/jadwal/model"
	"backend/internal/modules/nilai/dto"
	"backend/internal/modules/nilai/model"
	"backend/internal/modules/nilai/repository"
	soalmodel "backend/internal/modules/soal/model"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var jakartaLoc, _ = time.LoadLocation("Asia/Jakarta")

func init() {
	if jakartaLoc == nil {
		jakartaLoc = time.FixedZone("WIB", 7*60*60)
	}
}

type ExportResult struct {
	ZipBytes  []byte
	NamaUjian string
}

type NilaiService interface {
	CreateNilai(req *dto.CreateNilaiRequest) (*dto.NilaiResponse, error)
	GetNilaiByID(id string) (*dto.NilaiResponse, error)
	GetAllNilai(page, pageSize int, idPeserta, idJadwal string) (*dto.NilaiListResponse, error)
	GetNilaiByPeserta(idPeserta string, page, pageSize int) (*dto.NilaiListResponse, error)
	GetNilaiByJadwal(idJadwal string, page, pageSize int) (*dto.NilaiListResponse, error)
	UpdateNilai(id string, req *dto.UpdateNilaiRequest) (*dto.NilaiResponse, error)
	DeleteNilai(id string) error
	RestoreNilai(id string) error
	MulaiUjian(idPeserta, idJadwal string) (*dto.NilaiResponse, bool, error)
	ExportNilaiByJadwal(idJadwal string) (*ExportResult, error)
}

type nilaiService struct {
	repo        repository.NilaiRepository
	jawabanRepo jawabanrepo.JawabanRepository
	db          *gorm.DB
}

func NewNilaiService(repo repository.NilaiRepository, jawabanRepo jawabanrepo.JawabanRepository, db *gorm.DB) NilaiService {
	return &nilaiService{
		repo:        repo,
		jawabanRepo: jawabanRepo,
		db:          db,
	}
}

const timeLayout = "2006-01-02 15:04:05"

func parseTime(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse(timeLayout, *s)
	if err != nil {
		return nil, errors.New("format waktu tidak valid, gunakan: 2006-01-02 15:04:05")
	}
	return &t, nil
}

func (s *nilaiService) CreateNilai(req *dto.CreateNilaiRequest) (*dto.NilaiResponse, error) {
	if req.Nilai < 0 || req.Nilai > 100 {
		return nil, errors.New("nilai harus di antara 0 dan 100")
	}

	exists, err := s.repo.CheckDuplicate(req.IDPeserta, req.IDJadwal)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("nilai untuk peserta dan jadwal ini sudah ada — gunakan endpoint update")
	}

	wktMulai, err := parseTime(req.WktMulai)
	if err != nil {
		return nil, err
	}
	aktivitasTerakhir, err := parseTime(req.AktivitasTerakhir)
	if err != nil {
		return nil, err
	}
	wktSelesai, err := parseTime(req.WktSelesai)
	if err != nil {
		return nil, err
	}

	nilai := &model.Nilai{
		IDPeserta:         req.IDPeserta,
		IDJadwal:          req.IDJadwal,
		Nilai:             req.Nilai,
		WktMulai:          wktMulai,
		AktivitasTerakhir: aktivitasTerakhir,
		WktSelesai:        wktSelesai,
	}

	if err := s.repo.Create(nilai); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByIDWithDetail(nilai.ID)
	if err != nil {
		return nil, err
	}
	return detailToResponse(created), nil
}

func (s *nilaiService) GetNilaiByID(id string) (*dto.NilaiResponse, error) {
	result, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return detailToResponse(result), nil
}

func (s *nilaiService) GetAllNilai(page, pageSize int, idPeserta, idJadwal string) (*dto.NilaiListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	results, total, err := s.repo.GetAllWithDetail(page, pageSize, idPeserta, idJadwal)
	if err != nil {
		return nil, err
	}

	responses := []dto.NilaiResponse{}
	for _, r := range results {
		responses = append(responses, *detailToResponse(&r))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.NilaiListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *nilaiService) GetNilaiByPeserta(idPeserta string, page, pageSize int) (*dto.NilaiListResponse, error) {
	return s.GetAllNilai(page, pageSize, idPeserta, "")
}

func (s *nilaiService) GetNilaiByJadwal(idJadwal string, page, pageSize int) (*dto.NilaiListResponse, error) {
	return s.GetAllNilai(page, pageSize, "", idJadwal)
}

func (s *nilaiService) UpdateNilai(id string, req *dto.UpdateNilaiRequest) (*dto.NilaiResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if req.Nilai != nil {
		if *req.Nilai < 0 || *req.Nilai > 100 {
			return nil, errors.New("nilai harus di antara 0 dan 100")
		}
		existing.Nilai = *req.Nilai
	}

	if req.IDPeserta != nil || req.IDJadwal != nil {
		newPeserta := existing.IDPeserta
		newJadwal  := existing.IDJadwal
		if req.IDPeserta != nil {
			newPeserta = *req.IDPeserta
		}
		if req.IDJadwal != nil {
			newJadwal = *req.IDJadwal
		}
		if newPeserta != existing.IDPeserta || newJadwal != existing.IDJadwal {
			exists, err := s.repo.CheckDuplicate(newPeserta, newJadwal)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("nilai untuk peserta dan jadwal ini sudah ada")
			}
		}
		existing.IDPeserta = newPeserta
		existing.IDJadwal  = newJadwal
	}

	if req.WktMulai != nil {
		t, err := parseTime(req.WktMulai)
		if err != nil {
			return nil, err
		}
		existing.WktMulai = t
	}

	if req.AktivitasTerakhir != nil {
		t, err := parseTime(req.AktivitasTerakhir)
		if err != nil {
			return nil, err
		}
		existing.AktivitasTerakhir = t
	}

	if req.WktSelesai != nil {
		t, err := parseTime(req.WktSelesai)
		if err != nil {
			return nil, err
		}
		existing.WktSelesai = t

		if t != nil {
			nilai, err := s.repo.HitungNilai(id)
			if err != nil {
				return nil, err
			}
			existing.Nilai = nilai
		}
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		return nil, err
	}
	return detailToResponse(updated), nil
}

func (s *nilaiService) DeleteNilai(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.jawabanRepo.SoftDeleteByNilaiID(tx, id); err != nil {
			return err
		}
		now := time.Now()
		return tx.Model(&model.Nilai{}).Where("id = ?", id).Update("deleted_at", now).Error
	})
}

func (s *nilaiService) RestoreNilai(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.jawabanRepo.RestoreByNilaiID(tx, id); err != nil {
			return err
		}
		return tx.Model(&model.Nilai{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NULL")).Error
	})
}

func (s *nilaiService) MulaiUjian(idPeserta, idJadwal string) (*dto.NilaiResponse, bool, error) {
	// 1. Cek apakah record sudah ada
	existing, err := s.repo.GetByPesertaAndJadwal(idPeserta, idJadwal)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	// 2. Jika sudah ada → cek wkt_selesai
	if existing != nil {
		if existing.WktSelesai != nil {
			return nil, false, errors.New("Ujian sudah pernah dilakukan")
		}
		// Resume: ambil detail (dengan JOIN) lalu return
		detail, err := s.repo.GetByIDWithDetail(existing.ID)
		if err != nil {
			return nil, false, err
		}
		return detailToResponse(detail), false, nil
	}

	// 3. Belum ada → transaction: insert nilai + bulk insert jawaban
	var newNilaiID string
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 3a. Get jadwal untuk dapatkan id_bank_soal.
		// SELECT eksplisit + cast ke int agar acak_soal/acak_opsi terbaca
		// walau kolom DB masih boolean.
		var jadwal jadwalmodel.Jadwal
		if err := tx.Table("jadwal").
			Select("id, id_bank_soal, nama_ujian, tingkat, wkt_mulai, wkt_selesai, durasi, acak_soal::int AS acak_soal, acak_opsi::int AS acak_opsi, created_at, updated_at, deleted_at").
			Where("id = ? AND deleted_at IS NULL", idJadwal).
			First(&jadwal).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("jadwal tidak ditemukan")
			}
			return err
		}

		// 3b. Insert nilai baru
		now := time.Now().In(jakartaLoc)
		nilai := &model.Nilai{
			IDPeserta:         idPeserta,
			IDJadwal:          idJadwal,
			Nilai:             0,
			WktMulai:          &now,
			AktivitasTerakhir: &now,
			WktSelesai:        nil,
		}
		if err := tx.Create(nilai).Error; err != nil {
			return err
		}
		newNilaiID = nilai.ID

		// 3c. Query soal by bank_soal — acak jika acak_soal=1, urut jika 0
		var soals []soalmodel.Soal
		soalQuery := tx.Where("id_bank_soal = ? AND deleted_at IS NULL", jadwal.IDBankSoal)
		if jadwal.AcakSoal == 1 {
			soalQuery = soalQuery.Order("RANDOM()")
		} else {
			soalQuery = soalQuery.Order("no_soal ASC")
		}
		if err := soalQuery.Find(&soals).Error; err != nil {
			return err
		}

		// 3d. Build & bulk insert jawaban kosong
		if len(soals) > 0 {
			jawabans := make([]jawabanmodel.Jawaban, len(soals))
			if jadwal.AcakSoal == 1 {
				noUrutSequence := rand.Perm(len(soals))
				for i, soal := range soals {
					jawabans[i] = jawabanmodel.Jawaban{
						IDNilai:   nilai.ID,
						IDSoal:    soal.ID,
						IDPeserta: idPeserta,
						NoUrut:    noUrutSequence[i] + 1,
						Jawaban:   nil,
						IsBenar:   nil,
					}
				}
			} else {
				for i, soal := range soals {
					jawabans[i] = jawabanmodel.Jawaban{
						IDNilai:   nilai.ID,
						IDSoal:    soal.ID,
						IDPeserta: idPeserta,
						NoUrut:    i + 1,
						Jawaban:   nil,
						IsBenar:   nil,
					}
				}
			}
			if err := s.jawabanRepo.BulkCreateWithTx(tx, jawabans); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, false, err
	}

	// 4. Ambil detail nilai yang baru di-insert (di luar transaction, read)
	created, err := s.repo.GetByIDWithDetail(newNilaiID)
	if err != nil {
		return nil, false, err
	}
	return detailToResponse(created), true, nil
}

func detailToResponse(r *repository.NilaiWithDetail) *dto.NilaiResponse {
	return &dto.NilaiResponse{
		ID:                r.ID,
		IDPeserta:         r.IDPeserta,
		NamaPeserta:       r.NamaPeserta,
		IDJadwal:          r.IDJadwal,
		NamaUjian:         r.NamaUjian,
		Nilai:             r.Nilai,
		WktMulai:          r.WktMulai,
		AktivitasTerakhir: r.AktivitasTerakhir,
		WktSelesai:        r.WktSelesai,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

func (s *nilaiService) ExportNilaiByJadwal(idJadwal string) (*ExportResult, error) {
	// 1. Ambil nama ujian dari tabel jadwal
	var namaUjian string
	if err := s.db.Table("jadwal").
		Select("nama_ujian").
		Where("id = ? AND deleted_at IS NULL", idJadwal).
		Scan(&namaUjian).Error; err != nil || namaUjian == "" {
		return nil, errors.New("jadwal tidak ditemukan")
	}

	// 2. Ambil nama mapel dari jadwal
	type JadwalDetail struct {
		IDBankSoal string `gorm:"column:id_bank_soal"`
	}
	var jadwalDetail JadwalDetail
	if err := s.db.Table("jadwal").
		Select("id_bank_soal").
		Where("id = ? AND deleted_at IS NULL", idJadwal).
		Scan(&jadwalDetail).Error; err != nil {
		return nil, err
	}

	var namaMapel string
	if err := s.db.Table("bank_soal").
		Select("mapel.nama_mapel").
		Joins("INNER JOIN mapel ON bank_soal.id_mapel = mapel.id").
		Where("bank_soal.id = ? AND bank_soal.deleted_at IS NULL", jadwalDetail.IDBankSoal).
		Scan(&namaMapel).Error; err != nil || namaMapel == "" {
		return nil, errors.New("mapel tidak ditemukan untuk jadwal ini")
	}

	// 3. Ambil semua kelas yang terdaftar pada jadwal ini
	type KelasRow struct {
		IDKelas   string `gorm:"column:id_kelas"`
		NamaKelas string `gorm:"column:nama_kelas"`
	}
	var kelasList []KelasRow
	if err := s.db.Table("jadwal_kelas").
		Select("jadwal_kelas.id_kelas, kelas.nama_kelas").
		Joins("INNER JOIN kelas ON jadwal_kelas.id_kelas = kelas.id").
		Where("jadwal_kelas.id_jadwal = ?", idJadwal).
		Scan(&kelasList).Error; err != nil {
		return nil, err
	}
	if len(kelasList) == 0 {
		return nil, errors.New("tidak ada kelas yang terdaftar pada jadwal ini")
	}

	// 4. Buat ZIP di memory
	var zipBuf bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuf)

	for _, kelas := range kelasList {
		// 4a. Ambil data nilai peserta untuk kelas ini
		rows, err := s.repo.GetByJadwalAndKelas(idJadwal, kelas.IDKelas)
		if err != nil {
			return nil, fmt.Errorf("gagal ambil data kelas %s: %w", kelas.NamaKelas, err)
		}

		// 4b. Buat file Excel
		xlsx := excelize.NewFile()
		sheet := "Nilai Ujian"
		xlsx.SetSheetName("Sheet1", sheet)

		// Header
		headers := []string{"No", "Nama Peserta", "Username", "Nilai", "Waktu Mulai", "Waktu Selesai"}
		for col, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(col+1, 1)
			xlsx.SetCellValue(sheet, cell, h)
		}

		// Data rows
		for i, row := range rows {
			rowNum := i + 2
			xlsx.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), i+1)
			xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), row.NamaPeserta)
			xlsx.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), row.Username)
			xlsx.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), row.Nilai)

			wktMulai := "-"
			if row.WktMulai != nil {
				wktMulai = *row.WktMulai
			}
			xlsx.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), wktMulai)

			wktSelesai := "-"
			if row.WktSelesai != nil {
				wktSelesai = *row.WktSelesai
			}
			xlsx.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), wktSelesai)
		}

		// 4c. Tulis Excel ke buffer lalu masukkan ke ZIP
		var xlsBuf bytes.Buffer
		if err := xlsx.Write(&xlsBuf); err != nil {
			return nil, fmt.Errorf("gagal tulis excel kelas %s: %w", kelas.NamaKelas, err)
		}

		// Nama file: {MAPEL}_{KELAS}.xlsx (spasi diganti underscore)
		safeMapel := strings.ReplaceAll(namaMapel, " ", "_")
		safeKelas := strings.ReplaceAll(kelas.NamaKelas, " ", "_")
		filename := fmt.Sprintf("%s_%s.xlsx", safeMapel, safeKelas)

		zipEntry, err := zipWriter.Create(filename)
		if err != nil {
			return nil, err
		}
		if _, err := zipEntry.Write(xlsBuf.Bytes()); err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return &ExportResult{
		ZipBytes:  zipBuf.Bytes(),
		NamaUjian: namaUjian,
	}, nil
}
