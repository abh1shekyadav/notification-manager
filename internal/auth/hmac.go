package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type HMACValidator struct {
	Secret []byte
}

func NewHMACValidator(secret string) *HMACValidator {
	return &HMACValidator{Secret: []byte(secret)}
}

func (v *HMACValidator) ValidateToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return v.Secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
