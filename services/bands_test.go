package services

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"goMusic/db"
	"goMusic/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetBands(t *testing.T) {
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

	rows := sqlmock.NewRows([]string{"id", "name", "nationality", "number_of_members", "date_formed", "age", "active"}).
		AddRow(1, "Coldplay", "British", 4, "1996-01-01", 27, true)
	mock.ExpectQuery("SELECT id, name, nationality, number_of_members, date_formed, age, active FROM bands").WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/bands", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	GetBands(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := []models.Band{
		{Id: 1, Name: "Coldplay", Nationality: "British", NumberOfMembers: 4, DateFormed: "1996-01-01", Age: 27, Active: true},
	}
	var actual []models.Band
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestGetBandByID(t *testing.T) {
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

	t.Run("Successful band retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "nationality", "number_of_members", "date_formed", "age", "active"}).
			AddRow(1, "Coldplay", "British", 4, "1996-01-01", 27, true)
		mock.ExpectQuery("SELECT id, name, nationality, number_of_members, date_formed, age, active FROM bands WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/bands/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		GetBandByID(rr, req, 1)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var band models.Band
		if err := json.NewDecoder(rr.Body).Decode(&band); err != nil {
			t.Fatal(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Band not found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "nationality", "number_of_members", "date_formed", "age", "active"})
		mock.ExpectQuery("SELECT id, name, nationality, number_of_members, date_formed, age, active FROM bands WHERE id = ?").
			WithArgs(999).
			WillReturnRows(rows)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/bands/999", nil)
		if err != nil {
			t.Fatal(err)
		}

		GetBandByID(rr, req, 999)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestPostBand(t *testing.T) {
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

	band := models.Band{
		Id:              1,
		Name:            "Coldplay",
		Nationality:     "British",
		NumberOfMembers: 4,
		DateFormed:      "1996-01-01",
		Age:             27,
		Active:          true,
	}

	bandJSON, _ := json.Marshal(band)
	req, err := http.NewRequest("POST", "/bands", bytes.NewBuffer(bandJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mock.ExpectExec("INSERT INTO bands").
		WithArgs(band.Id, band.Name, band.Nationality, band.NumberOfMembers, band.DateFormed, band.Age, band.Active).
		WillReturnResult(sqlmock.NewResult(1, 1))

	PostBand(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateBandByID(t *testing.T) {
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

func TestDeleteBandByID(t *testing.T) {
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

		mock.ExpectExec("DELETE FROM bands WHERE id = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		result := DeleteBandByID(rr, req, 1)

		if !result {
			t.Errorf("DeleteBandByID returned false, expected true")
		}

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
