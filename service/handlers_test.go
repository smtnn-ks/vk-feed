package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"vk-feed/types"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type mockDeps struct{}

func (m mockDeps) createUser(name, password string) (types.User, error) {
	return types.User{Id: 1, Name: name}, nil
}

var m mockDeps
var valid *validator.Validate = validator.New()

func newRequest(methon, path string, body any) *http.Request {
	content, _ := json.Marshal(body)
	var b bytes.Buffer
	b.Write(content)
	return httptest.NewRequest("POST", "/signup", &b)
}

func TestNewSignupHandler(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		req := newRequest("POST", "/signup", types.SignDto{Name: "mock_name", Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSignupHandler(m, valid)(rr, req)
		assert.Equal(t, 201, rr.Code)
		var user types.User
		assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &user))
		assert.Equal(t, types.User{Id: 1, Name: "mock_name"}, user)
	})
	t.Run("No body provided", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup", nil)
		rr := httptest.NewRecorder()
		newSignupHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("No name provided", func(t *testing.T) {
		req := newRequest("POST", "/signup", types.SignDto{Name: "mock_name"})
		rr := httptest.NewRecorder()
		newSignupHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("No password provided", func(t *testing.T) {
		req := newRequest("POST", "/signup", types.SignDto{Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSignupHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("name too short", func(t *testing.T) {
		req := newRequest("POST", "/signup", types.SignDto{Name: "aaaaaaaaaaaaaaaaaaaa", Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSignupHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("name too long", func(t *testing.T) {
		req := newRequest("POST", "/signup", types.SignDto{Name: "a", Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSignupHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("password too short", func(t *testing.T) {
		req := newRequest("POST", "/signup", types.SignDto{Name: "mock_name", Password: "a"})
		rr := httptest.NewRecorder()
		newSignupHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("password too long", func(t *testing.T) {
		req := newRequest("POST", "/signup", types.SignDto{Name: "mock_name", Password: "aaaaaaaaaaaaaaaaaaaa"})
		rr := httptest.NewRecorder()
		newSignupHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
}
