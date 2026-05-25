package service

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"math"
	"math/rand"
	"strings"

	"backend/internal/constants"
	"backend/internal/modules/jawaban/dto"
	"backend/internal/modules/jawaban/model"
	"backend/internal/modules/jawaban/repository"

	"gorm.io/gorm"
)

type JawabanService interface {
	CreateJawaban(req *dto.CreateJawabanRequest) (*dto.JawabanResponse, error)
	GetJawabanByID(id string) (*dto.JawabanResponse, error)
	GetAllJawaban(page, pageSize int, idNilai, idPeserta, idSoal string) (*dto.JawabanListResponse, error)
	GetJawabanByNilai(idNilai string) ([]dto.JawabanResponse, error)
	GetJawabanByPeserta(idPeserta string, page, pageSize int) (*dto.JawabanListResponse, error)
	UpdateJawaban(id string, req *dto.UpdateJawabanRequest) (*dto.JawabanResponse, error)
	DeleteJawaban(id string) error
	RestoreJawaban(id string) error
}

type jawabanService struct {
	repo repository.JawabanRepository
}

func NewJawabanService(repo repository.JawabanRepository) JawabanService {
	return &jawabanService{repo: repo}
}

func normalizeJawaban(j string) (string, error) {
	j = strings.ToUpper(strings.TrimSpace(j))
	switch j {
	case "A", "B", "C", "D", "E":
		return j, nil
	default:
		return "", errors.New("jawaban harus salah satu dari: A, B, C, D, E")
	}
}

// generateOpsiOrder menghasilkan urutan acak yang sama dengan soal service.
// Wajib identik dengan soalService.generateSeed + rand.Shuffle.
func generateOpsiOrder(pesertaID, soalID string) []string {
	hash := md5.Sum([]byte(pesertaID + "|" + soalID))
	seed := int64(binary.BigEndian.Uint64(hash[:8]))
	rng := rand.New(rand.NewSource(seed))

	order := []string{"A", "B", "C", "D", "E"}
	rng.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})
	return order
}

// reverseMapJawaban mengonversi posisi acak yang dipilih peserta
// ke posisi asli sebelum diacak, agar bisa dibandingkan dengan kunci di DB.
func reverseMapJawaban(submittedJawaban, pesertaID, soalID string) string {
	order := generateOpsiOrder(pesertaID, soalID)
	idx := int(submittedJawaban[0] - 'A') // "B" → 1
	if idx < 0 || idx >= len(order) {
		return submittedJawaban
	}
	return order[idx] // posisi asli
}

func randomizeOpsi(detail *repository.JawabanWithDetail) (opsiA, opsiB, opsiC, opsiD, opsiE, gambarA, gambarB, gambarC, gambarD, gambarE, kunci string) {
	opsiOrder := generateOpsiOrder(detail.IDPeserta, detail.IDSoal)

	opsiMap := map[string]string{
		"A": detail.OpsiA, "B": detail.OpsiB, "C": detail.OpsiC,
		"D": detail.OpsiD, "E": detail.OpsiE,
	}
	gambarMap := map[string]string{
		"A": detail.GambarA, "B": detail.GambarB, "C": detail.GambarC,
		"D": detail.GambarD, "E": detail.GambarE,
	}

	newKunciPos := 0
	for i, key := range opsiOrder {
		if key == detail.Kunci {
			newKunciPos = i
			break
		}
	}
	newKunci := string(rune('A' + newKunciPos))

	return opsiMap[opsiOrder[0]], opsiMap[opsiOrder[1]], opsiMap[opsiOrder[2]], opsiMap[opsiOrder[3]], opsiMap[opsiOrder[4]],
		gambarMap[opsiOrder[0]], gambarMap[opsiOrder[1]], gambarMap[opsiOrder[2]], gambarMap[opsiOrder[3]], gambarMap[opsiOrder[4]],
		newKunci
}

