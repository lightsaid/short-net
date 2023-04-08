package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lightsaid/short-net/models"
	"golang.org/x/exp/slog"
)

func main() {
	var serverPort = os.Getenv("HTTP_SERVER_PORT")

	// db, err := openDB(os.Getenv("DB_SOURCE"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// 设置日志
	setupLogger()

	db := setupDB()
	var user models.User
	db.First(&user)

	err := db.AutoMigrate(&models.User{}, &models.Link{})
	if err != nil {
		log.Println("ERROR  ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})

	slog.Info("starting server on ", "port", serverPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", serverPort), mux)
	if err != nil {
		log.Println("start error: ", err)
	}
}
