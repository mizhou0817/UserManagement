# UserManagement

This is a basic user management system built in golang and using PostgreSQL as DB.
The system provides RESTful APIs that allow for adding new users, updating existing user information, 
retrieving a user by ID, and listing all users. All data will be stored in postgresql database 
whose config (lines from 24 to 30 in the file "handlers.go") can be edited.
  
## Project structure
```bash

├── README.md
│   ├── middleware                  // API handlers
│   │   ├── handlers.go
│   ├── models                      // DB models
│   │   └── models.go
├── router
│   └── router.go                   // API routes
├── go.mod                                    // Go modules
├── go.sum
├── main.go

```

## PostgreSQL Table Initialization

```sql
CREATE TABLE users (
  age INT,
  email TEXT,
  name  TEXT,
  id SERIAL PRIMARY KEY
);
```

## Endpoints:

api_endpoint: http://localhost:8080

```sh
POST   api_endpoint/users               # add a new user
GET    api_endpoint/users/{id}          # retrieve a user by ID
PUT    api_endpoint/users/{id}          # update existing user info by ID
GET    api_endpoint/users               # list all users by default order
GET    api_endpoint/users?sortBy=asc    # list all users by ascending order
GET    api_endpoint/users?sortBy=desc   # list all users by descending order
```

### Create User
This endpoint adds a new user to the `users` relation.  
Sends a `POST` request to `api_endpoint/users`:
```sh
curl -X POST 'http://localhost:8080/users' -d '{"name": "Mi", "email": "mz558@cornell.edu", "age": 28}'
```
Response:  
```sh
{"id":1,"message":"User created successfully"}
```

### Get User
This endpoint retrieves a user info by ID.  
Sends a `GET` request to `api_endpoint/users/{id}`:
```sh
curl -X GET 'http://localhost:8080/users/1'
```
Response:
```sh
{"id":1,"name":"Mi","email":"mz558@cornell.edu","age":28}
```

### Update User
This endpoint updates the user info by ID.  
Send a `PUT` request to `api_endpoint/users/{id}`:
```sh
curl -X PUT 'http://localhost:8080/users/1' -d '{"email": "mizhou0817@outlook.com"}'
```
Response:
```sh
{"id":1,"message":"User updated successfully. Total rows/record affected 1"}
```

### Get ALL User
This endpoint retrieves all users' information.  
Sends a `GET` request to `api_endpoint/users`:
```sh
curl -X GET 'http://localhost:8080/users'
```
Response:
```sh
[{"id":1,"name":"Mi","email":"mizhou0817@outlook.com","age":28},{"id":2,"name":"Mi","email":"mizhou0817@gmail.com","age":20},{"id":3,"name":"Sisi","email":"sisi0422@outlook.com","age":24}]
```
 
Sends a `GET` request to `api_endpoint/users` with specifying asceding order:
```sh
curl -X GET 'http://localhost:8080/users?sortBy=asc'
```
Response:
```sh
[{"id":2,"name":"Mi","email":"mizhou0817@gmail.com","age":20},{"id":1,"name":"Mi","email":"mizhou0817@outlook.com","age":28},{"id":3,"name":"Sisi","email":"sisi0422@outlook.com","age":24}]
```
 
Sends a `GET` request to `api_endpoint/users` with specifying descending order:
```sh
curl -X GET 'http://localhost:8080/users?sortBy=desc'
```
Response:
```sh
[{"id":3,"name":"Sisi","email":"sisi0422@outlook.com","age":24},{"id":1,"name":"Mi","email":"mizhou0817@outlook.com","age":28},{"id":2,"name":"Mi","email":"mizhou0817@gmail.com","age":20}]
```