package services

import (
	"encoding/json"
	"goMusic/db"
	"goMusic/models"
	"goMusic/utils"
	viewModelSong "goMusic/viewModels"
	"net/http"
)

func GetSongs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.DB.Query("SELECT id, title, length, price FROM songs")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var songs []models.Song

	for rows.Next() {
		var song models.Song

		err := rows.Scan(&song.Id, &song.Title, &song.Length, &song.Price)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	songVMs, err := viewModelSong.GetSongViewModels(songs)
	if err != nil {
		http.Error(w, "Error generating view models: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(songVMs)
}

func GetSongByID(w http.ResponseWriter, id int) {
	w.Header().Set("Content-Type", "application/json")
	row, err := db.DB.Query("SELECT id, title, length, price FROM songs WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if row.Next() {
		var song models.Song

		err := row.Scan(&song.Id, &song.Title, &song.Length, &song.Price)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		songVM, err := viewModelSong.GetSongViewModel(song)
		if err != nil {
			http.Error(w, "Error generating view model: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(songVM)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "band not found"})
}

func PostSong(w http.ResponseWriter, r *http.Request) {
	var newSong models.Song

	if !utils.DecodeAndValidate(w, r, &newSong) {
		return
	}

	_, err := db.DB.Exec("INSERT INTO songs (title, length, price, album_id, artist_id, band_id) VALUES (?, ?, ?, ?, ?, ?)",
		newSong.Title, newSong.Length, newSong.Price, newSong.AlbumId, newSong.ArtistId, newSong.BandId)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func UpdateSongByID(w http.ResponseWriter, r *http.Request, id int) bool {
	var updatedSong models.Song

	if !utils.DecodeAndValidate(w, r, &updatedSong) {
		return false
	}

	success := utils.ExecuteWithTransaction(w,
		"UPDATE songs SET title = ?, length = ?, price = ?, artist_id = ?, album_id = ?, band_id = ? WHERE id = ?",
		updatedSong.Title, updatedSong.Length, updatedSong.Price,
		updatedSong.ArtistId, updatedSong.AlbumId, updatedSong.BandId, id)

	if !success {
		return false
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSong)
	return true
}

func DeleteSongByID(w http.ResponseWriter, id int) bool {
	success := utils.ExecuteWithTransaction(w,
		"DELETE FROM songs WHERE id = ?",
		id)

	if !success {
		return false
	}

	w.WriteHeader(http.StatusNoContent)
	return true
}
