package mongo

import (
	"errors"
	"fmt"
	"github.com/DAT4/backend-project/models/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(user *user.User) (err error) {
	q2 := AddOneQuery{
		Model:      &user,
		Filter:     nil,
		Collection: "users",
	}
	err = q2.Add()
	if err != nil {
		return err
	}
	return nil
}

func UserFromId(id string) (user user.User, err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	q := FindOneQuery{
		Model:      &user,
		Filter:     bson.M{"_id": _id},
		Collection: "users",
	}
	err = q.Find()
	fmt.Println(user)
	return
}

func Authenticate(u *user.User) error {
	var tmpUser user.User
	q := FindOneQuery{
		Model: &tmpUser,
		Filter: bson.M{
			"username": u.Username,
		},
		Collection: "users",
	}

	err := q.Find()
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
	q1 := FindOneQuery{
		Model:      &tmpUser,
		Filter:     bson.M{"username": u.Username},
		Options:    options.FindOne(),
		Collection: "users",
	}
	err = q1.Find()
	if err == nil {
		return errors.New("A user already exists with this name")
	}
	return nil
}
