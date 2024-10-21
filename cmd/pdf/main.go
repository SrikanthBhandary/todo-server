package main

import (
	"log"

	"github.com/srikanthbhandary/todo-server/entity"
	"github.com/srikanthbhandary/todo-server/utility"
)

func main() {
	// Mocked ToDo data
	var todos = []entity.ToDo{
		{ToDoID: 1, Title: "Title1", Description: "Description1"},
		{ToDoID: 2, Title: "Title2", Description: "Description2"},
		{ToDoID: 3, Title: "Title3", Description: "Description3"},
		{ToDoID: 4, Title: "Title4", Description: "Description4"},
	}

	// User details
	userName := "John Doe"
	userEmail := "john.doe@example.com"

	path := utility.NewPDFGenerator(".")
	reportPath, err := path.GenerateToDosReport(1, userName, userEmail, todos)
	if err != nil {
		log.Println("error Reported")
	}
	log.Println("PDF Generated successfully at", reportPath)

}
