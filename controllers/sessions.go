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

	cookie := &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}

func (s *Sessions) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "access_token",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}
