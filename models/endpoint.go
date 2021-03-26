package models

import (
	"net/http"
)

type Endpoint struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
	Login   bool
	Method  string
}
