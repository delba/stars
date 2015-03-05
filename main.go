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
	t, err := template.ParseFiles(viewPath("public.html"))
	handle(err)

	err = t.Execute(w, nil)
	handle(err)
}

func PrivateIndex(w http.ResponseWriter, r *http.Request) {
	var err error

	starredRepositories, err := github.GetFollowingStarred()
	sort.Sort(model.ByPopularity(starredRepositories))

	t, err := template.ParseFiles(viewPath("private.html"))
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

	err = github.SetClient(r.URL.Query()["code"][0])
	handle(err)

	err = github.SetCurrentUser()
	handle(err)

	http.Redirect(w, r, "/", 302)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	github.Client, github.CurrentUser = nil, nil
	http.Redirect(w, r, "/", 302)
}

func viewPath(file string) string {
	var _, __FILE__, _, _ = runtime.Caller(1)

	return path.Join(path.Dir(__FILE__), "views", file)
}
