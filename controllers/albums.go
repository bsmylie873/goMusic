package controllers

import (
	"goMusic/authentication"
	"goMusic/services"
	"net/http"
	"strconv"
)

func RegisterAlbumRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /albums", services.GetAlbums)
	mux.HandleFunc("GET /albums/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		services.GetAlbumByID(w, r, id)
	})

	mux.HandleFunc("POST /albums", authentication.AuthMiddleware(services.PostAlbum))
	mux.HandleFunc("PUT /albums/{id}", authentication.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.Atoi(r.PathValue("id"))
			services.UpdateAlbumByID(w, r, id)
		},
	))
	mux.HandleFunc("DELETE /albums/{id}", authentication.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.Atoi(r.PathValue("id"))
			services.DeleteAlbumByID(w, r, id)
		},
	))
}
