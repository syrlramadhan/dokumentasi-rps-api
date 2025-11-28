package dto

// ==================== Gemini AI Request DTOs ====================

// GeminiRequest - Main request structure for Gemini API
type GeminiRequest struct {
	Contents          []GeminiContent  `json:"contents"`
	GenerationConfig  *GeminiGenConfig `json:"generationConfig,omitempty"`
	SafetySettings    []GeminiSafety   `json:"safetySettings,omitempty"`
	SystemInstruction *GeminiContent   `json:"systemInstruction,omitempty"`
}

// GeminiContent - Content structure (role + parts)
type GeminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart - Part of content (text, image, etc.)
type GeminiPart struct {
	Text string `json:"text,omitempty"`
}

// GeminiGenConfig - Generation configuration
type GeminiGenConfig struct {
	Temperature      float64                `json:"temperature,omitempty"`
	TopP             float64                `json:"topP,omitempty"`
	TopK             int                    `json:"topK,omitempty"`
	MaxOutputTokens  int                    `json:"maxOutputTokens,omitempty"`
	ResponseMimeType string                 `json:"responseMimeType,omitempty"`
	ResponseSchema   map[string]interface{} `json:"responseSchema,omitempty"`
}

// GeminiSafety - Safety settings
type GeminiSafety struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// ==================== Gemini AI Response DTOs ====================

// GeminiResponse - Main response structure from Gemini API
type GeminiResponse struct {
	Candidates    []GeminiCandidate `json:"candidates"`
	UsageMetadata GeminiUsage       `json:"usageMetadata"`
	ModelVersion  string            `json:"modelVersion,omitempty"`
}

// GeminiCandidate - Response candidate
type GeminiCandidate struct {
	Content       GeminiContent        `json:"content"`
	FinishReason  string               `json:"finishReason"`
	Index         int                  `json:"index"`
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings,omitempty"`
}

// GeminiSafetyRating - Safety rating for response
type GeminiSafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

// GeminiUsage - Token usage metadata
type GeminiUsage struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// GeminiError - Error response from Gemini
type GeminiError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

// RPS Structured Output Schema
type RPSStructuredOutput struct {
	Identitas           RPSIdentitas           `json:"identitas"`
	CapaianPembelajaran RPSCapaianPembelajaran `json:"capaian_pembelajaran"`
	DeskripsiMataKuliah RPSDeskripsi           `json:"deskripsi_mata_kuliah"`
	RencanaMingguan     []RPSRencanaMingguan   `json:"rencana_mingguan"`
	RencanaPenilaian    RPSPenilaian           `json:"rencana_penilaian"`
	DaftarReferensi     RPSReferensi           `json:"daftar_referensi"`
}

type RPSIdentitas struct {
	NamaMataKuliah string `json:"nama_mata_kuliah"`
	KodeMataKuliah string `json:"kode_mata_kuliah"`
	SKS            int    `json:"sks"`
	Semester       string `json:"semester"`
	Prasyarat      string `json:"prasyarat"`
	DosenPengampu  string `json:"dosen_pengampu"`
}

type RPSCapaianPembelajaran struct {
	CPLProdi []string `json:"cpl_prodi"`
	CPMK     []string `json:"cpmk"`
	SubCPMK  []string `json:"sub_cpmk"`
}

type RPSDeskripsi struct {
	DeskripsiSingkat string   `json:"deskripsi_singkat"`
	BahanKajian      []string `json:"bahan_kajian"`
}

type RPSRencanaMingguan struct {
	Minggu             int      `json:"minggu"`
	Topik              string   `json:"topik"`
	SubTopik           []string `json:"sub_topik"`
	IndikatorCapaian   string   `json:"indikator_capaian"`
	MetodePembelajaran string   `json:"metode_pembelajaran"`
	WaktuMenit         int      `json:"waktu_menit"`
	Referensi          string   `json:"referensi"`
	BentukPenilaian    string   `json:"bentuk_penilaian"`
}

type RPSPenilaian struct {
	Komponen []RPSKomponenPenilaian `json:"komponen"`
}

type RPSKomponenPenilaian struct {
	Nama      string `json:"nama"`
	Bobot     int    `json:"bobot"`
	Teknik    string `json:"teknik"`
	Instrumen string `json:"instrumen"`
}

type RPSReferensi struct {
	Utama     []string `json:"utama"`
	Pendukung []string `json:"pendukung"`
}

// AI Generation Result
type AIGenerationResult struct {
	Result     *RPSStructuredOutput   `json:"result"`
	AIMetadata map[string]interface{} `json:"ai_metadata"`
}

// Note: GenerateRPSOptions is defined in generated_rps_dto.go
