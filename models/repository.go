package models

type Repository struct {
	ArchiveURL       string      `json:"archive_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BlobsURL         string      `json:"blobs_url"`
	BranchesURL      string      `json:"branches_url"`
	CloneURL         string      `json:"clone_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	CommentsURL      string      `json:"comments_url"`
	CommitsURL       string      `json:"commits_url"`
	CompareURL       string      `json:"compare_url"`
	ContentsURL      string      `json:"contents_url"`
	ContributorsURL  string      `json:"contributors_url"`
	CreatedAt        string      `json:"created_at"`
	DefaultBranch    string      `json:"default_branch"`
	Description      string      `json:"description"`
	DownloadsURL     string      `json:"downloads_url"`
	EventsURL        string      `json:"events_url"`
	Fork             bool        `json:"fork"`
	Forks            int         `json:"forks"`
	ForksCount       int         `json:"forks_count"`
	ForksURL         string      `json:"forks_url"`
	FullName         string      `json:"full_name"`
	GitCommitsURL    string      `json:"git_commits_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitURL           string      `json:"git_url"`
	HasDownloads     bool        `json:"has_downloads"`
	HasIssues        bool        `json:"has_issues"`
	HasPages         bool        `json:"has_pages"`
	HasWiki          bool        `json:"has_wiki"`
	Homepage         string      `json:"homepage"`
	HooksURL         string      `json:"hooks_url"`
	HTMLURL          string      `json:"html_url"`
	ID               int         `json:"id"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	IssuesURL        string      `json:"issues_url"`
	KeysURL          string      `json:"keys_url"`
	LabelsURL        string      `json:"labels_url"`
	Language         string      `json:"language"`
	LanguagesURL     string      `json:"languages_url"`
	MergesURL        string      `json:"merges_url"`
	MilestonesURL    string      `json:"milestones_url"`
	MirrorURL        interface{} `json:"mirror_url"`
	Name             string      `json:"name"`
	NetworkCount     int         `json:"network_count"`
	NotificationsURL string      `json:"notifications_url"`
	OpenIssues       int         `json:"open_issues"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	Organization     struct {
		AvatarURL         string `json:"avatar_url"`
		EventsURL         string `json:"events_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		GravatarID        string `json:"gravatar_id"`
		HTMLURL           string `json:"html_url"`
		ID                int    `json:"id"`
		Login             string `json:"login"`
		OrganizationsURL  string `json:"organizations_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"organization"`
	Owner struct {
		AvatarURL         string `json:"avatar_url"`
		EventsURL         string `json:"events_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		GravatarID        string `json:"gravatar_id"`
		HTMLURL           string `json:"html_url"`
		ID                int    `json:"id"`
		Login             string `json:"login"`
		OrganizationsURL  string `json:"organizations_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"owner"`
	Private             bool   `json:"private"`
	PullsURL            string `json:"pulls_url"`
	PushedAt            string `json:"pushed_at"`
	ReleasesURL         string `json:"releases_url"`
	Size                int    `json:"size"`
	SSHURL              string `json:"ssh_url"`
	StargazersCount     int    `json:"stargazers_count"`
	StargazersURL       string `json:"stargazers_url"`
	StatusesURL         string `json:"statuses_url"`
	SubscribersCount    int    `json:"subscribers_count"`
	SubscribersURL      string `json:"subscribers_url"`
	SubscriptionURL     string `json:"subscription_url"`
	SvnURL              string `json:"svn_url"`
	TagsURL             string `json:"tags_url"`
	TeamsURL            string `json:"teams_url"`
	TreesURL            string `json:"trees_url"`
	UpdatedAt           string `json:"updated_at"`
	URL                 string `json:"url"`
	Watchers            int    `json:"watchers"`
	WatchersCount       int    `json:"watchers_count"`
	FollowingStargazers []User `json:"following_stargazers"`
	Starred             bool   `json:"starred"`
}

func (r *Repository) FollowingStargazersLogins() []string {
	var logins []string

	for _, user := range r.FollowingStargazers {
		logins = append(logins, user.Login)
	}

	return logins
}

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