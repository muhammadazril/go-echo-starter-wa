package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobLog struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	JobID     string             `bson:"job_id" json:"job_id" validate:"required"`
	Status    string             `bson:"status" json:"status" validate:"required"`
	To        string             `bson:"to" json:"to" validate:"required"`
	Messages  string             `bson:"message" json:"message" validate:"required"`
	Error     string             `bson:"error_message" json:"error_message"`
	// UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
}
