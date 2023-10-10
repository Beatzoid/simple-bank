package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	// Hash random password
	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	// Hashed password should be the same as the original password
	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	// A random string should not be the same as the hashed password
	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// Hash random password
	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)

	// Two different hashed passwords should not be the same
	require.NotEqual(t, hashedPassword1, hashedPassword2)

}
