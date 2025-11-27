package services

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/repositories"
)

type TemplateVersionService interface {
	Create(req *dto.CreateTemplateVersionRequest) (*dto.TemplateVersionResponse, error)
	FindAll() ([]dto.TemplateVersionResponse, error)
	FindByID(id uuid.UUID) (*dto.TemplateVersionResponse, error)
	FindByTemplateID(templateID uuid.UUID) ([]dto.TemplateVersionResponse, error)
	FindLatestByTemplateID(templateID uuid.UUID) (*dto.TemplateVersionResponse, error)
	Update(id uuid.UUID, req *dto.UpdateTemplateVersionRequest) (*dto.TemplateVersionResponse, error)
	Delete(id uuid.UUID) error
}

type templateVersionService struct {
	repo repositories.TemplateVersionRepository
}

func NewTemplateVersionService(repo repositories.TemplateVersionRepository) TemplateVersionService {
	return &templateVersionService{repo: repo}
}

func (s *templateVersionService) Create(req *dto.CreateTemplateVersionRequest) (*dto.TemplateVersionResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	version := helper.ToTemplateVersionModel(req)
	if err := s.repo.Create(version); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateVersionResponse(version), nil
}

func (s *templateVersionService) FindAll() ([]dto.TemplateVersionResponse, error) {
	versions, err := s.repo.FindAll()
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateVersionResponseList(versions), nil
}

func (s *templateVersionService) FindByID(id uuid.UUID) (*dto.TemplateVersionResponse, error) {
	version, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateVersionResponse(version), nil
}

func (s *templateVersionService) FindByTemplateID(templateID uuid.UUID) ([]dto.TemplateVersionResponse, error) {
	versions, err := s.repo.FindByTemplateID(templateID)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateVersionResponseList(versions), nil
}

func (s *templateVersionService) FindLatestByTemplateID(templateID uuid.UUID) (*dto.TemplateVersionResponse, error) {
	version, err := s.repo.FindLatestByTemplateID(templateID)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateVersionResponse(version), nil
}

func (s *templateVersionService) Update(id uuid.UUID, req *dto.UpdateTemplateVersionRequest) (*dto.TemplateVersionResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	version, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	if req.Version != nil {
		version.Version = *req.Version
	}
	if req.Definition != nil {
		version.Definition = req.Definition
	}

	if err := s.repo.Update(version); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateVersionResponse(version), nil
}

func (s *templateVersionService) Delete(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return helper.WrapDatabaseError(err)
	}

	return s.repo.Delete(id)
}
