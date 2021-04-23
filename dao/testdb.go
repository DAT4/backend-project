package dao

import (
	"errors"
	"github.com/DAT4/backend-project/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestDB struct {
	users []models.User
}

func (t *TestDB) Create(u *models.User) (err error) {
	u.Id = primitive.NewObjectID()
	t.users = append(t.users, *u)
	return nil
}
func (t *TestDB) UserFromId(id string) (user models.User, err error) {
	for _, u := range t.users {
		if id == u.Id.String() {
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
func (t *TestDB) Authenticate(u *models.User) (err error) {
	for _, dbu := range t.users {
		if dbu.Username == u.Username {
			if dbu.Password == u.Password {
				return nil
			}
			return errors.New("wrong password")
		}
	}
	return errors.New("user not found")
}
