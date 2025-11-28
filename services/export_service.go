package services

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/syrlramadhan/dokumentasi-rps-api/dto"
)

type ExportService interface {
	ExportToPDF(rps *dto.RPSStructuredOutput) ([]byte, error)
	ExportToHTML(rps *dto.RPSStructuredOutput) (string, error)
}

type exportService struct{}

func NewExportService() ExportService {
	return &exportService{}
}

// ExportToPDF generates a PDF document from RPS data
func (s *exportService) ExportToPDF(rps *dto.RPSStructuredOutput) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape for better table display
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)

	// Add first page
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 12, "RENCANA PEMBELAJARAN SEMESTER (RPS)", "", 1, "C", false, 0, "")
	pdf.Ln(8)

	// ==================== IDENTITAS ====================
	s.addSectionTitle(pdf, "I. IDENTITAS MATA KULIAH")

	identitasData := [][]string{
		{"Nama Mata Kuliah", rps.Identitas.NamaMataKuliah},
		{"Kode Mata Kuliah", rps.Identitas.KodeMataKuliah},
		{"SKS", fmt.Sprintf("%d", rps.Identitas.SKS)},
		{"Semester", rps.Identitas.Semester},
		{"Dosen Pengampu", rps.Identitas.DosenPengampu},
	}
	s.addKeyValueTable(pdf, identitasData)

	// ==================== CAPAIAN PEMBELAJARAN ====================
	s.addSectionTitle(pdf, "II. CAPAIAN PEMBELAJARAN")

	// CPL Prodi
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "A. Capaian Pembelajaran Lulusan (CPL) Prodi", "", 1, "L", false, 0, "")
	s.addNumberedList(pdf, rps.CapaianPembelajaran.CPLProdi)

	// CPMK
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "B. Capaian Pembelajaran Mata Kuliah (CPMK)", "", 1, "L", false, 0, "")
	s.addNumberedList(pdf, rps.CapaianPembelajaran.CPMK)

	// Sub-CPMK
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "C. Sub-CPMK", "", 1, "L", false, 0, "")
	s.addNumberedList(pdf, rps.CapaianPembelajaran.SubCPMK)

	// ==================== DESKRIPSI MATA KULIAH ====================
	s.addSectionTitle(pdf, "III. DESKRIPSI MATA KULIAH")

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, s.sanitizeText(rps.DeskripsiMataKuliah.DeskripsiSingkat), "", "J", false)
	pdf.Ln(3)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "Bahan Kajian:", "", 1, "L", false, 0, "")
	s.addNumberedList(pdf, rps.DeskripsiMataKuliah.BahanKajian)

	// ==================== RENCANA PEMBELAJARAN MINGGUAN ====================
	pdf.AddPage()
	s.addSectionTitle(pdf, "IV. RENCANA PEMBELAJARAN MINGGUAN")

	s.addWeeklyPlanTable(pdf, rps.RencanaMingguan)

	// ==================== RENCANA PENILAIAN ====================
	pdf.AddPage()
	s.addSectionTitle(pdf, "V. RENCANA PENILAIAN")

	s.addAssessmentTable(pdf, rps.RencanaPenilaian.Komponen)

	// ==================== DAFTAR REFERENSI ====================
	s.addSectionTitle(pdf, "VI. DAFTAR REFERENSI")

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "A. Referensi Utama", "", 1, "L", false, 0, "")
	s.addNumberedList(pdf, rps.DaftarReferensi.Utama)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "B. Referensi Pendukung", "", 1, "L", false, 0, "")
	s.addNumberedList(pdf, rps.DaftarReferensi.Pendukung)

	// Generate PDF bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// Helper functions for PDF generation
func (s *exportService) addSectionTitle(pdf *gofpdf.Fpdf, title string) {
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(41, 128, 185) // Blue color
	pdf.SetTextColor(255, 255, 255) // White text
	pdf.CellFormat(0, 8, title, "1", 1, "L", true, 0, "")
	pdf.SetTextColor(0, 0, 0) // Reset to black
	pdf.Ln(3)
}

func (s *exportService) addKeyValueTable(pdf *gofpdf.Fpdf, data [][]string) {
	pdf.SetFont("Arial", "", 10)
	for _, row := range data {
		pdf.SetFont("Arial", "B", 10)
		pdf.SetFillColor(240, 240, 240)
		pdf.CellFormat(60, 7, row[0], "1", 0, "L", true, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 7, s.sanitizeText(row[1]), "1", 1, "L", false, 0, "")
	}
	pdf.Ln(3)
}

