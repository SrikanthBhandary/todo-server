# ToDo Application in Go [![Go](https://github.com/srikanthdoc/todo-server/actions/workflows/go.yml/badge.svg)](https://github.com/srikanthdoc/todo-server/actions/workflows/go.yml)

This ToDo application demonstrates my learning journey with Golang and this will help the entry level developers to understand the golang components, showcasing various concepts and features of the language. 

## Learning Path

Throughout the development of this application, I have focused on several key topics in Go:

### 1. **Basic Syntax and Structures**
   - Understanding Go's syntax, data types, and control structures.

### 2. **Functions and Methods**
   - Creating functions, using method receivers, and understanding the difference between value and pointer receivers.

### 3. **Concurrency**
   - Implementing goroutines and channels to handle concurrent tasks.
   - Building a worker pool for efficient job processing.

### 4. **HTTP Server and Routing**
   - Setting up an HTTP server using the `net/http` package.
   - Implementing a RESTful API with route handling using the Gorilla Mux router.

### 5. **Middleware**
   - Developing middleware functions for logging requests and applying additional processing before reaching handlers.

### 6. **Error Handling**
   - Implementing robust error handling strategies throughout the application.

### 7. **Testing**
   - Writing unit tests for various components of the application using the `testing` package.

### 8. **Email Notifications**
   - Integrating an email notification system to send updates when tasks are created.

### 9. **Graceful Shutdown**
   - Ensuring the application can handle shutdown signals and terminate gracefully.

### 10. **Context Management**
   - Using the `context` package for managing cancellation and timeouts across goroutines.

## Getting Started

To run the application, follow these steps:

1. Clone the repository.
2. Run `go mod tidy` to install dependencies.
3. Go to the fodler cmd/server
4. Execute the application with `go run .`.
5. Access the API at `http://localhost:9999`.

## Future Improvements
- Handle validation
- Enhance error handling and logging.
- Create worker to mimic to send an email when todo is added
- Front end to validate delete todo, edit todo

---

Feel free to modify any sections or add additional details as needed!

![image](./docs/images/todo_login.png)

![image](./docs/images/todo_home.png)
