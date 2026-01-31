package service

import (
	"time"

	"github.com/CashInvoice-Golang-Assignment/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *models.User, secret string, expiryMinutes int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Duration(expiryMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
