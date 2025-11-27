package services

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/repositories"
)

type TemplateService interface {
	Create(req *dto.CreateTemplateRequest) (*dto.TemplateResponse, error)
	FindAll() ([]dto.TemplateResponse, error)
	FindByID(id uuid.UUID) (*dto.TemplateResponse, error)
	FindByProgramID(programID uuid.UUID) ([]dto.TemplateResponse, error)
	FindActiveByProgramID(programID uuid.UUID) ([]dto.TemplateResponse, error)
	Update(id uuid.UUID, req *dto.UpdateTemplateRequest) (*dto.TemplateResponse, error)
	Delete(id uuid.UUID) error
}

type templateService struct {
	repo repositories.TemplateRepository
}

func NewTemplateService(repo repositories.TemplateRepository) TemplateService {
	return &templateService{repo: repo}
}

func (s *templateService) Create(req *dto.CreateTemplateRequest) (*dto.TemplateResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	template := helper.ToTemplateModel(req)
	if err := s.repo.Create(template); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateResponse(template), nil
}

func (s *templateService) FindAll() ([]dto.TemplateResponse, error) {
	templates, err := s.repo.FindAll()
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateResponseList(templates), nil
}

func (s *templateService) FindByID(id uuid.UUID) (*dto.TemplateResponse, error) {
	template, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateResponse(template), nil
}

func (s *templateService) FindByProgramID(programID uuid.UUID) ([]dto.TemplateResponse, error) {
	templates, err := s.repo.FindByProgramID(programID)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateResponseList(templates), nil
}

func (s *templateService) FindActiveByProgramID(programID uuid.UUID) ([]dto.TemplateResponse, error) {
	templates, err := s.repo.FindActiveByProgramID(programID)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateResponseList(templates), nil
}

func (s *templateService) Update(id uuid.UUID, req *dto.UpdateTemplateRequest) (*dto.TemplateResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	template, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	if req.ProgramID != nil {
		template.ProgramID = req.ProgramID
	}
	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.Description != nil {
		template.Description = req.Description
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}

	if err := s.repo.Update(template); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToTemplateResponse(template), nil
}

func (s *templateService) Delete(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return helper.WrapDatabaseError(err)
	}

	return s.repo.Delete(id)
}
