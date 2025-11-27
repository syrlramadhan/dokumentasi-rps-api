package repositories

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/models"
	"gorm.io/gorm"
)

type ProgramRepository interface {
	Create(program *models.Program) error
	FindAll() ([]models.Program, error)
	FindByID(id uuid.UUID) (*models.Program, error)
	FindByCode(code string) (*models.Program, error)
	Update(program *models.Program) error
	Delete(id uuid.UUID) error
}

type programRepository struct {
	db *gorm.DB
}

func NewProgramRepository(db *gorm.DB) ProgramRepository {
	return &programRepository{db: db}
}

func (r *programRepository) Create(program *models.Program) error {
	return r.db.Create(program).Error
}

func (r *programRepository) FindAll() ([]models.Program, error) {
	var programs []models.Program
	err := r.db.Find(&programs).Error
	return programs, err
}

func (r *programRepository) FindByID(id uuid.UUID) (*models.Program, error) {
	var program models.Program
	err := r.db.First(&program, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &program, nil
}

func (r *programRepository) FindByCode(code string) (*models.Program, error) {
	var program models.Program
	err := r.db.First(&program, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &program, nil
}

func (r *programRepository) Update(program *models.Program) error {
	return r.db.Save(program).Error
}

func (r *programRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Program{}, "id = ?", id).Error
}
