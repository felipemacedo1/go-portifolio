package utils

import (
	"fmt"
	"portfolio-backend/models"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ValidationResult contains validation results
type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors"`
}

// Validator provides validation functionality
type Validator struct {
	errors []ValidationError
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: []ValidationError{},
	}
}

// AddError adds a validation error
func (v *Validator) AddError(field, message, code string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

// IsValid returns true if no validation errors exist
func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}

// GetResult returns the validation result
func (v *Validator) GetResult() ValidationResult {
	return ValidationResult{
		IsValid: v.IsValid(),
		Errors:  v.errors,
	}
}

// Reset clears all validation errors
func (v *Validator) Reset() {
	v.errors = []ValidationError{}
}

// Field validation methods

// Required validates that a field is not empty
func (v *Validator) Required(field string, value interface{}) *Validator {
	if isEmptyValue(value) {
		v.AddError(field, "This field is required", "REQUIRED")
	}
	return v
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field string, value string, min int) *Validator {
	if len(value) < min {
		v.AddError(field, fmt.Sprintf("Must be at least %d characters long", min), "MIN_LENGTH")
	}
	return v
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field string, value string, max int) *Validator {
	if len(value) > max {
		v.AddError(field, fmt.Sprintf("Must be at most %d characters long", max), "MAX_LENGTH")
	}
	return v
}

// Length validates exact string length
func (v *Validator) Length(field string, value string, length int) *Validator {
	if len(value) != length {
		v.AddError(field, fmt.Sprintf("Must be exactly %d characters long", length), "INVALID_LENGTH")
	}
	return v
}

// Min validates minimum numeric value
func (v *Validator) Min(field string, value, min int) *Validator {
	if value < min {
		v.AddError(field, fmt.Sprintf("Must be at least %d", min), "MIN_VALUE")
	}
	return v
}

// Max validates maximum numeric value
func (v *Validator) Max(field string, value, max int) *Validator {
	if value > max {
		v.AddError(field, fmt.Sprintf("Must be at most %d", max), "MAX_VALUE")
	}
	return v
}

// Range validates numeric value within range
func (v *Validator) Range(field string, value, min, max int) *Validator {
	if value < min || value > max {
		v.AddError(field, fmt.Sprintf("Must be between %d and %d", min, max), "OUT_OF_RANGE")
	}
	return v
}

// Email validates email format
func (v *Validator) Email(field string, value string) *Validator {
	if value != "" && !IsValidEmail(value) {
		v.AddError(field, "Invalid email format", "INVALID_EMAIL")
	}
	return v
}

// URL validates URL format
func (v *Validator) URL(field string, value string) *Validator {
	if value != "" && !IsValidURL(value) {
		v.AddError(field, "Invalid URL format", "INVALID_URL")
	}
	return v
}

// GitHubUsername validates GitHub username format
func (v *Validator) GitHubUsername(field string, value string) *Validator {
	if value != "" && !IsValidGitHubUsername(value) {
		v.AddError(field, "Invalid GitHub username format", "INVALID_GITHUB_USERNAME")
	}
	return v
}

// Regex validates against a regular expression
func (v *Validator) Regex(field string, value string, pattern string, message string) *Validator {
	if value != "" {
		matched, err := regexp.MatchString(pattern, value)
		if err != nil || !matched {
			v.AddError(field, message, "REGEX_MISMATCH")
		}
	}
	return v
}

// OneOf validates that value is one of the allowed values
func (v *Validator) OneOf(field string, value string, allowed []string) *Validator {
	if value != "" && !Contains(allowed, value) {
		v.AddError(field, fmt.Sprintf("Must be one of: %s", strings.Join(allowed, ", ")), "INVALID_CHOICE")
	}
	return v
}

// Date validates date format and range
func (v *Validator) Date(field string, value time.Time) *Validator {
	if value.IsZero() {
		v.AddError(field, "Invalid date", "INVALID_DATE")
	}
	return v
}

// FutureDate validates that date is in the future
func (v *Validator) FutureDate(field string, value time.Time) *Validator {
	if !value.IsZero() && value.Before(time.Now()) {
		v.AddError(field, "Date must be in the future", "DATE_NOT_FUTURE")
	}
	return v
}

// PastDate validates that date is in the past
func (v *Validator) PastDate(field string, value time.Time) *Validator {
	if !value.IsZero() && value.After(time.Now()) {
		v.AddError(field, "Date must be in the past", "DATE_NOT_PAST")
	}
	return v
}

// DateRange validates date within a range
func (v *Validator) DateRange(field string, value, start, end time.Time) *Validator {
	if !value.IsZero() && (value.Before(start) || value.After(end)) {
		v.AddError(field, "Date is outside allowed range", "DATE_OUT_OF_RANGE")
	}
	return v
}

// Custom validation function type
type CustomValidationFunc func(value interface{}) bool

// Custom validates using a custom function
func (v *Validator) Custom(field string, value interface{}, fn CustomValidationFunc, message string) *Validator {
	if !fn(value) {
		v.AddError(field, message, "CUSTOM_VALIDATION")
	}
	return v
}

// Struct validation methods

// ValidateMeta validates meta content
func (v *Validator) ValidateMeta(meta *models.Meta) *Validator {
	v.Required("name", meta.Name).
		MinLength("name", meta.Name, 2).
		MaxLength("name", meta.Name, 100)

	v.Required("title", meta.Title).
		MinLength("title", meta.Title, 2).
		MaxLength("title", meta.Title, 200)

	v.MaxLength("location", meta.Location, 100)
	v.GitHubUsername("github", meta.GitHub)
	v.Email("email", meta.Email)
	v.URL("linkedin", meta.LinkedIn)
	v.URL("website", meta.Website)
	v.MaxLength("bio", meta.Bio, 500)

	return v
}

