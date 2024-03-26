package service

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"
	imgC "vk-feed/image-checker"
	"vk-feed/types"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
)

func newSignupHandler(d dependencies, valid *validator.Validate) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var dto types.SignDto
		if err := json.Unmarshal(content, &dto); err != nil {
			if typeError, ok := err.(*json.UnmarshalTypeError); ok {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(typeError.Error()))
				return
			}
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := valid.Struct(dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		user, err := d.createUser(dto.Name, dto.Password)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload, err := json.Marshal(user)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}

func newSigninHandler(d dependencies, valid *validator.Validate) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var dto types.SignDto
		if err := json.Unmarshal(content, &dto); err != nil {
			if typeError, ok := err.(*json.UnmarshalTypeError); ok {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(typeError.Error()))
				return
			}
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := valid.Struct(dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		token, err := d.signIn(dto.Name, dto.Password)
		if err != nil {
			if err == ErrWrongCreds {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload, err := json.Marshal(token)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}

func newCreateAdHandler(d dependencies, valid *validator.Validate) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			log.Error("Content-Length is 0")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var dto types.AdDto
		if err := json.Unmarshal(content, &dto); err != nil {
			if typeError, ok := err.(*json.UnmarshalTypeError); ok {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(typeError.Error()))
				return
			}
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := valid.Struct(dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		userIdString := r.Header.Get("userid")
		if userIdString == "" {
			log.Error("UserId is not provided, yet fell into handler")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userId, err := strconv.Atoi(userIdString)
		if err != nil {
			log.Error("userId is not of type int, yet fell into handler")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ad, err := d.createAd(dto, userId)
		if err != nil {
			if slices.Contains([]error{imgC.ErrNotImage, imgC.ErrUrlUnavailable, imgC.ErrImageTooBig}, err) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload, err := json.Marshal(ad)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}

func newGetAdsHanlder(d dependencies, _ *validator.Validate) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var params types.GetAdParams
		if sortByStr := r.PathValue("sort_by"); sortByStr != string(types.SORT_BY_PRICE) {
			params.SortBy = types.SORT_BY_DATE
		} else {
			params.SortBy = types.SORT_BY_PRICE
		}
		if orderByStr := r.PathValue("order_by"); orderByStr != string(types.ORDER_BY_DESC) {
			params.OrderBy = types.ORDER_BY_ASC
		} else {
			params.OrderBy = types.ORDER_BY_DESC
		}
		maxPriceStr := r.PathValue("max_price")
		if maxPrice, err := strconv.Atoi(maxPriceStr); err != nil {
			params.MaxPrice = 1e6
		} else {
			if maxPrice < 1 {
				maxPrice = 1
			} else if maxPrice > 1e6 {
				maxPrice = 1e6
			}
			params.MaxPrice = maxPrice
		}
		minPriceStr := r.PathValue("min_price")
		if minPrice, err := strconv.Atoi(minPriceStr); err != nil {
			params.MinPrice = 1
		} else {
			if minPrice < 1 {
				minPrice = 1
			} else if minPrice > 1e6 {
				minPrice = 1e6
			}
			params.MinPrice = minPrice
		}
		pageStr := r.PathValue("page")
		if page, err := strconv.Atoi(pageStr); err != nil {
			params.Page = 0
		} else {
			if page < 0 {
				page = 0
			}
			params.Page = page
		}
		userIdStr := r.Header.Get("userid")
		var userId int
		if userIdStr != "" {
			var err error
			userId, err = strconv.Atoi(userIdStr)
			if err != nil {
				log.Error("userid is not int, yet fell into handler")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		feed, err := d.getAds(userId, params)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload, err := json.Marshal(feed)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}
