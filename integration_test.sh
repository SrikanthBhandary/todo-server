#!/bin/bash
set -e

# Wait for the application to start
echo "Waiting for the application to be ready..."
sleep 15  # Adjust this time as necessary

# Create a user
echo "Creating a new user..."
response=$(curl -s -X POST http://host.docker.internal:8080/users \
    -H "Content-Type: application/json" \
    -d '{"user_name": "testuser16", "password": "mypassword15", "email": "testuser16@example.com"}')

# Assert the response message
expected_message="User created successfully"
if [[ $(echo "$response" | jq -r '.message') == "$expected_message" ]]; then
    echo "User creation assertion passed: $expected_message"
else
    echo "User creation assertion failed: Expected '$expected_message', got '$response'"
    exit 1
fi

# Login the user and extract the token
echo "Logging in the user..."
login_response=$(curl -s -X POST http://host.docker.internal:8080/login \
    -H "Content-Type: application/json" \
    -d '{"username": "testuser16", "password": "mypassword15"}')

token=$(echo $login_response | jq -r '.token')
echo "Received JWT token: $token"

# Create a ToDo
echo "Creating a ToDo..."
curl -s -X POST http://host.docker.internal:8080/todos \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $token" \
    -d '{"title": "My New Todo", "description": "This is a description of my new todo"}'

# Get ToDos
echo "Getting the ToDos..."
curl -s -X GET http://host.docker.internal:8080/todos \
    -H "Authorization: Bearer $token"

echo "Integration test completed successfully."
