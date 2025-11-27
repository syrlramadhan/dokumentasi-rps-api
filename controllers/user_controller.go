package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// Create godoc
// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Create User Request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /users [post]
func (c *UserController) Create(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	user, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create user", "CREATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse("User created successfully", user))
}

// FindAll godoc
// @Summary Get all users
// @Tags Users
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /users [get]
func (c *UserController) FindAll(ctx *gin.Context) {
	users, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch users", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Users fetched successfully", users))
}

// FindByID godoc
// @Summary Get user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /users/{id} [get]
func (c *UserController) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid user ID", "INVALID_ID", nil))
		return
	}

	user, err := c.service.FindByID(id)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("User not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch user", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("User fetched successfully", user))
}

// Update godoc
// @Summary Update user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "Update User Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /users/{id} [put]
func (c *UserController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid user ID", "INVALID_ID", nil))
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	user, err := c.service.Update(id, &req)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("User not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update user", "UPDATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("User updated successfully", user))
}

// Delete godoc
// @Summary Delete user
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /users/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid user ID", "INVALID_ID", nil))
		return
	}

	if err := c.service.Delete(id); err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("User not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete user", "DELETE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("User deleted successfully", nil))
}
