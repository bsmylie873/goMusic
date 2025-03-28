package services

import (
	"encoding/json"
	"goMusic/db"
	"goMusic/models"
	"goMusic/utils"
	viewModelArtist "goMusic/viewModels"
	"net/http"
)

func GetArtists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.DB.Query("SELECT id, first_name, last_name, nationality, birth_date, age, alive, sex_id, title_id, band_id FROM artists")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var artists []models.Artist

	for rows.Next() {
		var artist models.Artist

		err := rows.Scan(&artist.Id, &artist.FirstName, &artist.LastName, &artist.Nationality, &artist.BirthDate, &artist.Age, &artist.Alive, &artist.SexId, &artist.TitleId, &artist.BandId)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		artists = append(artists, artist)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	artistVMs, err := viewModelArtist.GetArtistViewModels(artists)
	if err != nil {
		http.Error(w, "Error generating view models: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(artistVMs)
}

func GetArtistByID(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	row, err := db.DB.Query("SELECT id, first_name, last_name, nationality, birth_date, age, alive, sex_id, title_id, band_id FROM artists WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if row.Next() {
		var artist models.Artist

		err := row.Scan(&artist.Id, &artist.FirstName, &artist.LastName, &artist.Nationality, &artist.BirthDate, &artist.Age, &artist.Alive, &artist.SexId, &artist.TitleId, &artist.BandId)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		artistVM, err := viewModelArtist.GetArtistViewModel(artist)
		if err != nil {
			http.Error(w, "Error generating view model: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(artistVM)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "band not found"})
}

func PostArtist(w http.ResponseWriter, r *http.Request) {
	var newArtist models.Artist

	if !utils.DecodeAndValidate(w, r, &newArtist) {
		return
	}
	success := utils.ExecuteWithTransaction(w,
		"INSERT INTO artists (first_name, last_name, nationality, birth_date, age, alive, sex_id, title_id, band_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		newArtist.FirstName, newArtist.LastName, newArtist.Nationality, newArtist.BirthDate, newArtist.Age, newArtist.Alive, newArtist.SexId, newArtist.TitleId, newArtist.BandId)

	if !success {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func UpdateArtistByID(w http.ResponseWriter, r *http.Request, id int) bool {
	var updatedArtist models.Artist

	if !utils.DecodeAndValidate(w, r, &updatedArtist) {
		return false
	}

	success := utils.ExecuteWithTransaction(w,
		"UPDATE artists SET first_name = ?, last_name = ?, nationality = ?, birth_date = ?, age = ?, alive = ?, sex_id = ?, title_id = ?, band_id = ? WHERE id = ?",
		updatedArtist.FirstName, updatedArtist.LastName, updatedArtist.Nationality, updatedArtist.BirthDate, updatedArtist.Age, updatedArtist.Alive, updatedArtist.SexId, updatedArtist.TitleId, updatedArtist.BandId, id)

	if !success {
		return false
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedArtist)
	return true
}

func DeleteArtistByID(w http.ResponseWriter, id int) bool {
	success := utils.ExecuteWithTransaction(w,
		"DELETE FROM artists WHERE id = ?",
		id)

	if !success {
		return false
	}

	w.WriteHeader(http.StatusNoContent)
	return true
}
