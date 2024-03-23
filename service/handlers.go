package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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
