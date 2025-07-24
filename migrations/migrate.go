package migrations

import (
	"subscription-aggregator/internal/model"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	return db.AutoMigrate(&model.Subscription{})
}
