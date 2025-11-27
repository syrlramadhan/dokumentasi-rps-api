package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
)

type ProgramController struct {
	service services.ProgramService
}

func NewProgramController(service services.ProgramService) *ProgramController {
	return &ProgramController{service: service}
}

// Create godoc
// @Summary Create a new program
// @Tags Programs
// @Accept json
// @Produce json
// @Param request body dto.CreateProgramRequest true "Create Program Request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /programs [post]
func (c *ProgramController) Create(ctx *gin.Context) {
	var req dto.CreateProgramRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	program, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create program", "CREATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse("Program created successfully", program))
}

// FindAll godoc
// @Summary Get all programs
// @Tags Programs
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /programs [get]
func (c *ProgramController) FindAll(ctx *gin.Context) {
	programs, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch programs", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Programs fetched successfully", programs))
}

// FindByID godoc
// @Summary Get program by ID
// @Tags Programs
// @Produce json
// @Param id path string true "Program ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /programs/{id} [get]
func (c *ProgramController) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid program ID", "INVALID_ID", nil))
		return
	}

	program, err := c.service.FindByID(id)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Program not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch program", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Program fetched successfully", program))
}

// Update godoc
// @Summary Update program
// @Tags Programs
// @Accept json
// @Produce json
// @Param id path string true "Program ID"
// @Param request body dto.UpdateProgramRequest true "Update Program Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /programs/{id} [put]
func (c *ProgramController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid program ID", "INVALID_ID", nil))
		return
	}

	var req dto.UpdateProgramRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	program, err := c.service.Update(id, &req)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Program not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update program", "UPDATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Program updated successfully", program))
}

// Delete godoc
// @Summary Delete program
// @Tags Programs
// @Produce json
// @Param id path string true "Program ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /programs/{id} [delete]
func (c *ProgramController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid program ID", "INVALID_ID", nil))
		return
	}

	if err := c.service.Delete(id); err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Program not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete program", "DELETE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Program deleted successfully", nil))
}
