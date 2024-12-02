package lib

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret []byte

func InitSecret() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatalf("JWT_SECRET is not set")
	}
	jwtSecret = []byte(secret)
}

func GenerateJWT(email string, password string) (*string, error) {
	claim := &jwt.MapClaims{
		"email":    email,
		"password": password,
		"exp":      time.Now().AddDate(0, 0, 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	singnedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Fatalf("Error signing token %v", err)
		return nil, err
	}
	return &singnedToken, nil
}

func ValidateJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("token is expired")
			}
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
