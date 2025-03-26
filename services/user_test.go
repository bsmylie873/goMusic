package services_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"goMusic/db"
	"goMusic/models"
	"goMusic/services"
	"goMusic/viewModels"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupMockDB(t *testing.T) sqlmock.Sqlmock {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}

	db.DB = mockDB
	return mock
}

func TestRegisterUser(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		mock := setupMockDB(t)
		defer db.DB.Close()

		mock.ExpectExec("INSERT INTO users").
			WithArgs("testuser", sqlmock.AnyArg(), "test@example.com").
			WillReturnResult(sqlmock.NewResult(1, 1))

		reqBody := viewModels.RegisterRequest{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		services.RegisterUser(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response viewModels.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, 1, response.User.Id)
		assert.Equal(t, "testuser", response.User.Username)
	})

	t.Run("invalid request body", func(t *testing.T) {
		defer db.DB.Close()

		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		services.RegisterUser(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("username already exists", func(t *testing.T) {
		mock := setupMockDB(t)
		defer db.DB.Close()

		mock.ExpectExec("INSERT INTO users").
			WithArgs("existinguser", sqlmock.AnyArg(), "test@example.com").
			WillReturnError(sql.ErrNoRows)

		reqBody := viewModels.RegisterRequest{
			Username: "existinguser",
			Password: "password123",
			Email:    "test@example.com",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		services.RegisterUser(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		mock := setupMockDB(t)
		defer db.DB.Close()

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		rows := sqlmock.NewRows([]string{"id", "username", "password", "email"}).
			AddRow(1, "testuser", string(hashedPassword), "test@example.com")

		mock.ExpectQuery("SELECT id, username, password, email FROM users WHERE username = ?").
			WithArgs("testuser").
			WillReturnRows(rows)

		reqBody := viewModels.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		services.LoginUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response viewModels.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "testuser", response.User.Username)
	})

	t.Run("user not found", func(t *testing.T) {
		mock := setupMockDB(t)
		defer db.DB.Close()

		mock.ExpectQuery("SELECT id, username, password, email FROM users WHERE username = ?").
			WithArgs("nonexistentuser").
			WillReturnError(sql.ErrNoRows)

		reqBody := viewModels.LoginRequest{
			Username: "nonexistentuser",
			Password: "password123",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		services.LoginUser(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid password", func(t *testing.T) {
		mock := setupMockDB(t)
		defer db.DB.Close()

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

		rows := sqlmock.NewRows([]string{"id", "username", "password", "email"}).
			AddRow(1, "testuser", string(hashedPassword), "test@example.com")

		mock.ExpectQuery("SELECT id, username, password, email FROM users WHERE username = ?").
			WithArgs("testuser").
			WillReturnRows(rows)

		reqBody := viewModels.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		services.LoginUser(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGetProfile(t *testing.T) {
	t.Run("successfully get profile", func(t *testing.T) {
		mock := setupMockDB(t)
		defer db.DB.Close()

		rows := sqlmock.NewRows([]string{"id", "username", "email"}).
			AddRow(1, "testuser", "test@example.com")

		mock.ExpectQuery("SELECT id, username, email FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/profile", nil)
		ctx := context.WithValue(req.Context(), "userID", 1)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		services.GetProfile(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var user models.User
		json.Unmarshal(w.Body.Bytes(), &user)
		assert.Equal(t, 1, user.Id)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("user not found", func(t *testing.T) {
		mock := setupMockDB(t)
		defer db.DB.Close()

		mock.ExpectQuery("SELECT id, username, email FROM users WHERE id = ?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest("GET", "/profile", nil)
		ctx := context.WithValue(req.Context(), "userID", 999)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		services.GetProfile(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/profile", nil)

		w := httptest.NewRecorder()

		services.GetProfile(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
