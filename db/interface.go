package db

import "vk-feed/types"

type DBConnection interface {
	CreateUser(name, password string) (int, error)
	GetUserByName(name string) (int, string, error)
	CreateAd(dto types.AdDto, userId int) (int, error)
	GetAds(userId int, params types.GetAdParams) ([]types.AdFeed, error)
}
