package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	var serverPort = os.Getenv("HTTP_SERVER_PORT")

	fmt.Println(">>>>>> ", serverPort)
	fmt.Println(">>>>>> ", os.Getenv("VIEW_PATH"))
	fmt.Println(">>>>>> ", os.Getenv("PUBLIC_PATH"))
	fmt.Println(">>>>>> ", os.Getenv("DB_SOURCE"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})

	log.Println("starting server on :", serverPort)
	err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), mux)
	if err != nil {
		log.Println("start error: ", err)
	}
}

// type TmpTest struct {
// 	gorm.Model
// 	Name string `json:"name"`
// }

// func main() {
// 	mux := setupRoute()
// 	var err error
// 	fmt.Println(">>>>>>>>>> ", os.Getenv("VIEW_PATH"))
// 	fmt.Println(">>>>>>>>>> ", os.Getenv("PUBLIC_PATH"))

// 	dsn := os.Getenv("DB_SOURCE") //"root:abc123@tcp(mysqldb:3306)/shortnet?charset=utf8mb4&parseTime=True&loc=Local"
// 	fmt.Println("dsn: ", dsn)
// 	db, err := openDB(dsn)
// 	if err != nil {
// 		log.Println("ERROR openDB error1 ", err)
// 		return
// 	}
// 	err = db.AutoMigrate(&TmpTest{})
// 	if err != nil {
// 		log.Println("ERROR AutoMigrate error ", err)
// 		return
// 	}

// 	err = db.Create(&TmpTest{Name: "zhangsan"}).Error
// 	if err != nil {
// 		log.Println(">>> Create: ", err)
// 		return
// 	}

// 	var tmp TmpTest
// 	err = db.First(&tmp).Error
// 	if err != nil {
// 		log.Println(">>> First: ", err)
// 	}

// 	// var tmp TmpTest
// 	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
// 		json.NewEncoder(w).Encode(&tmp)
// 	})

// 	log.Println("starting server on :4000")
// 	err = http.ListenAndServe(":4000", mux)
// 	if err != nil {
// 		log.Println("start error:  ", err)
// 	}
// }

// func openDB(dsn string) (*gorm.DB, error) {
// 	logConfig := logger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
// 		logger.Config{
// 			SlowThreshold:             time.Second, // 慢 SQL 阈值
// 			LogLevel:                  logger.Info, // 日志级别
// 			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
// 			Colorful:                  true,        // 彩色打印
// 		},
// 	)
// 	return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logConfig})
// }
