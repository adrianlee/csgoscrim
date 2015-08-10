package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	User_id    int       `json:"user_id"`
	Username   *string   `json:"username"`
	First_name *string   `json:"first_name"`
	Last_name  *string   `json:"last_name"`
	Steam_id   *string   `json:"steam_id"`
	Email      *string   `json:"email"`
	Created_on time.Time `json:"created_on"`
	Last_login time.Time `json:"last_login"`
}

func (u User) SayHello() string {
	// return "Hi, I am " + u.username + " my name is " + u.first_name
	return ""
}

type NullString struct {
	sql.NullString
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	s.String = strings.Trim(string(data), `"`)
	s.Valid = true
	return nil
}

var (
	Okay            = errors.New("200: OK")
	ErrUnauthorized = errors.New("401: Unauthorized")
	ErrNotFound     = errors.New("404: Not found")
	ErrValidation   = errors.New("422: Validation error")
	ErrInternal     = errors.New("500: Something went wrong on our end")
)

const (
	DB_user = "csgo"
	DB_pass = "123123123"
	DB_host = "csgo.cwiirrozzw4c.us-west-2.rds.amazonaws.com:5432/csgoscrim"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://"+DB_user+":"+DB_pass+"@"+DB_host)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/login", LoginHandler)

	// Resources
	r.HandleFunc("/users", UsersListHandler).Methods("GET")
	r.HandleFunc("/users/{id}", UsersGetHandler).Methods("GET")
	r.HandleFunc("/users/{id}", RootHandler).Methods("PUT")

	r.HandleFunc("/teams", RootHandler).Methods("GET")
	r.HandleFunc("/matches", RootHandler)
	r.HandleFunc("/tournaments", RootHandler)

	// Static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Login page")
}

func UsersListHandler(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.Username, &user.Steam_id, &user.First_name, &user.Last_name, &user.Email, &user.Created_on, &user.Last_login, &user.User_id)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		users = append(users, user)
	}

	b, err := json.Marshal(users)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%s", b)
}

func UsersGetHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	// QueryRow always return a non-nil value
	row := db.QueryRow("SELECT * FROM users WHERE user_id = $1", id)

	var user User
	err := row.Scan(&user.Username, &user.Steam_id, &user.First_name, &user.Last_name, &user.Email, &user.Created_on, &user.Last_login, &user.User_id)

	// Check for errors
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, http.StatusText(404), 404)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Marshal results
	b, err := json.Marshal(user)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%s", b)
}
