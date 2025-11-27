package dto

// API Response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// Pagination
type PaginationRequest struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
	Data       interface{}        `json:"data,omitempty"`
	Pagination PaginationResponse `json:"pagination"`
}

// Helper functions
func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, code string, details map[string]string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    code,
			Details: details,
		},
	}
}

func PaginatedSuccessResponse(message string, data interface{}, pagination PaginationResponse) PaginatedResponse {
	return PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}