// ValidateSkill validates skill data
func (v *Validator) ValidateSkill(skill *models.Skill) *Validator {
	v.Required("name", skill.Name).
		MinLength("name", skill.Name, 1).
		MaxLength("name", skill.Name, 50)

	v.Range("level", skill.Level, 0, 100)
	v.MaxLength("category", skill.Category, 50)
	v.Min("years_exp", skill.YearsExp, 0)

	return v
}

// ValidateExperience validates experience data
func (v *Validator) ValidateExperience(exp *models.Experience) *Validator {
	v.Required("company", exp.Company).
		MinLength("company", exp.Company, 1).
		MaxLength("company", exp.Company, 100)

	v.Required("position", exp.Position).
		MinLength("position", exp.Position, 1).
		MaxLength("position", exp.Position, 100)

	v.MaxLength("location", exp.Location, 100)
	v.MaxLength("description", exp.Description, 1000)
	v.URL("company_url", exp.CompanyURL)

	// Validate dates
	if !exp.StartDate.IsZero() {
		v.PastDate("start_date", exp.StartDate)
		
		if exp.EndDate != nil && !exp.EndDate.IsZero() {
			if exp.EndDate.Before(exp.StartDate) {
				v.AddError("end_date", "End date must be after start date", "INVALID_DATE_RANGE")
			}
		}
	}

	return v
}

// ValidateProject validates project data
func (v *Validator) ValidateProject(project *models.Project) *Validator {
	v.Required("name", project.Name).
		MinLength("name", project.Name, 1).
		MaxLength("name", project.Name, 100)

	v.MaxLength("description", project.Description, 500)
	v.MaxLength("long_description", project.LongDesc, 2000)
	v.URL("github_url", project.GitHubURL)
	v.URL("live_url", project.LiveURL)
	v.URL("demo_url", project.DemoURL)

	// Validate status
	validStatuses := []string{"completed", "in-progress", "planned", "archived"}
	v.OneOf("status", project.Status, validStatuses)

	// Validate dates
	if !project.StartDate.IsZero() && project.EndDate != nil && !project.EndDate.IsZero() {
		if project.EndDate.Before(project.StartDate) {
			v.AddError("end_date", "End date must be after start date", "INVALID_DATE_RANGE")
		}
	}

	return v
}

// ValidateEducation validates education data
func (v *Validator) ValidateEducation(edu *models.Education) *Validator {
	v.Required("institution", edu.Institution).
		MinLength("institution", edu.Institution, 1).
		MaxLength("institution", edu.Institution, 100)

	v.Required("degree", edu.Degree).
		MinLength("degree", edu.Degree, 1).
		MaxLength("degree", edu.Degree, 100)

	v.MaxLength("field", edu.Field, 100)
	v.MaxLength("description", edu.Description, 1000)
	v.URL("url", edu.URL)

	// Validate GPA
	if edu.GPA > 0 {
		if edu.GPA < 0 || edu.GPA > 4.0 {
			v.AddError("gpa", "GPA must be between 0 and 4.0", "INVALID_GPA")
		}
	}

	// Validate dates
	if !edu.StartDate.IsZero() && edu.EndDate != nil && !edu.EndDate.IsZero() {
		if edu.EndDate.Before(edu.StartDate) {
			v.AddError("end_date", "End date must be after start date", "INVALID_DATE_RANGE")
		}
	}

	return v
}

// ValidateContentUpdateRequest validates content update request
func (v *Validator) ValidateContentUpdateRequest(req *models.ContentUpdateRequest) *Validator {
	validTypes := []string{"meta", "skills", "experience", "projects", "education"}
	v.Required("type", req.Type).
		OneOf("type", req.Type, validTypes)

	v.Required("data", req.Data)

	return v
}

// ValidateGitHubSyncRequest validates GitHub sync request
func (v *Validator) ValidateGitHubSyncRequest(req *models.GitHubSyncRequest) *Validator {
	v.Required("username", req.Username).
		GitHubUsername("username", req.Username)

	return v
}

// Helper functions

// isEmptyValue checks if a value is considered empty
func isEmptyValue(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String()) == ""
	case reflect.Bool:
		return false // bool is never empty
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return false // numbers are never empty
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return false // numbers are never empty
	case reflect.Float32, reflect.Float64:
		return false // numbers are never empty
	case reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}

	return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
}

// ValidateQueryParams validates common query parameters
func ValidateQueryParams(page, limit string) (int, int, []ValidationError) {
	var errors []ValidationError
	
	pageNum := 1
	limitNum := 10

	if page != "" {
		p, err := strconv.Atoi(page)
		if err != nil || p < 1 {
			errors = append(errors, ValidationError{
				Field:   "page",
				Message: "Page must be a positive integer",
				Code:    "INVALID_PAGE",
			})
		} else {
			pageNum = p
		}
	}

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil || l < 1 || l > 100 {
			errors = append(errors, ValidationError{
				Field:   "limit",
				Message: "Limit must be between 1 and 100",
				Code:    "INVALID_LIMIT",
			})
		} else {
			limitNum = l
		}
	}

	return pageNum, limitNum, errors
}