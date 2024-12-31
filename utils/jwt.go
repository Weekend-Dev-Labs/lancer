package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type SessionClaims struct {
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type AdminClaims struct {
	AdminID string `json:"uid"`
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

func GetSessionInfo(token string, key string) (*SessionClaims, error) {
	var claims SessionClaims

	withClaims, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !withClaims.Valid {
		return nil, fmt.Errorf("invalid session token")
	}

	return &claims, nil
}

func GetAdminToken(claims *AdminClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

func GetAdminInfo(token string, key string) (*AdminClaims, error) {
	var claims AdminClaims

	withClaims, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !withClaims.Valid {
		return nil, fmt.Errorf("invalid admin token")
	}

	return &claims, nil
}
