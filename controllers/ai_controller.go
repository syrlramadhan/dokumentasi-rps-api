package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
	"gorm.io/datatypes"
)

type AIController struct {
	aiService              services.AIService
	generatedRPSService    services.GeneratedRPSService
	templateVersionService services.TemplateVersionService
	courseService          services.CourseService
}

func NewAIController(
	aiService services.AIService,
	generatedRPSService services.GeneratedRPSService,
	templateVersionService services.TemplateVersionService,
	courseService services.CourseService,
) *AIController {
	return &AIController{
		aiService:              aiService,
		generatedRPSService:    generatedRPSService,
		templateVersionService: templateVersionService,
		courseService:          courseService,
	}
}

// GenerateRPSWithAI - Generate RPS menggunakan OpenAI (Synchronous)
// @Summary Generate RPS with AI (Sync)
// @Description Generate RPS using OpenAI structured output - waits for completion
// @Tags AI
// @Accept json
// @Produce json
// @Param request body dto.GenerateRPSRequest true "Generate RPS Request"
// @Success 200 {object} dto.APIResponse{data=dto.GeneratedRPSResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/generate/sync [post]
func (ctrl *AIController) GenerateRPSWithAI(c *gin.Context) {
	var req dto.GenerateRPSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	// Validate
	if req.TemplateVersionID == uuid.Nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("template_version_id is required", "VALIDATION_ERROR", nil))
		return
	}
	if req.CourseID == uuid.Nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("course_id is required", "VALIDATION_ERROR", nil))
		return
	}

	// Get template version
	templateVersion, err := ctrl.templateVersionService.FindByID(req.TemplateVersionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Template version not found", "NOT_FOUND", nil))
		return
	}

	// Get course
	course, err := ctrl.courseService.FindByID(req.CourseID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Course not found", "NOT_FOUND", nil))
		return
	}

	// Create initial generated_rps record
	createReq := &dto.CreateGeneratedRPSRequest{
		TemplateVersionID: &req.TemplateVersionID,
		CourseID:          &req.CourseID,
	}

	generatedRPS, err := ctrl.generatedRPSService.Create(createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create job", "CREATE_ERROR", nil))
		return
	}

	// Update status to processing
	ctrl.generatedRPSService.UpdateStatus(generatedRPS.ID, "processing")

	// Parse template definition
	var templateDef map[string]interface{}
	if templateVersion.Definition != nil {
		defBytes, _ := json.Marshal(templateVersion.Definition)
		json.Unmarshal(defBytes, &templateDef)
	}
	templateDef["id"] = templateVersion.ID.String()

	// Build course data
	courseData := map[string]interface{}{
		"id":      course.ID.String(),
		"title":   course.Title,
		"code":    course.Code,
		"credits": course.Credits,
	}

	// Build options
	options := dto.GenerateRPSOptions{
		Language: "Indonesia",
		Tone:     "formal",
	}
	if req.Options != nil {
		if req.Options.Language != "" {
			options.Language = req.Options.Language
		}
		if req.Options.Tone != "" {
			options.Tone = req.Options.Tone
		}
		options.Overrides = req.Options.Overrides
	}

	// Call AI service
	aiResult, err := ctrl.aiService.GenerateRPS(c.Request.Context(), generatedRPS.ID.String(), courseData, templateDef, options)
	if err != nil {
		ctrl.generatedRPSService.UpdateStatus(generatedRPS.ID, "failed")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("AI generation failed", "AI_ERROR", map[string]string{"error": err.Error()}))
		return
	}

	// Update generated_rps with result
	resultJSON, _ := json.Marshal(aiResult.Result)
	metadataJSON, _ := json.Marshal(aiResult.AIMetadata)

	status := "done"
	updateReq := &dto.UpdateGeneratedRPSRequest{
		Status:     &status,
		Result:     datatypes.JSON(resultJSON),
		AIMetadata: datatypes.JSON(metadataJSON),
	}

	updatedRPS, err := ctrl.generatedRPSService.Update(generatedRPS.ID, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update result", "UPDATE_ERROR", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("RPS generated successfully", updatedRPS))
}

// GenerateRPSAsync - Generate RPS secara async (return job_id langsung)
// @Summary Generate RPS Async
// @Description Start async RPS generation, returns job_id immediately
// @Tags AI
// @Accept json
// @Produce json
// @Param request body dto.GenerateRPSRequest true "Generate RPS Request"
// @Success 202 {object} dto.APIResponse
// @Router /api/v1/generate [post]
func (ctrl *AIController) GenerateRPSAsync(c *gin.Context) {
	var req dto.GenerateRPSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", "INVALID_REQUEST", nil))
		return
	}

	// Validate
	if req.TemplateVersionID == uuid.Nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("template_version_id is required", "VALIDATION_ERROR", nil))
		return
	}
	if req.CourseID == uuid.Nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("course_id is required", "VALIDATION_ERROR", nil))
		return
	}

	// Create job with status "queued"
	createReq := &dto.CreateGeneratedRPSRequest{
		TemplateVersionID: &req.TemplateVersionID,
		CourseID:          &req.CourseID,
	}

	generatedRPS, err := ctrl.generatedRPSService.Create(createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create job", "CREATE_ERROR", nil))
		return
	}

	// Start async generation in goroutine
	go ctrl.processGenerationAsync(generatedRPS.ID, req)

	c.JSON(http.StatusAccepted, dto.SuccessResponse("RPS generation started", gin.H{
		"job_id": generatedRPS.ID,
		"status": "queued",
	}))
}

