package controllers

import (
	"goMusic/authentication"
	"goMusic/services"
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		RegisterUser(w, r)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		LoginUser(w, r)
	})

	mux.HandleFunc("/profile", authentication.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		GetProfile(w, r)
	}))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	services.RegisterUser(w, r)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	services.LoginUser(w, r)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	services.GetProfile(w, r)
}
