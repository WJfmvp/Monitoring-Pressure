package JWT

import (
	"errors"
	"time"

	"Monitoring-Pressure/models/users"

	jwtgo "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("Monitoring-Pressure-secret-key")

type MyClaims struct {
	UserID    int64          `json:"user_id"`
	Telephone string         `json:"telephone"`
	Role      users.UserRole `json:"role"`
	jwtgo.RegisteredClaims
}

func GenerateToken(userID int64, telephone string, role users.UserRole) (string, error) {
	claims := MyClaims{
		UserID:    userID,
		Telephone: telephone,
		Role:      role,
		RegisteredClaims: jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwtgo.NewNumericDate(time.Now()),
			NotBefore: jwtgo.NewNumericDate(time.Now()),
			Issuer:    "Monitoring-Pressure",
			Subject:   telephone,
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwtgo.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*MyClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
