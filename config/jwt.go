package config

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("sua-chave-secreta")

func GenerateToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
