package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/srikanthdoc/todo-server/entity"
)

func (s *Router) CreateToDo(w http.ResponseWriter, r *http.Request) {
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

	err = s.todoService.AddToDo(r.Context(), &todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create todo", "message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "ToDo created successfully"})
}

func (s *Router) GetAllToDos(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the context
	userID := r.Context().Value("userID").(int)

	todos, err := s.todoService.GetAllTodos(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to retrieve todos", "message": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (s *Router) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["todoID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid todo ID"})
		return
	}

	// Extract user ID from the context
	userID := r.Context().Value("userID").(int)

	todo, err := s.todoService.GetTodo(r.Context(), userID, todoID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "todo not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (s *Router) DeleteToDo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["todoID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid todo ID"})
		return
	}

	// Extract user ID from the context
	userID := r.Context().Value("userID").(int)

	err = s.todoService.DeleteToDo(r.Context(), userID, todoID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete todo", "message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Router) DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the context
	userID := r.Context().Value("userID").(int)

	err := s.todoService.DeleteAllTodos(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete todos", "message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
