package middle

import (
	"encoding/json"
	"fmt"
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/models"
	"io"
)

func UserFromJson(data io.ReadCloser) (user models.User, err error) {
	err = json.NewDecoder(data).Decode(&user)
	return
}

func UserFromToken(token string) (user models.User, err error) {
	id, err := extractClaims(token)
	fmt.Println(id)
	if err != nil {
		return
	}
	user, err = dao.UserFromId(id)
	return
}
