package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
	models "github.com/syrlramadhan/dokumentasi-rps-api/models/mongo"
	mongoRepo "github.com/syrlramadhan/dokumentasi-rps-api/repositories/mongo"
)

type AIService interface {
	GenerateRPS(ctx context.Context, generatedRPSID string, courseData map[string]interface{}, templateDef map[string]interface{}, options dto.GenerateRPSOptions) (*dto.AIGenerationResult, error)
	GetPromptByID(ctx context.Context, id string) (*models.AIPrompt, error)
	GetPromptsByGeneratedRPSID(ctx context.Context, generatedRPSID string) ([]models.AIPrompt, error)
	GetGenerationByRPSID(ctx context.Context, generatedRPSID string) (*models.AIGeneration, error)
	GetPromptStats(ctx context.Context) (*mongoRepo.AIPromptStats, error)
	GetAllPrompts(ctx context.Context, limit, offset int64) ([]models.AIPromptSummary, error)
	GetAllGenerations(ctx context.Context, limit, offset int64) ([]models.AIGeneration, error)
}

type aiService struct {
	apiKey             string
	model              string
	httpClient         *http.Client
	aiPromptRepo       mongoRepo.AIPromptRepository
	aiGenerationRepo   mongoRepo.AIGenerationRepository
	promptTemplateRepo mongoRepo.PromptTemplateRepository
}

func NewAIService(
	aiPromptRepo mongoRepo.AIPromptRepository,
	aiGenerationRepo mongoRepo.AIGenerationRepository,
	promptTemplateRepo mongoRepo.PromptTemplateRepository,
) AIService {
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-2.0-flash"
	}

	return &aiService{
		apiKey:             os.Getenv("GEMINI_API_KEY"),
		model:              model,
		httpClient:         &http.Client{Timeout: 120 * time.Second},
		aiPromptRepo:       aiPromptRepo,
		aiGenerationRepo:   aiGenerationRepo,
		promptTemplateRepo: promptTemplateRepo,
	}
}

// getGeminiAPIURL builds the Gemini API URL
func (s *aiService) getGeminiAPIURL() string {
	return fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", s.model, s.apiKey)
}

