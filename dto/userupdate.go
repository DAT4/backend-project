package dto

import (
	"encoding/json"
	"io"
)

type UserUpdate struct {
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

func UserUpdateFromJson(data io.ReadCloser) (filter UserUpdate, err error) {
	err = json.NewDecoder(data).Decode(&filter)
	return
}

func (u UserUpdate) ToJson(data io.Writer) error {
	return json.NewEncoder(data).Encode(u)
}
