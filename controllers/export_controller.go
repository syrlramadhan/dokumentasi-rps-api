package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
)

type ExportController struct {
	exportService       services.ExportService
	generatedRPSService services.GeneratedRPSService
}

func NewExportController(
	exportService services.ExportService,
	generatedRPSService services.GeneratedRPSService,
) *ExportController {
	return &ExportController{
		exportService:       exportService,
		generatedRPSService: generatedRPSService,
	}
}

// ExportToPDF exports RPS to PDF format
// @Summary Export RPS to PDF
// @Description Export a generated RPS to PDF document
// @Tags Export
// @Produce application/pdf
// @Param id path string true "Generated RPS ID (UUID)"
// @Success 200 {file} binary "PDF file"
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/export/{id}/pdf [get]
func (ctrl *ExportController) ExportToPDF(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid ID format", "INVALID_ID", nil))
		return
	}

	// Get generated RPS
	generatedRPS, err := ctrl.generatedRPSService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Generated RPS not found", "NOT_FOUND", nil))
		return
	}

	// Check if RPS is completed
	if generatedRPS.Status != "done" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("RPS generation is not completed yet", "NOT_READY", nil))
		return
	}

	// Parse result to RPSStructuredOutput
	var rpsData dto.RPSStructuredOutput
	if generatedRPS.Result != nil {
		resultBytes, err := json.Marshal(generatedRPS.Result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to parse RPS data", "PARSE_ERROR", nil))
			return
		}
		if err := json.Unmarshal(resultBytes, &rpsData); err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to parse RPS data", "PARSE_ERROR", nil))
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("RPS result is empty", "EMPTY_RESULT", nil))
		return
	}

	// Generate PDF
	pdfBytes, err := ctrl.exportService.ExportToPDF(&rpsData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to generate PDF", "EXPORT_ERROR", nil))
		return
	}

	// Set headers and send file
	filename := fmt.Sprintf("RPS_%s_%s.pdf", rpsData.Identitas.KodeMataKuliah, rpsData.Identitas.NamaMataKuliah)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ExportToHTML exports RPS to HTML format (can be converted to DOCX via browser)
// @Summary Export RPS to HTML
// @Description Export a generated RPS to HTML document (can be saved as DOCX from browser)
// @Tags Export
// @Produce text/html
// @Param id path string true "Generated RPS ID (UUID)"
// @Success 200 {string} string "HTML content"
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/export/{id}/html [get]
func (ctrl *ExportController) ExportToHTML(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid ID format", "INVALID_ID", nil))
		return
	}

	// Get generated RPS
	generatedRPS, err := ctrl.generatedRPSService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Generated RPS not found", "NOT_FOUND", nil))
		return
	}

	// Check if RPS is completed
	if generatedRPS.Status != "done" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("RPS generation is not completed yet", "NOT_READY", nil))
		return
	}

	// Parse result to RPSStructuredOutput
	var rpsData dto.RPSStructuredOutput
	if generatedRPS.Result != nil {
		resultBytes, err := json.Marshal(generatedRPS.Result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to parse RPS data", "PARSE_ERROR", nil))
			return
		}
		if err := json.Unmarshal(resultBytes, &rpsData); err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to parse RPS data", "PARSE_ERROR", nil))
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("RPS result is empty", "EMPTY_RESULT", nil))
		return
	}

	// Generate HTML
	htmlContent, err := ctrl.exportService.ExportToHTML(&rpsData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to generate HTML", "EXPORT_ERROR", nil))
		return
	}

	// Set headers and send HTML
	filename := fmt.Sprintf("RPS_%s_%s.html", rpsData.Identitas.KodeMataKuliah, rpsData.Identitas.NamaMataKuliah)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.String(http.StatusOK, htmlContent)
}

// ExportToHTMLPreview exports RPS to HTML for preview (inline, not download)
// @Summary Preview RPS as HTML
// @Description Preview a generated RPS as HTML in browser
// @Tags Export
// @Produce text/html
// @Param id path string true "Generated RPS ID (UUID)"
// @Success 200 {string} string "HTML content"
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/export/{id}/preview [get]
func (ctrl *ExportController) ExportToHTMLPreview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid ID format", "INVALID_ID", nil))
		return
	}

	// Get generated RPS
	generatedRPS, err := ctrl.generatedRPSService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Generated RPS not found", "NOT_FOUND", nil))
		return
	}

	// Check if RPS is completed
	if generatedRPS.Status != "done" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("RPS generation is not completed yet", "NOT_READY", nil))
		return
	}

	// Parse result to RPSStructuredOutput
	var rpsData dto.RPSStructuredOutput
	if generatedRPS.Result != nil {
		resultBytes, err := json.Marshal(generatedRPS.Result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to parse RPS data", "PARSE_ERROR", nil))
			return
		}
		if err := json.Unmarshal(resultBytes, &rpsData); err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to parse RPS data", "PARSE_ERROR", nil))
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("RPS result is empty", "EMPTY_RESULT", nil))
		return
	}

	// Generate HTML
	htmlContent, err := ctrl.exportService.ExportToHTML(&rpsData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to generate HTML", "EXPORT_ERROR", nil))
		return
	}

	// Set headers for inline display (preview)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}

// GetExportFormats returns available export formats
// @Summary Get available export formats
// @Description Get list of available export formats for RPS
// @Tags Export
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/export/formats [get]
func (ctrl *ExportController) GetExportFormats(c *gin.Context) {
	formats := []map[string]interface{}{
		{
			"format":      "pdf",
			"name":        "PDF Document",
			"description": "Export RPS ke format PDF untuk cetak atau distribusi",
			"endpoint":    "/api/v1/export/{id}/pdf",
			"mime_type":   "application/pdf",
		},
		{
			"format":      "html",
			"name":        "HTML Document",
			"description": "Export RPS ke format HTML (dapat dibuka di browser dan disimpan sebagai DOCX)",
			"endpoint":    "/api/v1/export/{id}/html",
			"mime_type":   "text/html",
		},
		{
			"format":      "preview",
			"name":        "HTML Preview",
			"description": "Preview RPS di browser sebelum download",
			"endpoint":    "/api/v1/export/{id}/preview",
			"mime_type":   "text/html",
		},
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Export formats retrieved successfully", formats))
}
