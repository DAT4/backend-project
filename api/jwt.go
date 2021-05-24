package api

import (
	"encoding/json"
	"github.com/DAT4/backend-project/dto"
	"github.com/DAT4/backend-project/middle"
	"net/http"
)

// TokenHandler godoc
// @Summary login
// @Description login with user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.User true "Login with user"
// @Success 200 {object} middle.TokenPair
// @Failure 500 {object} Response
// @Router /auth [post]
func (a *API) TokenHandler(w http.ResponseWriter, r *http.Request) {
	if user, err := dto.UserFromJson(r.Body); err != nil {
		handleError("auth_user", err, w, http.StatusInternalServerError)
	} else {
		filter := dto.UserFilter{Username: string(user.Username)}
		if users, err := find(r, filter, a.Game.Users.Find); err != nil {
			handleError("auth_user", err, w, http.StatusInternalServerError)
		} else {
			dbUser := users[0].(dto.User)
			err = user.Check(dbUser.Password)
			if tokens, err := middle.MakeToken(dbUser); err != nil {
				handleError("auth_make_token", err, w, http.StatusInternalServerError)
			} else {
				if err = json.NewEncoder(w).Encode(tokens); err != nil {
					handleError("auth_token_to_json", err, w, http.StatusInternalServerError)
				}
			}
		}
	}
}

// RefreshToken godoc
// @Summary get a new token
// @Description get a new token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.User true "Login with user"
// @Success 200 {object} middle.TokenPair
// @Failure 500 {object} Response
// @Router /auth [post]
func (a *API) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if token, err := middle.ExtractJWTToken(r, middle.REFRESH); err != nil {
		handleError("refresh_token_extract", err, w, http.StatusInternalServerError)
	} else {
		if u, err := middle.UserFromToken(token, a.Game.Users); err != nil {
			handleError("user_from_refresh_token", err, w, http.StatusInternalServerError)
		} else {
			user := u.(dto.User)
			filter := dto.UserFilter{Username: string(user.Username)}
			if users, err := find(r, filter, a.Game.Users.Find); err != nil {
				handleError("auth_user", err, w, http.StatusInternalServerError)
			} else {
				dbUser := users[0].(dto.User)
				if err = user.Check(dbUser.Password); err != nil {
					handleError("refresh_token_not_authorized", err, w, http.StatusInternalServerError)
				} else {
					if tokens, err := middle.RefreshToken(token, user); err != nil {
						handleError("refresh_token_get", err, w, http.StatusInternalServerError)
					} else {
						if err = json.NewEncoder(w).Encode(tokens); err != nil {
							handleError("refresh_token_to_json", err, w, http.StatusInternalServerError)
						}
					}
				}
			}
		}
	}
}