func (s *aiService) GenerateRPS(ctx context.Context, generatedRPSID string, courseData map[string]interface{}, templateDef map[string]interface{}, options dto.GenerateRPSOptions) (*dto.AIGenerationResult, error) {
	startTime := time.Now()

	// Check API Key
	if s.apiKey == "" {
		log.Println("❌ ERROR: GEMINI_API_KEY is empty!")
		return nil, fmt.Errorf("GEMINI_API_KEY is not set. Please set it in .env file")
	}
	log.Printf("✅ Gemini API Key detected (length: %d)", len(s.apiKey))
	log.Printf("✅ Using model: %s", s.model)

	// Set defaults
	if options.Language == "" {
		options.Language = "Indonesia"
	}
	if options.Tone == "" {
		options.Tone = "formal dan akademis"
	}

	// Create AI Generation record in MongoDB
	generation := &models.AIGeneration{
		GeneratedRPSID:    generatedRPSID,
		CourseID:          fmt.Sprintf("%v", courseData["id"]),
		CourseName:        fmt.Sprintf("%v", courseData["title"]),
		CourseCode:        fmt.Sprintf("%v", courseData["code"]),
		TemplateVersionID: fmt.Sprintf("%v", templateDef["id"]),
		Attempts:          []models.GenerationAttempt{},
		TotalAttempts:     0,
		FinalStatus:       "processing",
	}

	generation, err := s.aiGenerationRepo.Create(ctx, generation)
	if err != nil {
		log.Printf("Warning: failed to create AI generation record: %v", err)
		generation = &models.AIGeneration{ID: primitive.NewObjectID()}
	}

	// Build prompts
	systemPrompt := s.buildSystemPrompt()
	userPrompt := s.buildUserPrompt(courseData, templateDef, options)

	log.Printf("📝 System prompt length: %d chars", len(systemPrompt))
	log.Printf("📝 User prompt length: %d chars", len(userPrompt))

	// Create AI Prompt record
	aiPrompt := &models.AIPrompt{
		GeneratedRPSID: generatedRPSID,
		CourseID:       fmt.Sprintf("%v", courseData["id"]),
		TemplateID:     fmt.Sprintf("%v", templateDef["id"]),
		SystemPrompt:   systemPrompt,
		UserPrompt:     userPrompt,
		FullPrompt:     fmt.Sprintf("System: %s\n\nUser: %s", systemPrompt, userPrompt),
		Model:          s.model,
		Temperature:    0.7,
		MaxTokens:      8192,
		ResponseFormat: "json_object",
		CourseData:     courseData,
		TemplateData:   templateDef,
		Options:        map[string]interface{}{"language": options.Language, "tone": options.Tone, "overrides": options.Overrides},
		Status:         "pending",
	}

	// Build Gemini request with structured output
	reqBody := dto.GeminiRequest{
		SystemInstruction: &dto.GeminiContent{
			Parts: []dto.GeminiPart{{Text: systemPrompt}},
		},
		Contents: []dto.GeminiContent{
			{
				Role:  "user",
				Parts: []dto.GeminiPart{{Text: userPrompt}},
			},
		},
		GenerationConfig: &dto.GeminiGenConfig{
			Temperature:      0.7,
			TopP:             0.95,
			TopK:             40,
			MaxOutputTokens:  8192,
			ResponseMimeType: "application/json",
			ResponseSchema:   s.GetRPSJSONSchema(),
		},
		SafetySettings: []dto.GeminiSafety{
			{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "BLOCK_NONE"},
			{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "BLOCK_NONE"},
			{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "BLOCK_NONE"},
			{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "BLOCK_NONE"},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("❌ Failed to marshal request: %v", err)
		s.recordFailedAttempt(ctx, generation.ID, aiPrompt, "failed to marshal request", 0, time.Since(startTime).Milliseconds())
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("📤 Sending request to Gemini API...")

	// Make HTTP request
	apiURL := s.getGeminiAPIURL()
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("❌ Failed to create request: %v", err)
		s.recordFailedAttempt(ctx, generation.ID, aiPrompt, "failed to create request", 0, time.Since(startTime).Milliseconds())
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("❌ Failed to call Gemini API: %v", err)
		s.recordFailedAttempt(ctx, generation.ID, aiPrompt, err.Error(), 0, time.Since(startTime).Milliseconds())
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Failed to read response: %v", err)
		s.recordFailedAttempt(ctx, generation.ID, aiPrompt, "failed to read response", 0, time.Since(startTime).Milliseconds())
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("📥 Gemini Response Status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Gemini API Error Response: %s", string(body))
		var geminiErr dto.GeminiError
		json.Unmarshal(body, &geminiErr)
		errMsg := fmt.Sprintf("Gemini API error (status %d): %s", resp.StatusCode, geminiErr.Error.Message)
		s.recordFailedAttempt(ctx, generation.ID, aiPrompt, errMsg, 0, time.Since(startTime).Milliseconds())
		return nil, fmt.Errorf(errMsg)
	}

	// Parse Gemini response
	var geminiResp dto.GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		log.Printf("❌ Failed to parse Gemini response: %v", err)
		log.Printf("Raw response: %s", string(body))
		s.recordFailedAttempt(ctx, generation.ID, aiPrompt, "failed to parse Gemini response", 0, time.Since(startTime).Milliseconds())
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		log.Printf("❌ No candidates in Gemini response")
		s.recordFailedAttempt(ctx, generation.ID, aiPrompt, "no candidates in response", 0, time.Since(startTime).Milliseconds())
		return nil, fmt.Errorf("no candidates in Gemini response")
	}

	// Get response text
	var responseContent string
	if len(geminiResp.Candidates[0].Content.Parts) > 0 {
		responseContent = geminiResp.Candidates[0].Content.Parts[0].Text
	}

	if responseContent == "" {
		log.Printf("❌ Empty response content from Gemini")
		s.recordFailedAttempt(ctx, generation.ID, aiPrompt, "empty response content", 0, time.Since(startTime).Milliseconds())
		return nil, fmt.Errorf("empty response content from Gemini")
	}

	log.Printf("📄 Response content length: %d chars", len(responseContent))

	// Parse structured output
	var rpsResult dto.RPSStructuredOutput
	if err := json.Unmarshal([]byte(responseContent), &rpsResult); err != nil {
		log.Printf("❌ Failed to parse RPS output: %v", err)
		log.Printf("Response content: %s", responseContent[:min(500, len(responseContent))])
		s.recordFailedAttempt(ctx, generation.ID, aiPrompt, "failed to parse RPS output", geminiResp.UsageMetadata.TotalTokenCount, time.Since(startTime).Milliseconds())
		return nil, fmt.Errorf("failed to parse RPS structured output: %w", err)
	}

	requestDuration := time.Since(startTime).Milliseconds()

	log.Printf("✅ Generation successful!")
	log.Printf("📊 Tokens used: %d (prompt: %d, completion: %d)",
		geminiResp.UsageMetadata.TotalTokenCount,
		geminiResp.UsageMetadata.PromptTokenCount,
		geminiResp.UsageMetadata.CandidatesTokenCount)
	log.Printf("⏱️ Duration: %dms", requestDuration)

	// Update AI Prompt with success data
	aiPrompt.Response = responseContent
	aiPrompt.ParsedResponse = map[string]interface{}{"rps": rpsResult}
	aiPrompt.PromptTokens = geminiResp.UsageMetadata.PromptTokenCount
	aiPrompt.CompletionTokens = geminiResp.UsageMetadata.CandidatesTokenCount
	aiPrompt.TotalTokens = geminiResp.UsageMetadata.TotalTokenCount
	aiPrompt.RequestDurationMs = requestDuration
	aiPrompt.Status = "success"
	aiPrompt.FinishReason = geminiResp.Candidates[0].FinishReason

	// Save prompt to MongoDB
	savedPrompt, err := s.aiPromptRepo.Create(ctx, aiPrompt)
	if err != nil {
		log.Printf("Warning: failed to save AI prompt to MongoDB: %v", err)
	}

	// Update generation record with success
	attempt := models.GenerationAttempt{
		AttemptNumber: 1,
		Status:        "success",
		TokensUsed:    geminiResp.UsageMetadata.TotalTokenCount,
		DurationMs:    requestDuration,
		Timestamp:     time.Now(),
	}
	if savedPrompt != nil {
		attempt.PromptID = savedPrompt.ID
	}

	s.aiGenerationRepo.AddAttempt(ctx, generation.ID, attempt)

	// Convert result to map for storage
	resultMap := make(map[string]interface{})
	resultBytes, _ := json.Marshal(rpsResult)
	json.Unmarshal(resultBytes, &resultMap)

	s.aiGenerationRepo.UpdateFinalStatus(ctx, generation.ID, "success", resultMap)

	// Build AI metadata
	aiMetadata := map[string]interface{}{
		"model":               s.model,
		"prompt_tokens":       geminiResp.UsageMetadata.PromptTokenCount,
		"completion_tokens":   geminiResp.UsageMetadata.CandidatesTokenCount,
		"total_tokens":        geminiResp.UsageMetadata.TotalTokenCount,
		"temperature":         0.7,
		"generation_time_ms":  requestDuration,
		"finish_reason":       geminiResp.Candidates[0].FinishReason,
		"response_format":     "structured_output",
		"provider":            "google_gemini",
		"mongo_prompt_id":     "",
		"mongo_generation_id": generation.ID.Hex(),
	}

	if savedPrompt != nil {
		aiMetadata["mongo_prompt_id"] = savedPrompt.ID.Hex()
	}

	return &dto.AIGenerationResult{
		Result:     &rpsResult,
		AIMetadata: aiMetadata,
	}, nil
}

