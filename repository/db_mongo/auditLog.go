package dbmongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Admin_ID    string             `json:"admin_id" bson:"admin_id"`
	AuditLog_ID uint64             `json:"audit_log_id" bson:"audit_log_id"`
	Action      string             `json:"action" bson:"action"`
	EntityID    string             `json:"entity_id" bson:"entity_id"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	IPAddress   string             `json:"ip_address" bson:"ip_address"`
	Model
}
