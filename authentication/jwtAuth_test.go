package authentication

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	const testSecret = "test-secret-key"
	os.Setenv("JWT_SECRET_KEY", testSecret)
	defer os.Unsetenv("JWT_SECRET_KEY")

	userID := 1

	token, err := GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	t.Logf("Generated token: %s", token)

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(testSecret), nil
	})

	if err != nil {
		t.Logf("Error parsing token: %v", err)
	}

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, userID, claims.UserID)

	expiresAt := claims.ExpiresAt.Time
	expectedExpiry := time.Now().Add(24 * time.Hour)
	timeDiff := expectedExpiry.Sub(expiresAt)
	assert.Less(t, timeDiff.Abs(), 5*time.Minute) // Allow 5 minute tolerance
}

func TestAuthMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET_KEY")

	t.Run("Valid token", func(t *testing.T) {
		userID := 456
		token, err := GenerateToken(userID)
		assert.NoError(t, err)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxUserID := r.Context().Value("userID")
			assert.Equal(t, userID, ctxUserID)
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		middleware := AuthMiddleware(nextHandler)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		rr := httptest.NewRecorder()

		middleware := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should not be called with invalid token")
		}))
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Missing token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		middleware := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should not be called without token")
		}))
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Expired token", func(t *testing.T) {
		expirationTime := time.Now().Add(-1 * time.Hour) // 1 hour in the past
		claims := &Claims{
			UserID: 789,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rr := httptest.NewRecorder()

		middleware := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should not be called with expired token")
		})
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Token without Bearer prefix", func(t *testing.T) {
		userID := 123
		token, err := GenerateToken(userID)
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", token)
		rr := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxUserID := r.Context().Value("userID")
			assert.Equal(t, userID, ctxUserID)
			w.WriteHeader(http.StatusOK)
		})

		middleware := AuthMiddleware(nextHandler)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
