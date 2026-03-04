package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Sender    string             `bson:"sender" json:"sender"`
	Receiver  string             `bson:"receiver" json:"receiver"`
	Message   string             `bson:"message" json:"message"`
	IsRead    *bool              `bson:"is_read,omitempty" json:"is_read"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}
