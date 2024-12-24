package utils

import "github.com/golang-jwt/jwt/v5"

type SessionClaims struct {
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

func GetSessionToken(claims *SessionClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
