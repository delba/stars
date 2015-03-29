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
