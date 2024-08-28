package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID, ipAddress, jwtSecretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"ip":      ipAddress,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	secretKey := []byte(jwtSecretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(tokenBytes), nil
}

func ExtractUserIDFromToken(tokenString, jwtSecretKey string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return userID, nil
}
