package middle

import (
	"errors"
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/models"
	"regexp"
	"unicode"
)

func Validate(user models.User) error {
	var err error
	err = dao.UsernameTaken(&user)
	if err != nil {
		return err
	}
	err = validatePassword(user.Password)
	if err != nil {
		return err
	}
	err = validateUsername(user.Username)
	if err != nil {
		return err
	}
	err = validateEmail(user.Email)
	if err != nil {
		return err
	}
	return nil
}

func validatePassword(password models.Password) error {
	var upp, low, num, sym bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upp = true
		case unicode.IsLower(char):
			low = true
		case unicode.IsNumber(char):
			num = true
		case unicode.IsPunct(char):
			sym = true
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

func validateUsername(username models.Username) error {
	re, _ := regexp.Compile(`^[a-z]{4,20}$`)
	ok := re.MatchString(string(username))
	if !ok {
		return errors.New("username is invalid")
	}
	return nil
}

func validateEmail(email models.Email) error {
	re, _ := regexp.Compile(`^\w+@\w+\.\w+$`)
	ok := re.MatchString(string(email))
	if !ok {
		return errors.New("username is invalid")
	}
	return nil
}

func validateMac(mac models.Mac) error {
	re, _ := regexp.Compile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	ok := re.MatchString(string(mac))
	if !ok {
		return errors.New("mac address is invalid")
	}
	return nil
}

func validateIp(ip models.Ip) error {
	re, _ := regexp.Compile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	ok := re.MatchString(string(ip))
	if !ok {
		return errors.New("ip address is invalid")
	}
	return nil
}
