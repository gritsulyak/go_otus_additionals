package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gritsulyak/go_otus_additionals/go_mic/004_github_gitlab_demo/libs/logger"
	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler(db))
	mux.HandleFunc("/users", usersHandler(db))
	mux.HandleFunc("/users/", userByIDHandler(db))
	return mux
}

// GET /health
func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// GET /users       — список всех
// POST /users      — создать пользователя
func usersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listUsers(w, r, db)
		case http.MethodPost:
			createUser(w, r, db)
		default:
			http.Error(w, "not (Get or Post) method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// GET    /users/{id}  — получить по ID
// DELETE /users/{id}  — удалить по ID
func userByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {

			http.Error(w, "invalid id for /users/{id}", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			getUser(w, r, db, id)
		case http.MethodDelete:
			deleteUser(w, r, db, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func listUsers(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
	logger.Debug("listing users")
	rows, err := db.Query(`SELECT id, name FROM users ORDER BY id`)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("listUsers query: %v", err)
		return
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "rows error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func createUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		http.Error(w, "invalid body: name required", http.StatusBadRequest)
		return
	}

	var u User
	err := db.QueryRow(
		`INSERT INTO users (name) VALUES ($1) RETURNING id, name`, input.Name,
	).Scan(&u.ID, &u.Name)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("createUser insert: %v", err)
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func getUser(w http.ResponseWriter, _ *http.Request, db *sql.DB, id int) {
	var u User
	err := db.QueryRow(`SELECT id, name FROM users WHERE id = $1`, id).Scan(&u.ID, &u.Name)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("getUser query: %v", err)
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func deleteUser(w http.ResponseWriter, _ *http.Request, db *sql.DB, id int) {
	res, err := db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("deleteUser exec: %v", err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}

func main() {
	logger.Info("API starting")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error(err.Error())
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, NewRouter(db)))
}