func (ctrl *AIController) processGenerationAsync(jobID uuid.UUID, req dto.GenerateRPSRequest) {
	// Update to processing
	ctrl.generatedRPSService.UpdateStatus(jobID, "processing")

	// Get template version
	templateVersion, err := ctrl.templateVersionService.FindByID(req.TemplateVersionID)
	if err != nil {
		ctrl.markAsFailed(jobID, "Template version not found")
		return
	}

	// Get course
	course, err := ctrl.courseService.FindByID(req.CourseID)
	if err != nil {
		ctrl.markAsFailed(jobID, "Course not found")
		return
	}

	// Parse template definition
	var templateDef map[string]interface{}
	if templateVersion.Definition != nil {
		defBytes, _ := json.Marshal(templateVersion.Definition)
		json.Unmarshal(defBytes, &templateDef)
	}
	templateDef["id"] = templateVersion.ID.String()

	// Build course data
	courseData := map[string]interface{}{
		"id":      course.ID.String(),
		"title":   course.Title,
		"code":    course.Code,
		"credits": course.Credits,
	}

	// Build options
	options := dto.GenerateRPSOptions{
		Language: "Indonesia",
		Tone:     "formal",
	}
	if req.Options != nil {
		if req.Options.Language != "" {
			options.Language = req.Options.Language
		}
		if req.Options.Tone != "" {
			options.Tone = req.Options.Tone
		}
		options.Overrides = req.Options.Overrides
	}

	// Call AI (use background context since HTTP request already returned)
	aiResult, err := ctrl.aiService.GenerateRPS(context.Background(), jobID.String(), courseData, templateDef, options)
	if err != nil {
		ctrl.markAsFailed(jobID, err.Error())
		return
	}

	// Update with result
	resultJSON, _ := json.Marshal(aiResult.Result)
	metadataJSON, _ := json.Marshal(aiResult.AIMetadata)

	status := "done"
	updateReq := &dto.UpdateGeneratedRPSRequest{
		Status:     &status,
		Result:     datatypes.JSON(resultJSON),
		AIMetadata: datatypes.JSON(metadataJSON),
	}

	ctrl.generatedRPSService.Update(jobID, updateReq)
}

func (ctrl *AIController) markAsFailed(jobID uuid.UUID, errorMsg string) {
	metadataJSON, _ := json.Marshal(map[string]string{"error": errorMsg})
	status := "failed"
	updateReq := &dto.UpdateGeneratedRPSRequest{
		Status:     &status,
		AIMetadata: datatypes.JSON(metadataJSON),
	}
	ctrl.generatedRPSService.Update(jobID, updateReq)
}

// GetPromptByID - Get AI prompt by MongoDB ID
// @Summary Get AI prompt by ID
// @Description Get detailed AI prompt information from MongoDB
// @Tags AI Admin
// @Produce json
// @Param id path string true "Prompt ID (MongoDB ObjectID)"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/admin/ai/prompts/{id} [get]
func (ctrl *AIController) GetPromptByID(c *gin.Context) {
	id := c.Param("id")

	prompt, err := ctrl.aiService.GetPromptByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Prompt not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Prompt retrieved successfully", prompt))
}

// GetPromptsByGeneratedRPSID - Get all prompts for a generated RPS
// @Summary Get prompts by Generated RPS ID
// @Description Get all AI prompts associated with a generated RPS
// @Tags AI Admin
// @Produce json
// @Param generated_rps_id path string true "Generated RPS ID (UUID)"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/admin/ai/prompts/rps/{generated_rps_id} [get]
func (ctrl *AIController) GetPromptsByGeneratedRPSID(c *gin.Context) {
	generatedRPSID := c.Param("generated_rps_id")

	prompts, err := ctrl.aiService.GetPromptsByGeneratedRPSID(c.Request.Context(), generatedRPSID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to get prompts", "INTERNAL_ERROR", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Prompts retrieved successfully", prompts))
}

// GetGenerationByRPSID - Get AI generation record
// @Summary Get AI generation by RPS ID
// @Description Get AI generation record from MongoDB by Generated RPS ID
// @Tags AI Admin
// @Produce json
// @Param generated_rps_id path string true "Generated RPS ID (UUID)"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/admin/ai/generations/{generated_rps_id} [get]
func (ctrl *AIController) GetGenerationByRPSID(c *gin.Context) {
	generatedRPSID := c.Param("generated_rps_id")

	generation, err := ctrl.aiService.GetGenerationByRPSID(c.Request.Context(), generatedRPSID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Generation not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Generation retrieved successfully", generation))
}

// GetPromptStats - Get AI prompt statistics
// @Summary Get AI prompt statistics
// @Description Get aggregated statistics of AI prompts
// @Tags AI Admin
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/admin/ai/stats [get]
func (ctrl *AIController) GetPromptStats(c *gin.Context) {
	stats, err := ctrl.aiService.GetPromptStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to get stats", "INTERNAL_ERROR", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Stats retrieved successfully", stats))
}

// GetAllPrompts - Get all AI prompts with pagination
// @Summary Get all AI prompts
// @Description Get all AI prompts with pagination
// @Tags AI Admin
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/admin/ai/prompts [get]
func (ctrl *AIController) GetAllPrompts(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	prompts, err := ctrl.aiService.GetAllPrompts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to get prompts", "INTERNAL_ERROR", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Prompts retrieved successfully", prompts))
}

// GetAllGenerations - Get all AI generations with pagination
// @Summary Get all AI generations
// @Description Get all AI generations with pagination
// @Tags AI Admin
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/admin/ai/generations [get]
func (ctrl *AIController) GetAllGenerations(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	generations, err := ctrl.aiService.GetAllGenerations(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to get generations", "INTERNAL_ERROR", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Generations retrieved successfully", generations))
}
