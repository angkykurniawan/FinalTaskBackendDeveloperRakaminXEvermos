package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey_evermos_rakamin_2025"
	}
	return []byte(secret)
}

func GenerateJWT(userID uint, role string) (string, error) {
	claims := JWTClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

func ValidateToken(tokenString string) (uint, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return getSecret(), nil
		},
	)

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return 0, "", errors.New("token otorisasi tidak valid")
	}

	return claims.UserID, claims.Role, nil
}
