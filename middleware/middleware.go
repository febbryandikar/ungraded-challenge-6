package middleware

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"ungraded-challenge-6/entity"
)

type Claims struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(requiredRole string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.Message{
				Status:  "failed",
				Code:    http.StatusUnauthorized,
				Message: "Missing Authorization header",
			})
			log.Println("Missing Authorization header")
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})

		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.Message{
				Status:  "failed",
				Code:    http.StatusUnauthorized,
				Message: "Invalid Authorization header",
			})
			log.Println("Invalid Authorization header")
			return
		}

		if requiredRole != "" && claims.Role != requiredRole {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.Message{
				Status:  "failed",
				Code:    http.StatusUnauthorized,
				Message: "Invalid Authorization header",
			})
			log.Println("Invalid Authorization header")
			return
		}

		next(w, r, p)
	}
}
