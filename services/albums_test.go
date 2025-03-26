package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"goMusic/db"
	"goMusic/models"
	"goMusic/viewModels"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAlbums(t *testing.T) {
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

	rows := sqlmock.NewRows([]string{"id", "title", "price", "artist_id", "band_id"}).
		AddRow(1, "Parachutes", 9.99, nil, 1)

	mock.ExpectQuery("SELECT id, title, price, artist_id, band_id FROM albums").
		WillReturnRows(rows)

	bandRows := sqlmock.NewRows([]string{"name"}).
		AddRow("Coldplay")

	mock.ExpectQuery("SELECT name FROM bands WHERE id = ?").
		WithArgs(1).
		WillReturnRows(bandRows)

	req, err := http.NewRequest("GET", "/albums", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	GetAlbums(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var albums []viewModels.AlbumViewModel
	if err := json.Unmarshal(rr.Body.Bytes(), &albums); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(albums) != 1 || albums[0].Title != "Parachutes" {
		t.Errorf("Wrong album data: got %+v", albums)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAlbumByID(t *testing.T) {
	originalDB := db.DB
	defer func() {
		db.DB = originalDB
	}()

	t.Run("Album found with band", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer mockDB.Close()
		db.DB = mockDB

		albumRows := sqlmock.NewRows([]string{"id", "title", "price", "artist_id", "band_id"}).
			AddRow(1, "Parachutes", 9.99, nil, 1)
		mock.ExpectQuery("SELECT id, title, price, artist_id, band_id FROM albums WHERE id = ?").
			WithArgs(1).
			WillReturnRows(albumRows)

		bandRows := sqlmock.NewRows([]string{"name", "nationality", "number_of_members", "date_formed", "age", "active"}).
			AddRow("Coldplay", "British", 4, "1996-01-01", 27, true)
		mock.ExpectQuery("SELECT name, nationality, number_of_members, date_formed, age, active FROM bands WHERE id = ?").
			WithArgs(1).
			WillReturnRows(bandRows)

		songRows := sqlmock.NewRows([]string{"id", "title", "length", "price"})
		mock.ExpectQuery("SELECT id, title, length, price FROM songs").
			WillReturnRows(songRows)

		req := httptest.NewRequest("GET", "/albums/1", nil)
		res := httptest.NewRecorder()

		GetAlbumByID(res, req, 1)

		if status := res.Code; status != http.StatusOK {
			t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
		}

		var album viewModels.DetailedAlbumViewModel
		if err := json.Unmarshal(res.Body.Bytes(), &album); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if album.Title != "Parachutes" || album.Price != 9.99 {
			t.Errorf("Wrong album data: got %+v", album)
		}

		if album.Band == nil || album.Band.Name != "Coldplay" {
			t.Errorf("Missing or incorrect band data: %+v", album.Band)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})

	t.Run("Album found with artist", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer mockDB.Close()
		db.DB = mockDB

		artistID := 2
		albumRows := sqlmock.NewRows([]string{"id", "title", "price", "artist_id", "band_id"}).
			AddRow(2, "Kind of Blue", 12.99, artistID, nil)
		mock.ExpectQuery("SELECT id, title, price, artist_id, band_id FROM albums WHERE id = ?").
			WithArgs(2).
			WillReturnRows(albumRows)

		artistRows := sqlmock.NewRows([]string{
			"first_name", "last_name", "nationality", "birth_date",
			"age", "alive", "sex", "title", "band_name",
		}).AddRow(
			"Miles", "Davis", "American", "1926-05-26",
			65, false, "Male", "Mr.", nil,
		)
		mock.ExpectQuery("SELECT a.first_name, a.last_name, a.nationality").
			WithArgs(artistID).
			WillReturnRows(artistRows)

		songRows := sqlmock.NewRows([]string{"id", "title", "length", "price"})
		mock.ExpectQuery("SELECT id, title, length, price FROM songs").
			WillReturnRows(songRows)

		req := httptest.NewRequest("GET", "/albums/2", nil)
		res := httptest.NewRecorder()

		GetAlbumByID(res, req, 2)

		if status := res.Code; status != http.StatusOK {
			t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
		}

		var album viewModels.DetailedAlbumViewModel
		if err := json.Unmarshal(res.Body.Bytes(), &album); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if album.Title != "Kind of Blue" || album.Artist == nil {
			t.Errorf("Wrong album data: got %+v", album)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})

	t.Run("Album not found", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer mockDB.Close()
		db.DB = mockDB

		albumRows := sqlmock.NewRows([]string{"id", "title", "price", "artist_id", "band_id"})
		mock.ExpectQuery("SELECT id, title, price, artist_id, band_id FROM albums WHERE id = ?").
			WithArgs(999).
			WillReturnRows(albumRows)

		req := httptest.NewRequest("GET", "/albums/999", nil)
		res := httptest.NewRecorder()

		GetAlbumByID(res, req, 999)

		if status := res.Code; status != http.StatusNotFound {
			t.Errorf("Wrong status code: got %v want %v", status, http.StatusNotFound)
		}

		var response map[string]string
		if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if message, exists := response["message"]; !exists || message != "album not found" {
			t.Errorf("Wrong error message: got %v", response)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})

	t.Run("Database error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer mockDB.Close()
		db.DB = mockDB

		// Return error on query
		mock.ExpectQuery("SELECT id, title, price, artist_id, band_id FROM albums WHERE id = ?").
			WithArgs(1).
			WillReturnError(errors.New("database connection lost"))

		req := httptest.NewRequest("GET", "/albums/1", nil)
		res := httptest.NewRecorder()

		GetAlbumByID(res, req, 1)

		if status := res.Code; status != http.StatusInternalServerError {
			t.Errorf("Wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})
}

func TestPostAlbum(t *testing.T) {
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

	bandID := 1
	album := models.Album{
		Title:    "Parachutes",
		Price:    9.99,
		BandId:   &bandID,
		ArtistId: nil,
	}

	albumJSON, _ := json.Marshal(album)
	req, err := http.NewRequest("POST", "/albums", bytes.NewBuffer(albumJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO albums").
		WithArgs(album.Title, album.Price, album.ArtistId, album.BandId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	PostAlbum(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateAlbumByID(t *testing.T) {
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

	t.Run("Successful update", func(t *testing.T) {
		band := models.Band{
			Name:            "Coldplay",
			Nationality:     "British",
			NumberOfMembers: 4,
			DateFormed:      "1996-01-01",
			Age:             28,
			Active:          true,
		}

		bandJSON, _ := json.Marshal(band)
		req, err := http.NewRequest("PUT", "/bands/1", bytes.NewBuffer(bandJSON))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE bands SET").
			WithArgs(band.Name, band.Nationality, band.NumberOfMembers, band.DateFormed, band.Age, band.Active, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		result := UpdateBandByID(rr, req, 1)

		if !result {
			t.Errorf("UpdateBandByID returned false, expected true")
		}

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestDeleteAlbumByID(t *testing.T) {
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

	t.Run("Successful delete", func(t *testing.T) {
		rr := httptest.NewRecorder()

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM albums WHERE id = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		result := DeleteAlbumByID(rr, 1)

		if !result {
			t.Errorf("DeleteAlbumByID returned false, expected true")
		}

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
