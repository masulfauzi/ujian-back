package repository

import (
	"backend/internal/modules/user/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}
