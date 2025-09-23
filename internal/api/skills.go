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

// SkillHandler handles skill-related requests
type SkillHandler struct {
	db *database.DB
}

// NewSkillHandler creates a new skill handler
func NewSkillHandler(db *database.DB) *SkillHandler {
	return &SkillHandler{db: db}
}

// GetSkills retrieves all skills
func (h *SkillHandler) GetSkills(c *gin.Context) {
	collection := h.db.GetCollection("skills")
	
	opts := options.Find().SetSort(bson.D{{Key: "category", Value: 1}, {Key: "name", Value: 1}})
	cursor, err := collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve skills"})
		return
	}
	defer cursor.Close(context.Background())

	var skills []models.Skill
	if err = cursor.All(context.Background(), &skills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode skills"})
		return
	}

	if skills == nil {
		skills = []models.Skill{}
	}

	c.JSON(http.StatusOK, skills)
}

// CreateSkill creates a new skill
func (h *SkillHandler) CreateSkill(c *gin.Context) {
	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := h.db.GetCollection("skills")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	result, err := collection.InsertOne(ctx, skill)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create skill"})
		return
	}

	skill.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, skill)
}

// UpdateSkill updates an existing skill
func (h *SkillHandler) UpdateSkill(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := h.db.GetCollection("skills")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": skill}
	
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update skill"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	skill.ID = objID
	c.JSON(http.StatusOK, skill)
}

// DeleteSkill deletes a skill
func (h *SkillHandler) DeleteSkill(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	collection := h.db.GetCollection("skills")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skill"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted successfully"})
}