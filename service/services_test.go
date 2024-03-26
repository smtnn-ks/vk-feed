package service

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"testing"
	"time"
	imgC "vk-feed/image-checker"
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

func (m mockDBConnection) CreateAd(dto types.AdDto, userId int) (id int, err error) {
	if userId == 0 {
		return 0, sql.ErrNoRows
	} else {
		return 1, nil
	}
}

func (m mockDBConnection) GetAds(userId int, params types.GetAdParams) ([]types.AdFeed, error) {
	return []types.AdFeed{
		{
			Id:        1,
			Title:     "mock_title",
			Content:   "mock_content",
			ImageUrl:  "http://mocksite.com/image.jpg",
			Price:     6969,
			CreatedAt: time.Now(),
			AuthorId:  1,
			IsYours:   false,
		},
	}, nil
}

type mockIC struct{}

func (m mockIC) Check(ctx context.Context, url string) error {
	if url == "OK" {
		return nil
	} else {
		return imgC.ErrBadImage
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

func TestCreateAd(t *testing.T) {
	d := deps{client: mockDBConnection{}, ic: mockIC{}}
	t.Run("OK", func(t *testing.T) {
		dto := types.AdDto{
			Title:    "mock_title",
			Content:  "mock_content",
			ImageUrl: "OK",
			Price:    6969,
		}
		resAd := types.Ad{
			Id:       1,
			Title:    dto.Title,
			Content:  dto.Content,
			ImageUrl: dto.ImageUrl,
			Price:    dto.Price,
		}
		ad, err := d.createAd(dto, 1)
		assert.NoError(t, err)
		assert.Equal(t, resAd, ad)
	})
	t.Run("Bad image", func(t *testing.T) {
		dto := types.AdDto{
			Title:    "mock_title",
			Content:  "mock_content",
			ImageUrl: "NOT OK",
			Price:    6969,
		}
		_, err := d.createAd(dto, 1)
		assert.Equal(t, imgC.ErrBadImage, err)
	})
	t.Run("Bad user ID", func(t *testing.T) {
		dto := types.AdDto{
			Title:    "mock_title",
			Content:  "mock_content",
			ImageUrl: "OK",
			Price:    6969,
		}
		_, err := d.createAd(dto, 0)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

// getAds does not do much, no tests needed
