package repositories

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/models"
	"gorm.io/gorm"
)

type TemplateRepository interface {
	Create(template *models.Template) error
	FindAll() ([]models.Template, error)
	FindByID(id uuid.UUID) (*models.Template, error)
	FindByProgramID(programID uuid.UUID) ([]models.Template, error)
	FindActiveByProgramID(programID uuid.UUID) ([]models.Template, error)
	Update(template *models.Template) error
	Delete(id uuid.UUID) error
}

type templateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *templateRepository) FindAll() ([]models.Template, error) {
	var templates []models.Template
	err := r.db.Preload("Program").Preload("Creator").Find(&templates).Error
	return templates, err
}

func (r *templateRepository) FindByID(id uuid.UUID) (*models.Template, error) {
	var template models.Template
	err := r.db.Preload("Program").Preload("Creator").First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *templateRepository) FindByProgramID(programID uuid.UUID) ([]models.Template, error) {
	var templates []models.Template
	err := r.db.Preload("Program").Preload("Creator").Where("program_id = ?", programID).Find(&templates).Error
	return templates, err
}

func (r *templateRepository) FindActiveByProgramID(programID uuid.UUID) ([]models.Template, error) {
	var templates []models.Template
	err := r.db.Preload("Program").Preload("Creator").Where("program_id = ? AND is_active = ?", programID, true).Find(&templates).Error
	return templates, err
}

func (r *templateRepository) Update(template *models.Template) error {
	return r.db.Save(template).Error
}

func (r *templateRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Template{}, "id = ?", id).Error
}
