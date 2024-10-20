package worker

import (
	"log"
	"time"
)

type EmailSender interface {
	SendEmail(to []string, subject, body string) error
}

// Job is an interface that defines a single method, Process.
// Any type that implements this method can be considered a "Job".
// Process performs the task and returns an error if something goes wrong.
type Job interface {
	Process() error // Process defines the behavior of a job and returns an error if any occurs.
}

// Notification is a concrete type that implements the Job interface.
// It represents a notification that needs to be sent.
type Notification struct {
	title string // The title or message of the notification.
}

// Process implements the Job interface for Notification.
// It simulates sending a notification by logging the title and adding a delay to simulate work.
func (n *Notification) Process() error {
	log.Println("sending the notification ", n.title) // Log the notification being sent.
	// Simulate a delay in processing (e.g., sending the notification might take time).
	time.Sleep(1 * time.Second)
	return nil // Return nil to indicate the job was processed without error.
}

type EmailJob struct {
	emailSender EmailSender // Use EmailSender interface for flexibility.
	toAddress   []string    // Recipient email addresses.
	subject     string      // Subject of the email.
	body        string      // Body of the email.
}

// NewEmailJob creates a new EmailJob with the provided parameters.
func NewEmailJob(emailSender EmailSender, to []string, subject, body string) *EmailJob {
	return &EmailJob{
		emailSender: emailSender,
		toAddress:   to,
		subject:     subject,
		body:        body,
	}
}

// Process implements the Job interface for EmailJob.
func (ej *EmailJob) Process() error {
	// Use the injected emailSender to send the email.
	return ej.emailSender.SendEmail(ej.toAddress, ej.subject, ej.body)
}
