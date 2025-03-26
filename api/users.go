package api

import (
	"log"
	"net/http"
	"time"
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
