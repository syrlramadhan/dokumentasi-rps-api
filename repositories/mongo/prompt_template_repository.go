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

type PromptTemplateRepository interface {
	Create(ctx context.Context, template *models.PromptTemplate) (*models.PromptTemplate, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.PromptTemplate, error)
	FindByName(ctx context.Context, name string) (*models.PromptTemplate, error)
	FindAll(ctx context.Context) ([]models.PromptTemplate, error)
	FindActive(ctx context.Context) ([]models.PromptTemplate, error)
	FindByCategory(ctx context.Context, category string) ([]models.PromptTemplate, error)
	Update(ctx context.Context, template *models.PromptTemplate) (*models.PromptTemplate, error)
	IncrementUsage(ctx context.Context, id primitive.ObjectID) error
	UpdateSuccessRate(ctx context.Context, id primitive.ObjectID, successRate float64) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type promptTemplateRepository struct {
	collection *mongo.Collection
}

func NewPromptTemplateRepository(db *mongo.Database) PromptTemplateRepository {
	collection := db.Collection("prompt_templates")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "name", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "is_active", Value: 1}}},
	}

	collection.Indexes().CreateMany(ctx, indexes)

	return &promptTemplateRepository{collection: collection}
}

func (r *promptTemplateRepository) Create(ctx context.Context, template *models.PromptTemplate) (*models.PromptTemplate, error) {
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	template.UsageCount = 0
	template.SuccessRate = 0

	result, err := r.collection.InsertOne(ctx, template)
	if err != nil {
		return nil, err
	}

	template.ID = result.InsertedID.(primitive.ObjectID)
	return template, nil
}

func (r *promptTemplateRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.PromptTemplate, error) {
	var template models.PromptTemplate
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&template)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *promptTemplateRepository) FindByName(ctx context.Context, name string) (*models.PromptTemplate, error) {
	var template models.PromptTemplate
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&template)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *promptTemplateRepository) FindAll(ctx context.Context) ([]models.PromptTemplate, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []models.PromptTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *promptTemplateRepository) FindActive(ctx context.Context) ([]models.PromptTemplate, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []models.PromptTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *promptTemplateRepository) FindByCategory(ctx context.Context, category string) ([]models.PromptTemplate, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"category": category, "is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []models.PromptTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *promptTemplateRepository) Update(ctx context.Context, template *models.PromptTemplate) (*models.PromptTemplate, error) {
	template.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": template.ID}, template)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (r *promptTemplateRepository) IncrementUsage(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$inc": bson.M{"usage_count": 1},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *promptTemplateRepository) UpdateSuccessRate(ctx context.Context, id primitive.ObjectID, successRate float64) error {
	update := bson.M{
		"$set": bson.M{
			"success_rate": successRate,
			"updated_at":   time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *promptTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
