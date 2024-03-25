package service

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"
	"vk-feed/types"

	"github.com/golang-jwt/jwt/v5"
)

var ErrWrongCreds error = errors.New("wrong credentials")

func (d deps) createUser(name, password string) (types.User, error) {
	temp := sha512.Sum512([]byte(password))
	hashPassword := base64.StdEncoding.EncodeToString(temp[:])
	id, err := d.client.CreateUser(name, hashPassword)
	if err != nil {
		return types.User{}, err
	}
	return types.User{Id: id, Name: name}, nil
}

func (d deps) signIn(name, password string) (types.Token, error) {
	id, pass, err := d.client.GetUserByName(name)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Token{}, ErrWrongCreds
		}
		return types.Token{}, err
	}
	temp := sha512.Sum512([]byte(password))
	hashPassword := base64.StdEncoding.EncodeToString(temp[:])
	if pass != hashPassword {
		return types.Token{}, ErrWrongCreds
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().UTC().Add(time.Hour * 24).Unix(),
	}).SignedString(d.jwtSecret)
	if err != nil {
		return types.Token{}, err
	}
	return types.Token{Token: token}, nil
}

func (d deps) createAd(dto types.AdDto, userId int) (types.Ad, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelCtx()
	if err := d.ic.Check(ctx, dto.ImageUrl); err != nil {
		return types.Ad{}, err
	}
	id, err := d.client.CreateAd(dto, userId)
	if err != nil {
		return types.Ad{}, err
	}
	out := types.Ad{
		Id:       id,
		Title:    dto.Title,
		Content:  dto.Content,
		ImageUrl: dto.ImageUrl,
		Price:    dto.Price,
	}
	return out, nil
}
