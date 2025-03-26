package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/magicznykacpur/psst-backend/auth"
	"github.com/magicznykacpur/psst-backend/internal/database"
)

type userResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
}

func (cfg *ApiConfig) HandlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.DB.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve users")
		return
	}

	usersResponse := []userResponse{}
	for _, user := range users {
		usersResponse = append(usersResponse,
			userResponse{
				ID:        user.ID.String(),
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
				Username:  user.UserName,
			},
		)
	}

	log.Println("returning users...")
	respondWithJSON(w, http.StatusOK, usersResponse)
}

func (cfg *ApiConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	userId, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "id malformed")
		return
	}

	user, err := cfg.DB.GetUserByID(r.Context(), userId)
	if err != nil && err.Error() == "sql: no rows in result set" {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil && err.Error() != "sql: no rows in result set" {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := userResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Username:  user.UserName,
	}

	log.Println("returning user...")
	respondWithJSON(w, http.StatusOK, response)
}

type createUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"user_name"`
	Password string `json:"password"`
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request body bytes")
		return
	}

	var createUserRequest createUserRequest
	err = json.Unmarshal(requestBytes, &createUserRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal request body")
		return
	}

	if createUserRequest.Email == "" || createUserRequest.Username == "" || createUserRequest.Password == "" {
		respondWithError(w, http.StatusBadRequest, "request body invalid, must contain: email, username and password")
		return
	}

	hashedPassword, err := auth.HashPassword(createUserRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(),
		database.CreateUserParams{
			Email:          createUserRequest.Email,
			UserName:       createUserRequest.Username,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil && strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
		respondWithError(w, http.StatusBadRequest, "user with that email or username already exists")
		return
	}

	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	log.Printf("user with email: %s, and user name: %s created...", user.Email, user.UserName)

	respondWithJSON(w, http.StatusCreated,
		userResponse{
			ID:        user.ID.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.UserName,
		},
	)
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request body bytes")
		return
	}

	var loginUserRequest loginUserRequest
	err = json.Unmarshal(requestBytes, &loginUserRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal request body")
		return
	}

	if loginUserRequest.Email == "" || loginUserRequest.Password == "" {
		respondWithError(w, http.StatusBadRequest, "request body invalid, must contain: email and password")
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), loginUserRequest.Email)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		respondWithError(w, http.StatusUnauthorized, "email or password incorrect")
		return
	}

	if err != nil && !strings.Contains(err.Error(), "sql: no rows in result set") {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	passwordCorrect := auth.CheckPassword(user.HashedPassword, loginUserRequest.Password)
	if !passwordCorrect {
		respondWithError(w, http.StatusUnauthorized, "email or password incorrect")
		return
	}

	token, err := auth.CreateJWTToken(user.ID, os.Getenv("JWT_SECRET"), time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create jwt token")
		return
	}

	log.Printf("user %s logged in...", user.UserName)
	respondWithJSON(w, http.StatusOK, loginUserResponse{Token: token})
}
