package model

type StarredRepository struct {
	Repository string
	Users      []string
}

type StarredRepositories []*StarredRepository

func (sr *StarredRepositories) FindOrCreateByRepository(repo string) *StarredRepository {
	starredRepository := sr.FindByRepository(repo)

	if starredRepository == nil {
		starredRepository = &StarredRepository{Repository: repo}
		*sr = append(*sr, starredRepository)
	}

	return starredRepository
}

func (sr *StarredRepositories) FindByRepository(repo string) *StarredRepository {
	for _, starredRepository := range *sr {
		if starredRepository.Repository == repo {
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
	return len(c[i].Users) > len(c[j].Users)
}
