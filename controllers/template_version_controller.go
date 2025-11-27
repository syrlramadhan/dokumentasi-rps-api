package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
)

type TemplateVersionController struct {
	service services.TemplateVersionService
}

func NewTemplateVersionController(service services.TemplateVersionService) *TemplateVersionController {
	return &TemplateVersionController{service: service}
}

// Create godoc
// @Summary Create a new template version
// @Tags Template Versions
// @Accept json
// @Produce json
// @Param request body dto.CreateTemplateVersionRequest true "Create Template Version Request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /template-versions [post]
func (c *TemplateVersionController) Create(ctx *gin.Context) {
	var req dto.CreateTemplateVersionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	version, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create template version", "CREATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse("Template version created successfully", version))
}

// FindAll godoc
// @Summary Get all template versions
// @Tags Template Versions
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /template-versions [get]
func (c *TemplateVersionController) FindAll(ctx *gin.Context) {
	versions, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch template versions", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Template versions fetched successfully", versions))
}

// FindByID godoc
// @Summary Get template version by ID
// @Tags Template Versions
// @Produce json
// @Param id path string true "Template Version ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /template-versions/{id} [get]
func (c *TemplateVersionController) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid template version ID", "INVALID_ID", nil))
		return
	}

	version, err := c.service.FindByID(id)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Template version not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch template version", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Template version fetched successfully", version))
}

// FindByTemplateID godoc
// @Summary Get template versions by template ID
// @Tags Template Versions
// @Produce json
// @Param template_id path string true "Template ID"
// @Success 200 {object} dto.APIResponse
// @Router /template-versions/template/{template_id} [get]
func (c *TemplateVersionController) FindByTemplateID(ctx *gin.Context) {
	templateID, err := uuid.Parse(ctx.Param("template_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid template ID", "INVALID_ID", nil))
		return
	}

	versions, err := c.service.FindByTemplateID(templateID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch template versions", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Template versions fetched successfully", versions))
}

// FindLatestByTemplateID godoc
// @Summary Get latest template version by template ID
// @Tags Template Versions
// @Produce json
// @Param template_id path string true "Template ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /template-versions/template/{template_id}/latest [get]
func (c *TemplateVersionController) FindLatestByTemplateID(ctx *gin.Context) {
	templateID, err := uuid.Parse(ctx.Param("template_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid template ID", "INVALID_ID", nil))
		return
	}

	version, err := c.service.FindLatestByTemplateID(templateID)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Template version not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch template version", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Latest template version fetched successfully", version))
}

// Update godoc
// @Summary Update template version
// @Tags Template Versions
// @Accept json
// @Produce json
// @Param id path string true "Template Version ID"
// @Param request body dto.UpdateTemplateVersionRequest true "Update Template Version Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /template-versions/{id} [put]
func (c *TemplateVersionController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid template version ID", "INVALID_ID", nil))
		return
	}

	var req dto.UpdateTemplateVersionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	version, err := c.service.Update(id, &req)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Template version not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update template version", "UPDATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Template version updated successfully", version))
}

// Delete godoc
// @Summary Delete template version
// @Tags Template Versions
// @Produce json
// @Param id path string true "Template Version ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /template-versions/{id} [delete]
func (c *TemplateVersionController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid template version ID", "INVALID_ID", nil))
		return
	}

	if err := c.service.Delete(id); err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Template version not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete template version", "DELETE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Template version deleted successfully", nil))
}
