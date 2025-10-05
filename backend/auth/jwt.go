package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type FamilyCalendarClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateFamilyCalendarJWT(userID uint) (string, error) {
	claims := FamilyCalendarClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "family-calendar-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}
