package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"`
	PublicKey string             `bson:"publicKey" json:"publicKey"`
	Photo     string             `bson:"photo" json:"photo"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
