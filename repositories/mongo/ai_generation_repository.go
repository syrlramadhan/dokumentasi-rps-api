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

type AIGenerationRepository interface {
	Create(ctx context.Context, generation *models.AIGeneration) (*models.AIGeneration, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.AIGeneration, error)
	FindByGeneratedRPSID(ctx context.Context, generatedRPSID string) (*models.AIGeneration, error)
	FindAll(ctx context.Context, limit, offset int64) ([]models.AIGeneration, error)
	FindByStatus(ctx context.Context, status string) ([]models.AIGeneration, error)
	Update(ctx context.Context, generation *models.AIGeneration) (*models.AIGeneration, error)
	AddAttempt(ctx context.Context, id primitive.ObjectID, attempt models.GenerationAttempt) error
	UpdateFinalStatus(ctx context.Context, id primitive.ObjectID, status string, result map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type aiGenerationRepository struct {
	collection *mongo.Collection
}

func NewAIGenerationRepository(db *mongo.Database) AIGenerationRepository {
	collection := db.Collection("ai_generations")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "generated_rps_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "final_status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}

	collection.Indexes().CreateMany(ctx, indexes)

	return &aiGenerationRepository{collection: collection}
}

func (r *aiGenerationRepository) Create(ctx context.Context, generation *models.AIGeneration) (*models.AIGeneration, error) {
	generation.CreatedAt = time.Now()
	generation.StartedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, generation)
	if err != nil {
		return nil, err
	}

	generation.ID = result.InsertedID.(primitive.ObjectID)
	return generation, nil
}

func (r *aiGenerationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.AIGeneration, error) {
	var generation models.AIGeneration
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&generation)
	if err != nil {
		return nil, err
	}
	return &generation, nil
}

func (r *aiGenerationRepository) FindByGeneratedRPSID(ctx context.Context, generatedRPSID string) (*models.AIGeneration, error) {
	var generation models.AIGeneration
	err := r.collection.FindOne(ctx, bson.M{"generated_rps_id": generatedRPSID}).Decode(&generation)
	if err != nil {
		return nil, err
	}
	return &generation, nil
}

func (r *aiGenerationRepository) FindAll(ctx context.Context, limit, offset int64) ([]models.AIGeneration, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var generations []models.AIGeneration
	if err := cursor.All(ctx, &generations); err != nil {
		return nil, err
	}
	return generations, nil
}

func (r *aiGenerationRepository) FindByStatus(ctx context.Context, status string) ([]models.AIGeneration, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"final_status": status}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var generations []models.AIGeneration
	if err := cursor.All(ctx, &generations); err != nil {
		return nil, err
	}
	return generations, nil
}

func (r *aiGenerationRepository) Update(ctx context.Context, generation *models.AIGeneration) (*models.AIGeneration, error) {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": generation.ID}, generation)
	if err != nil {
		return nil, err
	}
	return generation, nil
}

func (r *aiGenerationRepository) AddAttempt(ctx context.Context, id primitive.ObjectID, attempt models.GenerationAttempt) error {
	update := bson.M{
		"$push": bson.M{"attempts": attempt},
		"$inc":  bson.M{"total_attempts": 1, "total_tokens_used": attempt.TokensUsed, "total_duration_ms": attempt.DurationMs},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *aiGenerationRepository) UpdateFinalStatus(ctx context.Context, id primitive.ObjectID, status string, result map[string]interface{}) error {
	update := bson.M{
		"$set": bson.M{
			"final_status": status,
			"final_result": result,
			"completed_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *aiGenerationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
