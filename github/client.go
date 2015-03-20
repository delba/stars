package github

import (
	"net/http"

	"golang.org/x/oauth2"
)

const (
	clientID     = "aa78f0f464d4a73010a0"
	clientSecret = "842311922d9dd09ee074f63cf0218f9db2c75056"

	authURL     = "https://github.com/login/oauth/authorize"
	tokenURL    = "https://github.com/login/oauth/access_token"
	redirectURL = "http://localhost:8080/callback"
)

var config = oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	Endpoint: oauth2.Endpoint{
		AuthURL:  authURL,
		TokenURL: tokenURL,
	},
	RedirectURL: redirectURL,
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
