package controllers

import (
	"goMusic/services"
	"net/http"
	"regexp"
	"strconv"
)

func RegisterBandRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/bands", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetBands(w, r)
		case http.MethodPost:
			PostBand(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	bandIDPattern := regexp.MustCompile(`^/bands/(.+)$`)
	mux.HandleFunc("/bands/", func(w http.ResponseWriter, r *http.Request) {
		matches := bandIDPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil || len(matches) < 2 {
			http.NotFound(w, r)
			return
		}
		id := matches[1]
		num, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid band ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			GetBandByID(w, r, num)
		case http.MethodPut:
			UpdateBandByID(w, r, num)
		case http.MethodDelete:
			DeleteBandByID(w, r, num)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func GetBands(w http.ResponseWriter, r *http.Request) {
	services.GetBands(w, r)
}

func GetBandByID(w http.ResponseWriter, r *http.Request, id int) {
	services.GetBandByID(w, r, id)
}

func PostBand(w http.ResponseWriter, r *http.Request) {
	services.PostBand(w, r)
}

func UpdateBandByID(w http.ResponseWriter, r *http.Request, id int) {
	services.UpdateBandByID(w, r, id)
}

func DeleteBandByID(w http.ResponseWriter, r *http.Request, id int) {
	services.DeleteBandByID(w, r, id)
}
