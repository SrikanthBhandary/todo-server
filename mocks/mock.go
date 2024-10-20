package mocks

import (
	"fmt"
	"log"
)

// MockEmailSender is a mock implementation of the EmailSender interface for testing purposes.
type MockEmailSender struct {
	SentEmails []string // Holds the details of sent emails for verification.
}

// SendEmail simulates sending an email by storing the email details in the SentEmails slice.
func (m *MockEmailSender) SendEmail(to []string, subject, body string) error {
	// Simulate the email sending by appending to the SentEmails slice.
	emailDetails := fmt.Sprintf("To: %s, Subject: %s, Body: %s", to[0], subject, body)
	m.SentEmails = append(m.SentEmails, emailDetails)
	log.Println("sending...", emailDetails)
	return nil // Return nil to indicate success.
}
