
### CreateUser


    curl -X POST http://localhost:8080/users \
    -H "Content-Type: application/json" \
    -d '{"user_name": "testuser16", "password": "mypassword15", "email": "testuser16@example.com"}'


### Get User

    curl -X GET http://localhost:8080/users/3

### Login User

    curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"username": "testuser16", "password": "mypassword15"}'

### Create ToDo

    curl -X POST http://localhost:8080/todos \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk0MjE3MzAsInN1YiI6MTd9.YCXldEymIbglOe6___TaxFLu7NvflDNew-_xOBTghwY" \
    -d '{
    "title": "My New Todo",
    "description": "This is a description of my new todo"
    }'

### Get ToDO
    curl -X GET http://localhost:8080/todos \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk0MjE3MzAsInN1YiI6MTd9.YCXldEymIbglOe6___TaxFLu7NvflDNew-_xOBTghwY"