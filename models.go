package main

import (
	"encoding/json"
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
	Validation func() bool
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
	for _,ip := range user.Ips{
		err = ip.validate()
		if err != nil {
			return err
		}
	}
	for _,mac:= range user.Macs{
		err = mac.validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *User) authenticate() error {
	q := FindOneQuery{
		Model: user,
		Filter: bson.M{
			"username": user.Username,
			"password": user.Password,
		},
		Collection: "users",
	}

	err := q.find()
	if err != nil {
		return err
	}
	return nil
}
