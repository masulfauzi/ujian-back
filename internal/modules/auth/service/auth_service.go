package service

import (
	"errors"
	"time"

	"backend/configs"
	"backend/internal/constants"
	"backend/internal/modules/auth/dto"
	"backend/internal/modules/auth/repository"
	"backend/internal/modules/user/model"
	"backend/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GenerateToken(user *model.User) (string, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	existing, _ := s.repo.GetUserByEmail(req.Email)
	if existing != nil {
		return nil, errors.New(constants.ErrDuplicateEmail)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     constants.RoleUser,
		Status:   constants.StatusActive,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
		Role:  user.Role,
	}, nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if !utils.VerifyPassword(user.Password, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
		Role:  user.Role,
	}, nil
}

func (s *authService) GenerateToken(user *model.User) (string, error) {
	jwtConfig := configs.GetJWTConfig()

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(jwtConfig.Expired).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	jwtConfig := configs.GetJWTConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
