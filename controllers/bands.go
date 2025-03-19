package controllers

import (
	"goMusic/authentication"
	"goMusic/services"
	"net/http"
	"strconv"
)

func RegisterBandRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /bands", services.GetBands)
	mux.HandleFunc("GET /bands/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		services.GetBandByID(w, r, id)
	})

	mux.HandleFunc("POST /bands", authentication.AuthMiddleware(services.PostBand))
	mux.HandleFunc("PUT /bands/{id}", authentication.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.Atoi(r.PathValue("id"))
			services.UpdateBandByID(w, r, id)
		},
	))
	mux.HandleFunc("DELETE /bands/{id}", authentication.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.Atoi(r.PathValue("id"))
			services.DeleteBandByID(w, r, id)
		},
	))
}
