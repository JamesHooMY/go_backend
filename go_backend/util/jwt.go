package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(ID uint, name string) (token string, err error) {
	timeNow := time.Now()
	claims := Claims{
		ID:   ID,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(timeNow),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(viper.GetString("jwt.secretKey")))
}

func ParseJwtToken(token string) (claims Claims, err error) {
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secretKey")), nil
	})
	if err != nil {
		return Claims{}, err
	}

	if !jwtToken.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return claims, nil
}
