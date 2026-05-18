package service

import (
	"errors"
	"fmt"
	"math"

	"backend/internal/config"
	"backend/internal/constants"
	"backend/internal/modules/soal/dto"
	"backend/internal/modules/soal/model"
	"backend/internal/modules/soal/repository"
	"backend/internal/utils"

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

	var gambarSoalName string
	if req.GambarSoal != nil {
		filename, err := utils.SaveImage(req.GambarSoal, "soal")
		if err != nil {
			return nil, errors.New("gagal upload gambar_soal: " + err.Error())
		}
		gambarSoalName = filename
	}

	var gambarAName string
	if req.GambarA != nil {
		filename, err := utils.SaveImage(req.GambarA, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			return nil, errors.New("gagal upload gambar_a: " + err.Error())
		}
		gambarAName = filename
	}

	var gambarBName string
	if req.GambarB != nil {
		filename, err := utils.SaveImage(req.GambarB, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			utils.DeleteImage(gambarAName, "opsi")
			return nil, errors.New("gagal upload gambar_b: " + err.Error())
		}
		gambarBName = filename
	}

	var gambarCName string
	if req.GambarC != nil {
		filename, err := utils.SaveImage(req.GambarC, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			utils.DeleteImage(gambarAName, "opsi")
			utils.DeleteImage(gambarBName, "opsi")
			return nil, errors.New("gagal upload gambar_c: " + err.Error())
		}
		gambarCName = filename
	}

	var gambarDName string
	if req.GambarD != nil {
		filename, err := utils.SaveImage(req.GambarD, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			utils.DeleteImage(gambarAName, "opsi")
			utils.DeleteImage(gambarBName, "opsi")
			utils.DeleteImage(gambarCName, "opsi")
			return nil, errors.New("gagal upload gambar_d: " + err.Error())
		}
		gambarDName = filename
	}

	var gambarEName string
	if req.GambarE != nil {
		filename, err := utils.SaveImage(req.GambarE, "opsi")
		if err != nil {
			utils.DeleteImage(gambarSoalName, "soal")
			utils.DeleteImage(gambarAName, "opsi")
			utils.DeleteImage(gambarBName, "opsi")
			utils.DeleteImage(gambarCName, "opsi")
			utils.DeleteImage(gambarDName, "opsi")
			return nil, errors.New("gagal upload gambar_e: " + err.Error())
		}
		gambarEName = filename
	}

	soal := &model.Soal{
		IdBankSoal: req.IdBankSoal,
		NoSoal:     req.NoSoal,
		Soal:       req.Soal,
		GambarSoal: gambarSoalName,
		OpsiA:      req.OpsiA,
		OpsiB:      req.OpsiB,
		OpsiC:      req.OpsiC,
		OpsiD:      req.OpsiD,
		OpsiE:      req.OpsiE,
		GambarA:    gambarAName,
		GambarB:    gambarBName,
		GambarC:    gambarCName,
		GambarD:    gambarDName,
		GambarE:    gambarEName,
		Kunci:      req.Kunci,
	}

	if err := s.repo.Create(soal); err != nil {
		utils.DeleteImage(gambarSoalName, "soal")
		utils.DeleteImage(gambarAName, "opsi")
		utils.DeleteImage(gambarBName, "opsi")
		utils.DeleteImage(gambarCName, "opsi")
		utils.DeleteImage(gambarDName, "opsi")
		utils.DeleteImage(gambarEName, "opsi")
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

	if req.GambarSoal != nil {
		newFilename, err := utils.SaveImage(req.GambarSoal, "soal")
		if err != nil {
			return nil, errors.New("gagal upload gambar_soal: " + err.Error())
		}
		if soal.GambarSoal != "" {
			utils.DeleteImage(soal.GambarSoal, "soal")
		}
		soal.GambarSoal = newFilename
	}

	if req.GambarA != nil {
		newFilename, err := utils.SaveImage(req.GambarA, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_a: " + err.Error())
		}
		if soal.GambarA != "" {
			utils.DeleteImage(soal.GambarA, "opsi")
		}
		soal.GambarA = newFilename
	}

	if req.GambarB != nil {
		newFilename, err := utils.SaveImage(req.GambarB, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_b: " + err.Error())
		}
		if soal.GambarB != "" {
			utils.DeleteImage(soal.GambarB, "opsi")
		}
		soal.GambarB = newFilename
	}

	if req.GambarC != nil {
		newFilename, err := utils.SaveImage(req.GambarC, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_c: " + err.Error())
		}
		if soal.GambarC != "" {
			utils.DeleteImage(soal.GambarC, "opsi")
		}
		soal.GambarC = newFilename
	}

	if req.GambarD != nil {
		newFilename, err := utils.SaveImage(req.GambarD, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_d: " + err.Error())
		}
		if soal.GambarD != "" {
			utils.DeleteImage(soal.GambarD, "opsi")
		}
		soal.GambarD = newFilename
	}

	if req.GambarE != nil {
		newFilename, err := utils.SaveImage(req.GambarE, "opsi")
		if err != nil {
			return nil, errors.New("gagal upload gambar_e: " + err.Error())
		}
		if soal.GambarE != "" {
			utils.DeleteImage(soal.GambarE, "opsi")
		}
		soal.GambarE = newFilename
	}

	soal.NoSoal = req.NoSoal
	soal.Soal = req.Soal
	soal.OpsiA = req.OpsiA
	soal.OpsiB = req.OpsiB
	soal.OpsiC = req.OpsiC
	soal.OpsiD = req.OpsiD
	soal.OpsiE = req.OpsiE
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

	utils.DeleteImage(soal.GambarSoal, "soal")
	utils.DeleteImage(soal.GambarA, "opsi")
	utils.DeleteImage(soal.GambarB, "opsi")
	utils.DeleteImage(soal.GambarC, "opsi")
	utils.DeleteImage(soal.GambarD, "opsi")
	utils.DeleteImage(soal.GambarE, "opsi")

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
		NoSoal:     soal.NoSoal,
		Soal:       soal.Soal,
		GambarSoal: s.buildImageURL(soal.GambarSoal, "soal"),
		OpsiA:      soal.OpsiA,
		OpsiB:      soal.OpsiB,
		OpsiC:      soal.OpsiC,
		OpsiD:      soal.OpsiD,
		OpsiE:      soal.OpsiE,
		GambarA:    s.buildImageURL(soal.GambarA, "opsi"),
		GambarB:    s.buildImageURL(soal.GambarB, "opsi"),
		GambarC:    s.buildImageURL(soal.GambarC, "opsi"),
		GambarD:    s.buildImageURL(soal.GambarD, "opsi"),
		GambarE:    s.buildImageURL(soal.GambarE, "opsi"),
		Kunci:      soal.Kunci,
		CreatedAt:  soal.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  soal.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *soalService) buildImageURL(filename, folder string) string {
	if filename == "" {
		return ""
	}
	cfg := config.GetUploadConfig()
	return fmt.Sprintf("%s/%s/%s", cfg.ImageBaseURL, folder, filename)
}
