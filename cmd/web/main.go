package main

import (
	"os"
	"sync"

	"github.com/lightsaid/short-net/models"
)

type application struct {
	shortID uint
	wg      sync.WaitGroup
}

func main() {
	var serverPort = os.Getenv("HTTP_SERVER_PORT")

	// 设置日志
	setupLogger()

	// 连接、迁移 db
	db := setupDB()
	err := db.AutoMigrate(&models.User{}, &models.Link{})
	fatalOnError(err, "db.AutoMigrate failed")

	var app = application{

		wg: sync.WaitGroup{},
	}

	err = app.serve()
	fatalOnError(err, "app.serve failed")
}
