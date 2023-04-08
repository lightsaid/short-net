package dbrepo

import (
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testRepo Repository

// 宿主机访问 mysql dsn
const dsn = "root:abc123@tcp(127.0.0.1:3307)/shortnet?charset=utf8mb4&parseTime=True&loc=Local"

func TestMain(m *testing.M) {
	db, err := connectDB()
	if err != nil {
		log.Fatal("connect database failed: ", err)
	}
	testRepo = NewRepository(db)

	os.Exit(m.Run())
}

func connectDB() (*gorm.DB, error) {
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
