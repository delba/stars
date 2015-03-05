package main

import (
	"net/http"
	"os"
	"path"
	"runtime"
	"sort"
	"text/template"

	"github.com/delba/stars/github"
	"github.com/delba/stars/model"
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

	http.HandleFunc("/", Index)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/callback", Callback)

	http.ListenAndServe(":"+port, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	if github.CurrentUser == nil {
		PublicIndex(w, r)
	} else {
		PrivateIndex(w, r)
	}
}

func PublicIndex(w http.ResponseWriter, r *http.Request) {
	_, filename, _, _ := runtime.Caller(1)

	t, err := template.ParseFiles(path.Join(path.Dir(filename), "views", "public.html"))
	handle(err)

	err = t.Execute(w, nil)
	handle(err)
}

func PrivateIndex(w http.ResponseWriter, r *http.Request) {
	var err error

	starredRepositories, err := github.GetFollowingStarred()
	sort.Sort(model.ByPopularity(starredRepositories))

	_, filename, _, _ := runtime.Caller(1)

	t, err := template.ParseFiles(path.Join(path.Dir(filename), "views", "private.html"))
	handle(err)

	err = t.Execute(w, starredRepositories)
	handle(err)
}

func Login(w http.ResponseWriter, r *http.Request) {
	url := github.GetAuthURL()

	http.Redirect(w, r, url, 302)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	var err error

	github.Client, err = github.GetClient(r)
	handle(err)

	github.CurrentUser, err = github.GetCurrentUser()
	handle(err)

	http.Redirect(w, r, "/", 302)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	github.Client, github.CurrentUser = nil, nil
	http.Redirect(w, r, "/", 302)
}
