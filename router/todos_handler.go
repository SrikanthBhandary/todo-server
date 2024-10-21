package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/srikanthbhandary/todo-server/entity"
	"github.com/srikanthbhandary/todo-server/utility"
	"github.com/srikanthbhandary/todo-server/worker"
)

func (rt *Router) CreateToDo(w http.ResponseWriter, r *http.Request) {
	var todo entity.ToDo // Make sure you have a ToDo struct in your service/entity package
	err := json.NewDecoder(r.Body).Decode(&todo)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body", "message": err.Error()})
		return
	}

	// Extract user ID from the context
	userID := r.Context().Value("userID").(int)
	todo.UserID = userID // Associate the todo with the logged-in user

	err = rt.todoService.AddToDo(r.Context(), &todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create todo", "message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "ToDo created successfully"})
}

func (rt *Router) GetAllToDos(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the context
	userID := r.Context().Value("userID").(int)

	todos, err := rt.todoService.GetAllTodos(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to retrieve todos", "message": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (rt *Router) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["todoID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid todo ID"})
		return
	}

	// Extract user ID from the context
	userID := r.Context().Value("userID").(int)

	todo, err := rt.todoService.GetTodo(r.Context(), userID, todoID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "todo not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (rt *Router) DeleteToDo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["todoID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid todo ID"})
		return
	}

	// Extract user ID from the context
	userID := r.Context().Value("userID").(int)

	err = rt.todoService.DeleteToDo(r.Context(), userID, todoID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete todo", "message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *Router) DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the context
	userID := r.Context().Value("userID").(int)

	err := rt.todoService.DeleteAllTodos(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete todos", "message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *Router) DownloadToDos(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	user, err := rt.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	todos, err := rt.todoService.GetAllTodos(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error fetching todos", http.StatusInternalServerError)
		return
	}

	// Enqueue the PDF generation job
	pdfJob := &worker.PDFJob{
		UserID:    userID,
		UserName:  user.UserName,
		Email:     user.Email,
		Todos:     todos,
		Generator: utility.NewPDFGenerator(rt.Config.PDFOutputPath),
		WebSocket: rt.WorkerPool.WebSocket, // Use the active WebSocket connection
	}
	rt.WorkerPool.EnqueueJob(pdfJob)

	// Inform the client the job is queued
	w.Write([]byte("PDF generation started. You'll be notified when it's ready for download."))
}
