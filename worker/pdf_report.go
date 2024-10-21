package worker

import (
	"fmt"

	"github.com/srikanthbhandary/todo-server/entity"
	"github.com/srikanthbhandary/todo-server/utility"
)

type PDFJob struct {
	UserID    int
	UserName  string
	Email     string
	Todos     []entity.ToDo
	Generator *utility.PDFGenerator
	WebSocket *entity.WebSocketConnection // Assuming you have a WebSocket connection
}

func (pj *PDFJob) Process() error {
	// Generate the PDF report
	outputPath, err := pj.Generator.GenerateToDosReport(pj.UserID, pj.UserName, pj.Email, pj.Todos)
	if err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Notify frontend via WebSocket
	message := fmt.Sprintf("PDF generation complete. Download from: /download/%s", outputPath)
	pj.WebSocket.SendMessage([]byte(message)) // Notify via WebSocket

	return nil
}
