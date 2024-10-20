package router

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/srikanthbhandary/todo-server/config"
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
	Config      *config.Config
}

type Option func(*Router)

// WithConfig returns an Option that sets the Config for the Router
func WithConfig(cfg *config.Config) Option {
	return func(rt *Router) {
		rt.Config = cfg
	}
}

func NewRouter(todoSvc service.ToDoService, userSvc service.UserService,
	jwtService service.JWTValidator, rateLimiter RateLimiter,
	wp *worker.WorkerPool, emailSender worker.EmailSender,
	options ...Option) *Router {
	r := mux.NewRouter()
	rt := &Router{
		todoService: todoSvc,
		userService: userSvc,
		jwtService:  jwtService,
		rateLimiter: rateLimiter,
		Router:      r,
		WorkerPool:  wp,
		EmailSender: emailSender,
	}

	for _, option := range options {
		option(rt) // Apply each option to the Router
	}
	return rt
}

func (rt *Router) InitRoutes() {
	rt.Router.HandleFunc("/", rt.ServeHTML).Methods("GET")

	// User endpoints
	rt.Router.HandleFunc("/users", rt.CreateUser).Methods("POST")
	rt.Router.HandleFunc("/users/{id}", rt.GetUserByID).Methods("GET")
	rt.Router.HandleFunc("/login", rt.LoginUser).Methods("POST")

	rt.Router.Use(LoggingMiddleware) // Apply any other middleware as needed

	protectedRouter := rt.Router.PathPrefix("/todos").Subrouter()
	protectedRouter.Use(rt.JWTMiddleware)          // Apply JWT middleware to this subrouter
	protectedRouter.Use(rt.JRateLimiterMiddleware) // Apply JWT middleware to this subrouter

	// ToDo endpoints (protected)
	protectedRouter.HandleFunc("", rt.GetAllToDos).Methods("GET")            // /todos
	protectedRouter.HandleFunc("", rt.CreateToDo).Methods("POST")            // /todos for creating a todo
	protectedRouter.HandleFunc("/{todoID}", rt.GetTodo).Methods("GET")       // /todos/{todoID}
	protectedRouter.HandleFunc("/{todoID}", rt.DeleteToDo).Methods("DELETE") // /todos/{todoID}
	protectedRouter.HandleFunc("", rt.DeleteAllTodos).Methods("DELETE")      // /todos for deleting all todos

}

func (rt *Router) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		clearString := strings.ReplaceAll(tokenString, "Bearer ", "")
		userID, err := rt.jwtService.ValidateToken(clearString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Store user ID in context for later use
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rt *Router) JRateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve user ID from context, set by JWTMiddleware
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		// Check if the request is allowed by the rate limiter
		allowed, err := rt.rateLimiter.AllowRequest(strconv.Itoa(userID))
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
func (rt *Router) ServeHTML(w http.ResponseWriter, r *http.Request) {
	htmlFile := filepath.Join(rt.Config.HtmlAssetsPath, "index.html")
	file, err := os.ReadFile(htmlFile)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(file)
}
