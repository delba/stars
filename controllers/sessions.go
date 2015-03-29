package controllers

import (
	"net/http"

	"github.com/delba/stars/github"
)

type Sessions struct{}

func (s *Sessions) Login(w http.ResponseWriter, r *http.Request) {
	url := github.AuthURL()

	http.Redirect(w, r, url, 302)
}

func (s *Sessions) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]
	accessToken, err := github.GetAccessToken(code)
	handle(err)

	setCookie("access_token", accessToken, w)

	http.Redirect(w, r, "/", 302)
}

func (s *Sessions) Logout(w http.ResponseWriter, r *http.Request) {
	deleteCookie("access_token", w)

	http.Redirect(w, r, "/", 302)
}
