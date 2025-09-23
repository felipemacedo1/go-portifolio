package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/felipemacedo1/b/internal/config"
	"github.com/felipemacedo1/b/internal/database"
	"github.com/felipemacedo1/b/internal/github"
	"github.com/felipemacedo1/b/internal/models"
)

// ProfileHandler handles profile-related requests
type ProfileHandler struct {
	db     *database.DB
	config *config.Config
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(db *database.DB, cfg *config.Config) *ProfileHandler {
	return &ProfileHandler{
		db:     db,
		config: cfg,
	}
}

// GetProfile retrieves the user's profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}
	
	collection := h.db.GetCollection("profiles")
	
	var profile models.Profile
	err := collection.FindOne(context.Background(), bson.M{"login": h.config.GitHubUser}).Decode(&profile)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// SyncProfile synchronizes profile data from GitHub
func (h *ProfileHandler) SyncProfile(c *gin.Context) {
	githubClient := github.NewClient(h.config.GitHubToken, h.config.GitHubUser)
	
	profile, err := githubClient.GetProfile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile from GitHub"})
		return
	}

	if h.db == nil {
		// Return the profile data even if we can't save it
		c.JSON(http.StatusOK, gin.H{
			"message": "Profile fetched from GitHub (database not available for saving)",
			"profile": profile,
		})
		return
	}

	collection := h.db.GetCollection("profiles")
	
	// Upsert the profile
	filter := bson.M{"login": profile.Login}
	update := bson.M{
		"$set": profile,
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err = collection.UpdateOne(ctx, filter, update, &options.UpdateOptions{
		Upsert: &[]bool{true}[0],
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile synchronized successfully",
		"profile": profile,
	})
}