package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(secret string, userId uint, expireIn time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(expireIn).Unix()
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseJWT(secret, tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenInvalidType
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
