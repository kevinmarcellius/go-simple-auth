package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kevinmarcellius/go-simple-auth/internal/model"
)

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func GenerateJWT(user model.User, jwtKey string) (string, string, error) {
	// Generate access token
	accessToken, err := generateAccessTokenFromUser(user, jwtKey)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken(user, jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateAccessTokenFromUser(user model.User, jwtKey string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func generateRefreshToken(user model.User, jwtKey string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func ValidateRefreshToken(tokenString string, jwtKey string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}

func GenerateNewAccessToken(claims *Claims, jwtKey string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
