package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateJWTToken(id uuid.UUID, secret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "psst",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			Subject:   id.String(),
		})

	return token.SignedString([]byte(secret))
}

func ValidateJWT(token, secret string) (uuid.UUID, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	userId, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userId, nil
}
