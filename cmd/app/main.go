package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User1 struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	slog.Info("Welcome to Notein Global Server!")

	var (
		dbDriver = "mysql"
		dbSource = "my-test:noteinin@tcp(notein-debug-405205:us-central1:my-test-go-sql)/notein_user"
	)

	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		slog.Error("main", "connecting to MySQL error", err.Error())
	}
	defer db.Close()
	log.Println("connected to MySQL successfully.")

	err = db.Ping()
	if err != nil {
		slog.Error("main", "database ping error", err.Error())
	}
	log.Println("db ping successfully, Just test 2, 3")

	uh := userHandler{
		ctx: context.Background(),
		db:  db,
	}

	// err = uh.createDataFileAndTable()
	// if err != nil {
	// 	slog.Error("main", "createDataFileAndTable error", err.Error())
	// }
	// log.Println("create database file and table successfully.")

	http.Handle("/users", uh)
	http.ListenAndServe(":8080", nil)
}

type userHandler struct {
	ctx context.Context
	db  *sql.DB
}

func (uh userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uh.getUsers(w, r)
	case "POST":
		uh.createUser(w, r)
	default:
	}
}

var (
	insertUserQuery = "INSERT INTO users (id, name) VALUES (?, ?);"
	getUserQuery    = "SELECT id, name from users;"
)

func (uh userHandler) createDataFileAndTable() error {
	createDBQuery := "CREATE DATABASE IF NOT EXISTS notein_user;"
	_, err := uh.db.ExecContext(uh.ctx, createDBQuery)
	if err != nil {
		return err
	}

	useDBQuery := "USE notein_user;"
	_, err = uh.db.ExecContext(uh.ctx, useDBQuery)
	if err != nil {
		return err
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id INT PRIMARY KEY,
			name VARCHAR(255)
		);
	`
	_, err = uh.db.ExecContext(uh.ctx, createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func (uh userHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var u User1
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	stmt, err := uh.db.PrepareContext(uh.ctx, insertUserQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(uh.ctx, &u.ID, &u.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user added"))
}

func (uh userHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := uh.db.QueryContext(uh.ctx, getUserQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var users []User1
	for rows.Next() {
		var u User1
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(users)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
