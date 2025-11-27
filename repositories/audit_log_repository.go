package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/models"
	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(log *models.AuditLog) error
	FindAll() ([]models.AuditLog, error)
	FindByID(id int64) (*models.AuditLog, error)
	FindByUserID(userID uuid.UUID) ([]models.AuditLog, error)
	FindByAction(action string) ([]models.AuditLog, error)
	FindByTargetType(targetType string) ([]models.AuditLog, error)
	FindByTargetID(targetID uuid.UUID) ([]models.AuditLog, error)
	FindByDateRange(start, end time.Time) ([]models.AuditLog, error)
	Delete(id int64) error
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditLogRepository) FindAll() ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Preload("User").Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByID(id int64) (*models.AuditLog, error) {
	var log models.AuditLog
	err := r.db.Preload("User").First(&log, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *auditLogRepository) FindByUserID(userID uuid.UUID) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByAction(action string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Preload("User").Where("action = ?", action).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByTargetType(targetType string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Preload("User").Where("target_type = ?", targetType).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByTargetID(targetID uuid.UUID) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Preload("User").Where("target_id = ?", targetID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByDateRange(start, end time.Time) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Preload("User").Where("created_at BETWEEN ? AND ?", start, end).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) Delete(id int64) error {
	return r.db.Delete(&models.AuditLog{}, "id = ?", id).Error
}
