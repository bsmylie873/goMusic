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

func TestGetArtists(t *testing.T) {
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

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nationality", "birth_date", "age", "alive", "sex_id", "title_id", "band_id"}).
		AddRow(1, "Ed", "Sheeran", "British", "1991-02-17", 32, true, nil, nil, nil)

	mock.ExpectQuery("SELECT id, first_name, alst_name, nationality, birth_date, age, alive, sex_id, title_id, band_id FROM artists").
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/artists", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	GetArtists(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetArtistByID(t *testing.T) {
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

	t.Run("Artist found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nationality", "birth_date", "age", "alive", "sex_id", "title_id", "band_id"}).
			AddRow(1, "Ed", "Sheeran", "British", "1991-02-17", 32, true, nil, nil, nil)

		mock.ExpectQuery("SELECT id, first_name, last_name, nationality, birth_date, age, alive, sex_id, title_id, band_id FROM artists WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/artists/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		GetArtistByID(rr, req, 1)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Artist not found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nationality", "birth_date", "age", "alive", "sex_id", "title_id", "band_id"})

		mock.ExpectQuery("SELECT id, first_name, last_name, nationality, birth_date, age, alive, sex_id, title_id, band_id FROM artists WHERE id = ?").
			WithArgs(999).
			WillReturnRows(rows)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/artists/999", nil)
		if err != nil {
			t.Fatal(err)
		}

		GetArtistByID(rr, req, 999)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestPostArtist(t *testing.T) {
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

	artist := models.Artist{
		Id:          1,
		FirstName:   "Ed",
		LastName:    "Sheeran",
		Nationality: "British",
		BirthDate:   "1991-02-17",
		Age:         32,
		Alive:       true,
	}

	artistJSON, _ := json.Marshal(artist)
	req, err := http.NewRequest("POST", "/artists", bytes.NewBuffer(artistJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mock.ExpectExec("INSERT INTO artist").
		WithArgs(artist.Id, artist.FirstName, artist.LastName, artist.Nationality, artist.BirthDate, artist.Age, artist.Alive, artist.SexId, artist.TitleId, artist.BandId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	PostArtist(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateArtistByID(t *testing.T) {
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
		artist := models.Artist{
			FirstName:   "Ed",
			LastName:    "Sheeran",
			Nationality: "British",
			BirthDate:   "1991-02-17",
			Age:         32,
			Alive:       true,
		}

		artistJSON, _ := json.Marshal(artist)
		req, err := http.NewRequest("PUT", "/artists/1", bytes.NewBuffer(artistJSON))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mock.ExpectExec("UPDATE artists SET").
			WithArgs(artist.FirstName, artist.LastName, artist.Nationality, artist.BirthDate, artist.Age, artist.Alive, artist.SexId, artist.TitleId, artist.BandId, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		result := UpdateArtistByID(rr, req, 1)

		if !result {
			t.Errorf("UpdateArtistByID returned false, expected true")
		}

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestDeleteArtistByID(t *testing.T) {
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
		req, err := http.NewRequest("DELETE", "/artists/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		mock.ExpectExec("DELETE FROM artists WHERE id = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		result := DeleteArtistByID(rr, req, 1)

		if !result {
			t.Errorf("DeleteArtistByID returned false, expected true")
		}

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
