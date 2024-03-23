package service

import (
	"testing"
	"vk-feed/types"

	"github.com/stretchr/testify/assert"
)

type mockDBConnection struct{}

func (m mockDBConnection) CreateUser(name string, password string) (id int, err error) {
	return 1, nil
}

func TestCreateUser(t *testing.T) {
	d := deps{client: mockDBConnection{}}
	user, err := d.createUser("mock_username", "mock_password")
	assert.Equal(t, user, types.User{Id: 1, Name: "mock_username"})
	assert.Equal(t, err, nil)
}
