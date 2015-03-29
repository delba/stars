package github

/*
TODO Links
link := res.Header.Get("Link")
<https://api.github.com/user/following?page=2>; rel="next", <https://api.github.com/user/following?page=2>; rel="last"
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"

	"github.com/delba/stars/models"
)

const baseURL = "https://api.github.com"

func GetFollowing() ([]models.User, error) {
	var users []models.User
	var err error

	res, err := Client.Get(baseURL + "/user/following")
	if err != nil {
		return users, err
	}
	defer res.Body.Close()

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return users, err
	}

	err = json.Unmarshal(contents, &users)
	if err != nil {
		return users, err
	}

	return users, err
}

func GetStarredForUser(u models.User) ([]models.Repository, error) {
	var repositories []models.Repository
	var err error

	res, err := Client.Get(baseURL + "/users/" + u.Login + "/starred")
	if err != nil {
		return repositories, err
	}
	defer res.Body.Close()

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return repositories, err
	}

	err = json.Unmarshal(contents, &repositories)
	if err != nil {
		return repositories, err
	}

	return repositories, err
}

func GetFollowingStarred() (models.Repositories, error) {
	var repositories models.Repositories

	following, err := GetFollowing()
	if err != nil {
		fmt.Println("Error:", err)
		return repositories, err
	}

	c := make(chan map[models.User][]models.Repository)

	for _, user := range following {
		go func(u models.User, c chan map[models.User][]models.Repository) {
			repositories, err := GetStarredForUser(u)
			if err != nil {
				fmt.Println(err)
			}
			c <- map[models.User][]models.Repository{u: repositories}
		}(user, c)
	}

	for range following {
		for user, repos := range <-c {
			for _, repo := range repos {
				rp := repositories.FindOrAddRepository(repo)
				rp.FollowingStargazers = append(rp.FollowingStargazers, user)
			}
		}
	}

	sort.Sort(models.ByPopularity(repositories))

	return repositories, err
}

func IsStarringRepository(r *models.Repository) bool {
	res, err := Client.Get(baseURL + "/user/starred/" + r.FullName)
	if err != nil {
		fmt.Println(err)
	}

	if res.StatusCode == 204 {
		return true
	} else {
		return false
	}
}

func StarRepository(fullName string) {
	URL, _ := url.Parse(baseURL + "/user/starred/" + fullName)

	request := &http.Request{
		Method:        "PUT",
		URL:           URL,
		ContentLength: 0,
	}

	res, err := Client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.StatusCode)
}

func UnstarRepository(fullName string) {
	URL, _ := url.Parse(baseURL + "/user/starred/" + fullName)

	request := &http.Request{
		Method: "DELETE",
		URL:    URL,
	}

	res, err := Client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.StatusCode)
}

func IsFollowingUser(u models.User) bool {
	res, err := Client.Get(baseURL + "/user/following/" + u.Login)
	if err != nil {
		fmt.Println(err)
	}

	if res.StatusCode == 204 {
		return true
	} else {
		return false
	}
}

func FollowUser(username string) {
	URL, _ := url.Parse(baseURL + "/user/following/" + username)

	request := &http.Request{
		Method: "PUT",
		URL:    URL,
	}

	res, err := Client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.StatusCode)
}

func UnfollowUser(username string) {
	URL, _ := url.Parse(baseURL + "/user/following/" + username)

	request := &http.Request{
		Method: "DELETE",
		URL:    URL,
	}

	res, err := Client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.StatusCode)
}
