package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect DB
	dsn := os.Getenv("DATABASE_URL")
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	} else {
		log.Println("✅ Successfully connected to database")
	}

	// Router
	r := chi.NewRouter()

	r.Get("/notes", getNotes)
	r.Post("/notes", createNote)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.QueryContext(context.Background(), "SELECT id, title FROM notes")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		notes = append(notes, n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func createNote(w http.ResponseWriter, r *http.Request) {
	var n Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err := db.QueryRowContext(context.Background(),
		"INSERT INTO notes (title) VALUES ($1) RETURNING id", n.Title).Scan(&n.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}
