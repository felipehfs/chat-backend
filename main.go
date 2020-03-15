package main

import (
	"github.com/felipehfs/api/chat/config"
	"github.com/gorilla/mux"
)

func main() {
	mongo, err := config.CreateMongoDB()
	if err != nil {
		panic(err)
	}
	defer mongo.Close()
	r := mux.NewRouter().StrictSlash(true)
	server := config.NewServer(mongo, r)
	server.Run(":8080")
}
