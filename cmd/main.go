package main

import (
	"net/http"
	"os"
	"vk-feed/db"
	"vk-feed/service"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Warn("PORT is not specified. Default port of 6969 will be used.")
		port = "6969"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not specified")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not specified")
	}
	dbConn, err := db.Init(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("database connected")
	defer dbConn.Client.Close()

	service.Register(dbConn, []byte(jwtSecret))

	log.Infof("server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
