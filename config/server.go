package config

import (
	"net/http"

	"github.com/felipehfs/api/chat/controllers"
	"github.com/felipehfs/api/chat/repositories"
	"github.com/felipehfs/api/chat/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

// Server .
type Server struct {
	DB    *mgo.Session
	Route *mux.Router
}

// NewServer instantiates the server
func NewServer(db *mgo.Session, route *mux.Router) *Server {
	return &Server{DB: db, Route: route}
}

// Run the server on port described
func (s *Server) Run(port string) {

	userDao := repositories.NewUserDAO(s.DB)
	ws := services.NewWebsocketClient()

	userHandler := controllers.NewUserHandler(userDao)
	s.Route.Handle("/register", http.HandlerFunc(userHandler.Register)).Methods("POST")
	s.Route.Handle("/login", http.HandlerFunc(userHandler.Login)).Methods("POST")
	s.Route.Handle("/users/{id}/avatar", http.HandlerFunc(userHandler.UpdateAvatar)).Methods("PUT")

	dir := "/src/statics"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	handler := cors.Default().Handler(s.Route)

	http.Handle("/", handler)
	http.Handle("/chat", ws)

	http.ListenAndServe(port, nil)
}
