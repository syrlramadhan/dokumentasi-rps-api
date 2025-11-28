package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PromptTemplate represents reusable prompt templates
type PromptTemplate struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Version     int                `bson:"version" json:"version"`

	// Prompt content
	SystemPrompt       string `bson:"system_prompt" json:"system_prompt"`
	UserPromptTemplate string `bson:"user_prompt_template" json:"user_prompt_template"`

	// Variables that can be replaced in template
	Variables []PromptVariable `bson:"variables" json:"variables"`

	// Model settings
	DefaultModel       string  `bson:"default_model" json:"default_model"`
	DefaultTemperature float64 `bson:"default_temperature" json:"default_temperature"`
	DefaultMaxTokens   int     `bson:"default_max_tokens" json:"default_max_tokens"`

	// JSON Schema for structured output
	OutputSchema map[string]interface{} `bson:"output_schema,omitempty" json:"output_schema,omitempty"`

	// Metadata
	Category string   `bson:"category" json:"category"` // rps_generation, review, summary
	Tags     []string `bson:"tags" json:"tags"`
	IsActive bool     `bson:"is_active" json:"is_active"`

	// Usage stats
	UsageCount  int     `bson:"usage_count" json:"usage_count"`
	SuccessRate float64 `bson:"success_rate" json:"success_rate"`
	AvgTokens   int     `bson:"avg_tokens" json:"avg_tokens"`

	// Audit
	CreatedBy string    `bson:"created_by" json:"created_by"`
	UpdatedBy string    `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// PromptVariable defines a variable in the prompt template
type PromptVariable struct {
	Name         string `bson:"name" json:"name"`
	Description  string `bson:"description" json:"description"`
	Required     bool   `bson:"required" json:"required"`
	DefaultValue string `bson:"default_value,omitempty" json:"default_value,omitempty"`
	Type         string `bson:"type" json:"type"` // string, number, array, object
}
