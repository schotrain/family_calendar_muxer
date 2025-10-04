package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type FamilyCalendarClaims struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	jwt.RegisteredClaims
}

func GenerateFamilyCalendarJWT(email, givenName, familyName string) (string, error) {
	claims := FamilyCalendarClaims{
		Email:      email,
		GivenName:  givenName,
		FamilyName: familyName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "family-calendar-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}