func (s *aiService) recordFailedAttempt(ctx context.Context, generationID primitive.ObjectID, prompt *models.AIPrompt, errorMsg string, tokens int, duration int64) {
	prompt.Status = "failed"
	prompt.ErrorMessage = errorMsg
	prompt.TotalTokens = tokens
	prompt.RequestDurationMs = duration

	savedPrompt, _ := s.aiPromptRepo.Create(ctx, prompt)

	attempt := models.GenerationAttempt{
		AttemptNumber: 1,
		Status:        "failed",
		TokensUsed:    tokens,
		DurationMs:    duration,
		ErrorMessage:  errorMsg,
		Timestamp:     time.Now(),
	}
	if savedPrompt != nil {
		attempt.PromptID = savedPrompt.ID
	}

	s.aiGenerationRepo.AddAttempt(ctx, generationID, attempt)
	s.aiGenerationRepo.UpdateFinalStatus(ctx, generationID, "failed", nil)
}

func (s *aiService) GetPromptByID(ctx context.Context, id string) (*models.AIPrompt, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid prompt ID: %w", err)
	}
	return s.aiPromptRepo.FindByID(ctx, objectID)
}

func (s *aiService) GetPromptsByGeneratedRPSID(ctx context.Context, generatedRPSID string) ([]models.AIPrompt, error) {
	return s.aiPromptRepo.FindByGeneratedRPSID(ctx, generatedRPSID)
}

