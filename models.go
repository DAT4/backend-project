package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
)

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username Username
	Password Password
	Email    Email
	Macs     []Mac
	Ips      []Ip
}

type FindOneQuery struct {
	Model      interface{}
	Filter     bson.M
	Options    *options.FindOneOptions
	Collection string
}

type AddOneQuery struct {
	Model      interface{}
	Filter     bson.M
	Collection string
}

func (user *User) fromJson(data io.ReadCloser) error {
	err := json.NewDecoder(data).Decode(user)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) validate() error {
	var err error
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

func (user *User) authenticate() error {
	var tmpUser User
	q := FindOneQuery{
		Model: &tmpUser,
		Filter: bson.M{
			"username": user.Username,
		},
		Collection: "users",
	}

	err := q.find()
	if err != nil {
		return err
	}

	ok := user.check(tmpUser.Password)
	if !ok {
		return errors.New("password incorrect")
	}
	return nil
}
