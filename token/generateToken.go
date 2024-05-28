package token

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
	"ungraded-challenge-6/entity"
	"ungraded-challenge-6/middleware"
)

func GenerateToken(user entity.User) (string, error) {
	claims := &middleware.Claims{
		user.Email,
		user.FullName,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
