package services

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"goMusic/db"
	"goMusic/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSongs(t *testing.T) {
	originalDB := db.DB
	defer func() {
		db.DB = originalDB
	}()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db.DB = mockDB

	rows := sqlmock.NewRows([]string{"id", "title", "length", "price", "album_id"}).
		AddRow(1, "Yellow", 4.31, 1.29, 1)

	mock.ExpectQuery("SELECT id, title, length, price, album_id FROM songs").
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/songs", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	GetSongs(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSongByID(t *testing.T) {
	originalDB := db.DB
	defer func() {
		db.DB = originalDB
	}()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db.DB = mockDB

	t.Run("Song found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "length", "price", "album_id"}).
			AddRow(1, "Yellow", 4.31, 1.29, 1)

		mock.ExpectQuery("SELECT id, title, length, price, album_id FROM songs WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		albumRows := sqlmock.NewRows([]string{"id", "title", "price"}).
			AddRow(1, "Parachutes", 9.99)

		mock.ExpectQuery("SELECT id, title, price FROM albums WHERE id = ?").
			WithArgs(1).
			WillReturnRows(albumRows)

		_, err := http.NewRequest("GET", "/songs/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		GetSongByID(rr, 1)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Song not found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "length", "price", "album_id"})

		mock.ExpectQuery("SELECT id, title, length, price, album_id FROM songs WHERE id = ?").
			WithArgs(999).
			WillReturnRows(rows)

		_, err := http.NewRequest("GET", "/songs/999", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		GetSongByID(rr, 999)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestPostSong(t *testing.T) {
	originalDB := db.DB
	defer func() {
		db.DB = originalDB
	}()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db.DB = mockDB

	albumId := 1
	song := models.Song{
		Title:   "Yellow",
		Length:  431,
		Price:   1.29,
		AlbumId: &albumId,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO songs").
		WithArgs(song.Title, song.Length, song.Price, song.AlbumId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	songJSON, _ := json.Marshal(song)
	req, err := http.NewRequest("POST", "/songs", bytes.NewBuffer(songJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	PostSong(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateSongByID(t *testing.T) {
	originalDB := db.DB
	defer func() {
		db.DB = originalDB
	}()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db.DB = mockDB

	albumId := 1
	song := models.Song{
		Title:   "Yellow (2023 Remix)",
		Length:  445,
		Price:   1.49,
		AlbumId: &albumId,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE songs SET").
		WithArgs(song.Title, song.Length, song.Price, song.AlbumId, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	songJSON, _ := json.Marshal(song)
	req, err := http.NewRequest("PUT", "/songs/1", bytes.NewBuffer(songJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	result := UpdateSongByID(rr, req, 1)

	if !result {
		t.Errorf("UpdateSongByID returned false, expected true")
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteSongByID(t *testing.T) {
	originalDB := db.DB
	defer func() {
		db.DB = originalDB
	}()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db.DB = mockDB

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM songs WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	rr := httptest.NewRecorder()

	result := DeleteSongByID(rr, 1)

	if !result {
		t.Errorf("DeleteSongByID returned false, expected true")
	}

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
