package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeyLength = 32

// JWTMaker is a JSON web token maker
type JWTMaker struct {
	Secret string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeyLength)
	}
	return &JWTMaker{
		Secret: secretKey,
	}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.Secret))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) ValidateToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, invalidTokenErr
		}
		return []byte(maker.Secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, expiredTokenErr) {
			return nil, expiredTokenErr
		}
		return nil, invalidTokenErr
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, expiredTokenErr
	}

	return payload, nil
}
