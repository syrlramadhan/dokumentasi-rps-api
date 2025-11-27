package services

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/repositories"
)

type ProgramService interface {
	Create(req *dto.CreateProgramRequest) (*dto.ProgramResponse, error)
	FindAll() ([]dto.ProgramResponse, error)
	FindByID(id uuid.UUID) (*dto.ProgramResponse, error)
	FindByCode(code string) (*dto.ProgramResponse, error)
	Update(id uuid.UUID, req *dto.UpdateProgramRequest) (*dto.ProgramResponse, error)
	Delete(id uuid.UUID) error
}

type programService struct {
	repo repositories.ProgramRepository
}

func NewProgramService(repo repositories.ProgramRepository) ProgramService {
	return &programService{repo: repo}
}

func (s *programService) Create(req *dto.CreateProgramRequest) (*dto.ProgramResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	program := helper.ToProgramModel(req)
	if err := s.repo.Create(program); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToProgramResponse(program), nil
}

func (s *programService) FindAll() ([]dto.ProgramResponse, error) {
	programs, err := s.repo.FindAll()
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToProgramResponseList(programs), nil
}

func (s *programService) FindByID(id uuid.UUID) (*dto.ProgramResponse, error) {
	program, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToProgramResponse(program), nil
}

func (s *programService) FindByCode(code string) (*dto.ProgramResponse, error) {
	program, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToProgramResponse(program), nil
}

func (s *programService) Update(id uuid.UUID, req *dto.UpdateProgramRequest) (*dto.ProgramResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	program, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	if req.Code != nil {
		program.Code = *req.Code
	}
	if req.Name != nil {
		program.Name = *req.Name
	}

	if err := s.repo.Update(program); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToProgramResponse(program), nil
}

func (s *programService) Delete(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return helper.WrapDatabaseError(err)
	}

	return s.repo.Delete(id)
}
