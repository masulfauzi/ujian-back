package seeders

import (
	"backend/internal/constants"
	"backend/internal/modules/user/model"
	"backend/internal/utils"
	"log"

	"gorm.io/gorm"
)

// SeedAdmin membuat akun admin default. Aman dijalankan berulang: dilewati jika username sudah ada.
func SeedAdmin(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Where("username = ?", "admin").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("Akun admin sudah ada, seeding admin dilewati")
		return nil
	}

	hashedPassword, err := utils.HashPassword("klalen")
	if err != nil {
		return err
	}

	return db.Create(&model.User{
		Name:     "Administrator",
		Username: "admin",
		Email:    "admin@ujian.local",
		Password: hashedPassword,
		Role:     constants.RoleAdmin,
		Status:   constants.StatusActive,
	}).Error
}
