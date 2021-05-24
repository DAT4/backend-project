package base

import (
	"github.com/DAT4/backend-project/dto"
)

type DB interface {
	Insert(i dto.Object) (err error)
	Update(id string, u dto.Update) (o dto.Object, err error)
	Delete(id string) (err error)
	FindOne(id string) (o dto.Object, err error)
	Find(f dto.Filter) (o []dto.Object, err error)
}