func (s *aiService) GetGenerationByRPSID(ctx context.Context, generatedRPSID string) (*models.AIGeneration, error) {
	return s.aiGenerationRepo.FindByGeneratedRPSID(ctx, generatedRPSID)
}

func (s *aiService) GetPromptStats(ctx context.Context) (*mongoRepo.AIPromptStats, error) {
	return s.aiPromptRepo.GetStats(ctx)
}

func (s *aiService) GetAllPrompts(ctx context.Context, limit, offset int64) ([]models.AIPromptSummary, error) {
	return s.aiPromptRepo.FindAll(ctx, limit, offset)
}

func (s *aiService) GetAllGenerations(ctx context.Context, limit, offset int64) ([]models.AIGeneration, error) {
	return s.aiGenerationRepo.FindAll(ctx, limit, offset)
}

// GetRPSJSONSchema returns the JSON Schema for Gemini structured output
// Note: Gemini doesn't support "additionalProperties" and uses different format
func (s *aiService) GetRPSJSONSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"identitas": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"nama_mata_kuliah": map[string]interface{}{"type": "string"},
					"kode_mata_kuliah": map[string]interface{}{"type": "string"},
					"sks":              map[string]interface{}{"type": "integer"},
					"semester":         map[string]interface{}{"type": "string"},
					"prasyarat":        map[string]interface{}{"type": "string"},
					"dosen_pengampu":   map[string]interface{}{"type": "string"},
				},
				"required": []string{"nama_mata_kuliah", "kode_mata_kuliah", "sks", "semester", "prasyarat", "dosen_pengampu"},
			},
			"capaian_pembelajaran": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cpl_prodi": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					"cpmk":      map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					"sub_cpmk":  map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
				},
				"required": []string{"cpl_prodi", "cpmk", "sub_cpmk"},
			},
			"deskripsi_mata_kuliah": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"deskripsi_singkat": map[string]interface{}{"type": "string"},
					"bahan_kajian":      map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
				},
				"required": []string{"deskripsi_singkat", "bahan_kajian"},
			},
			"rencana_mingguan": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"minggu":              map[string]interface{}{"type": "integer"},
						"topik":               map[string]interface{}{"type": "string"},
						"sub_topik":           map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
						"indikator_capaian":   map[string]interface{}{"type": "string"},
						"metode_pembelajaran": map[string]interface{}{"type": "string"},
						"waktu_menit":         map[string]interface{}{"type": "integer"},
						"referensi":           map[string]interface{}{"type": "string"},
						"bentuk_penilaian":    map[string]interface{}{"type": "string"},
					},
					"required": []string{"minggu", "topik", "sub_topik", "indikator_capaian", "metode_pembelajaran", "waktu_menit", "referensi", "bentuk_penilaian"},
				},
			},
			"rencana_penilaian": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"komponen": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"nama":      map[string]interface{}{"type": "string"},
								"bobot":     map[string]interface{}{"type": "integer"},
								"teknik":    map[string]interface{}{"type": "string"},
								"instrumen": map[string]interface{}{"type": "string"},
							},
							"required": []string{"nama", "bobot", "teknik", "instrumen"},
						},
					},
				},
				"required": []string{"komponen"},
			},
			"daftar_referensi": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"utama":     map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					"pendukung": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
				},
				"required": []string{"utama", "pendukung"},
			},
		},
		"required": []string{"identitas", "capaian_pembelajaran", "deskripsi_mata_kuliah", "rencana_mingguan", "rencana_penilaian", "daftar_referensi"},
	}
}

