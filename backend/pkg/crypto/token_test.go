package crypto

import (
	"ocserv-bakend/pkg/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	secret := "my-secret-key"
	config.Set(&config.Config{JWTSecret: secret})
	userID := "12345"
	expire := time.Now().Add(time.Hour).Unix()

	tokenString, err := GenerateAccessToken(userID, expire, true)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims["sub"])
	assert.Equal(t, true, claims["isAdmin"])
}
