package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var d deps = deps{jwtSecret: []byte("some-jwt-secret")}

func mockFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuthMiddleware(t *testing.T) {
	// generating jwt token to parse it in tests
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1",
		"exp": time.Now().UTC().Add(time.Hour * 24).Unix(),
	}).SignedString(d.jwtSecret)
	wrongToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1",
		"exp": time.Now().UTC().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte("wrong-secret"))
	expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1",
		"exp": time.Now().UTC().Add(time.Hour * -1).Unix(),
	}).SignedString(d.jwtSecret)

	t.Run("OK if mandatory", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		authMiddleware(d, mockFunc, false)(rr, req)
		assert.Equal(t, 200, rr.Code)
		assert.Equal(t, "1", req.Header.Get("userid"))
	})
	t.Run("OK if optional", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		authMiddleware(d, mockFunc, true)(rr, req)
		assert.Equal(t, 200, rr.Code)
		assert.Equal(t, "1", req.Header.Get("userid"))
	})
	t.Run("Wrong token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup", nil)
		req.Header.Set("Authorization", "Bearer "+wrongToken)
		rr := httptest.NewRecorder()
		authMiddleware(d, mockFunc, false)(rr, req)
		assert.Equal(t, 401, rr.Code)
	})
	t.Run("Wrong token, but it's optional so OK", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup", nil)
		req.Header.Set("Authorization", "Bearer "+wrongToken)
		rr := httptest.NewRecorder()
		authMiddleware(d, mockFunc, true)(rr, req)
		assert.Equal(t, 200, rr.Code)
		assert.Equal(t, "", req.Header.Get("userid"))
	})
	t.Run("Expired token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		rr := httptest.NewRecorder()
		authMiddleware(d, mockFunc, false)(rr, req)
		assert.Equal(t, 401, rr.Code)
	})
	t.Run("Expired token, but it's optional so OK", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		rr := httptest.NewRecorder()
		authMiddleware(d, mockFunc, true)(rr, req)
		assert.Equal(t, 200, rr.Code)
		assert.Equal(t, "", req.Header.Get("userid"))
	})
}
