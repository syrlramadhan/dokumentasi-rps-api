package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/helper"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
)

type CourseController struct {
	service services.CourseService
}

func NewCourseController(service services.CourseService) *CourseController {
	return &CourseController{service: service}
}

// Create godoc
// @Summary Create a new course
// @Tags Courses
// @Accept json
// @Produce json
// @Param request body dto.CreateCourseRequest true "Create Course Request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /courses [post]
func (c *CourseController) Create(ctx *gin.Context) {
	var req dto.CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	course, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create course", "CREATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse("Course created successfully", course))
}

// FindAll godoc
// @Summary Get all courses
// @Tags Courses
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /courses [get]
func (c *CourseController) FindAll(ctx *gin.Context) {
	courses, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch courses", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Courses fetched successfully", courses))
}

// FindByID godoc
// @Summary Get course by ID
// @Tags Courses
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /courses/{id} [get]
func (c *CourseController) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid course ID", "INVALID_ID", nil))
		return
	}

	course, err := c.service.FindByID(id)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Course not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch course", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Course fetched successfully", course))
}

// FindByProgramID godoc
// @Summary Get courses by program ID
// @Tags Courses
// @Produce json
// @Param program_id path string true "Program ID"
// @Success 200 {object} dto.APIResponse
// @Router /courses/program/{program_id} [get]
func (c *CourseController) FindByProgramID(ctx *gin.Context) {
	programID, err := uuid.Parse(ctx.Param("program_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid program ID", "INVALID_ID", nil))
		return
	}

	courses, err := c.service.FindByProgramID(programID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch courses", "FETCH_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Courses fetched successfully", courses))
}

// Update godoc
// @Summary Update course
// @Tags Courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param request body dto.UpdateCourseRequest true "Update Course Request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /courses/{id} [put]
func (c *CourseController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid course ID", "INVALID_ID", nil))
		return
	}

	var req dto.UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	if err := helper.ValidateStruct(&req); err != nil {
		errors := helper.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", "VALIDATION_ERROR", errors))
		return
	}

	course, err := c.service.Update(id, &req)
	if err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Course not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update course", "UPDATE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Course updated successfully", course))
}

// Delete godoc
// @Summary Delete course
// @Tags Courses
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /courses/{id} [delete]
func (c *CourseController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid course ID", "INVALID_ID", nil))
		return
	}

	if err := c.service.Delete(id); err != nil {
		if helper.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse("Course not found", "NOT_FOUND", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete course", "DELETE_ERROR", nil))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Course deleted successfully", nil))
}
