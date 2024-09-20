package dbmongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mohon dicek 🌞🙏

type AuditLog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Admin_ID    string             `bson:"admin_id"`
	AuditLog_ID uint64             `bson:"audit_log_id"`
	Action      string             `bson:"action"`
	EntityID    string             `bson:"entity_id"`
	Timestamp   time.Time          `bson:"timestamp"`
	IPAddress   string             `bson:"ip_address"`
	Model
}
