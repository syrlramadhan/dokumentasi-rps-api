package services

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/repositories"
)

type CourseService interface {
	Create(req *dto.CreateCourseRequest) (*dto.CourseResponse, error)
	FindAll() ([]dto.CourseResponse, error)
	FindByID(id uuid.UUID) (*dto.CourseResponse, error)
	FindByProgramID(programID uuid.UUID) ([]dto.CourseResponse, error)
	Update(id uuid.UUID, req *dto.UpdateCourseRequest) (*dto.CourseResponse, error)
	Delete(id uuid.UUID) error
}

type courseService struct {
	repo repositories.CourseRepository
}

func NewCourseService(repo repositories.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) Create(req *dto.CreateCourseRequest) (*dto.CourseResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	course := helper.ToCourseModel(req)
	if err := s.repo.Create(course); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToCourseResponse(course), nil
}

func (s *courseService) FindAll() ([]dto.CourseResponse, error) {
	courses, err := s.repo.FindAll()
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToCourseResponseList(courses), nil
}

func (s *courseService) FindByID(id uuid.UUID) (*dto.CourseResponse, error) {
	course, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToCourseResponse(course), nil
}

func (s *courseService) FindByProgramID(programID uuid.UUID) ([]dto.CourseResponse, error) {
	courses, err := s.repo.FindByProgramID(programID)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToCourseResponseList(courses), nil
}

func (s *courseService) Update(id uuid.UUID, req *dto.UpdateCourseRequest) (*dto.CourseResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	course, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	if req.ProgramID != nil {
		course.ProgramID = req.ProgramID
	}
	if req.Code != nil {
		course.Code = *req.Code
	}
	if req.Title != nil {
		course.Title = *req.Title
	}
	if req.Credits != nil {
		course.Credits = req.Credits
	}

	if err := s.repo.Update(course); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToCourseResponse(course), nil
}

func (s *courseService) Delete(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return helper.WrapDatabaseError(err)
	}

	return s.repo.Delete(id)
}
