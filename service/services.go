package service

import (
	"crypto/sha512"
	"encoding/base64"
	"vk-feed/types"
)

func (d deps) createUser(name, password string) (types.User, error) {
	temp := sha512.Sum512([]byte(name))
	hashPassword := base64.StdEncoding.EncodeToString(temp[:])
	id, err := d.client.CreateUser(name, hashPassword)
	if err != nil {
		return types.User{}, err
	}
	return types.User{Id: id, Name: name}, nil
}
