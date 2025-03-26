package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	password := "password123"
	
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotNil(t, hashedPassword)

	matching := CheckPassword(hashedPassword, password)
	assert.True(t, matching)
}