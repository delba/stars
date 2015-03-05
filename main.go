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
	var viewFile string
	var data interface{}
	var err error

	if github.CurrentUser == nil {
		viewFile = viewPath("public.html")
		data = nil
	} else {
		viewFile = viewPath("private.html")
		data, err = github.GetFollowingStarred()
		handle(err)
		sort.Sort(model.ByPopularity(data.(model.StarredRepositories)))
	}

	t, err := template.ParseFiles(viewFile)
	handle(err)

	err = t.Execute(w, data)
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
