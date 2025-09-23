package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository represents a GitHub repository
type Repository struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	GitHubID    int64              `json:"github_id" bson:"github_id"`
	Name        string             `json:"name" bson:"name"`
	FullName    string             `json:"full_name" bson:"full_name"`
	Description string             `json:"description" bson:"description"`
	HTMLURL     string             `json:"html_url" bson:"html_url"`
	Language    string             `json:"language" bson:"language"`
	Stars       int                `json:"stars" bson:"stars"`
	Forks       int                `json:"forks" bson:"forks"`
	Private     bool               `json:"private" bson:"private"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Topics      []string           `json:"topics" bson:"topics"`
}

// Profile represents the user's GitHub profile
type Profile struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	GitHubID        int64              `json:"github_id" bson:"github_id"`
	Login           string             `json:"login" bson:"login"`
	Name            string             `json:"name" bson:"name"`
	Bio             string             `json:"bio" bson:"bio"`
	AvatarURL       string             `json:"avatar_url" bson:"avatar_url"`
	Location        string             `json:"location" bson:"location"`
	Company         string             `json:"company" bson:"company"`
	Blog            string             `json:"blog" bson:"blog"`
	Email           string             `json:"email" bson:"email"`
	PublicRepos     int                `json:"public_repos" bson:"public_repos"`
	Followers       int                `json:"followers" bson:"followers"`
	Following       int                `json:"following" bson:"following"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
	LastSyncedAt    time.Time          `json:"last_synced_at" bson:"last_synced_at"`
}

// Skill represents a skill/technology
type Skill struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Category string             `json:"category" bson:"category"`
	Level    string             `json:"level" bson:"level"` // beginner, intermediate, advanced, expert
	Icon     string             `json:"icon,omitempty" bson:"icon,omitempty"`
}

// Project represents a portfolio project
type Project struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Technologies []string          `json:"technologies" bson:"technologies"`
	GitHubURL   string             `json:"github_url,omitempty" bson:"github_url,omitempty"`
	LiveURL     string             `json:"live_url,omitempty" bson:"live_url,omitempty"`
	ImageURL    string             `json:"image_url,omitempty" bson:"image_url,omitempty"`
	Featured    bool               `json:"featured" bson:"featured"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Experience represents work experience
type Experience struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Company     string             `json:"company" bson:"company"`
	Position    string             `json:"position" bson:"position"`
	Description string             `json:"description" bson:"description"`
	StartDate   time.Time          `json:"start_date" bson:"start_date"`
	EndDate     *time.Time         `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Current     bool               `json:"current" bson:"current"`
	Location    string             `json:"location,omitempty" bson:"location,omitempty"`
	Technologies []string          `json:"technologies" bson:"technologies"`
}