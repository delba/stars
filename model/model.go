package model

import "github.com/octokit/go-octokit/octokit"

type StarredRepository struct {
	Repository *octokit.Repository
	Users      []*octokit.User
}

func (s *StarredRepository) UserLogins() []string {
	var userLogins []string

	for _, user := range s.Users {
		userLogins = append(userLogins, user.Login)
	}

	return userLogins
}

type StarredRepositories []*StarredRepository

func (sr *StarredRepositories) FindOrCreateByRepository(repository octokit.Repository) *StarredRepository {
	starredRepository := sr.FindByRepository(repository)

	if starredRepository == nil {
		starredRepository = &StarredRepository{Repository: &repository}
		*sr = append(*sr, starredRepository)
	}

	return starredRepository
}

func (sr *StarredRepositories) FindByRepository(repository octokit.Repository) *StarredRepository {
	for _, starredRepository := range *sr {
		if starredRepository.Repository.ID == repository.ID {
			return starredRepository
		}
	}

	return nil
}

type ByPopularity StarredRepositories

func (c ByPopularity) Len() int {
	return len(c)
}

func (c ByPopularity) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByPopularity) Less(i, j int) bool {
	if len(c[i].Users) == len(c[j].Users) {
		return c[i].Repository.FullName < c[j].Repository.FullName
	} else {
		return len(c[i].Users) > len(c[j].Users)
	}
}
