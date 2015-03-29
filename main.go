package main

import (
	"net/http"
	"os"

	"github.com/delba/stars/controllers"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var sessions controllers.Sessions
	var stars controllers.Stars

	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/":         stars.Index,
		"/star/":    stars.Star,
		"/login":    sessions.Login,
		"/logout":   sessions.Logout,
		"/callback": sessions.Callback,
	}

	for path, handler := range routes {
		http.HandleFunc(path, handler)
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.Handle("/favico.ico", fs)

	http.ListenAndServe(":"+port, nil)
}
