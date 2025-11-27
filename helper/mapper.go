package helper

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	"github.com/syrlramadhan/dokumentasi-rps-api/models"
)

// User Mapper
func ToUserResponse(user *models.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}
}

func ToUserResponseList(users []models.User) []dto.UserResponse {
	result := make([]dto.UserResponse, len(users))
	for i, user := range users {
		result[i] = *ToUserResponse(&user)
	}
	return result
}

func ToUserModel(req *dto.CreateUserRequest) *models.User {
	return &models.User{
		ID:          uuid.New(),
		Username:    req.Username,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Role:        req.Role,
	}
}

// Program Mapper
func ToProgramResponse(program *models.Program) *dto.ProgramResponse {
	if program == nil {
		return nil
	}
	return &dto.ProgramResponse{
		ID:   program.ID,
		Code: program.Code,
		Name: program.Name,
	}
}

func ToProgramResponseList(programs []models.Program) []dto.ProgramResponse {
	result := make([]dto.ProgramResponse, len(programs))
	for i, program := range programs {
		result[i] = *ToProgramResponse(&program)
	}
	return result
}

func ToProgramModel(req *dto.CreateProgramRequest) *models.Program {
	return &models.Program{
		ID:   uuid.New(),
		Code: req.Code,
		Name: req.Name,
	}
}

// Course Mapper
func ToCourseResponse(course *models.Course) *dto.CourseResponse {
	if course == nil {
		return nil
	}
	return &dto.CourseResponse{
		ID:        course.ID,
		ProgramID: course.ProgramID,
		Code:      course.Code,
		Title:     course.Title,
		Credits:   course.Credits,
		CreatedAt: course.CreatedAt,
		Program:   ToProgramResponse(course.Program),
	}
}

func ToCourseResponseList(courses []models.Course) []dto.CourseResponse {
	result := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		result[i] = *ToCourseResponse(&course)
	}
	return result
}

func ToCourseModel(req *dto.CreateCourseRequest) *models.Course {
	return &models.Course{
		ID:        uuid.New(),
		ProgramID: req.ProgramID,
		Code:      req.Code,
		Title:     req.Title,
		Credits:   req.Credits,
	}
}

// Template Mapper
func ToTemplateResponse(template *models.Template) *dto.TemplateResponse {
	if template == nil {
		return nil
	}
	return &dto.TemplateResponse{
		ID:          template.ID,
		ProgramID:   template.ProgramID,
		Name:        template.Name,
		Description: template.Description,
		CreatedBy:   template.CreatedBy,
		CreatedAt:   template.CreatedAt,
		IsActive:    template.IsActive,
		Program:     ToProgramResponse(template.Program),
		Creator:     ToUserResponse(template.Creator),
	}
}

func ToTemplateResponseList(templates []models.Template) []dto.TemplateResponse {
	result := make([]dto.TemplateResponse, len(templates))
	for i, template := range templates {
		result[i] = *ToTemplateResponse(&template)
	}
	return result
}

func ToTemplateModel(req *dto.CreateTemplateRequest) *models.Template {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	return &models.Template{
		ID:          uuid.New(),
		ProgramID:   req.ProgramID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
		IsActive:    isActive,
	}
}

// TemplateVersion Mapper
func ToTemplateVersionResponse(version *models.TemplateVersion) *dto.TemplateVersionResponse {
	if version == nil {
		return nil
	}
	return &dto.TemplateVersionResponse{
		ID:         version.ID,
		TemplateID: version.TemplateID,
		Version:    version.Version,
		Definition: version.Definition,
		CreatedBy:  version.CreatedBy,
		CreatedAt:  version.CreatedAt,
		Template:   ToTemplateResponse(version.Template),
		Creator:    ToUserResponse(version.Creator),
	}
}

func ToTemplateVersionResponseList(versions []models.TemplateVersion) []dto.TemplateVersionResponse {
	result := make([]dto.TemplateVersionResponse, len(versions))
	for i, version := range versions {
		result[i] = *ToTemplateVersionResponse(&version)
	}
	return result
}

func ToTemplateVersionModel(req *dto.CreateTemplateVersionRequest) *models.TemplateVersion {
	return &models.TemplateVersion{
		ID:         uuid.New(),
		TemplateID: req.TemplateID,
		Version:    req.Version,
		Definition: req.Definition,
		CreatedBy:  req.CreatedBy,
	}
}

// GeneratedRPS Mapper
func ToGeneratedRPSResponse(rps *models.GeneratedRPS) *dto.GeneratedRPSResponse {
	if rps == nil {
		return nil
	}
	return &dto.GeneratedRPSResponse{
		ID:                rps.ID,
		TemplateVersionID: rps.TemplateVersionID,
		CourseID:          rps.CourseID,
		GeneratedBy:       rps.GeneratedBy,
		Status:            rps.Status,
		Result:            rps.Result,
		ExportedFileURL:   rps.ExportedFileURL,
		AIMetadata:        rps.AIMetadata,
		CreatedAt:         rps.CreatedAt,
		UpdatedAt:         rps.UpdatedAt,
		TemplateVersion:   ToTemplateVersionResponse(rps.TemplateVersion),
		Course:            ToCourseResponse(rps.Course),
		Generator:         ToUserResponse(rps.Generator),
	}
}

func ToGeneratedRPSResponseList(rpsList []models.GeneratedRPS) []dto.GeneratedRPSResponse {
	result := make([]dto.GeneratedRPSResponse, len(rpsList))
	for i, rps := range rpsList {
		result[i] = *ToGeneratedRPSResponse(&rps)
	}
	return result
}

func ToGeneratedRPSModel(req *dto.CreateGeneratedRPSRequest) *models.GeneratedRPS {
	return &models.GeneratedRPS{
		ID:                uuid.New(),
		TemplateVersionID: req.TemplateVersionID,
		CourseID:          req.CourseID,
		GeneratedBy:       req.GeneratedBy,
		Status:            "queued",
	}
}

// AuditLog Mapper
func ToAuditLogResponse(log *models.AuditLog) *dto.AuditLogResponse {
	if log == nil {
		return nil
	}
	return &dto.AuditLogResponse{
		ID:         log.ID,
		UserID:     log.UserID,
		Action:     log.Action,
		TargetType: log.TargetType,
		TargetID:   log.TargetID,
		Payload:    log.Payload,
		CreatedAt:  log.CreatedAt,
		User:       ToUserResponse(log.User),
	}
}

func ToAuditLogResponseList(logs []models.AuditLog) []dto.AuditLogResponse {
	result := make([]dto.AuditLogResponse, len(logs))
	for i, log := range logs {
		result[i] = *ToAuditLogResponse(&log)
	}
	return result
}

func ToAuditLogModel(req *dto.CreateAuditLogRequest) *models.AuditLog {
	return &models.AuditLog{
		UserID:     req.UserID,
		Action:     req.Action,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Payload:    req.Payload,
	}
}
