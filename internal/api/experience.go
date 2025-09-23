package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/felipemacedo1/b/internal/database"
	"github.com/felipemacedo1/b/internal/models"
)

// ExperienceHandler handles experience-related requests
type ExperienceHandler struct {
	db *database.DB
}

// NewExperienceHandler creates a new experience handler
func NewExperienceHandler(db *database.DB) *ExperienceHandler {
	return &ExperienceHandler{db: db}
}

// GetExperience retrieves all work experience
func (h *ExperienceHandler) GetExperience(c *gin.Context) {
	collection := h.db.GetCollection("experience")
	
	// Sort by current first, then by start_date descending
	opts := options.Find().SetSort(bson.D{
		{Key: "current", Value: -1},
		{Key: "start_date", Value: -1},
	})
	
	cursor, err := collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve experience"})
		return
	}
	defer cursor.Close(context.Background())

	var experience []models.Experience
	if err = cursor.All(context.Background(), &experience); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode experience"})
		return
	}

	if experience == nil {
		experience = []models.Experience{}
	}

	c.JSON(http.StatusOK, experience)
}

// CreateExperience creates a new work experience entry
func (h *ExperienceHandler) CreateExperience(c *gin.Context) {
	var experience models.Experience
	if err := c.ShouldBindJSON(&experience); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := h.db.GetCollection("experience")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	result, err := collection.InsertOne(ctx, experience)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create experience"})
		return
	}

	experience.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, experience)
}

// UpdateExperience updates an existing work experience entry
func (h *ExperienceHandler) UpdateExperience(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid experience ID"})
		return
	}

	var experience models.Experience
	if err := c.ShouldBindJSON(&experience); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := h.db.GetCollection("experience")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": experience}
	
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update experience"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experience not found"})
		return
	}

	experience.ID = objID
	c.JSON(http.StatusOK, experience)
}

// DeleteExperience deletes a work experience entry
func (h *ExperienceHandler) DeleteExperience(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid experience ID"})
		return
	}

	collection := h.db.GetCollection("experience")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete experience"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experience not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Experience deleted successfully"})
}