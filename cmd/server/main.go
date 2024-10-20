package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/srikanthdoc/todo-server/mocks"
	"github.com/srikanthdoc/todo-server/repository"
	"github.com/srikanthdoc/todo-server/router"
	"github.com/srikanthdoc/todo-server/service"
	"github.com/srikanthdoc/todo-server/worker"
)

const (
	port               = ":9999"
	numOfWorkers       = 3
	shutdownTimeoutSec = 5
	dsn                = "user=postgres password=password dbname=postgres host=localhost sslmode=disable" // Update with your actual DB connection details
	jwtSecret          = "random-secret-key-for-production"                                               // ToDo: Secure this key
)

func main() {
	db := initDB()
	rdb := initRedisDB()

	defer db.Close()  // Ensure database connection is closed on exit
	defer rdb.Close() // Closing the redis connection

	shutdown := setupSignalHandler()
	jobChannel := make(chan worker.Job, 10) // Buffered channel for jobs
	ctx, cancel := context.WithCancel(context.Background())

	pool := setupWorkerPool(ctx, jobChannel)

	ratelLimiter := router.NewRedisRateLimiter(ctx, rdb, 100, 10*time.Second)

	userRepo := repository.NewPostgresUserRepository(db)
	todoRepo := repository.NewPostgresToDoRepository(db)

	userService := service.NewUserService(userRepo)
	todoService := service.NewTodoService(todoRepo)
	jwtService := service.NewJWTService(jwtSecret)

	emailSender := &mocks.MockEmailSender{}

	todoHandler := setupServer(todoService, userService, jwtService, ratelLimiter, pool, emailSender)

	srv := startHTTPServer(todoHandler)

	// Wait for shutdown signal
	<-shutdown

	shutdownServer(srv, pool, cancel)
}

// initDB initializes the database connection.
func initDB() *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the database: %s", err)
	}
	return db
}

// initRedisDB initializes the redis connection
func initRedisDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", //ToDo move all these to config
	})

	return rdb
}

// setupSignalHandler sets up a channel to listen for OS interrupt signals.
func setupSignalHandler() chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	return shutdown
}

// setupWorkerPool initializes the worker pool.
func setupWorkerPool(ctx context.Context, jobChannel chan worker.Job) *worker.WorkerPool {
	pool := worker.NewWorkerPool(numOfWorkers, jobChannel)
	pool.Init(ctx)
	return pool
}

// setupServer initializes the HTTP server with the router and services.
func setupServer(todoService service.ToDoService, userService service.UserService,
	jwtService service.JWTValidator, rateLimiter router.RateLimiter, pool *worker.WorkerPool,
	emailSender *mocks.MockEmailSender) *router.Router {
	todoHandler := router.NewRouter(todoService, userService, jwtService, rateLimiter, pool, emailSender)
	todoHandler.InitRoutes()
	return todoHandler
}

// startHTTPServer starts the HTTP server.
func startHTTPServer(todoHandler *router.Router) *http.Server {
	srv := &http.Server{
		Addr:    port,
		Handler: todoHandler.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("server started on", port)
	return srv
}

// shutdownServer handles graceful shutdown of the server.
func shutdownServer(srv *http.Server, pool *worker.WorkerPool, cancel context.CancelFunc) {
	log.Println("shutting down the server...")

	// Signal the worker pool to stop
	cancel() // Cancel the context for workers
	pool.Stop()

	// Wait for the worker pool to finish processing jobs
	pool.Wait()

	// Create a context with a timeout to allow ongoing connections to complete
	ctxShutdown, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeoutSec*time.Second)
	defer shutdownCancel() // Ensure this context is cancelled

	// Attempt graceful shutdown of the server
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}

	log.Println("server gracefully stopped")
}
