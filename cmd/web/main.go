package main

import (
	"html/template"
	"sync"
	"time"

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
}

type application struct {
	env           envConfig
	shortID       uint
	templateCache map[string]*template.Template
	wg            sync.WaitGroup
	mutex         sync.RWMutex
}

// func main() {
// 	fs := http.FileServer(http.Dir("./static"))
// 	http.Handle("/static/", http.StripPrefix("/static/", fs))

// 	log.Print("Listening on :4000...")
// 	err := http.ListenAndServe(":4000", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	// 加载配置
	envConf := setupConfig()

	// 设置日志
	setupLogger()

	// 连接、迁移 db
	db := setupDB()

	conn, err := db.DB()
	fatalOnError(err, "Get *sql.DB failed")

	// 设置最大连接数、空闲数，重用时间
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)
	defer conn.Close()

	// 执行迁移
	err = db.AutoMigrate(&models.User{}, &models.Link{})
	fatalOnError(err, "db.AutoMigrate failed")

	// 获取link表最后一条记录的short_hash值对应数值，用于生成后续短网址
	shortID := setupShortID(db)

	var app = application{
		env:     envConf,
		shortID: shortID,
		wg:      sync.WaitGroup{},
		mutex:   sync.RWMutex{},
	}

	app.genTemplateCache()

	err = app.serve()
	fatalOnError(err, "app.serve failed")
}