func (s *exportService) addNumberedList(pdf *gofpdf.Fpdf, items []string) {
	pdf.SetFont("Arial", "", 10)
	for i, item := range items {
		text := fmt.Sprintf("%d. %s", i+1, s.sanitizeText(item))
		pdf.MultiCell(0, 5, text, "", "L", false)
	}
	pdf.Ln(3)
}

func (s *exportService) addWeeklyPlanTable(pdf *gofpdf.Fpdf, plans []dto.RPSRencanaMingguan) {
	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(41, 128, 185) // Blue header
	pdf.SetTextColor(255, 255, 255)

	// Landscape A4: 297mm width - 30mm margins = 267mm usable
	headers := []string{"Minggu", "Topik", "Indikator", "Metode", "Waktu", "Penilaian"}
	widths := []float64{15, 70, 70, 45, 17, 50}

	for i, header := range headers {
		pdf.CellFormat(widths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(0, 0, 0)

	// Table content with alternating row colors
	pdf.SetFont("Arial", "", 8)
	for i, plan := range plans {
		// Alternating row colors
		if i%2 == 0 {
			pdf.SetFillColor(255, 255, 255)
		} else {
			pdf.SetFillColor(245, 245, 245)
		}

		// Calculate max lines needed for this row
		topik := s.wrapText(plan.Topik, 68)
		indikator := s.wrapText(plan.IndikatorCapaian, 68)
		metode := s.wrapText(plan.MetodePembelajaran, 43)
		penilaian := s.wrapText(plan.BentukPenilaian, 48)

		// Get max lines
		maxLines := s.maxInt(
			len(topik),
			s.maxInt(len(indikator), s.maxInt(len(metode), len(penilaian))),
		)
		if maxLines < 1 {
			maxLines = 1
		}

		rowHeight := float64(maxLines) * 5.0
		if rowHeight < 8 {
			rowHeight = 8
		}

		// Draw cells
		x := pdf.GetX()
		y := pdf.GetY()

		// Check if we need a new page
		if y+rowHeight > 190 {
			pdf.AddPage()
			y = pdf.GetY()
		}

		// Minggu
		pdf.SetXY(x, y)
		pdf.CellFormat(widths[0], rowHeight, fmt.Sprintf("%d", plan.Minggu), "1", 0, "C", true, 0, "")

		// Topik
		pdf.SetXY(x+widths[0], y)
		s.multiCellInTable(pdf, widths[1], rowHeight, strings.Join(topik, "\n"), true)

		// Indikator
		pdf.SetXY(x+widths[0]+widths[1], y)
		s.multiCellInTable(pdf, widths[2], rowHeight, strings.Join(indikator, "\n"), true)

		// Metode
		pdf.SetXY(x+widths[0]+widths[1]+widths[2], y)
		s.multiCellInTable(pdf, widths[3], rowHeight, strings.Join(metode, "\n"), true)

		// Waktu
		pdf.SetXY(x+widths[0]+widths[1]+widths[2]+widths[3], y)
		pdf.CellFormat(widths[4], rowHeight, fmt.Sprintf("%d", plan.WaktuMenit), "1", 0, "C", true, 0, "")

		// Penilaian
		pdf.SetXY(x+widths[0]+widths[1]+widths[2]+widths[3]+widths[4], y)
		s.multiCellInTable(pdf, widths[5], rowHeight, strings.Join(penilaian, "\n"), true)

		pdf.SetXY(x, y+rowHeight)
	}
}

func (s *exportService) multiCellInTable(pdf *gofpdf.Fpdf, width, height float64, text string, fill bool) {
	x := pdf.GetX()
	y := pdf.GetY()

	// Draw the cell border and fill
	if fill {
		r, g, b := pdf.GetFillColor()
		pdf.SetFillColor(int(r), int(g), int(b))
		pdf.Rect(x, y, width, height, "FD")
	} else {
		pdf.Rect(x, y, width, height, "D")
	}

	// Add text with padding
	pdf.SetXY(x+1, y+1)
	pdf.MultiCell(width-2, 5, text, "", "L", false)
}

func (s *exportService) wrapText(text string, maxWidth int) []string {
	text = s.sanitizeText(text)
	if text == "" {
		return []string{""}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		// Rough estimate: ~2.5 chars per mm for Arial 8pt
		if len(testLine) > maxWidth/2 {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func (s *exportService) maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *exportService) addAssessmentTable(pdf *gofpdf.Fpdf, komponen []dto.RPSKomponenPenilaian) {
	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(41, 128, 185)
	pdf.SetTextColor(255, 255, 255)

	headers := []string{"No", "Komponen", "Bobot (%)", "Teknik", "Instrumen"}
	widths := []float64{12, 50, 25, 85, 95}

	for i, header := range headers {
		pdf.CellFormat(widths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(0, 0, 0)

	// Table content
	pdf.SetFont("Arial", "", 9)
	totalBobot := 0
	for i, k := range komponen {
		// Alternating colors
		if i%2 == 0 {
			pdf.SetFillColor(255, 255, 255)
		} else {
			pdf.SetFillColor(245, 245, 245)
		}

		pdf.CellFormat(widths[0], 7, fmt.Sprintf("%d", i+1), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[1], 7, s.sanitizeText(k.Nama), "1", 0, "L", true, 0, "")
		pdf.CellFormat(widths[2], 7, fmt.Sprintf("%d%%", k.Bobot), "1", 0, "C", true, 0, "")
		pdf.CellFormat(widths[3], 7, s.truncateText(k.Teknik, 100), "1", 0, "L", true, 0, "")
		pdf.CellFormat(widths[4], 7, s.truncateText(k.Instrumen, 110), "1", 0, "L", true, 0, "")
		pdf.Ln(-1)
		totalBobot += k.Bobot
	}

	// Total row
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(widths[0]+widths[1], 8, "TOTAL", "1", 0, "C", true, 0, "")
	pdf.CellFormat(widths[2], 8, fmt.Sprintf("%d%%", totalBobot), "1", 0, "C", true, 0, "")
	pdf.CellFormat(widths[3]+widths[4], 8, "", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)
}

func (s *exportService) sanitizeText(text string) string {
	// Replace special characters that might cause issues
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\t", " ")
	return text
}

func (s *exportService) truncateText(text string, maxLen int) string {
	text = s.sanitizeText(text)
	if len(text) > maxLen {
		return text[:maxLen-3] + "..."
	}
	return text
}

// ExportToHTML generates an HTML document from RPS data
func (s *exportService) ExportToHTML(rps *dto.RPSStructuredOutput) (string, error) {
	html := `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>RPS - ` + rps.Identitas.NamaMataKuliah + `</title>
    <style>
        * {
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 297mm;
            margin: 0 auto;
            padding: 20mm;
            line-height: 1.6;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            text-align: center;
            font-size: 22pt;
            margin-bottom: 30px;
            color: #2980b9;
            border-bottom: 3px solid #2980b9;
            padding-bottom: 15px;
        }
        h2 {
            font-size: 14pt;
            background: linear-gradient(135deg, #2980b9 0%, #3498db 100%);
            color: white;
            padding: 12px 15px;
            margin-top: 25px;
            margin-bottom: 15px;
            border-radius: 5px;
        }
        h3 {
            font-size: 12pt;
            margin-top: 15px;
            color: #2c3e50;
            border-left: 4px solid #2980b9;
            padding-left: 10px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 15px 0;
            font-size: 10pt;
        }
        th, td {
            border: 1px solid #bdc3c7;
            padding: 10px 12px;
            text-align: left;
        }
        th {
            background: linear-gradient(135deg, #2980b9 0%, #3498db 100%);
            color: white;
            font-weight: 600;
        }
        tbody tr:nth-child(even) {
            background-color: #f8f9fa;
        }
        tbody tr:hover {
            background-color: #e8f4f8;
        }
        ol, ul {
            margin: 10px 0;
            padding-left: 25px;
        }
        li {
            margin-bottom: 5px;
        }
        .info-table {
            width: auto;
            min-width: 500px;
        }
        .info-table td:first-child {
            width: 200px;
            font-weight: 600;
            background-color: #ecf0f1;
            color: #2c3e50;
        }
        .info-table td:last-child {
            min-width: 300px;
        }
        .weekly-table th {
            text-align: center;
        }
        .weekly-table td:first-child {
            text-align: center;
            font-weight: 600;
        }
        .weekly-table td:nth-child(5) {
            text-align: center;
        }
        .assessment-table td:first-child {
            text-align: center;
        }
        .assessment-table td:nth-child(3) {
            text-align: center;
            font-weight: 600;
        }
        .total-row {
            font-weight: bold;
            background-color: #ecf0f1 !important;
        }
        .total-row td:first-child {
            text-align: center;
        }
        p {
            text-align: justify;
            margin-bottom: 15px;
        }
        @media print {
            body { 
                padding: 10mm;
                background-color: white;
            }
            .container {
                box-shadow: none;
                padding: 0;
            }
            h2 { 
                page-break-after: avoid;
                -webkit-print-color-adjust: exact;
                print-color-adjust: exact;
            }
            table { 
                page-break-inside: avoid;
            }
            th {
                -webkit-print-color-adjust: exact;
                print-color-adjust: exact;
            }
        }
    </style>
</head>
<body>
    <div class="container">
    <h1>RENCANA PEMBELAJARAN SEMESTER (RPS)</h1>
    
    <h2>I. IDENTITAS MATA KULIAH</h2>
    <table class="info-table">
        <tr><td>Nama Mata Kuliah</td><td>` + rps.Identitas.NamaMataKuliah + `</td></tr>
        <tr><td>Kode Mata Kuliah</td><td>` + rps.Identitas.KodeMataKuliah + `</td></tr>
        <tr><td>SKS</td><td>` + fmt.Sprintf("%d", rps.Identitas.SKS) + `</td></tr>
        <tr><td>Semester</td><td>` + rps.Identitas.Semester + `</td></tr>
        <tr><td>Dosen Pengampu</td><td>` + rps.Identitas.DosenPengampu + `</td></tr>
    </table>
    
    <h2>II. CAPAIAN PEMBELAJARAN</h2>
    <h3>A. Capaian Pembelajaran Lulusan (CPL) Prodi</h3>
    <ol>` + s.generateListItems(rps.CapaianPembelajaran.CPLProdi) + `</ol>
    
    <h3>B. Capaian Pembelajaran Mata Kuliah (CPMK)</h3>
    <ol>` + s.generateListItems(rps.CapaianPembelajaran.CPMK) + `</ol>
    
    <h3>C. Sub-CPMK</h3>
    <ol>` + s.generateListItems(rps.CapaianPembelajaran.SubCPMK) + `</ol>
    
    <h2>III. DESKRIPSI MATA KULIAH</h2>
    <p>` + rps.DeskripsiMataKuliah.DeskripsiSingkat + `</p>
    <h3>Bahan Kajian:</h3>
    <ol>` + s.generateListItems(rps.DeskripsiMataKuliah.BahanKajian) + `</ol>
    
    <h2>IV. RENCANA PEMBELAJARAN MINGGUAN</h2>
    <table class="weekly-table">
        <thead>
            <tr>
                <th>Minggu</th>
                <th>Topik</th>
                <th>Sub Topik</th>
                <th>Indikator Capaian</th>
                <th>Metode</th>
                <th>Waktu (menit)</th>
                <th>Penilaian</th>
            </tr>
        </thead>
        <tbody>` + s.generateWeeklyRows(rps.RencanaMingguan) + `</tbody>
    </table>
    
    <h2>V. RENCANA PENILAIAN</h2>
    <table class="assessment-table">
        <thead>
            <tr>
                <th>No</th>
                <th>Komponen</th>
                <th>Bobot (%)</th>
                <th>Teknik</th>
                <th>Instrumen</th>
            </tr>
        </thead>
        <tbody>` + s.generateAssessmentRows(rps.RencanaPenilaian.Komponen) + `</tbody>
    </table>
    
    <h2>VI. DAFTAR REFERENSI</h2>
    <h3>A. Referensi Utama</h3>
    <ol>` + s.generateListItems(rps.DaftarReferensi.Utama) + `</ol>
    
    <h3>B. Referensi Pendukung</h3>
    <ol>` + s.generateListItems(rps.DaftarReferensi.Pendukung) + `</ol>
    </div>
</body>
</html>`

	return html, nil
}

func (s *exportService) generateListItems(items []string) string {
	var sb strings.Builder
	for _, item := range items {
		sb.WriteString("<li>" + item + "</li>")
	}
	return sb.String()
}

func (s *exportService) generateWeeklyRows(plans []dto.RPSRencanaMingguan) string {
	var sb strings.Builder
	for _, p := range plans {
		subTopiks := strings.Join(p.SubTopik, ", ")
		sb.WriteString(fmt.Sprintf(`
            <tr>
                <td>%d</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%d</td>
                <td>%s</td>
            </tr>`,
			p.Minggu, p.Topik, subTopiks, p.IndikatorCapaian,
			p.MetodePembelajaran, p.WaktuMenit, p.BentukPenilaian))
	}
	return sb.String()
}

func (s *exportService) generateAssessmentRows(komponen []dto.RPSKomponenPenilaian) string {
	var sb strings.Builder
	totalBobot := 0
	for i, k := range komponen {
		sb.WriteString(fmt.Sprintf(`
            <tr>
                <td>%d</td>
                <td>%s</td>
                <td>%d%%</td>
                <td>%s</td>
                <td>%s</td>
            </tr>`,
			i+1, k.Nama, k.Bobot, k.Teknik, k.Instrumen))
		totalBobot += k.Bobot
	}
	sb.WriteString(fmt.Sprintf(`
        <tr class="total-row">
            <td colspan="2">TOTAL</td>
            <td>%d%%</td>
            <td colspan="2"></td>
        </tr>`, totalBobot))
	return sb.String()
}
