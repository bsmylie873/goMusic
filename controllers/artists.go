package controllers

import (
	"goMusic/authentication"
	"goMusic/services"
	"net/http"
	"strconv"
)

func RegisterArtistRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /artists", services.GetArtists)
	mux.HandleFunc("GET /artists/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		services.GetArtistByID(w, r, id)
	})

	mux.HandleFunc("POST /artists", authentication.AuthMiddleware(services.PostArtist))
	mux.HandleFunc("PUT /artists/{id}", authentication.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.Atoi(r.PathValue("id"))
			services.UpdateArtistByID(w, r, id)
		},
	))
	mux.HandleFunc("DELETE /artists/{id}", authentication.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.Atoi(r.PathValue("id"))
			services.DeleteArtistByID(w, r, id)
		},
	))
}
