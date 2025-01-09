package token

import (
	"os"
	"time"

	"examples.com/auth-service/custom_errors"
	"github.com/golang-jwt/jwt/v5"
)

var SessionKey = []byte(os.Getenv("SESSION_KEY"))

func CreateToken(login string) (string, error) {
	claims := jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SessionKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SessionKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return custom_errors.ErrNotValidToken
	}
	return nil
}
