package service

import (
	"vk-feed/types"
)

type dependencies interface {
	createUser(name, password string) (types.User, error)
}
