package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Portfolio content structures
type Portfolio struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Meta      Meta               `bson:"meta" json:"meta"`
	Skills    Skills             `bson:"skills" json:"skills"`
	Experience []Experience      `bson:"experience" json:"experience"`
	Projects  []Project          `bson:"projects" json:"projects"`
	Education []Education        `bson:"education" json:"education"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Meta struct {
	Name     string `bson:"name" json:"name" validate:"required"`
	Title    string `bson:"title" json:"title" validate:"required"`
	Location string `bson:"location" json:"location"`
	GitHub   string `bson:"github" json:"github"`
	Email    string `bson:"email" json:"email" validate:"email"`
	LinkedIn string `bson:"linkedin" json:"linkedin"`
	Website  string `bson:"website" json:"website"`
	Bio      string `bson:"bio" json:"bio"`
}

type Skills struct {
	Backend    []Skill `bson:"backend" json:"backend"`
	Frontend   []Skill `bson:"frontend" json:"frontend"`
	Database   []Skill `bson:"database" json:"database"`
	DevOps     []Skill `bson:"devops" json:"devops"`
	Tools      []Skill `bson:"tools" json:"tools"`
	Languages  []Skill `bson:"languages" json:"languages"`
}

type Skill struct {
	Name        string   `bson:"name" json:"name" validate:"required"`
	Level       int      `bson:"level" json:"level" validate:"min=0,max=100"`
	Category    string   `bson:"category" json:"category"`
	Icon        string   `bson:"icon" json:"icon"`
	YearsExp    int      `bson:"years_exp" json:"years_exp"`
	Certifications []string `bson:"certifications" json:"certifications"`
}

type Experience struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Company     string            `bson:"company" json:"company" validate:"required"`
	Position    string            `bson:"position" json:"position" validate:"required"`
	Location    string            `bson:"location" json:"location"`
	StartDate   time.Time         `bson:"start_date" json:"start_date"`
	EndDate     *time.Time        `bson:"end_date,omitempty" json:"end_date,omitempty"`
	IsCurrent   bool              `bson:"is_current" json:"is_current"`
	Description string            `bson:"description" json:"description"`
	Achievements []string         `bson:"achievements" json:"achievements"`
	Technologies []string         `bson:"technologies" json:"technologies"`
	CompanyLogo string            `bson:"company_logo" json:"company_logo"`
	CompanyURL  string            `bson:"company_url" json:"company_url"`
}

type Project struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string            `bson:"name" json:"name" validate:"required"`
	Description  string            `bson:"description" json:"description"`
	LongDesc     string            `bson:"long_description" json:"long_description"`
	Technologies []string          `bson:"technologies" json:"technologies"`
	GitHubURL    string            `bson:"github_url" json:"github_url"`
	LiveURL      string            `bson:"live_url" json:"live_url"`
	DemoURL      string            `bson:"demo_url" json:"demo_url"`
	Images       []string          `bson:"images" json:"images"`
	Featured     bool              `bson:"featured" json:"featured"`
	Status       string            `bson:"status" json:"status"` // "completed", "in-progress", "planned"
	StartDate    time.Time         `bson:"start_date" json:"start_date"`
	EndDate      *time.Time        `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Category     string            `bson:"category" json:"category"`
	Highlights   []string          `bson:"highlights" json:"highlights"`
	Challenges   []string          `bson:"challenges" json:"challenges"`
	Stars        int               `bson:"stars" json:"stars"`
	Forks        int               `bson:"forks" json:"forks"`
	Language     string            `bson:"language" json:"language"`
	UpdatedAt    time.Time         `bson:"updated_at" json:"updated_at"`
}

type Education struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Institution  string            `bson:"institution" json:"institution" validate:"required"`
	Degree       string            `bson:"degree" json:"degree" validate:"required"`
	Field        string            `bson:"field" json:"field"`
	StartDate    time.Time         `bson:"start_date" json:"start_date"`
	EndDate      *time.Time        `bson:"end_date,omitempty" json:"end_date,omitempty"`
	GPA          float64           `bson:"gpa,omitempty" json:"gpa,omitempty"`
	Honors       []string          `bson:"honors" json:"honors"`
	Courses      []string          `bson:"courses" json:"courses"`
	Description  string            `bson:"description" json:"description"`
	Logo         string            `bson:"logo" json:"logo"`
	URL          string            `bson:"url" json:"url"`
}

// Content types for flexible content management
type Content struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      string            `bson:"type" json:"type" validate:"required"` // "skills", "experience", "projects", "education", "meta"
	Data      interface{}       `bson:"data" json:"data"`
	Version   int               `bson:"version" json:"version"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedBy string            `bson:"updated_by" json:"updated_by"`
}

// Cache structure for storing temporary data
type CacheEntry struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Key       string            `bson:"key" json:"key" validate:"required"`
	Value     interface{}       `bson:"value" json:"value"`
	ExpiresAt time.Time         `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
}