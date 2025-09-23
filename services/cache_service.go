package services

import (
	"context"
	"encoding/json"
	"portfolio-backend/config"
	"portfolio-backend/database"
	"portfolio-backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CacheService struct {
	collection *mongo.Collection
}

func NewCacheService() *CacheService {
	return &CacheService{
		collection: database.Database.Collection("cache"),
	}
}

// Get retrieves a cached value by key
func (cs *CacheService) Get(ctx context.Context, key string, target interface{}) error {
	var cacheEntry models.CacheEntry
	
	filter := bson.M{
		"key": key,
		"expires_at": bson.M{"$gt": time.Now()},
	}
	
	err := cs.collection.FindOne(ctx, filter).Decode(&cacheEntry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("cache miss: %s", key)
		}
		return err
	}

	// Convert the cached value to the target type
	jsonBytes, err := json.Marshal(cacheEntry.Value)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, target)
}

// Set stores a value in cache with TTL
func (cs *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	cacheEntry := models.CacheEntry{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}

	// Use upsert to replace existing entries
	filter := bson.M{"key": key}
	update := bson.M{"$set": cacheEntry}
	opts := options.Update().SetUpsert(true)

	_, err := cs.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// Delete removes a cached value
func (cs *CacheService) Delete(ctx context.Context, key string) error {
	filter := bson.M{"key": key}
	_, err := cs.collection.DeleteOne(ctx, filter)
	return err
}

// DeletePattern removes all cache entries matching a pattern
func (cs *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	filter := bson.M{"key": bson.M{"$regex": pattern}}
	_, err := cs.collection.DeleteMany(ctx, filter)
	return err
}

// Exists checks if a key exists in cache and is not expired
func (cs *CacheService) Exists(ctx context.Context, key string) bool {
	filter := bson.M{
		"key": key,
		"expires_at": bson.M{"$gt": time.Now()},
	}
	
	count, err := cs.collection.CountDocuments(ctx, filter)
	return err == nil && count > 0
}

// GetTTL returns the remaining TTL for a key
func (cs *CacheService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	var cacheEntry models.CacheEntry
	
	filter := bson.M{"key": key}
	err := cs.collection.FindOne(ctx, filter).Decode(&cacheEntry)
	if err != nil {
		return 0, err
	}

	remaining := cacheEntry.ExpiresAt.Sub(time.Now())
	if remaining < 0 {
		return 0, fmt.Errorf("key expired")
	}

	return remaining, nil
}

// Cleanup removes expired entries (called by background job)
func (cs *CacheService) Cleanup(ctx context.Context) error {
	filter := bson.M{"expires_at": bson.M{"$lt": time.Now()}}
	result, err := cs.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount > 0 {
		fmt.Printf("Cleaned up %d expired cache entries\n", result.DeletedCount)
	}

	return nil
}

// GetStats returns cache statistics
func (cs *CacheService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	totalCount, err := cs.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	activeCount, err := cs.collection.CountDocuments(ctx, bson.M{
		"expires_at": bson.M{"$gt": time.Now()},
	})
	if err != nil {
		return nil, err
	}

	expiredCount := totalCount - activeCount

	return map[string]interface{}{
		"total_entries":   totalCount,
		"active_entries":  activeCount,
		"expired_entries": expiredCount,
		"hit_rate":        calculateHitRate(ctx, cs.collection),
	}, nil
}

// Helper methods for common cache operations

// GetGitHubData retrieves GitHub data from cache
func (cs *CacheService) GetGitHubData(ctx context.Context, username string, dataType string, target interface{}) error {
	key := fmt.Sprintf("github:%s:%s", username, dataType)
	return cs.Get(ctx, key, target)
}

// SetGitHubData stores GitHub data in cache
func (cs *CacheService) SetGitHubData(ctx context.Context, username string, dataType string, data interface{}) error {
	key := fmt.Sprintf("github:%s:%s", username, dataType)
	return cs.Set(ctx, key, data, config.AppConfig.GitHubCacheTTL)
}

// GetContentData retrieves content data from cache
func (cs *CacheService) GetContentData(ctx context.Context, contentType string, target interface{}) error {
	key := fmt.Sprintf("content:%s", contentType)
	return cs.Get(ctx, key, target)
}

// SetContentData stores content data in cache
func (cs *CacheService) SetContentData(ctx context.Context, contentType string, data interface{}) error {
	key := fmt.Sprintf("content:%s", contentType)
	return cs.Set(ctx, key, data, config.AppConfig.ContentCacheTTL)
}

// InvalidateGitHubCache removes all GitHub cache entries for a user
func (cs *CacheService) InvalidateGitHubCache(ctx context.Context, username string) error {
	pattern := fmt.Sprintf("github:%s:.*", username)
	return cs.DeletePattern(ctx, pattern)
}

// InvalidateContentCache removes all content cache entries
func (cs *CacheService) InvalidateContentCache(ctx context.Context) error {
	pattern := "content:.*"
	return cs.DeletePattern(ctx, pattern)
}

// Warm cache with initial data
func (cs *CacheService) WarmCache(ctx context.Context) error {
	// This method can be called on startup to pre-populate cache
	// with frequently accessed data
	return nil
}

func calculateHitRate(ctx context.Context, collection *mongo.Collection) float64 {
	// This is a simplified calculation
	// In a real implementation, you'd want to track hits/misses separately
	return 0.85 // Placeholder
}

// Background cleanup job
func (cs *CacheService) StartCleanupJob() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := cs.Cleanup(ctx); err != nil {
				fmt.Printf("Cache cleanup error: %v\n", err)
			}
			cancel()
		}
	}()
}