package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims структура для JWT токена
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword - хеширование пароля с солью
func HashPassword(password string) (hash string, salt string, err error) {
	// Генерация уникальной соли
	saltBytes := make([]byte, 32)
	_, err = rand.Read(saltBytes)
	if err != nil {
		return "", "", err
	}
	salt = base64.StdEncoding.EncodeToString(saltBytes)

	// Комбинируем пароль с солью
	saltedPassword := password + salt

	// Хеширование с помощью bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	hash = string(hashedBytes)

	return hash, salt, nil
}

// CheckPassword - проверка пароля
func CheckPassword(password, hash, salt string) bool {
	saltedPassword := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltedPassword))
	return err == nil
}

// GenerateAccessToken - генерация Access Token (JWT)
func GenerateAccessToken(userID uint, email, role, secret string, expiration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "paving-tiles-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateAccessToken - валидация Access Token
func ValidateAccessToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// HashToken - хеширование токена для хранения в БД
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GenerateRefreshToken - генерация случайного Refresh Token
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
