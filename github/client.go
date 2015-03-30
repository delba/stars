package github

import (
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

var config = oauth2.Config{
	ClientID:     os.Getenv("STARS_CLIENT_ID"),
	ClientSecret: os.Getenv("STARS_CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
	RedirectURL: os.Getenv("STARS_REDIRECT_URL"),
	Scopes:      []string{"public_repo"},
}

var Client *http.Client

func AuthURL() string {
	return config.AuthCodeURL(randomString())
}

func GetAccessToken(code string) (string, error) {
	var accessToken string
	var err error

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return accessToken, err
	}

	accessToken = token.AccessToken

	return accessToken, err
}

func SetClient(r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		Client = nil
		return
	}

	token := &oauth2.Token{
		AccessToken: cookie.Value,
	}

	Client = config.Client(oauth2.NoContext, token)
}

func randomString() string {
	return "hello"
}
