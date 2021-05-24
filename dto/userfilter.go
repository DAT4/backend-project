package dto

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"io"
	"net/http"
)

type UserFilter struct {
	Id       string `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
}

func UserFilterFromJson(data io.ReadCloser) (filter UserFilter, err error) {
	err = json.NewDecoder(data).Decode(&filter)
	return
}

func (f UserFilter) ToJson(data io.Writer) error {
	return json.NewEncoder(data).Encode(f)
}

func (f *UserFilter) Decode(r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}
	defer r.Body.Close()
	return schema.NewDecoder().Decode(f, r.Form)
}
