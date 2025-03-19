package controllers

import (
	"goMusic/services"
	"net/http"
	"regexp"
	"strconv"
)

func RegisterAlbumRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/albums", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetAlbums(w, r)
		case http.MethodPost:
			PostAlbum(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	albumIDPattern := regexp.MustCompile(`^/albums/(.+)$`)
	mux.HandleFunc("/albums/", func(w http.ResponseWriter, r *http.Request) {
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
			GetAlbumByID(w, r, num)
		case http.MethodPut:
			UpdateAlbumByID(w, r, num)
		case http.MethodDelete:
			DeleteAlbumByID(w, r, num)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func GetAlbums(w http.ResponseWriter, r *http.Request) {
	services.GetAlbums(w, r)
}

func GetAlbumByID(w http.ResponseWriter, r *http.Request, id int) {
	services.GetAlbumByID(w, r, id)
}

func PostAlbum(w http.ResponseWriter, r *http.Request) {
	services.PostAlbum(w, r)
}

func UpdateAlbumByID(w http.ResponseWriter, r *http.Request, id int) {
	services.UpdateAlbumByID(w, r, id)
}

func DeleteAlbumByID(w http.ResponseWriter, r *http.Request, id int) {
	services.DeleteAlbumByID(w, r, id)
}
