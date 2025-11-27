package repositories

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/models"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) error
	FindAll() ([]models.Course, error)
	FindByID(id uuid.UUID) (*models.Course, error)
	FindByProgramID(programID uuid.UUID) ([]models.Course, error)
	FindByCode(code string) (*models.Course, error)
	Update(course *models.Course) error
	Delete(id uuid.UUID) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) FindAll() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Preload("Program").Find(&courses).Error
	return courses, err
}

func (r *courseRepository) FindByID(id uuid.UUID) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Program").First(&course, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) FindByProgramID(programID uuid.UUID) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Preload("Program").Where("program_id = ?", programID).Find(&courses).Error
	return courses, err
}

func (r *courseRepository) FindByCode(code string) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Program").First(&course, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Course{}, "id = ?", id).Error
}
