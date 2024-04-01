package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go-postgres/models" // models package where User schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"sort"
	"strconv" // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	_ "github.com/lib/pq" // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

const (
	host = "localhost"

	user     = "mizhou0817"
	password = "12345678"
	dbname   = "fund_test"
)

// create connection with postgres db
func createConnection() *sql.DB {

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname)

	// Open the connection
	conn, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Fail to load the connection config.  %v", err)
	}

	// check the connection
	err = conn.Ping()

	if err != nil {
		log.Fatalf("Unable to connect psql.  %v", err)
	}

	// return the connection
	return conn
}

// CreateUser create a user in the postgres db
func CreateUser(w http.ResponseWriter, r *http.Request) {

	// create an empty user of type models.User
	var user models.User

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	insertID := insertUser(user)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetUser will return a single user by its id
func GetUser(w http.ResponseWriter, r *http.Request) {
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getUser function with user id to retrieve a single user
	user, err := getUser(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(user)
}

// GetAllUser will return all the users
func GetAllUser(w http.ResponseWriter, r *http.Request) {

	// get all the users in the db
	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	sortBy := r.URL.Query().Get("sortBy")

	if len(sortBy) > 0 {
		sort.Slice(users, func(i, j int) bool {
			if sortBy == "desc" {
				if users[i].Name == users[j].Name {
					return users[i].Age > users[j].Age
				} else {
					return users[i].Name > users[j].Name
				}
			} else {
				if users[i].Name == users[j].Name {
					return users[i].Age < users[j].Age
				} else {
					return users[i].Name < users[j].Name
				}
			}
		})
	}

	// send all the users as response
	json.NewEncoder(w).Encode(users)
}

// UpdateUser update user's detail in the postgres db
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty user of type models.User
	var user models.User

	// decode the json request to user
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update user to update the user
	updatedRows := updateUser(int64(id), user)

	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// ------------------------- handler functions ----------------
// insert one user in the DB
func insertUser(user models.User) int64 {

	// create the postgres db connection
	conn := createConnection()

	// close the db connection
	defer conn.Close()

	// create the dynamic insert sql query
	// returning id will return the id of the inserted user
	sqlStatement := `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := conn.QueryRow(sqlStatement, user.Name, user.Email, user.Age).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// return the inserted id
	return id
}

// get one user from the DB by its userid
func getUser(id int64) (models.User, error) {
	// create the postgres db connection
	conn := createConnection()

	// close the db connection
	defer conn.Close()

	// create a user of models.User type
	var user models.User

	// create the dynamical select sql query
	sqlStatement := `SELECT * FROM users WHERE id=$1`

	// execute the sql statement
	row := conn.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
	err := row.Scan(&user.Age, &user.Email, &user.Name, &user.ID)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return user, err
}

// get one user from the DB by its userid
func getAllUsers() ([]models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var users []models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var user models.User

		// unmarshal the row object to user
		err = rows.Scan(&user.Age, &user.Email, &user.Name, &user.ID)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		users = append(users, user)

	}

	// return empty user on error
	return users, err
}

// update user in the DB
func updateUser(id int64, user models.User) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE users SET name=$2, email=$3, age=$4 WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, user.Name, user.Email, user.Age)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
