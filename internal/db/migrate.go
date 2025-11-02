package db

import (
	"fmt"

	"gorm.io/gorm"
)

func Run(db *gorm.DB, models ...any) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 0) Extensions
		if err := tx.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
			return fmt.Errorf("enable pgcrypto: %w", err)
		}
		if err := tx.Exec(`CREATE EXTENSION IF NOT EXISTS citext`).Error; err != nil {
			return fmt.Errorf("enable citext: %w", err)
		}

		// 1) AutoMigrate (pass parents before children in caller)
		if err := tx.AutoMigrate(models...); err != nil {
			return fmt.Errorf("automigrate: %w", err)
		}
		return nil
	})
}
