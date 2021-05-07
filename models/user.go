package models

import (
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

func (user *User) Check(hashedPassword Password) bool {
	bytePwd := []byte(user.Password)
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}
