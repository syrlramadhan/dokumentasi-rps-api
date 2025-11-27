package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/repositories"
)

type AuditLogService interface {
	Create(req *dto.CreateAuditLogRequest) (*dto.AuditLogResponse, error)
	FindAll() ([]dto.AuditLogResponse, error)
	FindByID(id int64) (*dto.AuditLogResponse, error)
	FindByUserID(userID uuid.UUID) ([]dto.AuditLogResponse, error)
	FindByAction(action string) ([]dto.AuditLogResponse, error)
	FindByTargetType(targetType string) ([]dto.AuditLogResponse, error)
	FindByTargetID(targetID uuid.UUID) ([]dto.AuditLogResponse, error)
	FindByDateRange(start, end time.Time) ([]dto.AuditLogResponse, error)
	Delete(id int64) error
}

type auditLogService struct {
	repo repositories.AuditLogRepository
}

func NewAuditLogService(repo repositories.AuditLogRepository) AuditLogService {
	return &auditLogService{repo: repo}
}

func (s *auditLogService) Create(req *dto.CreateAuditLogRequest) (*dto.AuditLogResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	log := helper.ToAuditLogModel(req)
	if err := s.repo.Create(log); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToAuditLogResponse(log), nil
}

func (s *auditLogService) FindAll() ([]dto.AuditLogResponse, error) {
	logs, err := s.repo.FindAll()
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToAuditLogResponseList(logs), nil
}

func (s *auditLogService) FindByID(id int64) (*dto.AuditLogResponse, error) {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToAuditLogResponse(log), nil
}

func (s *auditLogService) FindByUserID(userID uuid.UUID) ([]dto.AuditLogResponse, error) {
	logs, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToAuditLogResponseList(logs), nil
}

func (s *auditLogService) FindByAction(action string) ([]dto.AuditLogResponse, error) {
	logs, err := s.repo.FindByAction(action)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToAuditLogResponseList(logs), nil
}

func (s *auditLogService) FindByTargetType(targetType string) ([]dto.AuditLogResponse, error) {
	logs, err := s.repo.FindByTargetType(targetType)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToAuditLogResponseList(logs), nil
}

func (s *auditLogService) FindByTargetID(targetID uuid.UUID) ([]dto.AuditLogResponse, error) {
	logs, err := s.repo.FindByTargetID(targetID)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToAuditLogResponseList(logs), nil
}

func (s *auditLogService) FindByDateRange(start, end time.Time) ([]dto.AuditLogResponse, error) {
	logs, err := s.repo.FindByDateRange(start, end)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToAuditLogResponseList(logs), nil
}

func (s *auditLogService) Delete(id int64) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return helper.WrapDatabaseError(err)
	}

	return s.repo.Delete(id)
}
