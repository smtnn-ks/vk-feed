package service

import (
	"net/http"
	"vk-feed/db"
	imgC "vk-feed/image-checker"

	"github.com/go-playground/validator/v10"
)

type deps struct {
	client    db.DBConnection
	jwtSecret []byte
	ic        imgC.ImageChecker
}

func Register(conn db.DBConnection, jwtSecret []byte) {
	d := deps{
		client:    conn,
		jwtSecret: jwtSecret,
		ic:        imgC.IC{},
	}
	valid := validator.New()
	http.HandleFunc("POST /signup", newSignupHandler(d, valid))
	http.HandleFunc("POST /signin", newSigninHandler(d, valid))
	http.HandleFunc("POST /ads", authMiddleware(d, newCreateAdHandler(d, valid), false))
	// http.HandleFunc("GET /ads", func(w http.ResponseWriter, r *http.Request) {})
}
