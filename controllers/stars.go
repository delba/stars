package controllers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/delba/stars/github"
	"github.com/delba/stars/models"
	"github.com/julienschmidt/httprouter"
)

type Stars struct{}

func (s *Stars) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var viewFile string
	var data interface{}
	var err error

	github.SetClient(r)

	if isLoggedIn(r) {
		viewFile = viewPath("private.html")
		var user models.User
		err = user.FetchFollowingStarred()
		handle(err)
		data = user.FollowingStarred
		// data, err = fetchFromCache("data.json")
	} else {
		viewFile = viewPath("public.html")
	}

	layoutFile := viewPath("layout.html")

	t, err := template.ParseFiles(layoutFile, viewFile)
	handle(err)

	err = t.Execute(w, data)
}

func (s *Stars) Star(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	owner := ps.ByName("owner")
	repo := ps.ByName("repo")

	data := map[string]string{"owner": owner, "repo": repo}
	json.NewEncoder(w).Encode(data)
}

func (s *Stars) Unstar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	owner := ps.ByName("owner")
	repo := ps.ByName("repo")

	data := map[string]string{"owner": owner, "repo": repo}
	json.NewEncoder(w).Encode(data)
}

func viewPath(file string) string {
	return path.Join("views", file)
}

func cacheToFile(data []byte, file string) error {
	var err error

	err = ioutil.WriteFile(file, data, 0777)

	return err
}

func fetchFromCache(file string) (models.Repositories, error) {
	var err error
	repositories := models.Repositories{}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return repositories, err
	}

	json.Unmarshal(contents, &repositories)

	return repositories, err
}
