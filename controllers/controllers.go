package controllers

import "net/http"

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func setCookie(name string, value string, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
	}

	http.SetCookie(w, cookie)
}

func deleteCookie(name string, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   name,
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

func isLoggedIn(r *http.Request) bool {
	_, err := r.Cookie("access_token")

	if err == nil {
		return true
	} else {
		return false
	}
}
