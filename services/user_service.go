package services

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/repositories"
)

type UserService interface {
	Create(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	FindAll() ([]dto.UserResponse, error)
	FindByID(id uuid.UUID) (*dto.UserResponse, error)
	FindByUsername(username string) (*dto.UserResponse, error)
	Update(id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(id uuid.UUID) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	user := helper.ToUserModel(req)
	if err := s.repo.Create(user); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToUserResponse(user), nil
}

func (s *userService) FindAll() ([]dto.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToUserResponseList(users), nil
}

func (s *userService) FindByID(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToUserResponse(user), nil
}

func (s *userService) FindByUsername(username string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToUserResponse(user), nil
}

func (s *userService) Update(id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if err := helper.ValidateStruct(req); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.DisplayName != nil {
		user.DisplayName = req.DisplayName
	}
	if req.Role != nil {
		user.Role = *req.Role
	}

	if err := s.repo.Update(user); err != nil {
		return nil, helper.WrapDatabaseError(err)
	}

	return helper.ToUserResponse(user), nil
}

func (s *userService) Delete(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return helper.WrapDatabaseError(err)
	}

	return s.repo.Delete(id)
}
