package utils

import (
	"softeng-platform/internal/config"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtSecret   string
	jwtSecretOnce sync.Once
)

// initJWTSecret 初始化JWT密钥，只加载一次
func initJWTSecret() {
	jwtSecretOnce.Do(func() {
		cfg := config.LoadConfig()
		jwtSecret = cfg.JWTSecret
	})
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int, username, role string) (string, error) {
	initJWTSecret()

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	initJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
