package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	models "github.com/syrlramadhan/dokumentasi-rps-api/models/mongo"
)

type AIPromptRepository interface {
	Create(ctx context.Context, prompt *models.AIPrompt) (*models.AIPrompt, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.AIPrompt, error)
	FindByGeneratedRPSID(ctx context.Context, generatedRPSID string) ([]models.AIPrompt, error)
	FindAll(ctx context.Context, limit, offset int64) ([]models.AIPromptSummary, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.AIPromptSummary, error)
	FindByModel(ctx context.Context, model string) ([]models.AIPromptSummary, error)
	FindByStatus(ctx context.Context, status string) ([]models.AIPromptSummary, error)
	Update(ctx context.Context, prompt *models.AIPrompt) (*models.AIPrompt, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetStats(ctx context.Context) (*AIPromptStats, error)
}

type aiPromptRepository struct {
	collection *mongo.Collection
}

type AIPromptStats struct {
	TotalPrompts    int64   `json:"total_prompts"`
	TotalTokens     int64   `json:"total_tokens"`
	SuccessCount    int64   `json:"success_count"`
	FailedCount     int64   `json:"failed_count"`
	AvgResponseTime float64 `json:"avg_response_time_ms"`
	AvgTokensPerReq float64 `json:"avg_tokens_per_request"`
}

func NewAIPromptRepository(db *mongo.Database) AIPromptRepository {
	collection := db.Collection("ai_prompts")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "generated_rps_id", Value: 1}}},
		{Keys: bson.D{{Key: "course_id", Value: 1}}},
		{Keys: bson.D{{Key: "model", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}

	collection.Indexes().CreateMany(ctx, indexes)

	return &aiPromptRepository{collection: collection}
}

func (r *aiPromptRepository) Create(ctx context.Context, prompt *models.AIPrompt) (*models.AIPrompt, error) {
	prompt.CreatedAt = time.Now()
	prompt.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, prompt)
	if err != nil {
		return nil, err
	}

	prompt.ID = result.InsertedID.(primitive.ObjectID)
	return prompt, nil
}

func (r *aiPromptRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.AIPrompt, error) {
	var prompt models.AIPrompt
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&prompt)
	if err != nil {
		return nil, err
	}
	return &prompt, nil
}

func (r *aiPromptRepository) FindByGeneratedRPSID(ctx context.Context, generatedRPSID string) ([]models.AIPrompt, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"generated_rps_id": generatedRPSID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prompts []models.AIPrompt
	if err := cursor.All(ctx, &prompts); err != nil {
		return nil, err
	}
	return prompts, nil
}

func (r *aiPromptRepository) FindAll(ctx context.Context, limit, offset int64) ([]models.AIPromptSummary, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset).
		SetProjection(bson.M{
			"_id":                 1,
			"generated_rps_id":    1,
			"model":               1,
			"status":              1,
			"total_tokens":        1,
			"request_duration_ms": 1,
			"created_at":          1,
		})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []models.AIPromptSummary
	if err := cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}
	return summaries, nil
}

func (r *aiPromptRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.AIPromptSummary, error) {
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetProjection(bson.M{
			"_id":                 1,
			"generated_rps_id":    1,
			"model":               1,
			"status":              1,
			"total_tokens":        1,
			"request_duration_ms": 1,
			"created_at":          1,
		})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []models.AIPromptSummary
	if err := cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}
	return summaries, nil
}

func (r *aiPromptRepository) FindByModel(ctx context.Context, model string) ([]models.AIPromptSummary, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetProjection(bson.M{
			"_id":                 1,
			"generated_rps_id":    1,
			"model":               1,
			"status":              1,
			"total_tokens":        1,
			"request_duration_ms": 1,
			"created_at":          1,
		})

	cursor, err := r.collection.Find(ctx, bson.M{"model": model}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []models.AIPromptSummary
	if err := cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}
	return summaries, nil
}

func (r *aiPromptRepository) FindByStatus(ctx context.Context, status string) ([]models.AIPromptSummary, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetProjection(bson.M{
			"_id":                 1,
			"generated_rps_id":    1,
			"model":               1,
			"status":              1,
			"total_tokens":        1,
			"request_duration_ms": 1,
			"created_at":          1,
		})

	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []models.AIPromptSummary
	if err := cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}
	return summaries, nil
}

func (r *aiPromptRepository) Update(ctx context.Context, prompt *models.AIPrompt) (*models.AIPrompt, error) {
	prompt.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": prompt.ID}, prompt)
	if err != nil {
		return nil, err
	}
	return prompt, nil
}

func (r *aiPromptRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *aiPromptRepository) GetStats(ctx context.Context) (*AIPromptStats, error) {
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":               nil,
				"total_prompts":     bson.M{"$sum": 1},
				"total_tokens":      bson.M{"$sum": "$total_tokens"},
				"success_count":     bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "success"}}, 1, 0}}},
				"failed_count":      bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "failed"}}, 1, 0}}},
				"avg_response_time": bson.M{"$avg": "$request_duration_ms"},
				"avg_tokens":        bson.M{"$avg": "$total_tokens"},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &AIPromptStats{}, nil
	}

	result := results[0]
	return &AIPromptStats{
		TotalPrompts:    getInt64(result, "total_prompts"),
		TotalTokens:     getInt64(result, "total_tokens"),
		SuccessCount:    getInt64(result, "success_count"),
		FailedCount:     getInt64(result, "failed_count"),
		AvgResponseTime: getFloat64(result, "avg_response_time"),
		AvgTokensPerReq: getFloat64(result, "avg_tokens"),
	}, nil
}

// Helper functions
func getInt64(m bson.M, key string) int64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case int64:
			return val
		case int32:
			return int64(val)
		case int:
			return int64(val)
		case float64:
			return int64(val)
		}
	}
	return 0
}

func getFloat64(m bson.M, key string) float64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int64:
			return float64(val)
		case int32:
			return float64(val)
		case int:
			return float64(val)
		}
	}
	return 0
}
