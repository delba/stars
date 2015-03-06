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
	"sort"
)

func GetFollowing() ([]User, error) {
	var users []User
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

func GetStarredForUser(u User) ([]Repository, error) {
	var repositories []Repository
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

func GetFollowingStarred() (Repositories, error) {
	var repositories Repositories

	following, err := GetFollowing()
	if err != nil {
		return repositories, err
	}

	c := make(chan map[User][]Repository)

	for _, user := range following {
		go func(u User, c chan map[User][]Repository) {
			repositories, err := GetStarredForUser(u)
			if err != nil {
				fmt.Println(err)
			}
			c <- map[User][]Repository{u: repositories}
		}(user, c)
	}

	for range following {
		for user, repos := range <-c {
			for _, repo := range repos {
				fmt.Println(repo.FullName)
				fmt.Println(repo.Owner.AvatarURL)
				fmt.Println(repo.Organization.AvatarURL)
				rp := repositories.FindOrAddRepository(repo)
				rp.FollowingStargazers = append(rp.FollowingStargazers, user)
			}
		}
	}

	sort.Sort(ByPopularity(repositories))

	return repositories, err
}

func isStarringRepository(r Repository) bool {
	// GET /user/starred/:owner/:repo
	// 204: true, 404: false
	return false
}

func StarRepository(r Repository) {
	// PUT /user/starred/:owner/:repo
}

func UnstarRepository(r Repository) {
	// DELETE /user/starred/:owner/:repo
}

func isFollowingUser(u User) bool {
	// GET /user/following/:username
	// 204: true, 404: false
	return false
}

func FollowUser(u User) {
	// PUT /user/following/:username
}

func UnfollowUser(u User) {
	// DELETE /user/following/:username
}
