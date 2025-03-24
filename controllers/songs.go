package controllers

import (
	"goMusic/authentication"
	"goMusic/services"
	"net/http"
	"strconv"
)

func RegisterSongRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /songs", services.GetSongs)
	mux.HandleFunc("GET /songs/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		services.GetSongByID(w, id)
	})

	mux.HandleFunc("POST /songs", authentication.AuthMiddleware(services.PostSong))
	mux.HandleFunc("PUT /songs/{id}", authentication.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.Atoi(r.PathValue("id"))
			services.UpdateSongByID(w, r, id)
		},
	))
	mux.HandleFunc("DELETE /songs/{id}", authentication.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.Atoi(r.PathValue("id"))
			services.DeleteSongByID(w, id)
		},
	))
}
