package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/felipemacedo1/b/internal/models"
)

// Client represents a GitHub API client
type Client struct {
	Token  string
	User   string
	Client *http.Client
}

// NewClient creates a new GitHub API client
func NewClient(token, user string) *Client {
	return &Client{
		Token:  token,
		User:   user,
		Client: &http.Client{Timeout: 30 * time.Second},
	}
}

// GitHubRepo represents a repository from GitHub API
type GitHubRepo struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	HTMLURL     string    `json:"html_url"`
	Language    string    `json:"language"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	Private     bool      `json:"private"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Topics      []string  `json:"topics"`
}

// GitHubUser represents a user from GitHub API
type GitHubUser struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	Location  string    `json:"location"`
	Company   string    `json:"company"`
	Blog      string    `json:"blog"`
	Email     string    `json:"email"`
	PublicRepos int     `json:"public_repos"`
	Followers int       `json:"followers"`
	Following int       `json:"following"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// makeRequest makes an authenticated request to GitHub API
func (c *Client) makeRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	return c.Client.Do(req)
}

// GetProfile fetches the user's GitHub profile
func (c *Client) GetProfile() (*models.Profile, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", c.User)
	resp, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, err
	}

	return &models.Profile{
		GitHubID:      githubUser.ID,
		Login:         githubUser.Login,
		Name:          githubUser.Name,
		Bio:           githubUser.Bio,
		AvatarURL:     githubUser.AvatarURL,
		Location:      githubUser.Location,
		Company:       githubUser.Company,
		Blog:          githubUser.Blog,
		Email:         githubUser.Email,
		PublicRepos:   githubUser.PublicRepos,
		Followers:     githubUser.Followers,
		Following:     githubUser.Following,
		CreatedAt:     githubUser.CreatedAt,
		UpdatedAt:     githubUser.UpdatedAt,
		LastSyncedAt:  time.Now(),
	}, nil
}

// GetRepositories fetches the user's public repositories
func (c *Client) GetRepositories() ([]models.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?type=public&sort=updated&per_page=100", c.User)
	resp, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var githubRepos []GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&githubRepos); err != nil {
		return nil, err
	}

	var repositories []models.Repository
	for _, repo := range githubRepos {
		repository := models.Repository{
			GitHubID:    repo.ID,
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			HTMLURL:     repo.HTMLURL,
			Language:    repo.Language,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			Private:     repo.Private,
			CreatedAt:   repo.CreatedAt,
			UpdatedAt:   repo.UpdatedAt,
			Topics:      repo.Topics,
		}
		repositories = append(repositories, repository)
	}

	return repositories, nil
}