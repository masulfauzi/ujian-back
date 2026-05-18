package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/soal/dto"
	"backend/internal/modules/soal/model"
	"backend/internal/modules/soal/repository"

	"gorm.io/gorm"
)

type SoalService interface {
	CreateSoal(req *dto.CreateSoalRequest) (*dto.SoalResponse, error)
	GetSoalByID(id string) (*dto.SoalResponse, error)
	GetAllSoal(page, pageSize int) (*dto.SoalListResponse, error)
	GetSoalByBankSoal(bankSoalID string, page, pageSize int) (*dto.SoalListResponse, error)
	UpdateSoal(id string, req *dto.UpdateSoalRequest) (*dto.SoalResponse, error)
	DeleteSoal(id string) error
	RestoreSoal(id string) error
}

type soalService struct {
	repo repository.SoalRepository
}

func NewSoalService(repo repository.SoalRepository) SoalService {
	return &soalService{repo: repo}
}

func (s *soalService) CreateSoal(req *dto.CreateSoalRequest) (*dto.SoalResponse, error) {
	if err := s.validateKunci(req.Kunci, req.OpsiA, req.OpsiB, req.OpsiC, req.OpsiD, req.OpsiE); err != nil {
		return nil, err
	}

	soal := &model.Soal{
		IdBankSoal: req.IdBankSoal,
		Soal:       req.Soal,
		GambarSoal: req.GambarSoal,
		OpsiA:      req.OpsiA,
		OpsiB:      req.OpsiB,
		OpsiC:      req.OpsiC,
		OpsiD:      req.OpsiD,
		OpsiE:      req.OpsiE,
		GambarA:    req.GambarA,
		GambarB:    req.GambarB,
		GambarC:    req.GambarC,
		GambarD:    req.GambarD,
		GambarE:    req.GambarE,
		Kunci:      req.Kunci,
	}

	if err := s.repo.Create(soal); err != nil {
		return nil, err
	}

	return s.modelToResponse(soal), nil
}

func (s *soalService) GetSoalByID(id string) (*dto.SoalResponse, error) {
	soal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.modelToResponse(soal), nil
}

func (s *soalService) GetAllSoal(page, pageSize int) (*dto.SoalListResponse, error) {
	soals, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	responses := []dto.SoalResponse{}
	for _, soal := range soals {
		responses = append(responses, *s.modelToResponse(&soal))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.SoalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *soalService) GetSoalByBankSoal(bankSoalID string, page, pageSize int) (*dto.SoalListResponse, error) {
	soals, total, err := s.repo.GetByBankSoalID(bankSoalID, page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	responses := []dto.SoalResponse{}
	for _, soal := range soals {
		responses = append(responses, *s.modelToResponse(&soal))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.SoalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *soalService) UpdateSoal(id string, req *dto.UpdateSoalRequest) (*dto.SoalResponse, error) {
	soal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if err := s.validateKunci(req.Kunci, req.OpsiA, req.OpsiB, req.OpsiC, req.OpsiD, req.OpsiE); err != nil {
		return nil, err
	}

	soal.Soal = req.Soal
	soal.GambarSoal = req.GambarSoal
	soal.OpsiA = req.OpsiA
	soal.OpsiB = req.OpsiB
	soal.OpsiC = req.OpsiC
	soal.OpsiD = req.OpsiD
	soal.OpsiE = req.OpsiE
	soal.GambarA = req.GambarA
	soal.GambarB = req.GambarB
	soal.GambarC = req.GambarC
	soal.GambarD = req.GambarD
	soal.GambarE = req.GambarE
	soal.Kunci = req.Kunci

	if err := s.repo.Update(soal); err != nil {
		return nil, err
	}

	return s.modelToResponse(soal), nil
}

func (s *soalService) DeleteSoal(id string) error {
	soal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(soal.ID)
}

func (s *soalService) RestoreSoal(id string) error {
	return s.repo.Restore(id)
}

func (s *soalService) validateKunci(kunci string, opsiA, opsiB, opsiC, opsiD, opsiE string) error {
	validKeys := map[string]bool{
		"A": true, "B": true, "C": true, "D": true, "E": true,
	}

	if !validKeys[kunci] {
		return errors.New("kunci harus A, B, C, D, atau E")
	}

	switch kunci {
	case "D":
		if opsiD == "" {
			return errors.New("opsi D tidak boleh kosong jika kunci D")
		}
	case "E":
		if opsiE == "" {
			return errors.New("opsi E tidak boleh kosong jika kunci E")
		}
	}

	return nil
}

func (s *soalService) modelToResponse(soal *model.Soal) *dto.SoalResponse {
	return &dto.SoalResponse{
		ID:         soal.ID,
		IdBankSoal: soal.IdBankSoal,
		Soal:       soal.Soal,
		GambarSoal: soal.GambarSoal,
		OpsiA:      soal.OpsiA,
		OpsiB:      soal.OpsiB,
		OpsiC:      soal.OpsiC,
		OpsiD:      soal.OpsiD,
		OpsiE:      soal.OpsiE,
		GambarA:    soal.GambarA,
		GambarB:    soal.GambarB,
		GambarC:    soal.GambarC,
		GambarD:    soal.GambarD,
		GambarE:    soal.GambarE,
		Kunci:      soal.Kunci,
		CreatedAt:  soal.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  soal.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
