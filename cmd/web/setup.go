package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"golang.org/x/exp/slog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func fatalOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("msg: %s", err.Error())
	}
}

func setupLogger() {
	var logfile = os.Getenv("ACCESS_LOG")
	if logfile == "" {
		logfile = "./storage/log/access.log"
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(&lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})))
}

func setupDB() *gorm.DB {
	var times = 5
	var db *gorm.DB
	var err error
	mylog, colorful := setGormLog()
	for i := 1; i <= times; i++ {
		db, err = connectDB(mylog, colorful)
		if err != nil {
			fmt.Printf("第 %d 次尝试连接数据库失败\n", i)
			time.Sleep(2 * time.Second)
		} else {
			return db
		}
	}

	fatalOnError(err, "connect mysql failed")
	return db
}

func setGormLog() (*log.Logger, bool) {
	var colorful = true
	var logfile = os.Getenv("MYSQL_LOG")
	if logfile == "" {
		logfile = "./storage/log/db.log"
	}
	mode := os.Getenv("RUN_MODE")
	mylog := log.New(os.Stdout, "\r\n", log.LstdFlags)
	if mode == "prod" {
		colorful = true
		mylog = log.New(&lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    1, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}, "\r\n", log.LstdFlags)
	}

	return mylog, colorful
}

func connectDB(mylog *log.Logger, colorful bool) (*gorm.DB, error) {
	logConfig := logger.New(
		// io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		// log.New(os.Stdout, "\r\n", log.LstdFlags),
		mylog,
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  colorful,    // 彩色打印
		},
	)

	return gorm.Open(mysql.Open(os.Getenv("DB_SOURCE")), &gorm.Config{Logger: logConfig})
}
