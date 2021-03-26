package dao

import (
	"errors"
	"fmt"
	"github.com/DAT4/backend-project/models/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(user *user.User) (err error) {
	q2 := addOneQuery{
		Model:      &user,
		Filter:     nil,
		Collection: "users",
	}
	return q2.add()
}

func UserFromId(id string) (user user.User, err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	q := findOneQuery{
		Model:      &user,
		Filter:     bson.M{"_id": _id},
		Collection: "users",
	}
	err = q.find()
	fmt.Println(user)
	return
}

func Authenticate(u *user.User) error {
	var tmpUser user.User
	q := findOneQuery{
		Model: &tmpUser,
		Filter: bson.M{
			"username": u.Username,
		},
		Collection: "users",
	}

	err := q.find()
	if err != nil {
		return err
	}

	ok := u.Check(tmpUser.Password)
	if !ok {
		return errors.New("password incorrect")
	}
	u.Password = tmpUser.Password
	return nil
}

func UsernameTaken(u *user.User) (err error) {
	var tmpUser user.User
	q1 := findOneQuery{
		Model:      &tmpUser,
		Filter:     bson.M{"username": u.Username},
		Options:    options.FindOne(),
		Collection: "users",
	}
	err = q1.find()
	if err == nil {
		return errors.New("A user already exists with this name")
	}
	return nil
}
