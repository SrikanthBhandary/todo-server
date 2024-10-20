package router

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/srikanthbhandary/todo-server/service"
	"github.com/srikanthbhandary/todo-server/worker"
)

type Router struct {
	todoService service.ToDoService
	userService service.UserService
	jwtService  service.JWTValidator
	rateLimiter RateLimiter
	Router      *mux.Router
	WorkerPool  *worker.WorkerPool
	EmailSender worker.EmailSender
}

func NewRouter(todoSvc service.ToDoService,
	userSvc service.UserService,
	jwtService service.JWTValidator,
	rateLimiter RateLimiter,
	wp *worker.WorkerPool,
	emailSender worker.EmailSender) *Router {
	r := mux.NewRouter()
	return &Router{
		todoService: todoSvc,
		userService: userSvc,
		jwtService:  jwtService,
		rateLimiter: rateLimiter,
		Router:      r,
		WorkerPool:  wp,
		EmailSender: emailSender,
	}
}

func (s *Router) InitRoutes() {
	s.Router.HandleFunc("/", s.ServeHTML).Methods("GET")

	//prometheus endpoint
	s.Router.Handle("/metrics", promhttp.Handler())

	// User endpoints
	s.Router.HandleFunc("/users", s.CreateUser).Methods("POST")
	s.Router.HandleFunc("/users/{id}", s.GetUserByID).Methods("GET")
	s.Router.HandleFunc("/login", s.LoginUser).Methods("POST")

	s.Router.Use(LoggingMiddleware) // Apply any other middleware as needed

	protectedRouter := s.Router.PathPrefix("/todos").Subrouter()
	protectedRouter.Use(s.JWTMiddleware)          // Apply JWT middleware to this subrouter
	protectedRouter.Use(s.JRateLimiterMiddleware) // Apply JWT middleware to this subrouter

	// ToDo endpoints (protected)
	protectedRouter.HandleFunc("", s.GetAllToDos).Methods("GET")            // /todos
	protectedRouter.HandleFunc("", s.CreateToDo).Methods("POST")            // /todos for creating a todo
	protectedRouter.HandleFunc("/{todoID}", s.GetTodo).Methods("GET")       // /todos/{todoID}
	protectedRouter.HandleFunc("/{todoID}", s.DeleteToDo).Methods("DELETE") // /todos/{todoID}
	protectedRouter.HandleFunc("", s.DeleteAllTodos).Methods("DELETE")      // /todos for deleting all todos

}

func (s *Router) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		clearString := strings.ReplaceAll(tokenString, "Bearer ", "")
		userID, err := s.jwtService.ValidateToken(clearString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Store user ID in context for later use
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Router) JRateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve user ID from context, set by JWTMiddleware
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		// Check if the request is allowed by the rate limiter
		allowed, err := s.rateLimiter.AllowRequest(strconv.Itoa(userID))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
func (s *Router) ServeHTML(w http.ResponseWriter, r *http.Request) {
	//To-Do move it to config.
	htmlFile := "/Users/srikanth.bhandary/Desktop/todo/todo-server/static/html/index.html" // Path to your HTML file
	file, err := os.ReadFile(htmlFile)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(file)
}
