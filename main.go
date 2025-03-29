package main

import (
	_ "github.com/lib/pq"

	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/magicznykacpur/psst-backend/api"
	"github.com/magicznykacpur/psst-backend/env"
	"github.com/magicznykacpur/psst-backend/internal/database"
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

	mux.HandleFunc("POST /api/validity", apiConfig.HandlerCheckTokenValidity)

	mux.HandleFunc("GET /api/users", apiConfig.HandlerGetUsers)
	mux.HandleFunc("GET /api/users/{id}", apiConfig.HandlerGetUser)
	mux.HandleFunc("POST /api/users", apiConfig.HandlerCreateUser)

	mux.Handle("GET /api/users/me", apiConfig.IsLoggedIn(apiConfig.HandlerGetMe))
	mux.Handle("GET /api/users/to-chat-with", apiConfig.IsLoggedIn(apiConfig.HandlerGetUsersAvailableToChatWith))

	mux.HandleFunc("POST /api/login", apiConfig.HandlerLoginUser)

	mux.Handle("POST /api/messages", apiConfig.IsLoggedIn(apiConfig.HandlerCreateMessage))
	mux.Handle("GET /api/messages/{chat_id}", apiConfig.IsLoggedIn(apiConfig.HandlerGetAllMessagesFromChat))
	mux.Handle("DELETE /api/messages", apiConfig.IsLoggedIn(apiConfig.HandlerDeleteMessageFromChat))

	mux.Handle("POST /api/chats", apiConfig.IsLoggedIn(apiConfig.HandlerCreateChat))
	mux.Handle("GET /api/chats", apiConfig.IsLoggedIn(apiConfig.HandlerGetAllUsersChats))
	mux.Handle("DELETE /api/chats/{id}", apiConfig.IsLoggedIn(apiConfig.HandlerDeleteChat))

	server := http.Server{Handler: &mux, Addr: "localhost:" + apiConfig.Port}

	log.Printf("starting `psst` server at %s", server.Addr)
	server.ListenAndServe()
}
