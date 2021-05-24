package dao

import (
	"errors"
	"github.com/DAT4/backend-project/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestDB struct {
	users map[string]dto.User
}

func (t *TestDB) Insert(i dto.Object) (err error) {
	var user = i.(*dto.User)
	user.Id = primitive.NewObjectID().Hex()
	t.users[user.Id] = *user
	i = user
	return nil
}

func NewTestDB() *TestDB {
	return &TestDB{users: make(map[string]dto.User)}
}

func (t *TestDB) Update(id string, u dto.Update) (o dto.Object, err error) {
	update := u.(dto.UserUpdate)
	e, ok := t.users[id]
	if ok {
		if update.Username != "" {
			e.Username = dto.Username(update.Username)
		}
		if update.Role != "" {
			e.Role = update.Role
		}
		if update.Password != "" {
			e.Password = dto.Password(update.Password)
		}
		if update.Email != "" {
			e.Email = dto.Email(update.Email)
		}
		return e, nil
	}
	return nil, errors.New("no such user with them id, this")

}
func (t *TestDB) Delete(id string) (err error) {
	_, ok := t.users[id]
	if ok {
		delete(t.users, id)
		return nil
	}
	return errors.New("no such user with them id, this")
}
func (t *TestDB) FindOne(id string) (o dto.Object, err error) {
	e, ok := t.users[id]
	if ok {
		return e, nil
	}
	return nil, errors.New("no such user with them id, this")
}
func (t *TestDB) Find(f dto.Filter) (o []dto.Object, err error) {
	filter := f.(dto.UserFilter)
	for _, e := range t.users {
		if string(e.Username) != filter.Username && filter.Username != "" {
			continue
		}
		if filter.Role != e.Role && filter.Role != "" {
			continue
		}
		if filter.Email != string(e.Email) && filter.Email != "" {
			continue
		}
		o = append(o, e)
	}
	return o, nil
}
