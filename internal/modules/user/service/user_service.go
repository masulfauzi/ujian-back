package service

import (
	"errors"
	"strings"

	"backend/internal/constants"
	"backend/internal/modules/user/dto"
	"backend/internal/modules/user/model"
	"backend/internal/modules/user/repository"
	"backend/internal/utils"

	"gorm.io/gorm"
)

type UserService interface {
	Create(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetAll() ([]dto.UserResponse, error)
	GetByID(id string) (*dto.UserResponse, error)
	Update(id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(id string) error
	GetByEmail(email string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, errors.New("username tidak boleh kosong")
	}
	if existing, _ := s.repo.GetByUsername(username); existing != nil {
		return nil, errors.New("username sudah terdaftar")
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		// Kolom email unique & not null; guru yang login via sianida tidak wajib punya email
		email = username + "@sianida.local"
	}
	if existing, _ := s.repo.GetByEmail(email); existing != nil {
		return nil, errors.New(constants.ErrDuplicateEmail)
	}

	// Password kosong = login via sianida
	hashedPassword := ""
	if req.Password != "" {
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		hashedPassword = hashed
	}

	role := req.Role
	if role == "" {
		role = constants.RoleAdmin
	}

	user := &model.User{
		Username: username,
		Name:     req.Name,
		Email:    email,
		Password: hashedPassword,
		Role:     role,
		Status:   constants.StatusActive,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return s.modelToResponse(user), nil
}

func (s *userService) GetAll() ([]dto.UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, *s.modelToResponse(&user))
	}

	return responses, nil
}

func (s *userService) GetByID(id string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.modelToResponse(user), nil
}

func (s *userService) Update(id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		if existing, _ := s.repo.GetByEmail(req.Email); existing != nil {
			return nil, errors.New(constants.ErrDuplicateEmail)
		}
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.Password != nil {
		if *req.Password == "" {
			// Pindah ke login via sianida
			user.Password = ""
		} else {
			if len(*req.Password) < 6 {
				return nil, errors.New("password minimal 6 karakter")
			}
			hashed, err := utils.HashPassword(*req.Password)
			if err != nil {
				return nil, err
			}
			user.Password = hashed
		}
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.modelToResponse(user), nil
}

func (s *userService) Delete(id string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(user.ID)
}

func (s *userService) GetByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) modelToResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		Status:       user.Status,
		LoginSianida: user.Password == "",
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