func (s *aiService) buildSystemPrompt() string {
	return `Anda adalah asisten akademik ahli dalam menyusun Rencana Pembelajaran Semester (RPS) untuk perguruan tinggi di Indonesia.

Tugas Anda:
1. Membuat RPS yang lengkap, terstruktur, dan sesuai standar SNPT (Standar Nasional Pendidikan Tinggi)
2. Menggunakan bahasa Indonesia yang formal dan akademis
3. Memastikan setiap komponen RPS terisi dengan konten yang relevan dan berkualitas
4. Membuat rencana pembelajaran mingguan yang detail untuk 16 minggu (termasuk UTS di minggu 8 dan UAS di minggu 16)

Panduan penyusunan:
- CPL (Capaian Pembelajaran Lulusan) harus sesuai dengan profil lulusan program studi
- CPMK (Capaian Pembelajaran Mata Kuliah) harus mendukung CPL
- Sub-CPMK harus terukur dan dapat dicapai dalam satu atau beberapa pertemuan
- Metode pembelajaran harus bervariasi dan sesuai dengan karakteristik materi
- Penilaian harus mencakup aspek kognitif, afektif, dan psikomotorik`
}

func (s *aiService) buildUserPrompt(courseData map[string]interface{}, templateDef map[string]interface{}, options dto.GenerateRPSOptions) string {
	templateJSON, _ := json.MarshalIndent(templateDef, "", "  ")

	// Default values
	semester := "Ganjil 2024/2025"
	dosenPengampu := "Tim Dosen"
	prasyarat := "-"
	programStudi := ""
	fakultas := ""

	// Override with options
	if options.Semester != "" {
		semester = options.Semester
	}
	if options.DosenPengampu != "" {
		dosenPengampu = options.DosenPengampu
	}
	if options.Prasyarat != "" {
		prasyarat = options.Prasyarat
	}
	if options.ProgramStudi != "" {
		programStudi = options.ProgramStudi
	}
	if options.Fakultas != "" {
		fakultas = options.Fakultas
	}

	// Legacy overrides support
	if options.Overrides != nil {
		if sem, ok := options.Overrides["semester"].(string); ok && sem != "" {
			semester = sem
		}
		if dosen, ok := options.Overrides["dosen_pengampu"].(string); ok && dosen != "" {
			dosenPengampu = dosen
		}
		if prasy, ok := options.Overrides["prasyarat"].(string); ok && prasy != "" {
			prasyarat = prasy
		}
	}

	// Build program info section
	programInfo := ""
	if programStudi != "" || fakultas != "" {
		programInfo = "\n\n## INFORMASI PROGRAM STUDI"
		if programStudi != "" {
			programInfo += fmt.Sprintf("\n- Program Studi: %s", programStudi)
		}
		if fakultas != "" {
			programInfo += fmt.Sprintf("\n- Fakultas: %s", fakultas)
		}
	}

	return fmt.Sprintf(`Buatkan Rencana Pembelajaran Semester (RPS) untuk mata kuliah berikut:

## INFORMASI MATA KULIAH
- Nama Mata Kuliah: %v
- Kode Mata Kuliah: %v
- Jumlah SKS: %v
- Semester: %s
- Dosen Pengampu: %s
- Prasyarat: %s%s

## TEMPLATE STRUKTUR (untuk referensi)
%s

## INSTRUKSI KHUSUS
- Bahasa: %s
- Gaya penulisan: %s
- Jumlah pertemuan: 16 minggu (UTS minggu ke-8, UAS minggu ke-16)
- Waktu per pertemuan: 150 menit (3 SKS) atau sesuaikan dengan SKS
- PENTING: Gunakan nama dosen "%s" untuk field dosen_pengampu
- PENTING: Gunakan "%s" untuk field prasyarat

## KETENTUAN PENILAIAN
- Total bobot harus = 100%%
- Komponen minimal: Tugas, Kuis, UTS, UAS
- Bisa ditambah: Praktikum, Proyek, Presentasi

Buatkan RPS yang lengkap dan berkualitas.`,
		courseData["title"],
		courseData["code"],
		courseData["credits"],
		semester,
		dosenPengampu,
		prasyarat,
		programInfo,
		string(templateJSON),
		options.Language,
		options.Tone,
		dosenPengampu,
		prasyarat,
	)
}
