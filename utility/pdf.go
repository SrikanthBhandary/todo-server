package utility

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/srikanthbhandary/todo-server/entity"
)

// RandomFilename generates a random filename with the specified extension.
func RandomFilename(extension string) (string, error) {
	// Generate 8 random bytes (16 hex characters)
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert bytes to a hex string
	randomString := hex.EncodeToString(bytes)

	// Append the extension
	filename := fmt.Sprintf("%s.%s", randomString, extension)

	return filename, nil
}

type PDFGenerator struct {
	OutPutDirectory string
}

func NewPDFGenerator(opDirectory string) *PDFGenerator {
	return &PDFGenerator{OutPutDirectory: opDirectory}
}

func (pc *PDFGenerator) GenerateToDosReport(userID int, userName, email string, todos []entity.ToDo) (string, error) {
	// Create the PDF file
	fileName, err := RandomFilename(".pdf")
	if err != nil {
		return "", fmt.Errorf("error generating filename: %w", err)
	}

	outPutFilePath := filepath.Join(pc.OutPutDirectory, strconv.Itoa(userID)+"_"+fileName)
	file, err := os.Create(outPutFilePath)
	if err != nil {
		fmt.Println("Error creating PDF file:", err)
		return "", err
	}
	defer file.Close()

	// Write the PDF header
	file.WriteString("%PDF-1.4\n")

	// Define the body of the PDF
	pdfObjects := []string{
		"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n",                                                                              // Object 1: Catalog
		"2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n",                                                                      // Object 2: Pages
		"3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n", // Object 3: Page
	}

	// Start building the content stream for the PDF (title, user details, table, watermark)
	var contentStream string

	// 1. Title in color (RGB for teal color)
	contentStream += "BT /F1 16 Tf 0.0 0.5 0.5 rg 100 750 Td (ToDo List) Tj ET\n" // Teal title

	// 2. User details
	contentStream += fmt.Sprintf("BT /F1 12 Tf 0 0 0 rg 100 720 Td (User: %s) Tj ET\n", userName)
	contentStream += fmt.Sprintf("BT /F1 12 Tf 0 0 0 rg 100 700 Td (Email: %s) Tj ET\n", email)

	// 3. Watermark: Add semi-transparent text in the background (in gray)
	contentStream += "q 0.9 g /F1 40 Tf 200 400 Td 0.2 Tc (DRAFT) Tj Q\n" // 'q' and 'Q' start and end the graphics state to isolate the watermark

	// 4. Table headers
	contentStream += "BT /F1 12 Tf 0 0 0 rg 100 670 Td (ID) Tj ET\n"          // ID column
	contentStream += "BT /F1 12 Tf 0 0 0 rg 200 670 Td (Task) Tj ET\n"        // Task column
	contentStream += "BT /F1 12 Tf 0 0 0 rg 400 670 Td (Description) Tj ET\n" // Description column
	contentStream += "100 665 m 500 665 l S\n"                                // Draw a line under the headers

	// 5. Loop through todos and display them in a table format
	yPosition := 650
	for _, todo := range todos {
		// Print ID
		contentStream += fmt.Sprintf("BT /F1 12 Tf 0 0 0 rg 100 %d Td (%d) Tj ET\n", yPosition, todo.ToDoID)

		// Print Task
		contentStream += fmt.Sprintf("BT /F1 12 Tf 0 0 0 rg 200 %d Td (%s) Tj ET\n", yPosition, todo.Title)

		// Print Description
		contentStream += fmt.Sprintf("BT /F1 12 Tf 0 0 0 rg 400 %d Td (%s) Tj ET\n", yPosition, todo.Description)

		// Move to next line
		yPosition -= 20
	}

	// Object 4: Page contents (dynamically generated todos and user info)
	pdfObjects = append(pdfObjects, fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", len(contentStream), contentStream))
	pdfObjects = append(pdfObjects, "5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n") // Object 5: Font

	// Write all objects to the file
	for _, obj := range pdfObjects {
		file.WriteString(obj)
	}

	// Write the cross-reference table
	file.WriteString("xref\n0 6\n0000000000 65535 f \n0000000010 00000 n \n0000000067 00000 n \n0000000120 00000 n \n0000000278 00000 n \n0000000337 00000 n \n")

	// Write the trailer
	file.WriteString("trailer\n<< /Size 6 /Root 1 0 R >>\nstartxref\n392\n%%EOF\n")

	fmt.Println("PDF with colorful title, table, and watermark generated successfully!")
	return outPutFilePath, nil
}
