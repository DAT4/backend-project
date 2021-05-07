package dao

import (
	"github.com/DAT4/backend-project/models"
)

type DBase interface {
	Create(userIn models.User) (userOut models.User, err error)
	UserFromId(id string) (user models.User, err error)
	UsernameTaken(u *models.User) (err error)
	UserFromName(name string) (user models.User, err error)
}
