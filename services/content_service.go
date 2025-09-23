package services

import (
	"context"
	"fmt"
	"portfolio-backend/database"
	"portfolio-backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContentService struct {
	collection   *mongo.Collection
	cacheService *CacheService
}

func NewContentService() *ContentService {
	return &ContentService{
		collection:   database.Database.Collection("content"),
		cacheService: NewCacheService(),
	}
}

// GetPortfolio retrieves the complete portfolio data
func (cs *ContentService) GetPortfolio(ctx context.Context) (*models.Portfolio, error) {
	// Try cache first
	var portfolio models.Portfolio
	if err := cs.cacheService.GetContentData(ctx, "portfolio", &portfolio); err == nil {
		return &portfolio, nil
	}

	// Build portfolio from individual content pieces
	portfolio = models.Portfolio{
		UpdatedAt: time.Now(),
	}

	// Get meta information
	if meta, err := cs.GetMeta(ctx); err == nil {
		portfolio.Meta = *meta
	}

	// Get skills
	if skills, err := cs.GetSkills(ctx); err == nil {
		portfolio.Skills = *skills
	}

	// Get experience
	if experience, err := cs.GetExperience(ctx); err == nil {
		portfolio.Experience = experience
	}

	// Get projects
	if projects, err := cs.GetProjects(ctx); err == nil {
		portfolio.Projects = projects
	}

	// Get education
	if education, err := cs.GetEducation(ctx); err == nil {
		portfolio.Education = education
	}

	// Cache the complete portfolio
	cs.cacheService.SetContentData(ctx, "portfolio", portfolio)

	return &portfolio, nil
}

// GetMeta retrieves meta information
func (cs *ContentService) GetMeta(ctx context.Context) (*models.Meta, error) {
	var meta models.Meta
	
	// Try cache first
	if err := cs.cacheService.GetContentData(ctx, "meta", &meta); err == nil {
		return &meta, nil
	}

	// Get from database
	var content models.Content
	filter := bson.M{"type": "meta"}
	err := cs.collection.FindOne(ctx, filter).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return default meta if not found
			return &models.Meta{
				Name:     "Portfolio Owner",
				Title:    "Developer",
				Location: "Location",
				GitHub:   "github-username",
			}, nil
		}
		return nil, err
	}

	// Convert interface{} to Meta struct
	if err := convertToStruct(content.Data, &meta); err != nil {
		return nil, err
	}

	// Cache the result
	cs.cacheService.SetContentData(ctx, "meta", meta)

	return &meta, nil
}

// GetSkills retrieves skills information
func (cs *ContentService) GetSkills(ctx context.Context) (*models.Skills, error) {
	var skills models.Skills
	
	// Try cache first
	if err := cs.cacheService.GetContentData(ctx, "skills", &skills); err == nil {
		return &skills, nil
	}

	// Get from database
	var content models.Content
	filter := bson.M{"type": "skills"}
	err := cs.collection.FindOne(ctx, filter).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.Skills{}, nil
		}
		return nil, err
	}

	// Convert interface{} to Skills struct
	if err := convertToStruct(content.Data, &skills); err != nil {
		return nil, err
	}

	// Cache the result
	cs.cacheService.SetContentData(ctx, "skills", skills)

	return &skills, nil
}

// GetExperience retrieves experience information
func (cs *ContentService) GetExperience(ctx context.Context) ([]models.Experience, error) {
	var experience []models.Experience
	
	// Try cache first
	if err := cs.cacheService.GetContentData(ctx, "experience", &experience); err == nil {
		return experience, nil
	}

	// Get from database
	var content models.Content
	filter := bson.M{"type": "experience"}
	err := cs.collection.FindOne(ctx, filter).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Experience{}, nil
		}
		return nil, err
	}

	// Convert interface{} to Experience slice
	if err := convertToStruct(content.Data, &experience); err != nil {
		return nil, err
	}

	// Cache the result
	cs.cacheService.SetContentData(ctx, "experience", experience)

	return experience, nil
}

// GetProjects retrieves projects information
func (cs *ContentService) GetProjects(ctx context.Context) ([]models.Project, error) {
	var projects []models.Project
	
	// Try cache first
	if err := cs.cacheService.GetContentData(ctx, "projects", &projects); err == nil {
		return projects, nil
	}

	// Get from database
	var content models.Content
	filter := bson.M{"type": "projects"}
	err := cs.collection.FindOne(ctx, filter).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Project{}, nil
		}
		return nil, err
	}

	// Convert interface{} to Project slice
	if err := convertToStruct(content.Data, &projects); err != nil {
		return nil, err
	}

	// Cache the result
	cs.cacheService.SetContentData(ctx, "projects", projects)

	return projects, nil
}

