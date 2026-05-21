package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/peserta/dto"
	"backend/internal/modules/peserta/model"
	"backend/internal/modules/peserta/repository"
	"backend/internal/utils"

	"gorm.io/gorm"
)

func pesertaWithKelasToResponse(p *repository.PesertaWithKelas) *dto.PesertaResponse {
	return &dto.PesertaResponse{
		ID:        p.ID,
		Nama:      p.Nama,
		IDKelas:   p.IDKelas,
		NamaKelas: p.NamaKelas,
		Username:  p.Username,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type PesertaService interface {
	CreatePeserta(req *dto.CreatePesertaRequest) (*dto.PesertaResponse, error)
	GetPesertaByID(id string) (*dto.PesertaResponse, error)
	GetAllPeserta(page, pageSize int, idKelas string) (*dto.PesertaListResponse, error)
	UpdatePeserta(id string, req *dto.UpdatePesertaRequest) (*dto.PesertaResponse, error)
	DeletePeserta(id string) error
	RestorePeserta(id string) error
}

type pesertaService struct {
	repo repository.PesertaRepository
}

func NewPesertaService(repo repository.PesertaRepository) PesertaService {
	return &pesertaService{repo: repo}
}

func (s *pesertaService) CreatePeserta(req *dto.CreatePesertaRequest) (*dto.PesertaResponse, error) {
	existing, err := s.repo.GetByUsername(req.Username)
	if err == nil && existing != nil {
		return nil, errors.New("username sudah digunakan")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	peserta := &model.Peserta{
		Nama:     req.Nama,
		IDKelas:  req.IDKelas,
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := s.repo.Create(peserta); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByID(peserta.ID)
	if err != nil {
		return nil, err
	}

	return pesertaWithKelasToResponse(created), nil
}

func (s *pesertaService) GetPesertaByID(id string) (*dto.PesertaResponse, error) {
	peserta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return pesertaWithKelasToResponse(peserta), nil
}

func (s *pesertaService) GetAllPeserta(page, pageSize int, idKelas string) (*dto.PesertaListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	pesertaList, total, err := s.repo.GetAll(page, pageSize, idKelas)
	if err != nil {
		return nil, err
	}

	var responses []dto.PesertaResponse
	for _, p := range pesertaList {
		responses = append(responses, *pesertaWithKelasToResponse(&p))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PesertaListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *pesertaService) UpdatePeserta(id string, req *dto.UpdatePesertaRequest) (*dto.PesertaResponse, error) {
	existing, err := s.repo.GetRawByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if req.Username != existing.Username {
		taken, err := s.repo.GetByUsername(req.Username)
		if err == nil && taken != nil && taken.ID != id {
			return nil, errors.New("username sudah digunakan")
		}
	}

	peserta := &model.Peserta{
		ID:       existing.ID,
		Nama:     req.Nama,
		IDKelas:  req.IDKelas,
		Username: req.Username,
		Password: existing.Password,
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("gagal memproses password")
		}
		peserta.Password = hashedPassword
	}

	if err := s.repo.Update(peserta); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return pesertaWithKelasToResponse(updated), nil
}

func (s *pesertaService) DeletePeserta(id string) error {
	peserta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(peserta.ID)
}

func (s *pesertaService) RestorePeserta(id string) error {
	return s.repo.Restore(id)
}
