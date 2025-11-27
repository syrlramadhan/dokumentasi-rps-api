package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
)

type GeneratedRPSController struct {
	service services.GeneratedRPSService
}

func NewGeneratedRPSController(service services.GeneratedRPSService) *GeneratedRPSController {
	return &GeneratedRPSController{service: service}
}

// Create godoc
// @Summary Create a new generated RPS
// @Tags Generated RPS
// @Accept json
// @Produce json
// @Param request body dto.CreateGeneratedRPSRequest true "Create Generated RPS Request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /generated-rps [post]
func (c *GeneratedRPSController) Create(ctx *gin.Context) {
	var req dto.CreateGeneratedRPSRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	rps, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create generated RPS", "CREATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse("Generated RPS created successfully", rps))
}

// FindAll godoc
// @Summary Get all generated RPS
// @Tags Generated RPS
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /generated-rps [get]
func (c *GeneratedRPSController) FindAll(ctx *gin.Context) {
	rpsList, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch generated RPS", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Generated RPS fetched successfully", rpsList))
}

// FindByID godoc
// @Summary Get generated RPS by ID
// @Tags Generated RPS
// @Produce json
// @Param id path string true "Generated RPS ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /generated-rps/{id} [get]
func (c *GeneratedRPSController) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid generated RPS ID", "INVALID_ID", nil))
		return
	}

	rps, err := c.service.FindByID(id)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Generated RPS not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch generated RPS", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Generated RPS fetched successfully", rps))
}

// FindByCourseID godoc
// @Summary Get generated RPS by course ID
// @Tags Generated RPS
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} dto.APIResponse
// @Router /generated-rps/course/{course_id} [get]
func (c *GeneratedRPSController) FindByCourseID(ctx *gin.Context) {
	courseID, err := uuid.Parse(ctx.Param("course_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid course ID", "INVALID_ID", nil))
		return
	}

	rpsList, err := c.service.FindByCourseID(courseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch generated RPS", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Generated RPS fetched successfully", rpsList))
}

// FindByStatus godoc
// @Summary Get generated RPS by status
// @Tags Generated RPS
// @Produce json
// @Param status path string true "Status (queued|processing|done|failed)"
// @Success 200 {object} dto.APIResponse
// @Router /generated-rps/status/{status} [get]
func (c *GeneratedRPSController) FindByStatus(ctx *gin.Context) {
	status := ctx.Param("status")
	if status != "queued" && status != "processing" && status != "done" && status != "failed" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid status", "INVALID_STATUS", nil))
		return
	}

	rpsList, err := c.service.FindByStatus(status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch generated RPS", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Generated RPS fetched successfully", rpsList))
}

// Update godoc
// @Summary Update generated RPS
// @Tags Generated RPS
// @Accept json
// @Produce json
// @Param id path string true "Generated RPS ID"
// @Param request body dto.UpdateGeneratedRPSRequest true "Update Generated RPS Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /generated-rps/{id} [put]
func (c *GeneratedRPSController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid generated RPS ID", "INVALID_ID", nil))
		return
	}

	var req dto.UpdateGeneratedRPSRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	rps, err := c.service.Update(id, &req)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Generated RPS not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update generated RPS", "UPDATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Generated RPS updated successfully", rps))
}

// UpdateStatus godoc
// @Summary Update generated RPS status
// @Tags Generated RPS
// @Accept json
// @Produce json
// @Param id path string true "Generated RPS ID"
// @Param request body dto.UpdateGeneratedRPSStatusRequest true "Update Status Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /generated-rps/{id}/status [patch]
func (c *GeneratedRPSController) UpdateStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid generated RPS ID", "INVALID_ID", nil))
		return
	}

	var req dto.UpdateGeneratedRPSStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	if err := c.service.UpdateStatus(id, req.Status); err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Generated RPS not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update status", "UPDATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Status updated successfully", nil))
}

// Delete godoc
// @Summary Delete generated RPS
// @Tags Generated RPS
// @Produce json
// @Param id path string true "Generated RPS ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /generated-rps/{id} [delete]
func (c *GeneratedRPSController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid generated RPS ID", "INVALID_ID", nil))
		return
	}

	if err := c.service.Delete(id); err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Generated RPS not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete generated RPS", "DELETE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Generated RPS deleted successfully", nil))
}

// Export godoc
// @Summary Export generated RPS to PDF/DOCX
// @Tags Generated RPS
// @Produce application/pdf,application/vnd.openxmlformats-officedocument.wordprocessingml.document
// @Param id path string true "Generated RPS ID"
// @Param format query string false "Export format (pdf|docx)" default(pdf)
// @Success 200 {file} binary
// @Failure 404 {object} dto.APIResponse
// @Router /generated/{id}/export [get]
func (c *GeneratedRPSController) Export(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid generated RPS ID", "INVALID_ID", nil))
		return
	}

	format := ctx.DefaultQuery("format", "pdf")
	if format != "pdf" && format != "docx" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid format. Use 'pdf' or 'docx'", "INVALID_FORMAT", nil))
		return
	}

	rps, err := c.service.FindByID(id)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Generated RPS not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch generated RPS", "FETCH_ERROR", nil))
		return
	}

	// Check if exported file URL exists
	if rps.ExportedFileURL != nil && *rps.ExportedFileURL != "" {
		ctx.Redirect(http.StatusFound, *rps.ExportedFileURL)
		return
	}

	// TODO: Generate PDF/DOCX file on-the-fly or queue for generation
	ctx.JSON(http.StatusOK, dto.SuccessResponse("Export functionality - TODO: implement file generation", gin.H{
		"id":     rps.ID,
		"format": format,
		"status": rps.Status,
	}))
}

// CompleteGeneration godoc
// @Summary Complete generation (internal worker endpoint)
// @Tags Internal
// @Accept json
// @Produce json
// @Param request body dto.CompleteGenerationRequest true "Complete Generation Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /internal/complete_generation [post]
func (c *GeneratedRPSController) CompleteGeneration(ctx *gin.Context) {
	var req dto.CompleteGenerationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	// Update the generated RPS with result
	updateReq := &dto.UpdateGeneratedRPSRequest{
		Status:          &req.Status,
		Result:          req.Result,
		ExportedFileURL: req.ExportedFileURL,
		AIMetadata:      req.AIMetadata,
	}

	rps, err := c.service.Update(req.JobID, updateReq)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Generated RPS not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to complete generation", "UPDATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Generation completed successfully", rps))
}
