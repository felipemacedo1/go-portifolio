package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"portfolio-backend/config"
	"portfolio-backend/database"
	"portfolio-backend/models"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GitHubService struct {
	client       *http.Client
	cacheService *CacheService
	collection   *mongo.Collection
}

func NewGitHubService() *GitHubService {
	return &GitHubService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		cacheService: NewCacheService(),
		collection:   database.Database.Collection("github_data"),
	}
}

// GetProfile retrieves GitHub profile information
func (gs *GitHubService) GetProfile(ctx context.Context, username string) (*models.GitHubProfile, error) {
	// Try cache first
	var profile models.GitHubProfile
	if err := gs.cacheService.GetGitHubData(ctx, username, "profile", &profile); err == nil {
		return &profile, nil
	}

	// Fetch from GitHub API
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add GitHub token if available
	if config.AppConfig.GitHubToken != "" {
		req.Header.Set("Authorization", "token "+config.AppConfig.GitHubToken)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := gs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var apiProfile models.GitHubAPIProfile
	if err := json.NewDecoder(resp.Body).Decode(&apiProfile); err != nil {
		return nil, err
	}

	// Convert to internal model
	profile = models.GitHubProfile{
		Login:           apiProfile.Login,
		Name:            apiProfile.Name,
		AvatarURL:       apiProfile.AvatarURL,
		Bio:             apiProfile.Bio,
		Company:         apiProfile.Company,
		Location:        apiProfile.Location,
		Email:           apiProfile.Email,
		Blog:            apiProfile.Blog,
		TwitterUsername: apiProfile.TwitterUsername,
		PublicRepos:     apiProfile.PublicRepos,
		PublicGists:     apiProfile.PublicGists,
		Followers:       apiProfile.Followers,
		Following:       apiProfile.Following,
		CreatedAt:       apiProfile.CreatedAt,
		UpdatedAt:       apiProfile.UpdatedAt,
		LastFetched:     time.Now(),
	}

	// Cache the result
	gs.cacheService.SetGitHubData(ctx, username, "profile", profile)

	// Store in database for persistence
	gs.storeProfile(ctx, profile)

	return &profile, nil
}

// GetRepositories retrieves user's public repositories
func (gs *GitHubService) GetRepositories(ctx context.Context, username string) ([]models.GitHubRepository, error) {
	// Try cache first
	var repos []models.GitHubRepository
	if err := gs.cacheService.GetGitHubData(ctx, username, "repositories", &repos); err == nil {
		return repos, nil
	}

	// Fetch from GitHub API with pagination
	allRepos := []models.GitHubRepository{}
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=%d&sort=updated", username, page, perPage)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		if config.AppConfig.GitHubToken != "" {
			req.Header.Set("Authorization", "token "+config.AppConfig.GitHubToken)
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := gs.client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
		}

		var apiRepos []models.GitHubAPIRepository
		if err := json.NewDecoder(resp.Body).Decode(&apiRepos); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		if len(apiRepos) == 0 {
			break
		}

		// Convert to internal models
		for _, apiRepo := range apiRepos {
			repo := models.GitHubRepository{
				GitHubID:        apiRepo.ID,
				Name:            apiRepo.Name,
				FullName:        apiRepo.FullName,
				Description:     apiRepo.Description,
				Private:         apiRepo.Private,
				Fork:            apiRepo.Fork,
				HTMLURL:         apiRepo.HTMLURL,
				CloneURL:        apiRepo.CloneURL,
				Homepage:        apiRepo.Homepage,
				Language:        apiRepo.Language,
				Size:            apiRepo.Size,
				StargazersCount: apiRepo.StargazersCount,
				WatchersCount:   apiRepo.WatchersCount,
				ForksCount:      apiRepo.ForksCount,
				OpenIssuesCount: apiRepo.OpenIssuesCount,
				DefaultBranch:   apiRepo.DefaultBranch,
				Topics:          apiRepo.Topics,
				HasWiki:         apiRepo.HasWiki,
				HasPages:        apiRepo.HasPages,
				HasDownloads:    apiRepo.HasDownloads,
				Archived:        apiRepo.Archived,
				Disabled:        apiRepo.Disabled,
				PushedAt:        apiRepo.PushedAt,
				CreatedAt:       apiRepo.CreatedAt,
				UpdatedAt:       apiRepo.UpdatedAt,
				LastFetched:     time.Now(),
				Owner:           username,
			}

			// Fetch languages for each repo (with rate limiting consideration)
			if repo.Language != "" {
				languages, _ := gs.getRepositoryLanguages(ctx, username, repo.Name)
				repo.Languages = languages
			}

			allRepos = append(allRepos, repo)
		}

		// Check if there are more pages
		if len(apiRepos) < perPage {
			break
		}
		page++
	}

	// Cache the result
	gs.cacheService.SetGitHubData(ctx, username, "repositories", allRepos)

	// Store in database
	gs.storeRepositories(ctx, allRepos)

	return allRepos, nil
}

