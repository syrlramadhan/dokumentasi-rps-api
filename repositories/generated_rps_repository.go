package repositories

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/models"
	"gorm.io/gorm"
)

type GeneratedRPSRepository interface {
	Create(rps *models.GeneratedRPS) error
	FindAll() ([]models.GeneratedRPS, error)
	FindByID(id uuid.UUID) (*models.GeneratedRPS, error)
	FindByCourseID(courseID uuid.UUID) ([]models.GeneratedRPS, error)
	FindByGeneratedBy(userID uuid.UUID) ([]models.GeneratedRPS, error)
	FindByStatus(status string) ([]models.GeneratedRPS, error)
	Update(rps *models.GeneratedRPS) error
	UpdateStatus(id uuid.UUID, status string) error
	Delete(id uuid.UUID) error
}

type generatedRPSRepository struct {
	db *gorm.DB
}

func NewGeneratedRPSRepository(db *gorm.DB) GeneratedRPSRepository {
	return &generatedRPSRepository{db: db}
}

func (r *generatedRPSRepository) Create(rps *models.GeneratedRPS) error {
	return r.db.Create(rps).Error
}

func (r *generatedRPSRepository) FindAll() ([]models.GeneratedRPS, error) {
	var rpsList []models.GeneratedRPS
	err := r.db.Preload("TemplateVersion").Preload("Course").Preload("Generator").Find(&rpsList).Error
	return rpsList, err
}

func (r *generatedRPSRepository) FindByID(id uuid.UUID) (*models.GeneratedRPS, error) {
	var rps models.GeneratedRPS
	err := r.db.Preload("TemplateVersion").Preload("Course").Preload("Generator").First(&rps, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rps, nil
}

func (r *generatedRPSRepository) FindByCourseID(courseID uuid.UUID) ([]models.GeneratedRPS, error) {
	var rpsList []models.GeneratedRPS
	err := r.db.Preload("TemplateVersion").Preload("Course").Preload("Generator").Where("course_id = ?", courseID).Find(&rpsList).Error
	return rpsList, err
}

func (r *generatedRPSRepository) FindByGeneratedBy(userID uuid.UUID) ([]models.GeneratedRPS, error) {
	var rpsList []models.GeneratedRPS
	err := r.db.Preload("TemplateVersion").Preload("Course").Preload("Generator").Where("generated_by = ?", userID).Find(&rpsList).Error
	return rpsList, err
}

func (r *generatedRPSRepository) FindByStatus(status string) ([]models.GeneratedRPS, error) {
	var rpsList []models.GeneratedRPS
	err := r.db.Preload("TemplateVersion").Preload("Course").Preload("Generator").Where("status = ?", status).Find(&rpsList).Error
	return rpsList, err
}

func (r *generatedRPSRepository) Update(rps *models.GeneratedRPS) error {
	return r.db.Save(rps).Error
}

func (r *generatedRPSRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.GeneratedRPS{}).Where("id = ?", id).Update("status", status).Error
}

func (r *generatedRPSRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.GeneratedRPS{}, "id = ?", id).Error
}
