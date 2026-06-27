package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Manager struct {
	secret []byte
}

func New(secret string) *Manager {
	return &Manager{
		secret: []byte(secret),
	}
}

func (m *Manager) Generate(userID uuid.UUID) (string, error) {

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)

	return token.SignedString(m.secret)
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {

	token, err := jwtlib.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwtlib.Token) (any, error) {

			if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return m.secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
