package token

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

var testMaker Maker

func createMaker(t *testing.T) {
	maker, err := NewTokenMaker("eLWbxMhrsM1681390917098627557735")
	require.NoError(t, err)

	testMaker = maker
}

func createToken(t *testing.T, userID uint, duration time.Duration) (*Payload, string) {
	createMaker(t)

	require.NotEmpty(t, testMaker)

	token, payload, err := testMaker.GenToken(userID, duration)
	require.NoError(t, err)
	require.Equal(t, userID, payload.UserID)
	require.WithinDuration(t, payload.ExpiresAt.Time, time.Now().Add(duration), 2*time.Second)
	fmt.Println(token)
	parts := strings.Split(token, ".")
	require.True(t, len(parts) == 3)

	return payload, token
}

func TestGenToken(t *testing.T) {
	createToken(t, 100, time.Minute)
}

func TestVerifyToken(t *testing.T) {
	payload, token := createToken(t, 100, time.Minute)
	payload2, err := testMaker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, payload, payload2)
}

func TestExpiredToken(t *testing.T) {
	_, token := createToken(t, 100, -time.Minute)
	_, err := testMaker.VerifyToken(token)
	fmt.Println(err)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)
}
