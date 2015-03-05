package model

type StarredRepository struct {
	Repository string
	Users      []string
}

type ByPopularity []StarredRepository

func (c ByPopularity) Len() int {
	return len(c)
}

func (c ByPopularity) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByPopularity) Less(i, j int) bool {
	return len(c[i].Users) > len(c[j].Users)
}
