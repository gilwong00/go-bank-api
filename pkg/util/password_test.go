package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashAndValidatePassword(t *testing.T) {
	password := GetRandomString(6)
	hashedPassword1, err := Hashpassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = ValidatePassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := GetRandomString(6)
	err = ValidatePassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// check if the same password is hashed twice, 2 different hashes should be produced
	hashedPassword2, err := Hashpassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
