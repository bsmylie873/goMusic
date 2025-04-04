package services

import (
	"encoding/json"
	"goMusic/db"
	"goMusic/models"
	"goMusic/utils"
	viewModelBand "goMusic/viewModels"
	"net/http"
)

func GetBands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.DB.Query("SELECT id, name, nationality, number_of_members, date_formed, age, active FROM bands")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bands []models.Band

	for rows.Next() {
		var band models.Band

		err := rows.Scan(&band.Id, &band.Name, &band.Nationality, &band.NumberOfMembers, &band.DateFormed, &band.Age, &band.Active)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		bands = append(bands, band)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bandVMs, err := viewModelBand.GetBandViewModels(bands)
	if err != nil {
		http.Error(w, "Error generating view models: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bandVMs)
}

func GetBandByID(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	row, err := db.DB.Query("SELECT id, name, nationality, number_of_members, date_formed, age, active FROM bands WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if row.Next() {
		var band models.Band

		err := row.Scan(&band.Id, &band.Name, &band.Nationality, &band.NumberOfMembers, &band.DateFormed, &band.Age, &band.Active)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		bandVM, err := viewModelBand.GetBandViewModel(band)
		if err != nil {
			http.Error(w, "Error generating view model: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(bandVM)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "band not found"})
}

func PostBand(w http.ResponseWriter, r *http.Request) {
	var newBand models.Band

	if !utils.DecodeAndValidate(w, r, &newBand) {
		return
	}
	success := utils.ExecuteWithTransaction(w,
		"INSERT INTO bands (name, nationality, number_of_members, date_formed, age, active) VALUES (?, ?, ?, ?, ?, ?)",
		newBand.Name, newBand.Nationality, newBand.NumberOfMembers, newBand.DateFormed, newBand.Age, newBand.Active)

	if !success {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func UpdateBandByID(w http.ResponseWriter, r *http.Request, id int) bool {
	var updatedBand models.Band

	if !utils.DecodeAndValidate(w, r, &updatedBand) {
		return false
	}

	success := utils.ExecuteWithTransaction(w,
		"UPDATE bands SET name = ?, nationality = ?, number_of_members = ?, date_formed = ?, age = ?, active = ? WHERE id = ?",
		updatedBand.Name, updatedBand.Nationality, updatedBand.NumberOfMembers, updatedBand.DateFormed, updatedBand.Age, updatedBand.Active, id)

	if !success {
		return false
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBand)
	return true
}

func DeleteBandByID(w http.ResponseWriter, id int) bool {
	success := utils.ExecuteWithTransaction(w,
		"DELETE FROM bands WHERE id = ?",
		id)

	if !success {
		return false
	}

	w.WriteHeader(http.StatusNoContent)
	return true
}
