package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lightsaid/short-net/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	log.Println("starting server on :", serverPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", serverPort), mux)
	if err != nil {
		log.Println("start error: ", err)
	}
}

func openDB(dsn string) (*gorm.DB, error) {
	logConfig := logger.New(
		// io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logConfig})
}
