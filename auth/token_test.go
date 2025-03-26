package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBearerToken(t *testing.T) {
	// test header withouth authorization
	headerNoAuth := http.Header{}
	_, err := GetBearerToken(headerNoAuth)
	require.Error(t, err)
	assert.Equal(t, "authorization header not provided", err.Error())

	// test header malformed
	headerMalformed := http.Header{}
	headerMalformed.Set("Authorization", "Bearer token something")
	_, err = GetBearerToken(headerMalformed)
	require.Error(t, err)
	assert.Equal(t, "authorization header malformed", err.Error())

	// test bearer malformed
	headerBearerMalformed := http.Header{}
	headerBearerMalformed.Set("Authorization", "Bearerrr token")
	_, err = GetBearerToken(headerBearerMalformed)
	require.Error(t, err)
	assert.Equal(t, "bearer token malformed", err.Error())

	// test bearer valid
	headerBearerValid := http.Header{}
	headerBearerValid.Set("Authorization", "Bearer myBearerToken")
	bearerToken, err := GetBearerToken(headerBearerValid)
	require.NoError(t, err)
	assert.Equal(t, "myBearerToken", bearerToken)
}
