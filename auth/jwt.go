package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// Claims mirrors the payload issued by the backend's /api/auth/login.
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}

// ParseToken validates the token stored in the browser cookie and returns
// its claims, so handlers/templates can show "logged in as ..." and the
// role-based navigation without calling the API again.
func ParseToken(secret, tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("auth: unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("auth: invalid token")
	}
	return claims, nil
}
