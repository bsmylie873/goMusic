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

	bandRows := sqlmock.NewRows([]string{"name", "nationality", "number_of_members", "date_formed", "age", "active"}).
		AddRow("Coldplay", "British", 4, "1996-01-01", 27, true)

	mock.ExpectQuery("SELECT name, nationality, number_of_members, date_formed, age, active FROM bands WHERE id = ?").
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

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAlbumByID(t *testing.T) {
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

	t.Run("Album found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "price", "artist_id", "band_id"}).
			AddRow(1, "Parachutes", 9.99, nil, 1)
		mock.ExpectQuery("SELECT id, title, price, artist_id, band_id FROM albums WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		bandRows := sqlmock.NewRows([]string{"name", "nationality", "number_of_members", "date_formed", "age", "active"}).
			AddRow("Coldplay", "British", 4, "1996-01-01", 27, true)

		mock.ExpectQuery("SELECT name, nationality, number_of_members, date_formed, age, active FROM bands WHERE id = ?").
			WithArgs(1).
			WillReturnRows(bandRows)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/albums/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		GetAlbumByID(rr, req, 1)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Album not found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "price", "artist_id", "band_id"})
		mock.ExpectQuery("SELECT id, title, price, artist_id, band_id FROM albums WHERE id = ?").
			WithArgs(999).
			WillReturnRows(rows)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/albums/999", nil)
		if err != nil {
			t.Fatal(err)
		}

		GetAlbumByID(rr, req, 999)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
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
		ID:       1,
		Title:    "Parachutes",
		Price:    9.99,
		BandID:   &bandID,
		ArtistID: nil,
	}

	albumJSON, _ := json.Marshal(album)
	req, err := http.NewRequest("POST", "/albums", bytes.NewBuffer(albumJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mock.ExpectExec("INSERT INTO albums").
		WithArgs(album.ID, album.Title, album.Price, album.ArtistID, album.BandID).
		WillReturnResult(sqlmock.NewResult(1, 1))

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

		mock.ExpectExec("UPDATE bands SET").
			WithArgs(band.Name, band.Nationality, band.NumberOfMembers, band.DateFormed, band.Age, band.Active, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

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
		req, err := http.NewRequest("DELETE", "/bands/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		mock.ExpectExec("DELETE FROM albums WHERE id = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		result := DeleteAlbumByID(rr, req, 1)

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
