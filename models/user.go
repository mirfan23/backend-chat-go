package models

import "time"

type User struct {
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"-"`
	PublicKey string    `bson:"publicKey" json:"publicKey"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
