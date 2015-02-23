package main

import (
	"fmt"
	"os"

	"github.com/octokit/go-octokit/octokit"
)

var client *octokit.Client
var stars []octokit.Repository

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	client = octokit.NewClient(nil)
}

func main() {
	username := os.Args[1]

	user, err := getUser(username)
	handle(err)

	following, err := getFollowing(user)
	handle(err)

	for _, user := range following {
		starred, err := getStarred(user)
		handle(err)

		stars = append(stars, starred...)
	}

	fmt.Println(stars)
}

func getUser(username string) (*octokit.User, error) {
	url, err := octokit.UserURL.Expand(octokit.M{"user": username})
	if err != nil {
		return nil, err
	}

	user, result := client.Users(url).One()
	if result.HasError() {
		return nil, result.Err
	}

	return user, nil
}

func getFollowing(user *octokit.User) ([]octokit.User, error) {
	url, err := user.FollowingURL.Expand(nil)
	if err != nil {
		return nil, err
	}

	following, result := client.Users(url).All()
	if result.HasError() {
		return nil, result.Err
	}

	return following, nil
}

func getStarred(user octokit.User) ([]octokit.Repository, error) {
	url, err := user.StarredURL.Expand(nil)
	if err != nil {
		return nil, err
	}

	starred, result := client.Repositories(url).All()
	if result.HasError() {
		return nil, result.Err
	}

	return starred, nil
}
