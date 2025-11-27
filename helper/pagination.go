package helper

import "math"

// CalculateTotalPages calculates total pages based on total items and limit
func CalculateTotalPages(totalItems int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int(math.Ceil(float64(totalItems) / float64(limit)))
}

// CalculateOffset calculates offset for pagination
func CalculateOffset(page, limit int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * limit
}

// GetDefaultPagination returns default pagination values
func GetDefaultPagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
