package dao

import (
	"errors"
	"fmt"
	"github.com/DAT4/backend-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(user *models.User) (err error) {
	q2 := addOneQuery{
		model:      &user,
		filter:     nil,
		collection: "users",
	}
	return q2.add()
}

func UserFromId(id string) (user models.User, err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	q := findOneQuery{
		model:      &user,
		filter:     bson.M{"_id": _id},
		collection: "users",
	}
	err = q.find()
	fmt.Println(user)
	return
}

func Authenticate(u *models.User) error {
	var tmpUser models.User
	q := findOneQuery{
		model: &tmpUser,
		filter: bson.M{
			"username": u.Username,
		},
		collection: "users",
	}

	err := q.find()
	if err != nil {
		return err
	}

	//TODO This logic should be in the middle package
	ok := u.Check(tmpUser.Password)
	if !ok {
		return errors.New("password incorrect")
	}
	u.Id = tmpUser.Id
	u.Password = tmpUser.Password
	return nil
}

func UsernameTaken(u *models.User) (err error) {
	var tmpUser models.User
	q1 := findOneQuery{
		model:      &tmpUser,
		filter:     bson.M{"username": u.Username},
		options:    options.FindOne(),
		collection: "users",
	}
	err = q1.find()
	if err == nil {
		return errors.New("A user already exists with this name")
	}
	return nil
}
