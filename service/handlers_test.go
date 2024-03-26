package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	imgC "vk-feed/image-checker"
	"vk-feed/types"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type mockDeps struct{}

func (m mockDeps) createUser(name, password string) (types.User, error) {
	return types.User{Id: 1, Name: name}, nil
}

func (m mockDeps) signIn(name, password string) (types.Token, error) {
	if name == "mock_name" && password == "mock_password" {
		return types.Token{Token: "mock_token"}, nil
	} else {
		return types.Token{}, ErrWrongCreds
	}
}

func (m mockDeps) createAd(dto types.AdDto, userId int) (types.Ad, error) {
	if dto.ImageUrl != "http://mocksite.com/image.jpg" {
		return types.Ad{}, imgC.ErrBadImage
	} else if userId == 0 {
		return types.Ad{}, sql.ErrNoRows
	} else {
		return types.Ad{
			Id:       1,
			Title:    dto.Title,
			Content:  dto.Content,
			ImageUrl: dto.ImageUrl,
			Price:    dto.Price,
		}, nil
	}
}

var outParams types.GetAdParams

func (m mockDeps) getAds(userId int, params types.GetAdParams) ([]types.AdFeed, error) {
	outParams = params
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

var m mockDeps
var valid *validator.Validate = validator.New()

func newRequest(method, path string, body any) *http.Request {
	content, _ := json.Marshal(body)
	var b bytes.Buffer
	b.Write(content)
	return httptest.NewRequest(method, path, &b)
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
	t.Run("name too long", func(t *testing.T) {
		req := newRequest("POST", "/signup", types.SignDto{Name: "aaaaaaaaaaaaaaaaaaaa", Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSignupHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("name too short", func(t *testing.T) {
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

func TestNewSigninHandler(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		req := newRequest("POST", "/signin", types.SignDto{Name: "mock_name", Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 201, rr.Code)
		var token types.Token
		assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &token))
	})
	t.Run("No body provided", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signin", nil)
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("No name provided", func(t *testing.T) {
		req := newRequest("POST", "/signin", types.SignDto{Name: "mock_name"})
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("No password provided", func(t *testing.T) {
		req := newRequest("POST", "/signin", types.SignDto{Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("name too short", func(t *testing.T) {
		req := newRequest("POST", "/signin", types.SignDto{Name: "aaaaaaaaaaaaaaaaaaaa", Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("name too long", func(t *testing.T) {
		req := newRequest("POST", "/signin", types.SignDto{Name: "a", Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("password too short", func(t *testing.T) {
		req := newRequest("POST", "/signin", types.SignDto{Name: "mock_name", Password: "a"})
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("password too long", func(t *testing.T) {
		req := newRequest("POST", "/signin", types.SignDto{Name: "mock_name", Password: "aaaaaaaaaaaaaaaaaaaa"})
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("user not found", func(t *testing.T) {
		req := newRequest("POST", "/signin", types.SignDto{Name: "wrong_name", Password: "mock_password"})
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 404, rr.Code)
	})
	t.Run("wrong password", func(t *testing.T) {
		req := newRequest("POST", "/signin", types.SignDto{Name: "mock_name", Password: "wrong_password"})
		rr := httptest.NewRecorder()
		newSigninHandler(m, valid)(rr, req)
		assert.Equal(t, 404, rr.Code)
	})
}

func TestNewCreateAdHandler(t *testing.T) {
	mockImageUrl := "http://mocksite.com/image.jpg"
	t.Run("OK", func(t *testing.T) {
		dto := types.AdDto{
			Title:    "mock_title",
			Content:  "mock_content",
			ImageUrl: mockImageUrl,
			Price:    6969,
		}
		req := newRequest("POST", "/ads", dto)
		req.Header.Add("userid", "1")
		rr := httptest.NewRecorder()
		newCreateAdHandler(m, valid)(rr, req)
		ad := types.Ad{
			Id:       1,
			Title:    dto.Title,
			Content:  dto.Content,
			ImageUrl: dto.ImageUrl,
			Price:    dto.Price,
		}
		assert.Equal(t, 201, rr.Code)
		var out types.Ad
		assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &out))
		assert.Equal(t, ad, out)
	})
	t.Run("No body provided", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/ads", nil)
		req.Header.Add("userid", "1")
		rr := httptest.NewRecorder()
		newCreateAdHandler(m, valid)(rr, req)
		assert.Equal(t, 400, rr.Code)
	})
	t.Run("Partial DTO provided", func(t *testing.T) {
		cases := []struct {
			name string
			in   types.AdDto
		}{
			{
				name: "no title",
				in: types.AdDto{
					Content:  "mock_content",
					ImageUrl: mockImageUrl,
					Price:    6969,
				},
			},
			{
				name: "no content",
				in: types.AdDto{
					Title:    "mock_title",
					ImageUrl: mockImageUrl,
					Price:    6969,
				},
			},
			{
				name: "no image url",
				in: types.AdDto{
					Title:   "mock_title",
					Content: "mock_content",
					Price:   6969,
				},
			},
			{
				name: "no price",
				in: types.AdDto{
					Title:    "mock_title",
					Content:  "mock_content",
					ImageUrl: mockImageUrl,
				},
			},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				req := newRequest("POST", "/ads", c.in)
				req.Header.Add("userid", "1")
				rr := httptest.NewRecorder()
				newCreateAdHandler(m, valid)(rr, req)
				assert.Equal(t, 400, rr.Code, c.name)
			})
		}
	})
	t.Run("validation test", func(t *testing.T) {
		cases := []struct {
			name string
			in   types.AdDto
		}{
			{
				name: "title too short",
				in: types.AdDto{
					Title:    "a",
					Content:  "mock_content",
					ImageUrl: mockImageUrl,
					Price:    6969,
				},
			},
			{
				name: "title too long",
				in: types.AdDto{
					Title:    strings.Repeat("a", 256),
					Content:  "mock_content",
					ImageUrl: mockImageUrl,
					Price:    6969,
				},
			},
			{
				name: "content too short",
				in: types.AdDto{
					Title:    "mock_title",
					Content:  "a",
					ImageUrl: mockImageUrl,
					Price:    6969,
				},
			},
			{
				name: "content too long",
				in: types.AdDto{
					Title:    "mock_title",
					Content:  strings.Repeat("a", 1001),
					ImageUrl: mockImageUrl,
					Price:    6969,
				},
			},
			{
				name: "ImageUrl is not url",
				in: types.AdDto{
					Title:    "mock_title",
					Content:  "mock_content",
					ImageUrl: "a",
					Price:    6969,
				},
			},
			{
				name: "price too low",
				in: types.AdDto{
					Title:    "mock_title",
					Content:  "mock_content",
					ImageUrl: mockImageUrl,
					Price:    0,
				},
			},
			{
				name: "price too high",
				in: types.AdDto{
					Title:    "mock_title",
					Content:  "mock_content",
					ImageUrl: mockImageUrl,
					Price:    1e7,
				},
			},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				req := newRequest("POST", "/ads", c.in)
				req.Header.Add("userid", "1")
				rr := httptest.NewRecorder()
				newCreateAdHandler(m, valid)(rr, req)
				assert.Equal(t, 400, rr.Code, c.name)
			})
		}
	})
}

func TestNewGetAdsHandler(t *testing.T) {
	t.Run("no params", func(t *testing.T) {
		var m mockDeps
		req := httptest.NewRequest("GET", "/ads", nil)
		rr := httptest.NewRecorder()
		newGetAdsHanlder(m, valid)(rr, req)
		assert.Equal(t, 200, rr.Code)
		defaultParams := types.GetAdParams{
			Page:     0,
			MinPrice: 1,
			MaxPrice: 1e6,
			SortBy:   types.SORT_BY_DATE,
			OrderBy:  types.ORDER_BY_ASC,
		}
		assert.Equal(t, defaultParams, outParams)
	})
	type mockParams struct {
		page     string
		minPrice string
		maxPrice string
		sortBy   string
		orderBy  string
	}
	cases := []struct {
		name string
		in   mockParams
		out  types.GetAdParams
	}{
		{
			name: "correct params",
			in: mockParams{
				page:     "1",
				minPrice: "100",
				maxPrice: "10000",
				sortBy:   "price",
				orderBy:  "desc",
			},
			out: types.GetAdParams{
				Page:     1,
				MinPrice: 100,
				MaxPrice: 10000,
				SortBy:   types.SORT_BY_PRICE,
				OrderBy:  types.ORDER_BY_DESC,
			},
		},
		{
			name: "bad params",
			in: mockParams{
				page:     "-1",
				minPrice: "0",
				maxPrice: "100000000000000000000000000000000000",
				sortBy:   "foo",
				orderBy:  "bar",
			},
			out: types.GetAdParams{
				Page:     0,
				MinPrice: 1,
				MaxPrice: 1e6,
				SortBy:   types.SORT_BY_DATE,
				OrderBy:  types.ORDER_BY_ASC,
			},
		},
		{
			name: "bad types",
			in: mockParams{
				page:     "monke",
				minPrice: "monke",
				maxPrice: "monke",
				sortBy:   "monke",
				orderBy:  "monke",
			},
			out: types.GetAdParams{
				Page:     0,
				MinPrice: 1,
				MaxPrice: 1e6,
				SortBy:   types.SORT_BY_DATE,
				OrderBy:  types.ORDER_BY_ASC,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/ads", nil)
			req.SetPathValue("page", c.in.page)
			req.SetPathValue("min_price", c.in.minPrice)
			req.SetPathValue("max_price", c.in.maxPrice)
			req.SetPathValue("sort_by", c.in.sortBy)
			req.SetPathValue("order_by", c.in.orderBy)
			rr := httptest.NewRecorder()
			newGetAdsHanlder(m, valid)(rr, req)
			assert.Equal(t, 200, rr.Code)
			assert.Equal(t, c.out, outParams)
		})
	}
}
