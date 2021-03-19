package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DAT4/backend-project/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username Username
	Password Password
	Email    Email
	Macs     []Mac
	Ips      []Ip
}

func UserFromJson(data io.ReadCloser) (user User, err error) {
	err = json.NewDecoder(data).Decode(&user)
	return
}

func UserFromToken(r *http.Request) (user User, err error) {
	token, err := ExtractJWTToken(r)
	if err != nil {
		return
	}
	id, err := ExtractClaims(token)
	if err != nil {
		return
	}
	user, err = userFromId(id)
	return
}

func userFromId(id string) (user User, err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	q := database.FindOneQuery{
		Model:      &user,
		Filter:     bson.M{"_id": _id},
		Collection: "users",
	}
	err = q.Find()
	fmt.Println(user)
	return
}

func (user *User) Create() (err error) {
	q2 := database.AddOneQuery{
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

func (user *User) Validate() error {
	var err error
	err = user.UsernameTaken()
	if err != nil {
		return err
	}
	err = user.Password.validate()
	if err != nil {
		return err
	}
	err = user.Username.validate()
	if err != nil {
		return err
	}
	err = user.Email.validate()
	if err != nil {
		return err
	}
	if len(user.Ips) == 0 {
		return errors.New("ip is missing")
	}
	for _, ip := range user.Ips {
		fmt.Println(ip)
		err = ip.validate()
		if err != nil {
			return err
		}
	}
	if len(user.Macs) == 0 {
		return errors.New("mac address is missing")
	}
	for _, mac := range user.Macs {
		err = mac.validate()
		if err != nil {
			return err
		}
	}
	fmt.Println("Dont checking")
	return nil
}

func (user *User) Authenticate() (User, error) {
	var tmpUser User
	q := database.FindOneQuery{
		Model: &tmpUser,
		Filter: bson.M{
			"username": user.Username,
		},
		Collection: "users",
	}

	err := q.Find()
	if err != nil {
		return tmpUser, err
	}

	ok := user.check(tmpUser.Password)
	if !ok {
		return tmpUser, errors.New("password incorrect")
	}
	return tmpUser, nil
}

func (user *User) HashAndSalt() error {
	bytePwd := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = Password(hash)
	return nil
}

func (user *User) check(hashedPassword Password) bool {
	bytePwd := []byte(user.Password)
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (user *User) UsernameTaken() (err error) {
	var tmpUser User
	q1 := database.FindOneQuery{
		Model:      &tmpUser,
		Filter:     bson.M{"username": user.Username},
		Options:    options.FindOne(),
		Collection: "users",
	}
	err = q1.Find()
	if err == nil {
		return errors.New("A user already exists with this name")
	}
	return nil
}

