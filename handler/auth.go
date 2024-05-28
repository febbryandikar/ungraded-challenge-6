package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"ungraded-challenge-6/entity"
	"ungraded-challenge-6/token"
	"ungraded-challenge-6/utility"
)

type NewAuthHandler struct {
	*sql.DB
}

func (h *NewAuthHandler) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Error while parsing body",
		})
		log.Println("Error while parsing body:", err)
		return
	}

	valid, errMessage := utility.ValidateUser(user)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: errMessage,
		})
		log.Println("Error while parsing body:", err)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while hashing password",
		})
		log.Println("Error while hashing password:", err)
		return
	}
	user.Password = string(hashedPass)

	var exists bool
	err = h.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`, user.Email).Scan(&exists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while checking user",
		})
		log.Println("Error while checking user:", err)
		return
	}

	if exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "User already exists",
		})
		log.Println("User already exists")
		return
	}

	_, err = h.Exec(`INSERT INTO users (email, password, full_name, age, occupation, role) VALUES (?, ?, ?, ?, ?, ?)`, user.Email, user.Password, user.FullName, user.Age, user.Occupation, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while inserting user",
		})
		log.Println("Error while inserting user:", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusCreated,
		Message: "User created successfully",
	})
}

func (h *NewAuthHandler) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var login entity.User
	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Error while parsing body",
		})
		log.Println("Error while parsing body:", err)
		return
	}

	row, err := h.Query(`SELECT email, password, full_name, age, occupation, role FROM users WHERE email = ?`, login.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while checking user",
		})
		log.Println("Error while checking user:", err)
		return
	}
	defer row.Close()

	found := false
	for row.Next() {
		found = true

		err = row.Scan(&user.Email, &user.Password, &user.FullName, &user.Age, &user.Occupation, &user.Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.Message{
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Message: "Error while checking user",
			})
			log.Println("Error while checking user:", err)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.Message{
				Status:  "failed",
				Code:    http.StatusUnauthorized,
				Message: "Passwords do not match",
			})
			log.Println("Password do not match")
			return
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "User not found",
		})
		log.Println("User not found")
		return
	}

	tokenString, err := token.GenerateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while generating token",
		})
		log.Println("Error while generating token:", err)
		return
	}

	data := map[string]string{
		"token": tokenString,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusCreated,
		Message: "User login successfully",
		Data:    data,
	})
}
