// TODO Refactor!!!
package main

import (
	"net/http"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/octokit/go-octokit/octokit"
	"golang.org/x/oauth2"
)

const (
	ClientID     = "aa78f0f464d4a73010a0"
	ClientSecret = "842311922d9dd09ee074f63cf0218f9db2c75056"

	AuthURL     = "https://github.com/login/oauth/authorize"
	TokenURL    = "https://github.com/login/oauth/access_token"
	RedirectURL = "http://localhost:8080/callback"
)

var config = oauth2.Config{
	ClientID:     ClientID,
	ClientSecret: ClientSecret,
	Endpoint: oauth2.Endpoint{
		AuthURL:  AuthURL,
		TokenURL: TokenURL,
	},
	RedirectURL: RedirectURL,
}

var client *octokit.Client
var currentUser *octokit.User

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
	if currentUser == nil {
		PublicIndex(w, r)
	} else {
		PrivateIndex(w, r)
	}
}

func PublicIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(path.Join("views", "public.html"))
	handle(err)

	err = t.Execute(w, nil)
	handle(err)
}

type StarredRepository struct {
	Repository string
	Users      []string
}

type ByPopularity []StarredRepository

func (c ByPopularity) Len() int {
	return len(c)
}

func (c ByPopularity) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByPopularity) Less(i, j int) bool {
	return len(c[i].Users) > len(c[j].Users)
}

func PrivateIndex(w http.ResponseWriter, r *http.Request) {
	var err error

	followingStarred, err := GetFollowingStarred()

	var starredRepositories []StarredRepository

	for repository, users := range followingStarred {
		starredRepositories = append(starredRepositories, StarredRepository{
			Repository: repository,
			Users:      users,
		})
	}

	sort.Sort(ByPopularity(starredRepositories))

	t, err := template.ParseFiles(path.Join("views", "private.html"))
	handle(err)

	err = t.Execute(w, starredRepositories)
	handle(err)
}

func Login(w http.ResponseWriter, r *http.Request) {
	url := config.AuthCodeURL(randomString())

	http.Redirect(w, r, url, 302)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	var err error

	client, err = GetClient(r)
	handle(err)

	currentUser, err = GetCurrentUser()
	handle(err)

	http.Redirect(w, r, "/", 302)
}

func GetClient(r *http.Request) (*octokit.Client, error) {
	var client *octokit.Client
	var err error
	code := r.URL.Query()["code"][0]
	token, err := config.Exchange(nil, code)
	if err != nil {
		return client, err
	}

	client = octokit.NewClient(octokit.TokenAuth{token.AccessToken})

	return client, err
}

func GetCurrentUser() (*octokit.User, error) {
	var user *octokit.User
	var err error

	url, err := octokit.CurrentUserURL.Expand(nil)
	if err != nil {
		return user, err
	}

	user, result := client.Users(url).One()
	if result.HasError() {
		return user, err
	}

	return user, err
}

func Logout(w http.ResponseWriter, r *http.Request) {
	client, currentUser = nil, nil
	http.Redirect(w, r, "/", 302)
}

func GetFollowing() ([]octokit.User, error) {
	var following []octokit.User

	followingURL, err := currentUser.FollowingURL.Expand(nil)
	if err != nil {
		return following, err
	}

	users, result := client.Users(followingURL).All()
	if result.HasError() {
		return following, err
	}
	following = append(following, users...)

	// TODO goroutine
	for {
		if result.NextPage == nil {
			break
		}
		nextPageURL, _ := result.NextPage.Expand(nil)
		users, result = client.Users(nextPageURL).All()
		if result.HasError() {
			break
		}
		following = append(following, users...)
	}

	return following, err
}

func GetFollowingStarred() (map[string][]string, error) {
	var err error

	c := make(chan map[octokit.User][]octokit.Repository)
	result := make(map[string][]string)
	// var result map[octokit.Repository][]octokit.User

	following, err := GetFollowing()
	if err != nil {
		return result, err
	}

	for _, user := range following {
		go GetStarredRepositories(user, c)
	}

	for range following {
		for user, repos := range <-c {
			for _, repo := range repos {
				result[repo.FullName] = append(result[repo.FullName], user.Login)
			}
		}
	}

	return result, err
}

func GetStarredRepositories(u octokit.User, c chan map[octokit.User][]octokit.Repository) {
	url, err := u.StarredURL.Expand(nil)
	if err != nil {
		panic(err)
	}

	repositories, result := client.Repositories(url).All()
	if result.HasError() {
		panic(result.Err)
	}

	c <- map[octokit.User][]octokit.Repository{u: repositories}
}

func randomString() string {
	return "hello"
}
