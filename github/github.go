package github

import "github.com/octokit/go-octokit/octokit"

var client = octokit.NewClient(nil)

func GetFollowingStarred(username string) ([]octokit.Repository, error) {
	var repositories []octokit.Repository

	user, err := GetUser(username)
	if err != nil {
		return repositories, err
	}

	following, err := GetFollowing(user)
	if err != nil {
		return repositories, err
	}

	for _, user := range following {
		starred, err := GetStarred(user)
		if err != nil {
			return repositories, err
		}

		repositories = append(repositories, starred...)
	}

	return repositories, err
}

func GetUser(username string) (*octokit.User, error) {
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

func GetFollowing(user *octokit.User) ([]octokit.User, error) {
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

func GetStarred(user octokit.User) ([]octokit.Repository, error) {
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
