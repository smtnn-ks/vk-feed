package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"vk-feed/db"
	"vk-feed/service"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT is not specified. Default port of 6969 will be used.")
		port = "6969"
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not specified")
	}
	dbConn, err := db.Init(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("database connected")
	defer dbConn.Client.Close()

	// http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello, World!"))
	// })

	service.Register(dbConn)

	fmt.Printf("server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
