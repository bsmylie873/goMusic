package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"goMusic/db"
	"goMusic/models"
	"goMusic/utils"
	viewModelAlbum "goMusic/viewModels"
	"net/http"
	"time"
)

func GetAlbums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.DB.Query("SELECT id, title, price, artist_id, band_id FROM albums")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var albums []models.Album

	for rows.Next() {
		var album models.Album
		var artistID, bandID sql.NullInt64

		err := rows.Scan(&album.ID, &album.Title, &album.Price, &artistID, &bandID)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if artistID.Valid {
			id := int(artistID.Int64)
			album.ArtistID = &id
		}

		if bandID.Valid {
			id := int(bandID.Int64)
			album.BandID = &id
		}

		albums = append(albums, album)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	albumVMs, err := viewModelAlbum.GetAlbumViewModels(albums)
	if err != nil {
		http.Error(w, "Error generating view models: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(albumVMs)
}

func GetAlbumByID(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	row, err := db.DB.Query("SELECT id, title, price, artist_id, band_id FROM albums WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if row.Next() {
		var album models.Album
		var artistID, bandID sql.NullInt64

		err := row.Scan(&album.ID, &album.Title, &album.Price, &artistID, &bandID)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if artistID.Valid {
			id := int(artistID.Int64)
			album.ArtistID = &id
		}

		if bandID.Valid {
			id := int(bandID.Int64)
			album.BandID = &id
		}

		albumVM, err := viewModelAlbum.GetAlbumViewModel(album)
		if err != nil {
			http.Error(w, "Error generating view model: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(albumVM)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "album not found"})
}

func PostAlbum(w http.ResponseWriter, r *http.Request) {
	var newAlbum models.Album

	if !utils.DecodeAndValidate(w, r, &newAlbum) {
		return
	}

	_, err := db.DB.Exec("INSERT INTO albums (title, price, artist_id, band_id) VALUES (?, ?, ?, ?)",
		newAlbum.Title, newAlbum.Price, newAlbum.ArtistID, newAlbum.BandID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func UpdateAlbumByID(w http.ResponseWriter, r *http.Request, id int) bool {
	var updatedAlbum models.Album

	if !utils.DecodeAndValidate(w, r, &updatedAlbum) {
		return false
	}

	_, err := db.DB.Exec("UPDATE albums SET title = ?, price = ?, artist_id = ?, band_id = ? WHERE id = ?",
		updatedAlbum.Title, updatedAlbum.Price, updatedAlbum.ArtistID, updatedAlbum.BandID, id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return false
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedAlbum)
	return true
}

func DeleteAlbumByID(w http.ResponseWriter, r *http.Request, id int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return false
	}

	_, err = tx.Exec("DELETE FROM albums WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return false
	}
	w.WriteHeader(http.StatusNoContent)
	return true
}
