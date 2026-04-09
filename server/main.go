package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mauroloaiza/groovy/server/internal/db"
	"github.com/mauroloaiza/groovy/server/internal/library"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()

	musicDir := os.Getenv("MUSIC_DIR")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok","service":"groovy"}`)
	})

	r.Mount("/api/library", library.Router(database, musicDir))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Groovy API  →  http://localhost:%s", port)
	log.Printf("DB          →  %s", dbPath())
	log.Printf("Music dir   →  %s", musicDir)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func dbPath() string {
	if p := os.Getenv("DB_PATH"); p != "" {
		return p
	}
	return "groovy.db"
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
