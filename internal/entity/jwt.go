package entity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const JWTExpire = time.Hour * 3

type UnexpectedSigningMethodError struct {
	alg string
}

func (e *UnexpectedSigningMethodError) Error() string {
	return fmt.Sprintf("unexpected signing method: %s", e.alg)
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserUUID string
}

func BuildJWTString(userUUID string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTExpire)),
		},
		UserUUID: userUUID,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWTString(tokenString string, secret []byte) (*jwt.Token, *JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if alg := token.Method.Alg(); alg != jwt.SigningMethodHS256.Name {
			return nil, &UnexpectedSigningMethodError{alg}
		}

		return secret, nil
	})
	if err != nil {
		return nil, nil, err
	}

	return token, claims, nil
}
