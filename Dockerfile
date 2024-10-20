# Stage 1: Build the Go app
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o todo-server ./cmd/server

# Stage 2: Final image
FROM alpine:3.18

# Set the working directory inside the final container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/todo-server .

# Copy the static HTML files (if applicable)
COPY ./static ./static

# Copy the config file (you can replace this with a volume in docker-compose)
COPY config.yaml /app/config.yaml

# Expose the port your Go app will run on
EXPOSE 8080

# Command to run the Go application
CMD ["./todo-server"]