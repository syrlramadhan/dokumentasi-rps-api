package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
)

type TemplateController struct {
	service services.TemplateService
}

func NewTemplateController(service services.TemplateService) *TemplateController {
	return &TemplateController{service: service}
}

// Create godoc
// @Summary Create a new template
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body dto.CreateTemplateRequest true "Create Template Request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /templates [post]
func (c *TemplateController) Create(ctx *gin.Context) {
	var req dto.CreateTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	template, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create template", "CREATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse("Template created successfully", template))
}

// FindAll godoc
// @Summary Get all templates
// @Tags Templates
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /templates [get]
func (c *TemplateController) FindAll(ctx *gin.Context) {
	templates, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch templates", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Templates fetched successfully", templates))
}

// FindByID godoc
// @Summary Get template by ID
// @Tags Templates
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /templates/{id} [get]
func (c *TemplateController) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid template ID", "INVALID_ID", nil))
		return
	}

	template, err := c.service.FindByID(id)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Template not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch template", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Template fetched successfully", template))
}

// FindByProgramID godoc
// @Summary Get templates by program ID
// @Tags Templates
// @Produce json
// @Param program_id path string true "Program ID"
// @Success 200 {object} dto.APIResponse
// @Router /templates/program/{program_id} [get]
func (c *TemplateController) FindByProgramID(ctx *gin.Context) {
	programID, err := uuid.Parse(ctx.Param("program_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid program ID", "INVALID_ID", nil))
		return
	}

	templates, err := c.service.FindByProgramID(programID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch templates", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Templates fetched successfully", templates))
}

// FindActiveByProgramID godoc
// @Summary Get active templates by program ID
// @Tags Templates
// @Produce json
// @Param program_id path string true "Program ID"
// @Success 200 {object} dto.APIResponse
// @Router /templates/program/{program_id}/active [get]
func (c *TemplateController) FindActiveByProgramID(ctx *gin.Context) {
	programID, err := uuid.Parse(ctx.Param("program_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid program ID", "INVALID_ID", nil))
		return
	}

	templates, err := c.service.FindActiveByProgramID(programID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch templates", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Active templates fetched successfully", templates))
}

// Update godoc
// @Summary Update template
// @Tags Templates
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Param request body dto.UpdateTemplateRequest true "Update Template Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /templates/{id} [put]
func (c *TemplateController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid template ID", "INVALID_ID", nil))
		return
	}

	var req dto.UpdateTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	template, err := c.service.Update(id, &req)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Template not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update template", "UPDATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Template updated successfully", template))
}

// Delete godoc
// @Summary Delete template
// @Tags Templates
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /templates/{id} [delete]
func (c *TemplateController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid template ID", "INVALID_ID", nil))
		return
	}

	if err := c.service.Delete(id); err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Template not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete template", "DELETE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Template deleted successfully", nil))
}
