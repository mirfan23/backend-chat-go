package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomId    string             `bson:"roomId" json:"roomId"`
	Sender    string             `bson:"sender" json:"sender"`
	Text      string             `bson:"text" json:"text"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	IsRead    *bool              `bson:"is_read,omitempty" json:"is_read"`
}
