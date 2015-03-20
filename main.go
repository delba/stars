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
	http.HandleFunc("/star/", Star)
	http.Handle("/public", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":"+port, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	var viewFile string
	var data interface{}
	var err error

	github.SetClient(r)

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

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "access_token",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}

func Star(w http.ResponseWriter, r *http.Request) {
	fullName := r.URL.Path[6:]

	switch r.Method {
	case "GET":
		github.StarRepository(fullName)
	case "DELETE":
		fmt.Println("Unstar repo")
	}
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
