package services

import (
	"encoding/json"
	"goMusic/authentication"
	"goMusic/db"
	"goMusic/models"
	"goMusic/viewModels"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req viewModels.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	result, err := db.DB.Exec(
		"INSERT INTO users (username, password, email) VALUES (?, ?, ?)",
		req.Username, string(hashedPassword), req.Email,
	)
	if err != nil {
		http.Error(w, "Username or email already exists", http.StatusBadRequest)
		return
	}

	id, _ := result.LastInsertId()

	token, err := authentication.GenerateToken(int(id))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Id:       int(id),
		Username: req.Username,
		Email:    req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(viewModels.AuthResponse{
		Token: token,
		User:  user,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req viewModels.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	var hashedPassword string
	err = db.DB.QueryRow(
		"SELECT id, username, password, email FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.Id, &user.Username, &hashedPassword, &user.Email)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := authentication.GenerateToken(user.Id)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(viewModels.AuthResponse{
		Token: token,
		User:  user,
	})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Invalid user context", http.StatusInternalServerError)
		return
	}

	var user models.User
	err := db.DB.QueryRow(
		"SELECT id, username, email FROM users WHERE id = ?",
		userID,
	).Scan(&user.Id, &user.Username, &user.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
