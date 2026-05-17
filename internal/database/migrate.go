package database

import (
	"backend/internal/modules/user/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
	)
}
