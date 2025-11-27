package repositories

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/models"
	"gorm.io/gorm"
)

type TemplateVersionRepository interface {
	Create(version *models.TemplateVersion) error
	FindAll() ([]models.TemplateVersion, error)
	FindByID(id uuid.UUID) (*models.TemplateVersion, error)
	FindByTemplateID(templateID uuid.UUID) ([]models.TemplateVersion, error)
	FindLatestByTemplateID(templateID uuid.UUID) (*models.TemplateVersion, error)
	Update(version *models.TemplateVersion) error
	Delete(id uuid.UUID) error
}

type templateVersionRepository struct {
	db *gorm.DB
}

func NewTemplateVersionRepository(db *gorm.DB) TemplateVersionRepository {
	return &templateVersionRepository{db: db}
}

func (r *templateVersionRepository) Create(version *models.TemplateVersion) error {
	return r.db.Create(version).Error
}

func (r *templateVersionRepository) FindAll() ([]models.TemplateVersion, error) {
	var versions []models.TemplateVersion
	err := r.db.Preload("Template").Preload("Creator").Find(&versions).Error
	return versions, err
}

func (r *templateVersionRepository) FindByID(id uuid.UUID) (*models.TemplateVersion, error) {
	var version models.TemplateVersion
	err := r.db.Preload("Template").Preload("Creator").First(&version, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *templateVersionRepository) FindByTemplateID(templateID uuid.UUID) ([]models.TemplateVersion, error) {
	var versions []models.TemplateVersion
	err := r.db.Preload("Template").Preload("Creator").Where("template_id = ?", templateID).Order("version DESC").Find(&versions).Error
	return versions, err
}

func (r *templateVersionRepository) FindLatestByTemplateID(templateID uuid.UUID) (*models.TemplateVersion, error) {
	var version models.TemplateVersion
	err := r.db.Preload("Template").Preload("Creator").Where("template_id = ?", templateID).Order("version DESC").First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *templateVersionRepository) Update(version *models.TemplateVersion) error {
	return r.db.Save(version).Error
}

func (r *templateVersionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.TemplateVersion{}, "id = ?", id).Error
}
