package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/felipehfs/api/chat/models"
	"github.com/felipehfs/api/chat/repositories"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserHandler represents the handle
type UserHandler struct {
	Dao *repositories.UserDAO
}

// NewUserHandler .
func NewUserHandler(dao *repositories.UserDAO) *UserHandler {
	return &UserHandler{
		Dao: dao,
	}
}

// Register adds new users for api
func (handler UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = bson.NewObjectIdWithTime(time.Now())

	if err := handler.Dao.Register(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := handler.Dao.FindById(user.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Usuário criado com sucesso!",
		"data":    created,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateAvatar updates the avatars
func (handler UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseMultipartForm(5 << 20)
	file, fileHandler, err := r.FormFile("avatar")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	destination, err := os.OpenFile("/src/statics/"+fileHandler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.Copy(destination, file)

	id := bson.ObjectIdHex(vars["id"])

	_, err = handler.Dao.UpdateAvatar(id, fileHandler.Filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Arquivo carregado com sucesso!",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login .
func (handler UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := handler.Dao.Login(user)

	if result.AvatarURL != "" {
		result.AvatarURL = os.Getenv("STATIC_URL") + result.AvatarURL
	}
	if err != nil {

		if err == mgo.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result.Password = ""

	json.NewEncoder(w).Encode(result)
}
