package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(secretKey string, expire_seconds int, id int) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirps",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire_seconds) * time.Second)),
			Subject:   strconv.Itoa(id),
		},
	)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
