package dao

import (
	"github.com/DAT4/backend-project/models"
)

type DBase interface {
	Create(u *models.User) (err error)
	UserFromId(id string) (user models.User, err error)
	UsernameTaken(u *models.User) (err error)
	Authenticate(u *models.User) (err error)
}