func (s *jawabanService) CreateJawaban(req *dto.CreateJawabanRequest) (*dto.JawabanResponse, error) {
	jawaban, err := normalizeJawaban(req.Jawaban)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.CheckDuplicate(req.IDNilai, req.IDSoal)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("jawaban untuk soal ini di attempt tersebut sudah ada — gunakan endpoint update")
	}

	kunci, err := s.repo.GetSoalKunci(req.IDSoal)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("soal tidak ditemukan")
		}
		return nil, err
	}
	if kunci == "" {
		return nil, errors.New("soal tidak ditemukan")
	}

	// Reverse-map jawaban dari posisi acak ke posisi asli sebelum dibandingkan kunci
	originalJawaban := reverseMapJawaban(jawaban, req.IDPeserta, req.IDSoal)
	isBenarBool := strings.EqualFold(originalJawaban, strings.TrimSpace(kunci))
	isBenarInt := 0
	if isBenarBool {
		isBenarInt = 1
	}

	row := &model.Jawaban{
		IDNilai:   req.IDNilai,
		IDSoal:    req.IDSoal,
		IDPeserta: req.IDPeserta,
		NoUrut:    req.NoUrut,
		Jawaban:   &jawaban,
		IsBenar:   &isBenarInt,
	}

	if err := s.repo.Create(row); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByIDWithDetail(row.ID)
	if err != nil {
		return nil, err
	}
	return detailToResponse(created), nil
}

func (s *jawabanService) GetJawabanByID(id string) (*dto.JawabanResponse, error) {
	result, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return detailToResponse(result), nil
}

func (s *jawabanService) GetAllJawaban(page, pageSize int, idNilai, idPeserta, idSoal string) (*dto.JawabanListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	results, total, err := s.repo.GetAllWithDetail(page, pageSize, idNilai, idPeserta, idSoal)
	if err != nil {
		return nil, err
	}

	responses := []dto.JawabanResponse{}
	for _, r := range results {
		responses = append(responses, *detailToResponse(&r))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.JawabanListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *jawabanService) GetJawabanByNilai(idNilai string) ([]dto.JawabanResponse, error) {
	resp, err := s.GetAllJawaban(1, 99999, idNilai, "", "")
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (s *jawabanService) GetJawabanByPeserta(idPeserta string, page, pageSize int) (*dto.JawabanListResponse, error) {
	return s.GetAllJawaban(page, pageSize, "", idPeserta, "")
}

func (s *jawabanService) UpdateJawaban(id string, req *dto.UpdateJawabanRequest) (*dto.JawabanResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	jawaban, err := normalizeJawaban(req.Jawaban)
	if err != nil {
		return nil, err
	}

	if req.IDNilai != existing.IDNilai || req.IDSoal != existing.IDSoal {
		exists, err := s.repo.CheckDuplicate(req.IDNilai, req.IDSoal)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("jawaban untuk soal ini di attempt tersebut sudah ada")
		}
	}

	kunci, err := s.repo.GetSoalKunci(req.IDSoal)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("soal tidak ditemukan")
		}
		return nil, err
	}
	if kunci == "" {
		return nil, errors.New("soal tidak ditemukan")
	}

	// Reverse-map jawaban dari posisi acak ke posisi asli sebelum dibandingkan kunci
	originalJawaban := reverseMapJawaban(jawaban, req.IDPeserta, req.IDSoal)
	isBenarBool := strings.EqualFold(originalJawaban, strings.TrimSpace(kunci))
	isBenarInt := 0
	if isBenarBool {
		isBenarInt = 1
	}

	existing.IDNilai   = req.IDNilai
	existing.IDSoal    = req.IDSoal
	existing.IDPeserta = req.IDPeserta
	existing.NoUrut    = req.NoUrut
	existing.Jawaban   = &jawaban
	existing.IsBenar   = &isBenarInt

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDWithDetail(id)
	if err != nil {
		return nil, err
	}
	return detailToResponse(updated), nil
}

func (s *jawabanService) DeleteJawaban(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(id)
}

func (s *jawabanService) RestoreJawaban(id string) error {
	return s.repo.Restore(id)
}

func detailToResponse(r *repository.JawabanWithDetail) *dto.JawabanResponse {
	opsiA, opsiB, opsiC, opsiD, opsiE, gambarA, gambarB, gambarC, gambarD, gambarE, kunci := randomizeOpsi(r)

	return &dto.JawabanResponse{
		ID:          r.ID,
		IDNilai:     r.IDNilai,
		IDSoal:      r.IDSoal,
		NoUrut:      r.NoUrut,
		NoSoal:      r.NoSoal,
		SoalText:    r.SoalText,
		Kunci:       kunci,
		OpsiA:       opsiA,
		OpsiB:       opsiB,
		OpsiC:       opsiC,
		OpsiD:       opsiD,
		OpsiE:       opsiE,
		GambarA:     gambarA,
		GambarB:     gambarB,
		GambarC:     gambarC,
		GambarD:     gambarD,
		GambarE:     gambarE,
		IDPeserta:   r.IDPeserta,
		NamaPeserta: r.NamaPeserta,
		Jawaban:     r.Jawaban,
		IsBenar:     r.IsBenar,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
