package models

import "gopkg.in/mgo.v2/bson"

// User represents the credentials
type User struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string        `json:"username,omitempty"`
	Password  string        `json:"password,omitempty"`
	AvatarURL string        `json:"avatar_url,omitempty" bson:"avatarURL,omitempty"`
}
