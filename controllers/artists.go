package controllers

import (
	"goMusic/services"
	"net/http"
	"regexp"
	"strconv"
)

func RegisterArtistRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/artists", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetArtists(w, r)
		case http.MethodPost:
			PostArtist(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	albumIDPattern := regexp.MustCompile(`^/artists/(.+)$`)
	mux.HandleFunc("/artists/", func(w http.ResponseWriter, r *http.Request) {
		matches := albumIDPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil || len(matches) < 2 {
			http.NotFound(w, r)
			return
		}
		id := matches[1]
		num, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid album ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			GetArtistByID(w, r, num)
		case http.MethodPut:
			UpdateArtistByID(w, r, num)
		case http.MethodDelete:
			DeleteArtistByID(w, r, num)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func GetArtists(w http.ResponseWriter, r *http.Request) {
	services.GetArtists(w, r)
}

func GetArtistByID(w http.ResponseWriter, r *http.Request, id int) {
	services.GetArtistByID(w, r, id)
}

func PostArtist(w http.ResponseWriter, r *http.Request) {
	services.PostArtist(w, r)
}

func UpdateArtistByID(w http.ResponseWriter, r *http.Request, id int) {
	services.UpdateArtistByID(w, r, id)
}

func DeleteArtistByID(w http.ResponseWriter, r *http.Request, id int) {
	services.DeleteArtistByID(w, r, id)
}
