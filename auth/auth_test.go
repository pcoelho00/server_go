package auth

import (
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateToken(t *testing.T) {
	// Set up the test data
	secretKey := "mySecretKey"
	expireSeconds := 3600
	id := 123

	// Call the CreateToken function
	token, err := CreateToken(secretKey, expireSeconds, id)
	if err != nil {
		t.Fatal(err)
	}

	// Parse the token to verify its validity
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check the token claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		t.Errorf("Token is invalid")
	}

	// Check the issuer claim
	issuer, ok := claims["iss"].(string)
	if !ok || issuer != "chirps" {
		t.Errorf("Expected issuer claim to be 'chirps', but got %q", issuer)
	}

	// Check the issued at claim
	issuedAt, ok := claims["iat"].(float64)
	if !ok || issuedAt == 0 {
		t.Errorf("Expected issued at claim to be a valid timestamp, but got %v", issuedAt)
	}

	// Check the expires at claim
	expiresAt, ok := claims["exp"].(float64)
	if !ok || expiresAt == 0 {
		t.Errorf("Expected expires at claim to be a valid timestamp, but got %v", expiresAt)
	}

	// Check the subject claim
	subject, ok := claims["sub"].(string)
	if !ok || subject != strconv.Itoa(id) {
		t.Errorf("Expected subject claim to be %q, but got %q", strconv.Itoa(id), subject)
	}
}
func TestCreateRefreshToken(t *testing.T) {
	length := 32

	// Call the CreateRefreshToken function
	token, err := CreateRefreshToken(length)
	if err != nil {
		t.Fatal(err)
	}

	// Check the token length
	if len(token) != length*2 {
		t.Errorf("Expected token length to be %d, but got %d", length*2, len(token))
	}
}
