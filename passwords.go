package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func (user *User) hashAndSalt() error {
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
