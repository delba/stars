package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/delba/stars/github"
	"github.com/julienschmidt/httprouter"
)

type Sessions struct{}

func (s *Sessions) New(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	url := github.AuthURL()

	http.Redirect(w, r, url, 302)
}

func (s *Sessions) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	code := r.URL.Query()["code"][0]
	accessToken, err := github.GetAccessToken(code)
	handle(err)

	setCookie("access_token", accessToken, w)

	http.Redirect(w, r, "/", 302)
}

func (s *Sessions) Destroy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	deleteCookie("access_token", w)

	data := map[string]string{"location": "/"}
	json.NewEncoder(w).Encode(data)
}
