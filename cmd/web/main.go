package main

import (
	"sync"

	"github.com/lightsaid/short-net/models"
)

type envConfig struct {
	DBPort         int    `mapstruct:"DB_PORT"`
	DBName         string `mapstruct:"DB_NAME"`
	DBPassword     string `mapstruct:"DB_PASSWORD"`
	HTTPServerPort int    `mapstruct:"HTTP_SERVER_PORT"`
	RunMode        string `mapstruct:"RUN_MODE"`
	MySQLLog       string `mapstruct:"MYSQL_LOG"`
	AccessLog      string `mapstruct:"ACCESS_LOG"`
	ViewPath       string `mapstruct:"VIEW_PATH"`
	PublicPath     string `mapstruct:"PUBLIC_PATH"`
}

type application struct {
	env     envConfig
	shortID uint
	wg      sync.WaitGroup
	mutex   sync.RWMutex
}

func main() {
	// 加载配置
	envConf := setupConfig()

	// 设置日志
	setupLogger()

	// 连接、迁移 db
	db := setupDB()
	err := db.AutoMigrate(&models.User{}, &models.Link{})
	fatalOnError(err, "db.AutoMigrate failed")

	// 获取link表最后一条记录的short_hash值对应数值，用于生成后续短网址
	shortID := setupShortID(db)

	var app = application{
		env:     envConf,
		shortID: shortID,
		wg:      sync.WaitGroup{},
		mutex:   sync.RWMutex{},
	}

	err = app.serve()
	fatalOnError(err, "app.serve failed")
}
