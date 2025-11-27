package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
)

type AuditLogController struct {
	service services.AuditLogService
}

func NewAuditLogController(service services.AuditLogService) *AuditLogController {
	return &AuditLogController{service: service}
}

// Create godoc
// @Summary Create a new audit log
// @Tags Audit Logs
// @Accept json
// @Produce json
// @Param request body dto.CreateAuditLogRequest true "Create Audit Log Request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /audit-logs [post]
func (c *AuditLogController) Create(ctx *gin.Context) {
	var req dto.CreateAuditLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	log, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create audit log", "CREATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse("Audit log created successfully", log))
}

// FindAll godoc
// @Summary Get all audit logs
// @Tags Audit Logs
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /audit-logs [get]
func (c *AuditLogController) FindAll(ctx *gin.Context) {
	logs, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch audit logs", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Audit logs fetched successfully", logs))
}

// FindByID godoc
// @Summary Get audit log by ID
// @Tags Audit Logs
// @Produce json
// @Param id path int true "Audit Log ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /audit-logs/{id} [get]
func (c *AuditLogController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid audit log ID", "INVALID_ID", nil))
		return
	}

	log, err := c.service.FindByID(id)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Audit log not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch audit log", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Audit log fetched successfully", log))
}

// FindByUserID godoc
// @Summary Get audit logs by user ID
// @Tags Audit Logs
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} dto.APIResponse
// @Router /audit-logs/user/{user_id} [get]
func (c *AuditLogController) FindByUserID(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid user ID", "INVALID_ID", nil))
		return
	}

	logs, err := c.service.FindByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch audit logs", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Audit logs fetched successfully", logs))
}

// FindByAction godoc
// @Summary Get audit logs by action
// @Tags Audit Logs
// @Produce json
// @Param action path string true "Action"
// @Success 200 {object} dto.APIResponse
// @Router /audit-logs/action/{action} [get]
func (c *AuditLogController) FindByAction(ctx *gin.Context) {
	action := ctx.Param("action")

	logs, err := c.service.FindByAction(action)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch audit logs", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Audit logs fetched successfully", logs))
}

// FindByDateRange godoc
// @Summary Get audit logs by date range
// @Tags Audit Logs
// @Produce json
// @Param start query string true "Start date (RFC3339)"
// @Param end query string true "End date (RFC3339)"
// @Success 200 {object} dto.APIResponse
// @Router /audit-logs/date-range [get]
func (c *AuditLogController) FindByDateRange(ctx *gin.Context) {
	startStr := ctx.Query("start")
	endStr := ctx.Query("end")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid start date format", "INVALID_DATE", nil))
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid end date format", "INVALID_DATE", nil))
		return
	}

	logs, err := c.service.FindByDateRange(start, end)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch audit logs", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Audit logs fetched successfully", logs))
}

// Delete godoc
// @Summary Delete audit log
// @Tags Audit Logs
// @Produce json
// @Param id path int true "Audit Log ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /audit-logs/{id} [delete]
func (c *AuditLogController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid audit log ID", "INVALID_ID", nil))
		return
	}

	if err := c.service.Delete(id); err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Audit log not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete audit log", "DELETE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Audit log deleted successfully", nil))
}
