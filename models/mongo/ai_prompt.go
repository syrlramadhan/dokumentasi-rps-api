package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AIPrompt represents a single AI prompt/completion record
type AIPrompt struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GeneratedRPSID string             `bson:"generated_rps_id" json:"generated_rps_id"`
	CourseID       string             `bson:"course_id" json:"course_id"`
	TemplateID     string             `bson:"template_id" json:"template_id"`

	// Prompt details
	SystemPrompt string `bson:"system_prompt" json:"system_prompt"`
	UserPrompt   string `bson:"user_prompt" json:"user_prompt"`
	FullPrompt   string `bson:"full_prompt" json:"full_prompt"`

	// Response details
	Response       string                 `bson:"response" json:"response"`
	ParsedResponse map[string]interface{} `bson:"parsed_response" json:"parsed_response"`

	// Model configuration
	Model            string  `bson:"model" json:"model"`
	Temperature      float64 `bson:"temperature" json:"temperature"`
	MaxTokens        int     `bson:"max_tokens" json:"max_tokens"`
	TopP             float64 `bson:"top_p,omitempty" json:"top_p,omitempty"`
	FrequencyPenalty float64 `bson:"frequency_penalty,omitempty" json:"frequency_penalty,omitempty"`
	PresencePenalty  float64 `bson:"presence_penalty,omitempty" json:"presence_penalty,omitempty"`

	// Usage statistics
	PromptTokens     int `bson:"prompt_tokens" json:"prompt_tokens"`
	CompletionTokens int `bson:"completion_tokens" json:"completion_tokens"`
	TotalTokens      int `bson:"total_tokens" json:"total_tokens"`

	// Timing
	RequestDurationMs int64 `bson:"request_duration_ms" json:"request_duration_ms"`

	// Status
	Status       string `bson:"status" json:"status"` // success, failed, timeout
	ErrorMessage string `bson:"error_message,omitempty" json:"error_message,omitempty"`

	// Metadata
	RequestID      string `bson:"request_id,omitempty" json:"request_id,omitempty"`
	FinishReason   string `bson:"finish_reason,omitempty" json:"finish_reason,omitempty"`
	ResponseFormat string `bson:"response_format" json:"response_format"`

	// Input context
	CourseData   map[string]interface{} `bson:"course_data" json:"course_data"`
	TemplateData map[string]interface{} `bson:"template_data" json:"template_data"`
	Options      map[string]interface{} `bson:"options" json:"options"`

	// Timestamps
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// AIPromptSummary is a lightweight version for listing
type AIPromptSummary struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	GeneratedRPSID    string             `bson:"generated_rps_id" json:"generated_rps_id"`
	Model             string             `bson:"model" json:"model"`
	Status            string             `bson:"status" json:"status"`
	TotalTokens       int                `bson:"total_tokens" json:"total_tokens"`
	RequestDurationMs int64              `bson:"request_duration_ms" json:"request_duration_ms"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
}
