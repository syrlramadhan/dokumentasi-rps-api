package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AIGeneration represents the full AI generation process with all attempts
type AIGeneration struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GeneratedRPSID    string             `bson:"generated_rps_id" json:"generated_rps_id"`
	CourseID          string             `bson:"course_id" json:"course_id"`
	CourseName        string             `bson:"course_name" json:"course_name"`
	CourseCode        string             `bson:"course_code" json:"course_code"`
	TemplateVersionID string             `bson:"template_version_id" json:"template_version_id"`

	// Generation attempts (bisa ada retry)
	Attempts      []GenerationAttempt `bson:"attempts" json:"attempts"`
	TotalAttempts int                 `bson:"total_attempts" json:"total_attempts"`

	// Final result
	FinalStatus string                 `bson:"final_status" json:"final_status"` // success, failed
	FinalResult map[string]interface{} `bson:"final_result,omitempty" json:"final_result,omitempty"`

	// Aggregated stats
	TotalTokensUsed int64   `bson:"total_tokens_used" json:"total_tokens_used"`
	TotalDurationMs int64   `bson:"total_duration_ms" json:"total_duration_ms"`
	TotalCost       float64 `bson:"total_cost" json:"total_cost"` // estimated cost in USD

	// User info
	RequestedBy string `bson:"requested_by,omitempty" json:"requested_by,omitempty"`

	// Timestamps
	StartedAt   time.Time `bson:"started_at" json:"started_at"`
	CompletedAt time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

// GenerationAttempt represents a single attempt in the generation process
type GenerationAttempt struct {
	AttemptNumber int                `bson:"attempt_number" json:"attempt_number"`
	PromptID      primitive.ObjectID `bson:"prompt_id" json:"prompt_id"`
	Status        string             `bson:"status" json:"status"`
	TokensUsed    int                `bson:"tokens_used" json:"tokens_used"`
	DurationMs    int64              `bson:"duration_ms" json:"duration_ms"`
	ErrorMessage  string             `bson:"error_message,omitempty" json:"error_message,omitempty"`
	Timestamp     time.Time          `bson:"timestamp" json:"timestamp"`
}
