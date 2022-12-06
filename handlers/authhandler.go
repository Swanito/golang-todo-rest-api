package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func ValidateToken(headers http.Header) jwt.MapClaims {
	if headers["Authorization"] != nil {
		token, _ := jwt.Parse(headers["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if token != nil && token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			return claims
		} else {
			return nil
		}
	} else {
		return nil
	}
}
