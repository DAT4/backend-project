package models

import (
	"errors"
	"regexp"
	"unicode"
)

type Username string
type Password string
type Email string
type Mac string
type Ip string

func (password Password) validate() error {
	var upp, low, num, sym bool

	for _, char := range password{
		switch {
		case unicode.IsUpper(char):
			upp = true
		case unicode.IsLower(char):
			low = true
		case unicode.IsNumber(char):
			num = true
		case unicode.IsPunct(char):
			sym = true
		default:
			return errors.New("password contains some shit")
		}
	}

	if !upp || !low || !num || !sym {
		return errors.New("password does now have required symbols")
	}

	if len(string(password)) > 20 || len(string(password)) < 4 {
		return errors.New("password has to be between 40 and 8 chars")
	}

	return nil
}

func (username Username) validate() error {
	re, _ := regexp.Compile(`^[a-z]{4,20}$`)
	ok := re.MatchString(string(username))
	if !ok {
		return errors.New("username is invalid")
	}
	return nil
}

func (email Email) validate() error {
	re, _ := regexp.Compile(`^s\d{6}@student\.dtu\.dk$`)
	ok := re.MatchString(string(email))
	if !ok {
		return errors.New("username is invalid")
	}
	return nil
}

func (mac Mac) validate() error{
	re, _ := regexp.Compile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	ok := re.MatchString(string(mac))
	if !ok {
		return errors.New("mac address is invalid")
	}
	return nil
}

func (ip Ip) validate() error{
	re, _ := regexp.Compile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	ok := re.MatchString(string(ip))
	if !ok {
		return errors.New("ip address is invalid")
	}
	return nil
}
