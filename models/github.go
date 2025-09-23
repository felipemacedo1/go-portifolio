package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GitHub API response structures
type GitHubProfile struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Login           string            `bson:"login" json:"login"`
	Name            string            `bson:"name" json:"name"`
	AvatarURL       string            `bson:"avatar_url" json:"avatar_url"`
	Bio             string            `bson:"bio" json:"bio"`
	Company         string            `bson:"company" json:"company"`
	Location        string            `bson:"location" json:"location"`
	Email           string            `bson:"email" json:"email"`
	Blog            string            `bson:"blog" json:"blog"`
	TwitterUsername string            `bson:"twitter_username" json:"twitter_username"`
	PublicRepos     int               `bson:"public_repos" json:"public_repos"`
	PublicGists     int               `bson:"public_gists" json:"public_gists"`
	Followers       int               `bson:"followers" json:"followers"`
	Following       int               `bson:"following" json:"following"`
	CreatedAt       time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `bson:"updated_at" json:"updated_at"`
	LastFetched     time.Time         `bson:"last_fetched" json:"last_fetched"`
}

type GitHubRepository struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	GitHubID        int64             `bson:"github_id" json:"github_id"`
	Name            string            `bson:"name" json:"name"`
	FullName        string            `bson:"full_name" json:"full_name"`
	Description     string            `bson:"description" json:"description"`
	Private         bool              `bson:"private" json:"private"`
	Fork            bool              `bson:"fork" json:"fork"`
	HTMLURL         string            `bson:"html_url" json:"html_url"`
	CloneURL        string            `bson:"clone_url" json:"clone_url"`
	Homepage        string            `bson:"homepage" json:"homepage"`
	Language        string            `bson:"language" json:"language"`
	Languages       map[string]int    `bson:"languages" json:"languages"`
	Size            int               `bson:"size" json:"size"`
	StargazersCount int               `bson:"stargazers_count" json:"stargazers_count"`
	WatchersCount   int               `bson:"watchers_count" json:"watchers_count"`
	ForksCount      int               `bson:"forks_count" json:"forks_count"`
	OpenIssuesCount int               `bson:"open_issues_count" json:"open_issues_count"`
	DefaultBranch   string            `bson:"default_branch" json:"default_branch"`
	Topics          []string          `bson:"topics" json:"topics"`
	HasWiki         bool              `bson:"has_wiki" json:"has_wiki"`
	HasPages        bool              `bson:"has_pages" json:"has_pages"`
	HasDownloads    bool              `bson:"has_downloads" json:"has_downloads"`
	Archived        bool              `bson:"archived" json:"archived"`
	Disabled        bool              `bson:"disabled" json:"disabled"`
	PushedAt        time.Time         `bson:"pushed_at" json:"pushed_at"`
	CreatedAt       time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `bson:"updated_at" json:"updated_at"`
	LastFetched     time.Time         `bson:"last_fetched" json:"last_fetched"`
	Owner           string            `bson:"owner" json:"owner"`
}

type GitHubContributions struct {
	ID                    primitive.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
	Username              string               `bson:"username" json:"username"`
	TotalContributions    int                  `bson:"total_contributions" json:"total_contributions"`
	ContributionCalendar  []ContributionWeek   `bson:"contribution_calendar" json:"contribution_calendar"`
	ContributionYears     []int                `bson:"contribution_years" json:"contribution_years"`
	LongestStreak         int                  `bson:"longest_streak" json:"longest_streak"`
	CurrentStreak         int                  `bson:"current_streak" json:"current_streak"`
	LastFetched           time.Time            `bson:"last_fetched" json:"last_fetched"`
}

type ContributionWeek struct {
	WeekStart string             `bson:"week_start" json:"week_start"`
	Days      []ContributionDay  `bson:"days" json:"days"`
}

type ContributionDay struct {
	Date  string `bson:"date" json:"date"`
	Count int    `bson:"count" json:"count"`
	Level int    `bson:"level" json:"level"` // 0-4 intensity level
}

type GitHubStats struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username             string            `bson:"username" json:"username"`
	TotalRepos           int               `bson:"total_repos" json:"total_repos"`
	TotalStars           int               `bson:"total_stars" json:"total_stars"`
	TotalForks           int               `bson:"total_forks" json:"total_forks"`
	TotalCommits         int               `bson:"total_commits" json:"total_commits"`
	TotalContributions   int               `bson:"total_contributions" json:"total_contributions"`
	MostUsedLanguages    []LanguageStat    `bson:"most_used_languages" json:"most_used_languages"`
	TopRepositories      []RepoStat        `bson:"top_repositories" json:"top_repositories"`
	RecentActivity       []ActivityStat    `bson:"recent_activity" json:"recent_activity"`
	ContributionStreak   int               `bson:"contribution_streak" json:"contribution_streak"`
	IssuesOpened         int               `bson:"issues_opened" json:"issues_opened"`
	PullRequestsOpened   int               `bson:"pull_requests_opened" json:"pull_requests_opened"`
	PullRequestsMerged   int               `bson:"pull_requests_merged" json:"pull_requests_merged"`
	LastFetched          time.Time         `bson:"last_fetched" json:"last_fetched"`
}

