package main

import (
	"net/http"
	"os"

	"github.com/delba/stars/controllers"
	"github.com/julienschmidt/httprouter"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	stars    controllers.Stars
	sessions controllers.Sessions
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := httprouter.New()
	router.GET("/", stars.Index)
	router.PUT("/star/:owner/:repo", stars.Star)
	router.DELETE("/star/:owner/:repo", stars.Unstar)
	router.GET("/login", sessions.New)
	router.GET("/callback", sessions.Create)
	router.DELETE("/logout", sessions.Destroy)

	fs := http.FileServer(http.Dir("public"))
	router.Handler("GET", "/public/*filepath", http.StripPrefix("/public/", fs))
	router.Handler("GET", "/favico.ico", fs)

	http.ListenAndServe(":"+port, router)
}
