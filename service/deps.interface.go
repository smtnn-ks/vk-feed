package service

import (
	"vk-feed/types"
)

type dependencies interface {
	createUser(name, password string) (types.User, error)
	signIn(name, password string) (types.Token, error)
	createAd(dto types.AdDto, userId int) (types.Ad, error)
}
