package api

import (
	"encoding/json"
	"errors"
	"github.com/DAT4/backend-project/dto"
	"github.com/DAT4/backend-project/middle"
	"github.com/gorilla/mux"
	"net/http"
)

// InsertUser godoc
// @Summary create a new user
// @Description create a new user
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body dto.User true "Create user"
// @Success 201 {object} dto.User{id=string}
// @Failure 500 {object} Response
// @Router /user [post]
func (a *API) InsertUser(w http.ResponseWriter, r *http.Request) {
	if user, err := dto.UserFromJson(r.Body); err != nil {
		handleError("user_from_json", err, w, http.StatusInternalServerError)
	} else {
		if out, err := a.Game.Users.Find(dto.UserFilter{Username: string(user.Username)}); len(out) > 0 || err != nil {
			handleError("users_find", err, w, http.StatusInternalServerError)
		} else {
			if err = middle.Validate(user); err != nil {
				handleError("user_validate", err, w, http.StatusInternalServerError)
			} else {
				if err = user.HashAndSalt(); err != nil {
					handleError("user_hash_n_salt", err, w, http.StatusInternalServerError)
				} else {
					if err = a.Game.Users.Insert(&user); err != nil {
						handleError("user_insert", err, w, http.StatusInternalServerError)
					}
				}
			}
		}
	}
}

// UpdateUser godoc
// @Summary update a user
// @Description update a user
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Param user body dto.UserUpdate true "Update user"
// @Success 200 {object} dto.User
// @Failure 500 {object} Response
// @Router /user/{id} [put]
func (a *API) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		handleError("user_get_id_from_vars", errors.New("no id provided"), w, http.StatusNotAcceptable)
		return
	}

	var user dto.User

	if u, err := middle.UserFromHeader(r, a.Game.Users); err != nil {
		handleError("user_from_token", err, w, http.StatusNotAcceptable)
		return
	} else {
		currentUser := u.(dto.User)
		if currentUser.Role != "admin" && currentUser.Id != id {
			handleError("user_is_admin", errors.New("you don't have the right to edit another user"), w, http.StatusNotAcceptable)
			return
		}
	}

	if u, err := a.Game.Users.FindOne(id); err != nil {
		handleError("user_find_one", err, w, http.StatusInternalServerError)
		return
	} else {
		user = u.(dto.User)
	}

	if updates, err := dto.UserUpdateFromJson(r.Body); err != nil {
		handleError("user_update_from_json", err, w, http.StatusNotAcceptable)
		return
	} else {

		if updates.Username != "" {
			if string(user.Username) != updates.Username {
				user.Username = dto.Username(updates.Username)
				if out, err := a.Game.Users.Find(dto.UserFilter{Username: string(user.Username)}); len(out) > 0 || err != nil {
					handleError("user_update", errors.New("username taken"), w, http.StatusInternalServerError)
					return
				}
			}
		}

		if updates.Password != "" {
			user.Password = dto.Password(updates.Password)
			err = user.Password.Validate()
			if err != nil {
				handleError("user_update", err, w, http.StatusNotAcceptable)
				return
			}
			err = user.HashAndSalt()
			if err != nil {
				handleError("user_hash_n_salt", err, w, http.StatusInternalServerError)
				return
			}
			updates.Password = string(user.Password)
		}

		if u, err := a.Game.Users.Update(id, updates); err != nil {
			handleError("user_update", err, w, http.StatusInternalServerError)
			return
		} else {
			if err = u.ToJson(w); err != nil {
				handleError("user_to_json", err, w, http.StatusInternalServerError)
				return
			}
		}
	}

}

// DeleteUser godoc
// @Summary delete a user
// @Description delete a user
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200 {object} Response
// @Failure default {object} Response
// @Router /user/{id} [delete]
func (a *API) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := deleteOne(r, a.Game.Users.Delete); err != nil {
		handleError("user_delete", err, w, http.StatusInternalServerError)
	} else {
		if err = json.NewEncoder(w).Encode(Response{
			Data:  "successfully deleted user",
			Error: "",
		}); err != nil {
			handleError("response_to_json", err, w, http.StatusInternalServerError)
		}
	}
}

// FindOneUser godoc
// @Summary find a user
// @Description find a user
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 200 {object} Response
// @Failure default {object} Response
// @Router /user/{id} [get]
func (a *API) FindOneUser(w http.ResponseWriter, r *http.Request) {
	if user, err := findOne(r, a.Game.Users.FindOne); err != nil {
		handleError("user_find_one", err, w, http.StatusInternalServerError)
	} else {
		if err = user.(dto.User).ToJson(w); err != nil {
			handleError("user_to_json", err, w, http.StatusInternalServerError)
		}
	}
}

// FindUsers godoc
// @Summary find users
// @Description find users
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param filter query dto.UserFilter true "User filter"
// @Success 200 {object} []dto.User
// @Failure default {object} Response
// @Router /user [get]
func (a *API) FindUsers(w http.ResponseWriter, r *http.Request) {
	var f dto.UserFilter
	if users, err := find(r, f, a.Game.Users.Find); err != nil {
		handleError("users_find", err, w, http.StatusInternalServerError)
	} else {
		var out []dto.User
		for _, e := range users {
			out = append(out, e.(dto.User))
		}
		if err = json.NewEncoder(w).Encode(out); err != nil {
			handleError("users_to_json", err, w, http.StatusInternalServerError)
		}
	}
}
