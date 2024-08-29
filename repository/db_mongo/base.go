package dbmongo

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	CreatedAt time.Time			`bson:"created_at"`
	UpdatedAt time.Time 		`bson:"updated_at"`
	DeletedAt gorm.DeletedAt	`gorm:"index" bson:"deleted_at"`
}