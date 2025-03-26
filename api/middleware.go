package api

import (
	"net/http"
	"os"

	"github.com/magicznykacpur/psst-backend/auth"
)

func (cfg *ApiConfig) IsLoggedIn(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userId, err := auth.ValidateJWT(bearerToken, os.Getenv("JWT_SECRET"))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		
		user, err := cfg.DB.GetUserByID(r.Context(), userId)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "user does not exist")
			return
		}

		r.Header.Set("User-ID", user.ID.String())

		next.ServeHTTP(w, r)
	})
}