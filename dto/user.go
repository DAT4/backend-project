package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"regexp"
	"unicode"
)

type Username string
type Password string
type Email string
type Mac string
type Ip string

type User struct {
	Id       string   `json:"-" bson:"_id,omitempty"`
	PlayerID int      `json:"-"`
	Username Username `json:"username"`
	Password Password `json:"password"`
	Email    Email    `json:"email"`
	Macs     []Mac    `json:"-"`
	Ips      []Ip     `json:"-"`
	Role     string   `json:"-"`
}

func (u *User) HashAndSalt() error {
	bytePwd := []byte(u.Password)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.Password = Password(hash)
	return nil
}

func (u *User) Check(hashedPassword Password) error {
	bytePwd := []byte(u.Password)
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		fmt.Println(err)
		return errors.New("password incorrect")
	}
	return nil
}

func (password Password) Validate() error {
	var upp, low, num bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upp = true
		case unicode.IsLower(char):
			low = true
		case unicode.IsNumber(char):
			num = true
		default:
			continue
		}
	}

	if !upp || !low || !num {
		return errors.New("password does now have required symbols")
	}

	if len(string(password)) > 20 || len(string(password)) < 4 {
		return errors.New("password has to be between 40 and 8 chars")
	}

	return nil
}

func (username Username) Validate() error {
	re, _ := regexp.Compile(`^[a-z]{4,20}$`)
	ok := re.MatchString(string(username))
	if !ok {
		return errors.New("username is invalid")
	}
	return nil
}

func (email Email) Validate() error {
	re, _ := regexp.Compile(`^\w+@\w+\.\w+$`)
	ok := re.MatchString(string(email))
	if !ok {
		return errors.New("email is invalid")
	}
	return nil
}

func UserFromJson(data io.ReadCloser) (user User, err error) {
	err = json.NewDecoder(data).Decode(&user)
	return
}

func (u User) ToJson(data io.Writer) error {
	return json.NewEncoder(data).Encode(u)
}
