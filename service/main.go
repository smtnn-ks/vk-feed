package service

import (
	"net/http"
	"vk-feed/db"

	"github.com/go-playground/validator/v10"
)

type deps struct {
	client    db.DBConnection
	validator *validator.Validate
}

func Register(conn db.DBConnection) {
	d := deps{
		client:    conn,
		validator: validator.New(),
	}
	http.HandleFunc("POST /signup", newSignupHandler(d))
	// http.HandleFunc("POST /signin", func(w http.ResponseWriter, r *http.Request) {})
	// http.HandleFunc("POST /ads", func(w http.ResponseWriter, r *http.Request) {})
	// http.HandleFunc("GET /ads", func(w http.ResponseWriter, r *http.Request) {})
}
