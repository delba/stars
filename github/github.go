package github

import (
	"net/http"

	"github.com/delba/stars/model"
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

var Client *octokit.Client
var CurrentUser *octokit.User

func GetAuthURL() string {
	return config.AuthCodeURL(randomString())
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

	user, result := Client.Users(url).One()
	if result.HasError() {
		return user, err
	}

	return user, err
}

func GetFollowing() ([]octokit.User, error) {
	var following []octokit.User

	followingURL, err := CurrentUser.FollowingURL.Expand(nil)
	if err != nil {
		return following, err
	}

	users, result := Client.Users(followingURL).All()
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
		users, result = Client.Users(nextPageURL).All()
		if result.HasError() {
			break
		}
		following = append(following, users...)
	}

	return following, err
}

func GetFollowingStarred() (model.StarredRepositories, error) {
	var starredRepositories model.StarredRepositories
	var err error

	c := make(chan map[octokit.User][]octokit.Repository)

	following, err := GetFollowing()
	if err != nil {
		return starredRepositories, err
	}

	for _, user := range following {
		go GetStarredRepositories(user, c)
	}

	for range following {
		for user, repos := range <-c {
			for _, repo := range repos {
				starredRepository := starredRepositories.FindOrCreateByRepository(repo)
				starredRepository.Users = append(starredRepository.Users, &user)
			}
		}
	}

	return starredRepositories, err
}

func GetStarredRepositories(u octokit.User, c chan map[octokit.User][]octokit.Repository) {
	url, err := u.StarredURL.Expand(nil)
	if err != nil {
		panic(err)
	}

	repositories, result := Client.Repositories(url).All()
	if result.HasError() {
		panic(result.Err)
	}

	c <- map[octokit.User][]octokit.Repository{u: repositories}
}

func randomString() string {
	return "hello"
}
