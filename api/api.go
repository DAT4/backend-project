package api

import (
	"errors"
	"github.com/DAT4/backend-project/dto"
	"github.com/DAT4/backend-project/middle"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	Game *middle.Game
}

func update(r *http.Request, fun func(id string, f dto.Update) (out dto.Object, err error)) (dto.Object, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, errors.New("no id provided")
	}
	updates, err := dto.UpdateFromJson(r.Body)
	if err != nil {
		return nil, err
	}
	return fun(id, updates)
}
func find(r *http.Request, f dto.Filter, fun func(f dto.Filter) (out []dto.Object, err error)) ([]dto.Object, error) {
	err := dto.FilterFromForm(r, f)
	if err != nil {
		return nil, err
	}
	return fun(f)
}
func findOne(r *http.Request, fun func(id string) (out dto.Object, err error)) (dto.Object, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return nil, errors.New("no id provided")
	}
	return fun(id)
}
func deleteOne(r *http.Request, fun func(id string) error) error {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return errors.New("no id provided")
	}
	return fun(id)
}
