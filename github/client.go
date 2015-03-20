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

func SetClient(code string) error {
	var err error

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return err
	}

	Client = config.Client(oauth2.NoContext, token)

	return err
}

func randomString() string {
	return "hello"
}
