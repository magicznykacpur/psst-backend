package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTValid(t *testing.T) {
	id, _ := uuid.Parse("eaa8cf1d-869b-4d97-9a14-320a0edf7853")
	secret := "my-secret-secret"
	expiresIn := time.Second

	token, err := CreateJWTToken(id, secret, expiresIn)
	require.NoError(t, err)

	retrievedId, err := ValidateJWT(token, secret)
	require.NoError(t, err)
	assert.Equal(t, id, retrievedId)
}

func TestJWTInvalidSecret(t *testing.T) {
	id, _ := uuid.Parse("eaa8cf1d-869b-4d97-9a14-320a0edf7853")
	secret := "my-secret-secret"
	expiresIn := time.Second

	token, err := CreateJWTToken(id, secret, expiresIn)
	require.NoError(t, err)

	_, err = ValidateJWT(token, secret + "-invalid")
	require.Error(t, err)
	assert.Equal(t, "token signature is invalid: signature is invalid", err.Error())

}

func TestJWTExpired(t *testing.T) {
	id, _ := uuid.Parse("eaa8cf1d-869b-4d97-9a14-320a0edf7853")
	secret := "my-secret-secret"
	expiresIn := time.Millisecond

	token, err := CreateJWTToken(id, secret, expiresIn)
	require.NoError(t, err)

	_, err = ValidateJWT(token, secret)
	require.Error(t, err)
	assert.Equal(t, "token has invalid claims: token is expired", err.Error())
}