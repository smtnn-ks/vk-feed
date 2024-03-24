package service

import (
	"net/http"
	"vk-feed/db"

	"github.com/go-playground/validator/v10"
)

type deps struct {
	client    db.DBConnection
	jwtSecret []byte
}

func Register(conn db.DBConnection, jwtSecret []byte) {
	d := deps{
		client:    conn,
		jwtSecret: jwtSecret,
	}
	valid := validator.New()
	http.HandleFunc("POST /signup", newSignupHandler(d, valid))
	http.HandleFunc("POST /signin", newSigninHandler(d, valid))
	// http.HandleFunc("POST /ads", func(w http.ResponseWriter, r *http.Request) {})
	// http.HandleFunc("GET /ads", func(w http.ResponseWriter, r *http.Request) {})
}