type LanguageStat struct {
	Name       string  `bson:"name" json:"name"`
	Bytes      int     `bson:"bytes" json:"bytes"`
	Percentage float64 `bson:"percentage" json:"percentage"`
}

type RepoStat struct {
	Name        string `bson:"name" json:"name"`
	FullName    string `bson:"full_name" json:"full_name"`
	Stars       int    `bson:"stars" json:"stars"`
	Forks       int    `bson:"forks" json:"forks"`
	Language    string `bson:"language" json:"language"`
	Description string `bson:"description" json:"description"`
	HTMLURL     string `bson:"html_url" json:"html_url"`
}

type ActivityStat struct {
	Type        string    `bson:"type" json:"type"` // "push", "create", "delete", "fork", "watch", etc.
	Repo        string    `bson:"repo" json:"repo"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	Description string    `bson:"description" json:"description"`
}

// GitHub API raw responses (for mapping)
type GitHubAPIProfile struct {
	Login             string    `json:"login"`
	ID                int64     `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type GitHubAPIRepository struct {
	ID               int64                  `json:"id"`
	NodeID           string                 `json:"node_id"`
	Name             string                 `json:"name"`
	FullName         string                 `json:"full_name"`
	Private          bool                   `json:"private"`
	Owner            GitHubAPIUser          `json:"owner"`
	HTMLURL          string                 `json:"html_url"`
	Description      string                 `json:"description"`
	Fork             bool                   `json:"fork"`
	URL              string                 `json:"url"`
	ArchiveURL       string                 `json:"archive_url"`
	AssigneesURL     string                 `json:"assignees_url"`
	BlobsURL         string                 `json:"blobs_url"`
	BranchesURL      string                 `json:"branches_url"`
	CollaboratorsURL string                 `json:"collaborators_url"`
	CommentsURL      string                 `json:"comments_url"`
	CommitsURL       string                 `json:"commits_url"`
	CompareURL       string                 `json:"compare_url"`
	ContentsURL      string                 `json:"contents_url"`
	ContributorsURL  string                 `json:"contributors_url"`
	DeploymentsURL   string                 `json:"deployments_url"`
	DownloadsURL     string                 `json:"downloads_url"`
	EventsURL        string                 `json:"events_url"`
	ForksURL         string                 `json:"forks_url"`
	GitCommitsURL    string                 `json:"git_commits_url"`
	GitRefsURL       string                 `json:"git_refs_url"`
	GitTagsURL       string                 `json:"git_tags_url"`
	GitURL           string                 `json:"git_url"`
	IssueCommentURL  string                 `json:"issue_comment_url"`
	IssueEventsURL   string                 `json:"issue_events_url"`
	IssuesURL        string                 `json:"issues_url"`
	KeysURL          string                 `json:"keys_url"`
	LabelsURL        string                 `json:"labels_url"`
	LanguagesURL     string                 `json:"languages_url"`
	MergesURL        string                 `json:"merges_url"`
	MilestonesURL    string                 `json:"milestones_url"`
	NotificationsURL string                 `json:"notifications_url"`
	PullsURL         string                 `json:"pulls_url"`
	ReleasesURL      string                 `json:"releases_url"`
	SSHURL           string                 `json:"ssh_url"`
	StarredURL       string                 `json:"starred_url"`
	StatusesURL      string                 `json:"statuses_url"`
	SubscribersURL   string                 `json:"subscribers_url"`
	SubscriptionURL  string                 `json:"subscription_url"`
	TagsURL          string                 `json:"tags_url"`
	TeamsURL         string                 `json:"teams_url"`
	TreesURL         string                 `json:"trees_url"`
	CloneURL         string                 `json:"clone_url"`
	MirrorURL        string                 `json:"mirror_url"`
	HooksURL         string                 `json:"hooks_url"`
	SvnURL           string                 `json:"svn_url"`
	Homepage         string                 `json:"homepage"`
	Language         string                 `json:"language"`
	ForksCount       int                    `json:"forks_count"`
	StargazersCount  int                    `json:"stargazers_count"`
	WatchersCount    int                    `json:"watchers_count"`
	Size             int                    `json:"size"`
	DefaultBranch    string                 `json:"default_branch"`
	OpenIssuesCount  int                    `json:"open_issues_count"`
	IsTemplate       bool                   `json:"is_template"`
	Topics           []string               `json:"topics"`
	HasIssues        bool                   `json:"has_issues"`
	HasProjects      bool                   `json:"has_projects"`
	HasWiki          bool                   `json:"has_wiki"`
	HasPages         bool                   `json:"has_pages"`
	HasDownloads     bool                   `json:"has_downloads"`
	Archived         bool                   `json:"archived"`
	Disabled         bool                   `json:"disabled"`
	Visibility       string                 `json:"visibility"`
	PushedAt         time.Time              `json:"pushed_at"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Permissions      map[string]interface{} `json:"permissions"`
}

type GitHubAPIUser struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}