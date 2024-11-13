package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLog struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AdminID   primitive.ObjectID `json:"admin_id" bson:"admin_id"`
	Action    string             `json:"action" bson:"action"`
	Entity    string             `json:"entity" bson:"entity"`
	EntityID  string             `json:"entity_id" bson:"entity_id"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	IPAddress string             `json:"ip_address" bson:"ip_address"`
}
