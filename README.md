# ToDo Application in Go 
[![Go](https://github.com/SrikanthBhandary/todo-server/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/SrikanthBhandary/todo-server/actions/workflows/go.yml) [![Integration Test](https://github.com/SrikanthBhandary/todo-server/actions/workflows/integration_test.yml/badge.svg)](https://github.com/SrikanthBhandary/todo-server/actions/workflows/integration_test.yml)

This ToDo application demonstrates my learning journey with Golang and this will help the entry level developers to understand the golang components, showcasing various concepts and features of the language. 

## Configuration
The application uses a configuration file to manage various settings, including database connection details, server ports, and Redis configurations.

Config File Structure
Create a config.yaml file in the root directory of your project with the following structure:


``` YAML
port: "9999"  # Server port
dsn: "host=localhost port=5432 user=srikanth password=password dbname=postgres sslmode=disable"  # Database connection string  
jwt_secret_key: "your_jwt_secret"  # JWT secret key
redis_address: "localhost:6379"  # Redis server address
html_assets_path: "./todo-server/static/html"  # Path to HTML files . Give the absolute path to the html templates
num_of_workers: 5  # Number of worker threads
smtp_host: "smtp.example.com"  # SMTP server address
smtp_port: 587  # SMTP server port
smtp_user_name: "your_email@example.com"  # SMTP username
smtp_password: "your_email_password"  # SMTP password
```


## Migrate Commands - To initate the DB with the tables


### Executing Migrations

To execute the migrations, you need to set the `POSTGRESQL_URL` environment variable to point to your PostgreSQL database. Use the following command:

```bash
export POSTGRESQL_URL='postgres://postgres:password@localhost:5432/postgres?sslmode=disable'
```

Replace `postgres:password@localhost:5432/postgres` with your actual database credentials and connection details.

Then, run the following command to apply the migrations:

```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

- `-database ${POSTGRESQL_URL}`: Specifies the database connection string.
- `-path db/migrations`: Indicates the directory containing the migration files.
- `up`: Applies the migrations to the database.

---

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
4. Execute the application with `go run .`. or just use go run cmd/server/main.go from root directory of the cloned repo.

5. Access the API at `http://localhost:8080`.

## Future Improvements
- Handle validation
- Enhance error handling and logging.
- Create worker to mimic to send an email when todo is added
- Front end to validate delete todo, edit todo

---

Feel free to modify any sections or add additional details as needed!

![image](./docs/images/todo_login.png)

![image](./docs/images/todo_home.png)
