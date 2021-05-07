package dao

import (
	"errors"
	"github.com/DAT4/backend-project/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestDB struct {
	users []models.User
}

func (t *TestDB) Create(userIn models.User) (userOut models.User, err error) {
	userIn.Id = primitive.NewObjectID()
	t.users = append(t.users, userIn)
	return userIn, nil
}
func (t *TestDB) UserFromId(id string) (user models.User, err error) {
	for _, u := range t.users {
		if id == u.Id.Hex() {
			return u, nil
		}
	}
	return models.User{}, errors.New("no users with this id")
}
func (t *TestDB) UsernameTaken(u *models.User) (err error) {
	for _, dbu := range t.users {
		if dbu.Username == u.Username {
			return errors.New("user with this username exist")
		}
	}
	return nil
}
func (t *TestDB) UserFromName(name string) (user models.User, err error) {
	for _, dbu := range t.users {
		if dbu.Username == models.Username(name) {
			return dbu, nil
		}
	}
	err = errors.New("user not found")
	return
}

func NewTestDB() *TestDB {
	return &TestDB{}

}
