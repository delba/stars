package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/delba/stars/github"
)

const baseURL = "https://api.github.com"

type User struct {
	AvatarURL           string      `json:"avatar_url"`
	Bio                 interface{} `json:"bio"`
	Blog                string      `json:"blog"`
	Company             string      `json:"company"`
	CreatedAt           string      `json:"created_at"`
	Email               string      `json:"email"`
	EventsURL           string      `json:"events_url"`
	Followers           int         `json:"followers"`
	FollowersURL        string      `json:"followers_url"`
	Following           int         `json:"following"`
	FollowingURL        string      `json:"following_url"`
	GistsURL            string      `json:"gists_url"`
	GravatarID          string      `json:"gravatar_id"`
	Hireable            bool        `json:"hireable"`
	HTMLURL             string      `json:"html_url"`
	ID                  int         `json:"id"`
	Location            string      `json:"location"`
	Login               string      `json:"login"`
	Name                string      `json:"name"`
	OrganizationsURL    string      `json:"organizations_url"`
	PublicGists         int         `json:"public_gists"`
	PublicRepos         int         `json:"public_repos"`
	ReceivedEventsURL   string      `json:"received_events_url"`
	ReposURL            string      `json:"repos_url"`
	SiteAdmin           bool        `json:"site_admin"`
	StarredURL          string      `json:"starred_url"`
	SubscriptionsURL    string      `json:"subscriptions_url"`
	Type                string      `json:"type"`
	UpdatedAt           string      `json:"updated_at"`
	URL                 string      `json:"url"`
	FollowingUsers      []User
	StarredRepositories []Repository
	FollowingStarred    []Repository
}

func (u *User) FetchFollowing() error {
	res, err := github.Client.Get(baseURL + "/user/following")
	if err != nil {
		return err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var users []User
	err = json.Unmarshal(contents, &users)
	if err != nil {
		return err
	}

	u.FollowingUsers = users

	return err
}

func (u *User) FetchStarred() error {
	res, err := github.Client.Get(baseURL + "/users/" + u.Login + "/starred")
	if err != nil {
		return err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var repositories []Repository
	err = json.Unmarshal(contents, &repositories)
	if err != nil {
		return err
	}

	u.StarredRepositories = repositories

	return err
}

func (u *User) FetchFollowingStarred() error {
	err := u.FetchFollowing()
	if err != nil {
		return err
	}

	c := make(chan map[*User][]Repository)

	for _, user := range u.FollowingUsers {
		go func(u User, c chan map[*User][]Repository) {
			err = u.FetchStarred()
			if err != nil {
				fmt.Println(err)
			}
			c <- map[*User][]Repository{&u: u.StarredRepositories}
		}(user, c)
	}

	var repositories Repositories

	for range u.FollowingUsers {
		for user, repos := range <-c {
			for _, repo := range repos {
				rp := repositories.FindOrAddRepository(repo)
				rp.FollowingStargazers = append(rp.FollowingStargazers, *user)
			}
		}
	}

	sort.Sort(ByPopularity(repositories))

	var repopo []Repository
	for _, repo := range repositories {
		repopo = append(repopo, *repo)
	}

	u.FollowingStarred = repopo

	return err
}