// GetEducation retrieves education information
func (cs *ContentService) GetEducation(ctx context.Context) ([]models.Education, error) {
	var education []models.Education
	
	// Try cache first
	if err := cs.cacheService.GetContentData(ctx, "education", &education); err == nil {
		return education, nil
	}

	// Get from database
	var content models.Content
	filter := bson.M{"type": "education"}
	err := cs.collection.FindOne(ctx, filter).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Education{}, nil
		}
		return nil, err
	}

	// Convert interface{} to Education slice
	if err := convertToStruct(content.Data, &education); err != nil {
		return nil, err
	}

	// Cache the result
	cs.cacheService.SetContentData(ctx, "education", education)

	return education, nil
}

// UpdateContent updates content by type
func (cs *ContentService) UpdateContent(ctx context.Context, contentType string, data interface{}, updatedBy string) error {
	now := time.Now()
	
	// Get existing content to increment version
	var existingContent models.Content
	filter := bson.M{"type": contentType}
	err := cs.collection.FindOne(ctx, filter).Decode(&existingContent)
	
	version := 1
	if err == nil {
		version = existingContent.Version + 1
	}

	// Create new content document
	content := models.Content{
		Type:      contentType,
		Data:      data,
		Version:   version,
		UpdatedAt: now,
		UpdatedBy: updatedBy,
	}

	if err == mongo.ErrNoDocuments {
		content.CreatedAt = now
		content.ID = primitive.NewObjectID()
		_, err = cs.collection.InsertOne(ctx, content)
	} else {
		content.CreatedAt = existingContent.CreatedAt
		update := bson.M{"$set": content}
		_, err = cs.collection.UpdateOne(ctx, filter, update)
	}

	if err != nil {
		return err
	}

	// Invalidate cache
	cs.cacheService.InvalidateContentCache(ctx)

	return nil
}

// GetContentHistory retrieves version history for content type
func (cs *ContentService) GetContentHistory(ctx context.Context, contentType string, limit int) ([]models.Content, error) {
	filter := bson.M{"type": contentType}
	opts := options.Find().
		SetSort(bson.D{{Key: "version", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := cs.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []models.Content
	err = cursor.All(ctx, &history)
	return history, err
}

// InitializeDefaultContent creates default content if none exists
func (cs *ContentService) InitializeDefaultContent(ctx context.Context) error {
	// Check if any content exists
	count, err := cs.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Content already exists
	}

	now := time.Now()

	// Default meta
	defaultMeta := models.Meta{
		Name:     "Felipe Macedo",
		Title:    "Desenvolvedor Full Cycle",
		Location: "São Paulo, Brasil",
		GitHub:   "felipemacedo1",
		Bio:      "Desenvolvedor apaixonado por tecnologia e inovação",
	}

	// Default skills
	defaultSkills := models.Skills{
		Backend: []models.Skill{
			{Name: "Java", Level: 90, Category: "backend"},
			{Name: "Go", Level: 80, Category: "backend"},
			{Name: "Spring Boot", Level: 85, Category: "backend"},
		},
		Frontend: []models.Skill{
			{Name: "JavaScript", Level: 80, Category: "frontend"},
			{Name: "React", Level: 80, Category: "frontend"},
		},
	}

	// Create content documents
	contents := []models.Content{
		{
			ID:        primitive.NewObjectID(),
			Type:      "meta",
			Data:      defaultMeta,
			Version:   1,
			CreatedAt: now,
			UpdatedAt: now,
			UpdatedBy: "system",
		},
		{
			ID:        primitive.NewObjectID(),
			Type:      "skills",
			Data:      defaultSkills,
			Version:   1,
			CreatedAt: now,
			UpdatedAt: now,
			UpdatedBy: "system",
		},
	}

	// Insert default content
	var documents []interface{}
	for _, content := range contents {
		documents = append(documents, content)
	}

	_, err = cs.collection.InsertMany(ctx, documents)
	return err
}

// convertToStruct converts interface{} to target struct using BSON
func convertToStruct(source interface{}, target interface{}) error {
	bytes, err := bson.Marshal(source)
	if err != nil {
		return err
	}
	return bson.Unmarshal(bytes, target)
}

// SearchContent performs text search on content
func (cs *ContentService) SearchContent(ctx context.Context, query string, contentTypes []string) ([]models.Content, error) {
	filter := bson.M{}
	
	if len(contentTypes) > 0 {
		filter["type"] = bson.M{"$in": contentTypes}
	}

	// Add text search if query provided
	if query != "" {
		filter["$or"] = []bson.M{
			{"type": bson.M{"$regex": query, "$options": "i"}},
			{"data.name": bson.M{"$regex": query, "$options": "i"}},
			{"data.title": bson.M{"$regex": query, "$options": "i"}},
			{"data.description": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	cursor, err := cs.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Content
	err = cursor.All(ctx, &results)
	return results, err
}