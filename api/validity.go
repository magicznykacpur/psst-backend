package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/magicznykacpur/psst-backend/auth"
)

type checkTokenRequest struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) HandlerCheckTokenValidity(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request bytes")
		return
	}

	var tokenRequest checkTokenRequest
	err = json.Unmarshal(requestBytes, &tokenRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal request")
		return
	}

	id, err := auth.ValidateJWT(tokenRequest.Token, os.Getenv("JWT_SECRET"))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token invalid")
		return
	}

	_, err = cfg.DB.GetUserByID(r.Context(), id)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		respondWithError(w, http.StatusUnauthorized, "user doesn't exist")
		return
	}

	w.WriteHeader(http.StatusOK)
}
