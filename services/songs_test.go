package services_test

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"goMusic/db"
	"goMusic/models"
	"goMusic/services"
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

	rows := sqlmock.NewRows([]string{"id", "title", "length", "price"}).
		AddRow(1, "Yellow", 431, 1.29)

	mock.ExpectQuery("SELECT id, title, length, price FROM songs").
		WillReturnRows(rows)

	albumRows := sqlmock.NewRows([]string{"id", "title", "price"}).
		AddRow(1, "Parachutes", 29.99)

	mock.ExpectQuery("SELECT a.id, a.title, a.price FROM albums a JOIN album_songs sa ON a.id = sa.album_id WHERE sa.song_id = ?").
		WithArgs(1).
		WillReturnRows(albumRows)

	artistRows := sqlmock.NewRows([]string{"id", "first_name", "last_name"}).
		AddRow(1, "Chris", "Martin")

	mock.ExpectQuery("SELECT a.id, a.first_name, a.last_name FROM artists a JOIN artist_songs sa ON a.id = sa.artist_id LEFT JOIN sexes s ON a.sex_id = s.id LEFT JOIN titles t ON a.title_id = t.id WHERE sa.song_id = ?").
		WithArgs(1).
		WillReturnRows(artistRows)

	bandRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Coldplay")

	mock.ExpectQuery("SELECT b.id, b.name FROM bands b JOIN band_songs sb ON b.id = sb.band_id WHERE sb.song_id = ?").
		WithArgs(1).
		WillReturnRows(bandRows)

	req, err := http.NewRequest("GET", "/songs", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	services.GetSongs(rr, req)

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
		rows := sqlmock.NewRows([]string{"id", "title", "length", "price"}).
			AddRow(1, "Yellow", 431, 1.29)

		mock.ExpectQuery("SELECT id, title, length, price FROM songs WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		albumRows := sqlmock.NewRows([]string{"id", "title", "price", "artist_id", "band_id"}).
			AddRow(1, "Album Title", 9.99, 2, 3)

		mock.ExpectQuery("SELECT a.id, a.title, a.price, a.artist_id, a.band_id FROM albums a JOIN album_songs sa ON a.id = sa.album_id WHERE sa.song_id = ?").
			WithArgs(1).
			WillReturnRows(albumRows)

		artistRows := sqlmock.NewRows([]string{"first_name", "last_name"}).
			AddRow("Chris", "Martin")

		mock.ExpectQuery("SELECT a.first_name, a.last_name FROM artists a WHERE a.id=?").
			WithArgs(2).
			WillReturnRows(artistRows)

		artistRows2 := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nationality", "birth_date", "age", "alive", "sex_id", "title_id", "band_id"}).
			AddRow(1, "Chris", "Martin", "British", "1977-03-02", 44, true, 1, 1, nil)

		mock.ExpectQuery("SELECT a.* FROM artists a JOIN artist_songs sa ON a.id = sa.artist_id WHERE sa.song_id = ?").
			WithArgs(1).
			WillReturnRows(artistRows2)

		sexRows := sqlmock.NewRows([]string{"name"}).
			AddRow("Male")

		mock.ExpectQuery("SELECT name FROM sexes WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sexRows)

		titleRows := sqlmock.NewRows([]string{"name"}).
			AddRow("Mr")

		mock.ExpectQuery("SELECT name FROM titles WHERE id = ?").
			WithArgs(1).
			WillReturnRows(titleRows)

		bandRows := sqlmock.NewRows([]string{"id", "name", "nationality", "number_of_members", "date_formed", "age", "active"}).
			AddRow(1, "Coldplay", "British", 4, "1996-01-16", 25, true)

		mock.ExpectQuery("SELECT b.* FROM bands b JOIN band_songs sb ON b.id = sb.band_id WHERE sb.song_id = ?").
			WithArgs(1).
			WillReturnRows(bandRows)

		_, err := http.NewRequest("GET", "/songs/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		services.GetSongByID(rr, 1)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Song not found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "length", "price"})

		mock.ExpectQuery("SELECT id, title, length, price FROM songs WHERE id = ?").
			WithArgs(999).
			WillReturnRows(rows)

		_, err := http.NewRequest("GET", "/songs/999", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		services.GetSongByID(rr, 999)

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

	song := models.Song{
		Title:  "Yellow",
		Length: 431,
		Price:  1.29,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO songs").
		WithArgs(song.Title, song.Length, song.Price, nil, nil, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	songJSON, _ := json.Marshal(song)
	req, err := http.NewRequest("POST", "/songs", bytes.NewBuffer(songJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	services.PostSong(rr, req)

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

	song := models.Song{
		Title:  "Yellow (2023 Remix)",
		Length: 445,
		Price:  1.49,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE songs SET").
		WithArgs(song.Title, song.Length, song.Price, nil, nil, nil, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	songJSON, _ := json.Marshal(song)
	req, err := http.NewRequest("PUT", "/songs/1", bytes.NewBuffer(songJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	result := services.UpdateSongByID(rr, req, 1)

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

	result := services.DeleteSongByID(rr, 1)

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
