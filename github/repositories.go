package github

type Repositories []*Repository

func (r *Repositories) FindOrAddRepository(repo Repository) *Repository {
	for _, rp := range *r {
		if rp.ID == repo.ID {
			return rp
		}
	}

	*r = append(*r, &repo)

	return &repo
}

type ByPopularity Repositories

func (r ByPopularity) Len() int {
	return len(r)
}

func (r ByPopularity) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ByPopularity) Less(i, j int) bool {
	iPopularity := len(r[i].FollowingStargazers)
	jPopularity := len(r[j].FollowingStargazers)

	if iPopularity == jPopularity {
		return r[i].FullName < r[j].FullName
	} else {
		return iPopularity > jPopularity
	}
}
