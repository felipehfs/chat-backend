package repositories

import (
	"github.com/felipehfs/api/chat/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserDAO manipulates the database
type UserDAO struct {
	Db *mgo.Session
}

// NewUserDAO instantiates
func NewUserDAO(db *mgo.Session) *UserDAO {
	return &UserDAO{
		Db: db,
	}
}

// Register .
func (dao UserDAO) Register(user models.User) error {
	return dao.Db.DB("chat-api").C("users").Insert(&user)
}

// FindById
func (dao UserDAO) FindById(id bson.ObjectId) (*models.User, error) {
	var user models.User
	err := dao.Db.DB("chat-api").C("users").Find(bson.M{
		"_id": id,
	}).One(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateAvatar sets the avatar filename
func (dao UserDAO) UpdateAvatar(id bson.ObjectId, filename string) (*models.User, error) {
	err := dao.Db.DB("chat-api").C("users").Update(bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"avatarURL": filename,
		},
	})

	if err != nil {
		return nil, err
	}

	var user models.User

	err = dao.Db.DB("chat-api").C("users").Find(bson.M{
		"_id": id,
	}).One(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Login .
func (dao UserDAO) Login(user models.User) (*models.User, error) {
	var search models.User
	err := dao.Db.DB("chat-api").C("users").Find(&bson.M{
		"username": user.Username,
	}).One(&search)

	if err != nil {
		return nil, err
	}

	return &search, nil
}
