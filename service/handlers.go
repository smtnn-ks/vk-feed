package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	imgC "vk-feed/image-checker"
	"vk-feed/types"

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
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var dto types.SignDto
		if err := json.Unmarshal(content, &dto); err != nil {
			log.Println(err)
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
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
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
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var dto types.SignDto
		if err := json.Unmarshal(content, &dto); err != nil {
			log.Println(err)
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
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload, err := json.Marshal(token)
		if err != nil {
			log.Println(err)
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
			log.Println("Content-Length is 0")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var dto types.AdDto
		if err := json.Unmarshal(content, &dto); err != nil {
			log.Println(err)
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
			log.Println("UserId is not provided, yet fell into handler")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userId, err := strconv.Atoi(userIdString)
		if err != nil {
			log.Println("userId is not of type int, yet fell into handler")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ad, err := d.createAd(dto, userId)
		if err != nil {
			if err == imgC.ErrBadImage {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("something wrong with the image"))
				return
			}
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload, err := json.Marshal(ad)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}
