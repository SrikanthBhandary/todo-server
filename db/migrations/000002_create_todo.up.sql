CREATE TABLE IF NOT EXISTS todos(
   todo_id serial PRIMARY KEY,
   user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
   title VARCHAR(50),
   description VARCHAR(50),
   email VARCHAR(300) UNIQUE NOT NULL
);
