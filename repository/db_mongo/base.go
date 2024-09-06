package dbmongo

import (
	"time"

)

type Model struct {
	CreatedAt time.Time			`bson:"created_at"`
	UpdatedAt time.Time 		`bson:"updated_at"`
	DeletedAt *time.Time		`bson:"deleted_at,omitempty"`
}