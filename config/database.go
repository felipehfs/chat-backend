package config

import (
	"os"

	"gopkg.in/mgo.v2"
)

// CreateMongoDB setups the mongoDB for us
func CreateMongoDB() (*mgo.Session, error) {
	return mgo.Dial(os.Getenv("MONGODB_URL"))
}
