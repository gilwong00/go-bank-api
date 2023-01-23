package db

import (
	"context"
	"go-bank-api/pkg/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.GetRandomOwner(),
		HashedPassword: "secret",
		FirstName:      util.GetRandomOwner(),
		LastName:       util.GetRandomOwner(),
		Email:          util.GetRandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := createTestUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Email, user2.Email)
}
