package testing

import (
	"encoding/base64"
	"test-task/pkg/utils"
	"testing"
	"time"

	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	userID := "6ab58fc3-6920-48a0-8851-a2f0650fa2a5"
	ipAddress := "192.168.0.1"
	jwtSecretKey := "test-secret-key"

	token, err := utils.GenerateAccessToken(userID, ipAddress, jwtSecretKey)
	assert.NoError(t, err, "expected no error when generating access token")
	assert.NotEmpty(t, token, "expected token to be non-empty")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil
	})
	assert.NoError(t, err, "expected no error when parsing token")
	assert.True(t, parsedToken.Valid, "expected token to be valid")

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok, "expected claims to be of type MapClaims")
	assert.Equal(t, userID, claims["user_id"], "expected user_id to match")
	assert.Equal(t, ipAddress, claims["ip"], "expected ip to match")

	exp := claims["exp"].(float64)
	assert.True(t, exp > float64(time.Now().Unix()), "expected expiration time to be in the future")
}

func TestGenerateRefreshToken(t *testing.T) {
	token, err := utils.GenerateRefreshToken()
	assert.NoError(t, err, "expected no error when generating refresh token")
	assert.NotEmpty(t, token, "expected token to be non-empty")

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	assert.NoError(t, err, "expected no error when decoding token")
	assert.Equal(t, 32, len(decodedToken), "expected decoded token length to be 32 bytes")
}

func TestExtractUserIDFromToken(t *testing.T) {
	userID := "6ab58fc3-6920-48a0-8851-a2f0650fa2a5"
	ipAddress := "192.168.0.1"
	jwtSecretKey := "test-secret-key"

	token, err := utils.GenerateAccessToken(userID, ipAddress, jwtSecretKey)
	assert.NoError(t, err, "expected no error when generating access token")

	extractedUserID, err := utils.ExtractUserIDFromToken(token, jwtSecretKey)
	assert.NoError(t, err, "expected no error when extracting user_id from token")
	assert.Equal(t, userID, extractedUserID, "expected extracted user_id to match")
}

func TestConvertStringToUUID(t *testing.T) {
	validUUID := "6ab58fc3-6920-48a0-8851-a2f0650fa2a5"
	invalidUUID := "invalid-uuid-string"

	uuidVal, err := utils.ConvertStringToUUID(validUUID)
	assert.NoError(t, err, "expected no error for valid UUID string")
	assert.Equal(t, validUUID, uuidVal.String(), "expected UUID to match the input string")

	_, err = utils.ConvertStringToUUID(invalidUUID)
	assert.Error(t, err, "expected error for invalid UUID string")
	assert.EqualError(t, err, "invalid UUID string", "expected error message to be 'invalid UUID string'")
}
