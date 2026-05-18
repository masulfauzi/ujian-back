package database

import (
	"backend/internal/database/seeders"
	"fmt"

	"gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) error {
	if err := seeders.SeedMapel(db); err != nil {
		return fmt.Errorf("failed to seed mapel: %w", err)
	}

	return nil
}
