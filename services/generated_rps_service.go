package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/repositories"
)

type GeneratedRPSService interface {
	Create(req *dto.CreateGeneratedRPSRequest) (*dto.GeneratedRPSResponse, error)
	FindAll() ([]dto.GeneratedRPSResponse, error)
	FindByID(id uuid.UUID) (*dto.GeneratedRPSResponse, error)
	FindByCourseID(courseID uuid.UUID) ([]dto.GeneratedRPSResponse, error)
	FindByGeneratedBy(userID uuid.UUID) ([]dto.GeneratedRPSResponse, error)
	FindByStatus(status string) ([]dto.GeneratedRPSResponse, error)
	Update(id uuid.UUID, req *dto.UpdateGeneratedRPSRequest) (*dto.GeneratedRPSResponse, error)
	UpdateStatus(id uuid.UUID, status string) error
	Delete(id uuid.UUID) error
}

type generatedRPSService struct {
	repo repositories.GeneratedRPSRepository
}

func NewGeneratedRPSService(repo repositories.GeneratedRPSRepository) GeneratedRPSService {
	return &generatedRPSService{repo: repo}
}

func (s *generatedRPSService) Create(req *dto.CreateGeneratedRPSRequest) (*dto.GeneratedRPSResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	rps := helper.ToGeneratedRPSModel(req)
	if err := s.repo.Create(rps); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToGeneratedRPSResponse(rps), nil
}

func (s *generatedRPSService) FindAll() ([]dto.GeneratedRPSResponse, error) {
	rpsList, err := s.repo.FindAll()
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToGeneratedRPSResponseList(rpsList), nil
}

func (s *generatedRPSService) FindByID(id uuid.UUID) (*dto.GeneratedRPSResponse, error) {
	rps, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToGeneratedRPSResponse(rps), nil
}

func (s *generatedRPSService) FindByCourseID(courseID uuid.UUID) ([]dto.GeneratedRPSResponse, error) {
	rpsList, err := s.repo.FindByCourseID(courseID)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToGeneratedRPSResponseList(rpsList), nil
}

func (s *generatedRPSService) FindByGeneratedBy(userID uuid.UUID) ([]dto.GeneratedRPSResponse, error) {
	rpsList, err := s.repo.FindByGeneratedBy(userID)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToGeneratedRPSResponseList(rpsList), nil
}

func (s *generatedRPSService) FindByStatus(status string) ([]dto.GeneratedRPSResponse, error) {
	rpsList, err := s.repo.FindByStatus(status)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToGeneratedRPSResponseList(rpsList), nil
}

func (s *generatedRPSService) Update(id uuid.UUID, req *dto.UpdateGeneratedRPSRequest) (*dto.GeneratedRPSResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	rps, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	if req.Status != nil {
		rps.Status = *req.Status
	}
	if req.Result != nil {
		rps.Result = req.Result
	}
	if req.ExportedFileURL != nil {
		rps.ExportedFileURL = req.ExportedFileURL
	}
	if req.AIMetadata != nil {
		rps.AIMetadata = req.AIMetadata
	}
	rps.UpdatedAt = time.Now()

	if err := s.repo.Update(rps); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToGeneratedRPSResponse(rps), nil
}

func (s *generatedRPSService) UpdateStatus(id uuid.UUID, status string) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return helper.WrapDatabaseError(err)
	}

	return s.repo.UpdateStatus(id, status)
}

func (s *generatedRPSService) Delete(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return helper.WrapDatabaseError(err)
	}

	return s.repo.Delete(id)
}
