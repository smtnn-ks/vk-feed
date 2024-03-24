package service

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"testing"
	"vk-feed/types"

	"github.com/stretchr/testify/assert"
)

type mockDBConnection struct{}

func (m mockDBConnection) CreateUser(name string, password string) (id int, err error) {
	return 1, nil
}

func (m mockDBConnection) GetUserByName(name string) (int, string, error) {
	if name == "mock_name" {
		temp := sha512.Sum512([]byte("mock_password"))
		hashPassword := base64.StdEncoding.EncodeToString(temp[:])
		return 1, hashPassword, nil
	} else {
		return 0, "", sql.ErrNoRows
	}
}

func TestCreateUser(t *testing.T) {
	d := deps{client: mockDBConnection{}, jwtSecret: []byte("mock_jwt_secret")}
	user, err := d.createUser("mock_username", "mock_password")
	assert.Equal(t, user, types.User{Id: 1, Name: "mock_username"})
	assert.Equal(t, err, nil)
}

func TestSignin(t *testing.T) {
	d := deps{client: mockDBConnection{}, jwtSecret: []byte("mock_jwt_secret")}
	t.Run("OK", func(t *testing.T) {
		_, err := d.signIn("mock_name", "mock_password")
		assert.NoError(t, err)
	})
	t.Run("Not found", func(t *testing.T) {
		_, err := d.signIn("wrong_name", "mock_password")
		assert.Equal(t, ErrWrongCreds, err)
	})
	t.Run("Wrong password", func(t *testing.T) {
		_, err := d.signIn("mock_name", "wrong_password")
		assert.Equal(t, ErrWrongCreds, err)
	})
}
