package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"text/template"

	"github.com/delba/stars/github"
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
	http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path[1:])
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.ListenAndServe(":"+port, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	var viewFile string
	var data interface{}
	var err error

	if github.Client == nil {
		viewFile = viewPath("public.html")
	} else {
		viewFile = viewPath("private.html")
		// data, err = github.GetFollowingStarred()
		data, err = fetchFromCache("data.json")
		handle(err)
	}

	layoutFile := viewPath("layout.html")

	t, err := template.ParseFiles(layoutFile, viewFile)
	handle(err)

	err = t.Execute(w, data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	url := github.AuthURL()

	http.Redirect(w, r, url, 302)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	err := github.SetClient(r.URL.Query()["code"][0])
	handle(err)

	http.Redirect(w, r, "/", 302)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	github.Client = nil
	http.Redirect(w, r, "/", 302)
}

func viewPath(file string) string {
	var _, __FILE__, _, _ = runtime.Caller(1)

	return path.Join(path.Dir(__FILE__), "views", file)
}

func cacheToFile(data []byte, file string) error {
	var err error

	err = ioutil.WriteFile(file, data, 0777)

	return err
}

func fetchFromCache(file string) (github.Repositories, error) {
	var err error
	repositories := github.Repositories{}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return repositories, err
	}

	json.Unmarshal(contents, &repositories)

	return repositories, err
}
