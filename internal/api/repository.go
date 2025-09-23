package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/felipemacedo1/b/internal/config"
	"github.com/felipemacedo1/b/internal/database"
	"github.com/felipemacedo1/b/internal/github"
	"github.com/felipemacedo1/b/internal/models"
)

// RepositoryHandler handles repository-related requests
type RepositoryHandler struct {
	db     *database.DB
	config *config.Config
}

// NewRepositoryHandler creates a new repository handler
func NewRepositoryHandler(db *database.DB, cfg *config.Config) *RepositoryHandler {
	return &RepositoryHandler{
		db:     db,
		config: cfg,
	}
}

// GetRepositories retrieves all repositories
func (h *RepositoryHandler) GetRepositories(c *gin.Context) {
	collection := h.db.GetCollection("repositories")
	
	// Sort by stars descending, then by updated_at descending
	opts := options.Find().SetSort(bson.D{
		{Key: "stars", Value: -1},
		{Key: "updated_at", Value: -1},
	})
	
	cursor, err := collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve repositories"})
		return
	}
	defer cursor.Close(context.Background())

	var repositories []models.Repository
	if err = cursor.All(context.Background(), &repositories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode repositories"})
		return
	}

	if repositories == nil {
		repositories = []models.Repository{}
	}

	c.JSON(http.StatusOK, repositories)
}

// SyncRepositories synchronizes repository data from GitHub
func (h *RepositoryHandler) SyncRepositories(c *gin.Context) {
	githubClient := github.NewClient(h.config.GitHubToken, h.config.GitHubUser)
	
	repositories, err := githubClient.GetRepositories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repositories from GitHub"})
		return
	}

	collection := h.db.GetCollection("repositories")
	
	// Clear existing repositories
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing repositories"})
		return
	}

	// Insert new repositories
	if len(repositories) > 0 {
		var docs []interface{}
		for _, repo := range repositories {
			docs = append(docs, repo)
		}
		
		_, err = collection.InsertMany(ctx, docs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save repositories"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Repositories synchronized successfully",
		"count":   len(repositories),
	})
}