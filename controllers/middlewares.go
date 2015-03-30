package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Middlewares struct{}

func (m *Middlewares) Authenticate(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if isLoggedIn(r) {
			next(w, r, ps)
			return
		}

		// TODO redirect
	}
}
