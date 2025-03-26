package main

import (
	_ "github.com/lib/pq"

	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/magicznykacpur/psst-backend/api"
	"github.com/magicznykacpur/psst-backend/internal/database"
	"github.com/magicznykacpur/psst-backend/internal/env"
)

func main() {
	env.LoadDotEnv()

	sqlDB, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("couldn't open database: %v", err)
	}

	apiConfig := api.ApiConfig{
		DB:   database.New(sqlDB),
		Port: os.Getenv("PORT"),
	}

	mux := http.ServeMux{}

	mux.HandleFunc("GET /api/users", apiConfig.HandlerGetUsers)
	mux.HandleFunc("GET /api/users/{id}", apiConfig.HandlerGetUser)
	mux.HandleFunc("POST /api/users", apiConfig.HandlerCreateUser)

	server := http.Server{Handler: &mux, Addr: "localhost:" + apiConfig.Port}

	log.Printf("starting `psst` server at %s", server.Addr)
	server.ListenAndServe()
}
