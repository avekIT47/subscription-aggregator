package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	UserID      string     `gorm:"type:varchar(255);index" json:"user_id"`
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	ServiceName string     `gorm:"index" json:"service_name"`
	Price       uint       `gorm:"index" json:"price"`
	StartDate   time.Time  `gorm:"index" json:"start_date"`
	EndDate     *time.Time `gorm:"index" json:"end_date,omitempty"`
}