// GetContributions retrieves contribution data (simplified version)
func (gs *GitHubService) GetContributions(ctx context.Context, username string) (*models.GitHubContributions, error) {
	// Try cache first
	var contributions models.GitHubContributions
	if err := gs.cacheService.GetGitHubData(ctx, username, "contributions", &contributions); err == nil {
		return &contributions, nil
	}

	// GitHub doesn't provide a direct API for contribution graph
	// We'll simulate based on repository activity and commits
	repos, err := gs.GetRepositories(ctx, username)
	if err != nil {
		return nil, err
	}

	// Calculate approximate contributions based on repository data
	totalContributions := 0
	for _, repo := range repos {
		if !repo.Fork {
			// Estimate contributions based on repo activity
			totalContributions += repo.StargazersCount + repo.ForksCount + 10 // Base contribution
		}
	}

	contributions = models.GitHubContributions{
		Username:           username,
		TotalContributions: totalContributions,
		ContributionYears:  []int{time.Now().Year(), time.Now().Year() - 1},
		LongestStreak:      30, // Placeholder
		CurrentStreak:      5,  // Placeholder
		LastFetched:        time.Now(),
	}

	// Cache the result
	gs.cacheService.SetGitHubData(ctx, username, "contributions", contributions)

	return &contributions, nil
}

// GetStats calculates aggregated GitHub statistics
func (gs *GitHubService) GetStats(ctx context.Context, username string) (*models.GitHubStats, error) {
	// Try cache first
	var stats models.GitHubStats
	if err := gs.cacheService.GetGitHubData(ctx, username, "stats", &stats); err == nil {
		return &stats, nil
	}

	// Get repositories to calculate stats
	repos, err := gs.GetRepositories(ctx, username)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	totalStars := 0
	totalForks := 0
	languageStats := make(map[string]int)
	topRepos := []models.RepoStat{}

	for _, repo := range repos {
		if !repo.Fork && !repo.Private {
			totalStars += repo.StargazersCount
			totalForks += repo.ForksCount

			// Count languages
			if repo.Language != "" {
				languageStats[repo.Language]++
			}

			// Add to top repos if it has stars
			if repo.StargazersCount > 0 {
				topRepos = append(topRepos, models.RepoStat{
					Name:        repo.Name,
					FullName:    repo.FullName,
					Stars:       repo.StargazersCount,
					Forks:       repo.ForksCount,
					Language:    repo.Language,
					Description: repo.Description,
					HTMLURL:     repo.HTMLURL,
				})
			}
		}
	}

	// Convert language stats to sorted list
	var mostUsedLanguages []models.LanguageStat
	for lang, count := range languageStats {
		percentage := float64(count) / float64(len(repos)) * 100
		mostUsedLanguages = append(mostUsedLanguages, models.LanguageStat{
			Name:       lang,
			Bytes:      count * 1000, // Approximate
			Percentage: percentage,
		})
	}

	stats = models.GitHubStats{
		Username:         username,
		TotalRepos:       len(repos),
		TotalStars:       totalStars,
		TotalForks:       totalForks,
		MostUsedLanguages: mostUsedLanguages,
		TopRepositories:  topRepos,
		LastFetched:      time.Now(),
	}

	// Cache the result
	gs.cacheService.SetGitHubData(ctx, username, "stats", stats)

	return &stats, nil
}

// SyncData forces a refresh of all GitHub data for a user
func (gs *GitHubService) SyncData(ctx context.Context, username string) error {
	// Invalidate cache
	gs.cacheService.InvalidateGitHubCache(ctx, username)

	// Fetch fresh data
	_, err := gs.GetProfile(ctx, username)
	if err != nil {
		return err
	}

	_, err = gs.GetRepositories(ctx, username)
	if err != nil {
		return err
	}

	_, err = gs.GetContributions(ctx, username)
	if err != nil {
		return err
	}

	_, err = gs.GetStats(ctx, username)
	return err
}

// Helper methods

func (gs *GitHubService) getRepositoryLanguages(ctx context.Context, username, repoName string) (map[string]int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", username, repoName)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if config.AppConfig.GitHubToken != "" {
		req.Header.Set("Authorization", "token "+config.AppConfig.GitHubToken)
	}

	resp, err := gs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch languages: %d", resp.StatusCode)
	}

	var languages map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, err
	}

	return languages, nil
}

func (gs *GitHubService) storeProfile(ctx context.Context, profile models.GitHubProfile) error {
	profile.ID = primitive.NewObjectID()
	filter := bson.M{"login": profile.Login}
	update := bson.M{"$set": profile}
	opts := options.Update().SetUpsert(true)

	_, err := gs.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (gs *GitHubService) storeRepositories(ctx context.Context, repos []models.GitHubRepository) error {
	if len(repos) == 0 {
		return nil
	}

	var operations []mongo.WriteModel
	for _, repo := range repos {
		repo.ID = primitive.NewObjectID()
		filter := bson.M{"github_id": repo.GitHubID}
		update := bson.M{"$set": repo}
		operation := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		operations = append(operations, operation)
	}

	_, err := gs.collection.BulkWrite(ctx, operations)
	return err
}

// CheckRateLimit checks GitHub API rate limit
func (gs *GitHubService) CheckRateLimit(ctx context.Context) (map[string]interface{}, error) {
	url := "https://api.github.com/rate_limit"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if config.AppConfig.GitHubToken != "" {
		req.Header.Set("Authorization", "token "+config.AppConfig.GitHubToken)
	}

	resp, err := gs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rateLimit map[string]interface{}
	if err := json.Unmarshal(body, &rateLimit); err != nil {
		return nil, err
	}

	return rateLimit, nil
}