package middle

import (
	"encoding/json"
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/models"
	"io"
)

func UserFromJson(data io.ReadCloser) (user models.User, err error) {
	err = json.NewDecoder(data).Decode(&user)
	return
}

func UserFromToken(token string, db dao.DBase) (user models.User, err error) {
	id, err := extractClaims(token)
	if err != nil {
		return
	}
	user, err = db.UserFromId(id)
	return
}
