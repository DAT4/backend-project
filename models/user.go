package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Username string
type Password string
type Email string
type Mac string
type Ip string

type User struct {
	Id       primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	PlayerID int                `json:"-"`
	Username Username           `json:"username"`
	Password Password           `json:"password"`
	Email    Email              `json:"email"`
	Macs     []Mac              `json:"-"`
	Ips      []Ip               `json:"-"`
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

func (user *User) Check(hashedPassword Password) bool {
	bytePwd := []byte(user.Password)
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
