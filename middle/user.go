package middle

import (
	"encoding/json"
	"github.com/DAT4/backend-project/dao/mongo"
	"github.com/DAT4/backend-project/models/user"
	"io"
	"net/http"
)

func UserFromJson(data io.ReadCloser) (user user.User, err error) {
	err = json.NewDecoder(data).Decode(&user)
	return
}

func UserFromToken(r *http.Request) (user user.User, err error) {
	token, err := ExtractJWTToken(r)
	if err != nil {
		return
	}
	id, err := ExtractClaims(token)
	if err != nil {
		return
	}
	user, err = mongo.UserFromId(id)
	return
}